package httpx

import (
	"net"
	"net/http"
	"net/netip"
)

// trustProxyHeader gates the CF-Connecting-IP read in ClientIP. Default off:
// the header is attacker-controlled on any path that skips the Cloudflare
// edge, so trusting it is an explicit per-deployment decision (#182).
var trustProxyHeader bool

// SetTrustProxyHeader enables (or disables) trusting CF-Connecting-IP as the
// client address. Call it exactly once at startup, before serving - it is a
// plain bool with no locking, so mutating it under live traffic is a data
// race. Wired from TRUST_PROXY_HEADER in cmd/api.
func SetTrustProxyHeader(v bool) { trustProxyHeader = v }

// ClientIP returns the originating client address - the single source of truth
// for rate-limit keys (via ClientIPKey) and session audit records. When proxy
// trust is enabled (SetTrustProxyHeader) it prefers Cloudflare's
// CF-Connecting-IP header, set at the edge in both deployment models (D20);
// otherwise (dev, smoke, tests) it uses the host part of RemoteAddr,
// preserving the direct-connection behavior.
//
// SECURITY: the header is only honored behind the TRUST_PROXY_HEADER opt-in,
// set where the deployment guarantees every request traverses the Cloudflare
// edge. In the prod compose topology that guarantee is structural
// (defense-in-depth): compose.prod.yml publishes no api port; the only paths
// in are Cloudflare -> cloudflared -> caddy (direct browser /api calls,
// headers forwarded unchanged) and the SvelteKit server's form actions, which
// forward the same edge-set header on their server-to-server calls
// ($lib/server/clientip.ts). Both hops live on the private compose network.
// Anywhere that topology does not hold (make dev on a LAN, a published port,
// a tunnel misconfiguration), the default-off flag keeps the unforgeable
// RemoteAddr as the key, so per-request header spoofing cannot mint fresh
// rate-limit buckets. Same trust posture the requestID middleware documents
// for X-Request-Id. X-Forwarded-For is deliberately not parsed: multi-valued,
// spoof-prone, and unnecessary here.
func ClientIP(r *http.Request) string {
	if trustProxyHeader {
		if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
			return ip
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// ClientIPKey returns the per-client rate-limit bucket: the IPv4 address, or
// for IPv6 the routed /64 prefix. A v6 client typically controls its whole
// /64, so per-address budgets could be rotated around for free; masking to
// the prefix restores "one client, one budget" without touching v4 behavior.
// Audit records (sessions.ip) keep the exact address via ClientIP.
func ClientIPKey(r *http.Request) string {
	ip := ClientIP(r)
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		// Unparseable (only possible in dev-shaped setups); the raw string
		// still buckets consistently.
		return ip
	}
	addr = addr.Unmap() // 4-mapped-in-6 is IPv4 in disguise
	if addr.Is4() {
		return addr.String()
	}
	// Prefix cannot fail here: a post-Unmap non-4 address is 128-bit, and
	// Prefix strips any zone rather than erroring.
	prefix, _ := addr.Prefix(64)
	return prefix.String()
}
