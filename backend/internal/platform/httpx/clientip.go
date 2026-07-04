package httpx

import (
	"net"
	"net/http"
	"net/netip"
)

// ClientIP returns the originating client address - the single source of truth
// for rate-limit keys (via ClientIPKey) and session audit records. It prefers
// Cloudflare's CF-Connecting-IP header, set at the edge in both deployment
// models (D20); when absent (dev, smoke, tests) it falls back to the host part
// of RemoteAddr, preserving the direct-connection behavior.
//
// SECURITY: trusting CF-Connecting-IP is correct only because the API is never
// directly exposed - compose.prod.yml publishes no api port; the only paths in
// are Cloudflare -> cloudflared -> caddy (direct browser /api calls, headers
// forwarded unchanged) and the SvelteKit server's form actions, which forward
// the same edge-set header on their server-to-server calls
// ($lib/server/clientip.ts). Both hops live on the private compose network, so the value
// cannot be client-forged in prod. Same trust posture the requestID middleware
// documents for X-Request-Id. X-Forwarded-For is deliberately not parsed:
// multi-valued, spoof-prone, and unnecessary here.
func ClientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
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
