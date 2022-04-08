package spider

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type Client struct {
	client *http.Client
	config *DownloadConfig
}

// Form used for post
type Form []Field
type Field struct {
	Key   string
	Value string
}

// NewClient would create a client with persisted cookie file.
func NewClient(config *DownloadConfig) *Client {
	// Create cookiejar.
	cookieFile := path.Join(config.DownloadPath, config.CookieFile)
	cookieJar, err := NewCookieJar(cookieFile)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Jar: cookieJar, Timeout: config.Timeout}
	return &Client{config: config, client: client}
}

// Get would perform http get.
func (s *Client) Get(link, referer string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, link, http.NoBody)
	if err != nil {
		return nil, err
	}

	return s.request(req, referer)
}

// FormPost would post with form-url-encoded.
func (s *Client) FormPost(link, referer string, form Form) (*http.Response, error) {
	// Prepare form data.
	values := make(url.Values, len(form))
	for _, field := range form {
		var value []string
		var ok bool
		if value, ok = values[field.Key]; ok {
			value = append(value, field.Value)
		} else {
			value = []string{field.Value}
		}
		values[field.Key] = value
	}

	req, err := http.NewRequest(http.MethodPost, link, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("illegal form post request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return s.request(req, referer)
}

// request would perform the final request and check.
func (s *Client) request(r *http.Request, referer string) (*http.Response, error) {
	r.Header.Set("User-Agent", s.config.UserAgent)
	if referer != "" {
		r.Header.Set("referer", referer)
	}

	var resp *http.Response

	for i := 0; i < s.config.Retry; i++ {
		var err error
		if resp, err = s.client.Do(r); err != nil {
			if IsTimeOut(err) {
				// Retry the timeout request.
				continue
			} else {
				return nil, err
			}
		}

		break
	}

	// Check http status code
	if resp.StatusCode != http.StatusOK {
		// Close this body manually.
		_ = resp.Body.Close()
		return nil, errors.New(resp.Status)
	} else {
		return resp, nil
	}
}

// CheckRedirect add extra check func for 302 or 301 redirect.
func (s *Client) CheckRedirect(checkFunc func(req *http.Request, via []*http.Request) error) {
	s.client.CheckRedirect = checkFunc
}
