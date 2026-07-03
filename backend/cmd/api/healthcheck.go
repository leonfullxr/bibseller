package main

import (
	"net/http"
	"time"
)

// healthcheck probes the running server's readiness endpoint and reports the
// result as an exit code. It exists because the prod image is distroless (no
// shell, no curl), so the compose healthcheck runs the binary itself:
// `/api -healthcheck`. Readiness (DB reachable) rather than liveness is what
// dependents care about (issue #134).
func healthcheck(base string) int {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(base + "/api/readyz")
	if err != nil {
		return 1
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
