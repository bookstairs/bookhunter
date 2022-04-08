package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
)

var ErrNeedSignin = errors.New("need user account to download books")

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
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	var builder strings.Builder
	builder.WriteString(url)

	// Join request path.
	for _, path := range paths {
		path = strings.TrimRight(path, "/")
		if !strings.HasPrefix(path, "/") {
			builder.WriteString("/")
		}
		builder.WriteString(path)
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
