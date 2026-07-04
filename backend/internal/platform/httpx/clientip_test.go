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

func TestClientIPKey(t *testing.T) {
	key := func(ip string) string {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "10.0.0.5:1234"
		r.Header.Set("CF-Connecting-IP", ip)
		return ClientIPKey(r)
	}

	// IPv4: the address is the bucket.
	if got := key("203.0.113.7"); got != "203.0.113.7" {
		t.Errorf("v4 key = %q, want the address itself", got)
	}
	// IPv6: everything in one routed /64 shares a bucket; a neighboring /64
	// does not.
	a := key("2001:db8:1:1::1")
	b := key("2001:db8:1:1:ffff:ffff:ffff:ffff")
	c := key("2001:db8:1:2::1")
	if a != b {
		t.Errorf("same /64 got different keys: %q vs %q", a, b)
	}
	if a == c {
		t.Errorf("different /64s share a key: %q", a)
	}
	// 4-mapped-in-6 is IPv4 in disguise: keyed as the v4 address.
	if got := key("::ffff:203.0.113.7"); got != "203.0.113.7" {
		t.Errorf("4-in-6 key = %q, want %q", got, "203.0.113.7")
	}
	// Unparseable header value still yields a usable (if odd) bucket.
	if got := key("not-an-ip"); got != "not-an-ip" {
		t.Errorf("unparseable key = %q, want the raw string", got)
	}
}
