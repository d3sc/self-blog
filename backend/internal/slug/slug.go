package slug

import (
	"regexp"
	"strings"
	"unicode"
)

var non = regexp.MustCompile(`[^a-z0-9]+`)

func Make(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			b.WriteRune(r)
		}
	}
	return strings.Trim(non.ReplaceAllString(b.String(), "-"), "-")
}
