package ids

import "testing"

func TestNewIsV7AndUnique(t *testing.T) {
	a, b := New(), New()
	if a.Version() != 7 {
		t.Errorf("version = %d, want 7", a.Version())
	}
	if a == b {
		t.Error("two generated ids are equal")
	}
}
