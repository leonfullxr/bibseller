/**
 * Forward the Cloudflare-set client address on a form action's
 * server-to-server API call.
 *
 * Direct browser /api calls reach the Go API through Caddy with
 * CF-Connecting-IP intact, but form actions hop through this server
 * (browser -> caddy -> web -> api); without forwarding, the API keys its
 * per-IP rate limits and sessions.ip audit on the web container's address -
 * one shared budget for every user (#133). Deliberately used only by the
 * throttled auth actions (login, register, forgot, reset, verify/resend);
 * catalog loads are unmetered and don't need it.
 *
 * Trust: the inbound value was set by Cloudflare at the edge (D20), and both
 * hops (caddy -> web, web -> api) run on the private compose network, so it
 * cannot be client-forged in prod. In dev the header is absent, nothing is
 * forwarded, and the API falls back to RemoteAddr.
 */
export function clientIPHeader(request: Request): Record<string, string> {
	const ip = request.headers.get('cf-connecting-ip');
	return ip ? { 'CF-Connecting-IP': ip } : {};
}
