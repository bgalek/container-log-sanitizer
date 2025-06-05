package redactor

import (
	"os"
	"regexp"
	"strings"
)

type Redactor struct {
	patterns []*regexp.Regexp
}

// NewRedactor creates a Redactor from env or config file
func NewRedactor() (*Redactor, error) {
	var regexes []string
	if env := os.Getenv("SANITIZE_PATTERNS"); env != "" {
		for _, pat := range strings.Split(env, ";") {
			pat = strings.TrimSpace(pat)
			if pat != "" {
				// Unescape double backslashes for env var input
				pat = strings.ReplaceAll(pat, "\\", "\\")
				regexes = append(regexes, pat)
			}
		}
	}
	var compiled []*regexp.Regexp
	for _, r := range regexes {
		if re, err := regexp.Compile(r); err == nil {
			compiled = append(compiled, re)
		}
	}
	return &Redactor{patterns: compiled}, nil
}

// RedactLine applies all patterns to a line
func (r *Redactor) RedactLine(line string) string {
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, "***")
	}
	return line
}
