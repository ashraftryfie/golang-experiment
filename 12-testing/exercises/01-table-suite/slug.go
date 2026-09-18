package slug

import (
	"regexp"
	"strings"
)

var (
	nonAlphaNumRegex  = regexp.MustCompile(`[^a-z0-9-]+`)
	multipleDashRegex = regexp.MustCompile(`-+`)
)

// Slugify cleans a string for use in URLs.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	s = nonAlphaNumRegex.ReplaceAllString(s, "")
	s = multipleDashRegex.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
