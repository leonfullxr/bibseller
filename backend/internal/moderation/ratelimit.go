package moderation

import "time"

// Per-account budget on the safety write endpoints (report, block). The
// limiter itself lives in internal/platform/ratelimit (#122).
const (
	reportRateMax    = 20
	reportRateWindow = time.Minute
)
