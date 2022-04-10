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
	config *Config
}

// Query used for get
type Query struct {
	Key   string
	Value string
}

// Form used for post
type Form []Field
type Field struct {
	Key   string
	Value string
}

// NewClient would create a client with persisted cookie file.
func NewClient(config *Config) *Client {
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
func (c *Client) Get(link, referer string, params ...*Query) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, link, http.NoBody)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		query := req.URL.Query()
		for _, param := range params {
			query.Add(param.Key, param.Value)
		}
		req.URL.RawQuery = query.Encode()
	}

	return c.request(req, referer)
}

// FormPost would post with form-url-encoded.
func (c *Client) FormPost(link, referer string, form Form) (*http.Response, error) {
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

	return c.request(req, referer)
}

// request would perform the final request and check.
func (c *Client) request(r *http.Request, referer string) (*http.Response, error) {
	r.Header.Set("User-Agent", c.config.UserAgent)
	if referer != "" {
		r.Header.Set("referer", referer)
	}

	var resp *http.Response

	err := c.Retry(func() (err error) {
		resp, err = c.client.Do(r)
		return err
	})
	if err != nil {
		return nil, err
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
func (c *Client) CheckRedirect(checkFunc func(req *http.Request, via []*http.Request) error) {
	c.client.CheckRedirect = checkFunc
}
