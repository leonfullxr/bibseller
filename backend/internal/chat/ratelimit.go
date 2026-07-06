package chat

import "time"

// Message-send budgets. Polling (GET) is unmetered; only writes are capped.
// Per-account is the primary anti-spam control; the per-IP cap is deliberately
// generous so it stops a single source flooding writes across many accounts
// without tripping a shared NAT of legitimate users. The limiter itself lives
// in internal/platform/ratelimit (#122).
const (
	msgRateMax    = 30
	ipRateMax     = 120
	msgRateWindow = time.Minute
)
