package httpx

import (
	"net"
	"net/http"
)

// ClientIP returns the originating client address - the single source of truth
// for rate-limit keys and session audit records. It prefers Cloudflare's
// CF-Connecting-IP header, set at the edge in both deployment models (D20);
// when absent (dev, smoke, tests) it falls back to the host part of
// RemoteAddr, preserving the direct-connection behavior.
//
// SECURITY: trusting CF-Connecting-IP is correct only because the API is never
// directly exposed - compose.prod.yml publishes no api port; the only path in
// is Cloudflare -> cloudflared -> caddy, and Caddy forwards headers unchanged
// (docs/DEPLOYMENT.md). Same trust posture the requestID middleware documents
// for X-Request-Id. X-Forwarded-For is deliberately not parsed: multi-valued,
// spoof-prone, and unnecessary here.
func ClientIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
