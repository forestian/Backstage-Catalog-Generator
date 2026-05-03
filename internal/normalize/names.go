package normalize

import (
	"regexp"
	"strings"
)

var nonAlphaNum = regexp.MustCompile(`[^a-z0-9-]`)
var multipleDashes = regexp.MustCompile(`-+`)

// EntityName converts a string into a DNS-like lowercase identifier suitable
// for Backstage metadata.name.
func EntityName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphaNum.ReplaceAllString(s, "")
	s = multipleDashes.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
