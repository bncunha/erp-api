package helper

import "testing"

func TestParseInt64(t *testing.T) {
	if got := ParseInt64("42"); got != 42 {
		t.Fatalf("expected 42, got %d", got)
	}
	if got := ParseInt64("invalid"); got != 0 {
		t.Fatalf("expected 0 for invalid input, got %d", got)
	}
}
