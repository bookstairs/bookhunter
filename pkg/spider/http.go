package spider

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"
)

const (
	HTTP  = "http://"
	HTTPS = "https://"
)

// DecodeResponse would parse the http response into a json based content.
func DecodeResponse(resp *http.Response, data any) (err error) {
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(data)

	return
}

// Extension the file extension from the link or file name.
func Extension(link string) (string, bool) {
	end := strings.LastIndex(link, "?")
	if end >= 0 {
		link = link[:end]
	}
	start := strings.LastIndex(link, ".") + 1

	filename := strings.ToLower(link[start:])
	if strings.Contains(filename, "/") {
		return "", false
	}
	return filename, true
}

// Filename parse the file name from Content-Disposition header.
// If there is no such head, we would return blank string.
func Filename(resp *http.Response) (name string) {
	if disposition := resp.Header.Get("Content-Disposition"); disposition != "" {
		if _, params, err := mime.ParseMediaType(disposition); err == nil {
			if filename, ok := params["filename"]; ok {
				name = filename
			}
		}
	}

	return
}

// GenerateURL would remove the "/" suffix and add schema prefix to url.
func GenerateURL(base string, paths ...string) string {
	// Remove suffix
	l := strings.TrimRight(base, "/")

	// Add schema prefix
	if !strings.HasPrefix(l, HTTPS) && !strings.HasPrefix(l, HTTP) {
		l = HTTP + l
	}

	var builder strings.Builder
	builder.WriteString(l)

	// Join request path.
	for _, p := range paths {
		p = strings.TrimRight(p, "/")
		if !strings.HasPrefix(p, "/") {
			builder.WriteString("/")
		}
		builder.WriteString(p)
	}

	return builder.String()
}
