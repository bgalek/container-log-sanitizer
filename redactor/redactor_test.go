package redactor

import (
	"os"
	"testing"
)

func TestRedactLine(t *testing.T) {
	os.Setenv("SANITIZE_PATTERNS", `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,};token=[^\s]+`)
	r, err := NewRedactor()
	if err != nil {
		t.Fatalf("failed to create redactor: %v", err)
	}
	cases := []struct {
		in, want string
	}{
		{"user@example.com", "***"},
		{"token=abc123", "***"},
		{"no match here", "no match here"},
	}
	for _, c := range cases {
		got := r.RedactLine(c.in)
		if got != c.want {
			t.Errorf("RedactLine(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
