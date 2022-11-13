package naming

import (
	"net/url"
	"strings"
	"unicode"
)

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Extension the file extension from the link.
func Extension(link string) (string, bool) {
	u, err := url.Parse(link)
	if err != nil {
		return "", false
	}

	path := u.Path
	start := strings.LastIndex(path, ".") + 1
	ext := path[start:]

	if isLetter(ext) {
		return strings.ToLower(ext), true
	}
	return "", false
}
