package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	ready := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/readyz" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ready.Close()
	if got := healthcheck(ready.URL); got != 0 {
		t.Errorf("ready server: exit = %d, want 0", got)
	}

	notReady := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer notReady.Close()
	if got := healthcheck(notReady.URL); got != 1 {
		t.Errorf("503 server: exit = %d, want 1", got)
	}

	down := httptest.NewServer(nil)
	down.Close() // connection refused
	if got := healthcheck(down.URL); got != 1 {
		t.Errorf("unreachable server: exit = %d, want 1", got)
	}
}
