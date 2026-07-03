package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.5:1234"
	if got := ClientIP(r); got != "10.0.0.5" {
		t.Errorf("header absent: ClientIP = %q, want %q", got, "10.0.0.5")
	}

	r.Header.Set("CF-Connecting-IP", "203.0.113.7")
	if got := ClientIP(r); got != "203.0.113.7" {
		t.Errorf("header present: ClientIP = %q, want %q", got, "203.0.113.7")
	}
}
