package naming

import (
	"mime"
	"net/url"
	"strings"
	"unicode"

	"github.com/go-resty/resty/v2"
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

// Filename parse the file name from Content-Disposition header.
// If there is no such head, we would return blank string.
func Filename(resp *resty.Response) (name string) {
	header := resp.Header()
	if disposition := header.Get("Content-Disposition"); disposition != "" {
		if _, params, err := mime.ParseMediaType(disposition); err == nil {
			if filename, ok := params["filename"]; ok {
				name = EscapeFilename(filename)
			}
		}
	}

	return
}
