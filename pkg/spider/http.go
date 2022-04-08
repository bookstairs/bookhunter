package spider

import (
	"encoding/json"
	"fmt"
	"mime"
	"net"
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

// GenerateUrl would remove the "/" suffix and add schema prefix to url.
func GenerateUrl(base string, paths ...string) string {
	// Remove suffix
	url := strings.TrimRight(base, "/")

	// Add schema prefix
	if !strings.HasPrefix(url, HTTPS) && !strings.HasPrefix(url, HTTP) {
		url = HTTP + url
	}

	var builder strings.Builder
	builder.WriteString(url)

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

// WrapTimeOut would convert the timeout error with a better prefix in error message.
func WrapTimeOut(err error) error {
	if timeoutErr, ok := err.(net.Error); ok && timeoutErr.Timeout() {
		return fmt.Errorf("timeout %v", timeoutErr)
	}

	return err
}
