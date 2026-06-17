package auth

import (
	"strings"
	"testing"
)

func TestValidateDisplayName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"trims surrounding space", "  Ana  ", "Ana", false},
		{"too short", "a", "", true},
		{"min length ok", "ab", "ab", false},
		{"blank after trim", "   ", "", true},
		{"too long", strings.Repeat("x", 51), "", true},
		{"max length ok", strings.Repeat("x", 50), strings.Repeat("x", 50), false},
		{"counts runes not bytes", "áé", "áé", false}, // 2 runes, 4 bytes
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDisplayName(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateDisplayName(%q) err = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ValidateDisplayName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
