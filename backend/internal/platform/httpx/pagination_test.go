package httpx

import (
	"net/url"
	"testing"
)

func TestParseLimit(t *testing.T) {
	cases := []struct {
		name    string
		raw     string // "" means the param is absent
		set     bool
		want    int32
		wantErr bool
	}{
		{name: "absent uses default", set: false, want: DefaultPageSize},
		{name: "valid", set: true, raw: "10", want: 10},
		{name: "max is allowed", set: true, raw: "100", want: MaxPageSize},
		{name: "zero rejected", set: true, raw: "0", wantErr: true},
		{name: "over max rejected", set: true, raw: "101", wantErr: true},
		{name: "non-numeric rejected", set: true, raw: "abc", wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q := url.Values{}
			if c.set {
				q.Set("limit", c.raw)
			}
			got, err := ParseLimit(q)
			if (err != nil) != c.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, c.wantErr)
			}
			if got != c.want {
				t.Errorf("got %d, want %d", got, c.want)
			}
		})
	}
}
