# Tech notes — near-term considerations

Non-binding engineering notes to revisit (binding decisions live in
[CONTEXT.md](CONTEXT.md)). Mirrored in the tracking issue
[#21](https://github.com/leonfullxr/bibseller/issues/21).

## Auth resolution middleware

Session → user resolution is now a single middleware (`auth.ResolveUser`,
`backend/internal/auth/middleware.go`) applied to the `/api/v1` sub-mux, with
handlers reading `auth.UserFromContext`. It runs the lookup on every
cookie-bearing request, including the (cacheable) catalog. That's fine at v1
scale; if it shows up in profiling, scope resolution to the routes that need
it. The trigger to reconsider the whole approach is many protected routes
landing in M4/M5.

## Session staleness across tabs

`locals.user` is refreshed per request via `GET /auth/me`, and a profile
rename syncs `locals.user` within the same action so the nav updates
immediately. But other already-open tabs keep their stale `user` until their
next navigation. Acceptable for v1; revisit if a session-refresh hook exists.

## Rate limiting is in-process and per-instance

`auth.RateLimit` is an in-process per-IP fixed-window limiter
(`backend/internal/auth/ratelimit.go`), wired only in `cmd/api` so the
shared-IP test suite never trips it. It is correct for a single API instance
(v1). Horizontal scaling of the API is the trigger to move the counter to a
shared store (e.g. Redis) — which would be the first "new infra" per
ARCHITECTURE.md.

## requestID trusts inbound X-Request-Id

`httpx.requestID` honors an inbound `X-Request-Id` for log correlation. That's
correct **only** behind our own reverse proxy. If that port ever becomes
directly client-facing, a client can forge correlation ids — sanitize
(length/charset cap) or stop honoring the header at that point. (Caveat noted
at the middleware.) Correctness/security item, not over-engineering.

## Stacked-PR workflow

Feature branches chain (feature → feature → main), not always straight to
`main`. Check a PR's `baseRefName` before merging so the work actually reaches
`main`; promote the tip branch with its own `→ main` PR when the stack is done.
