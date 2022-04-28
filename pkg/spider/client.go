package spider

import (
	"io"
	"path"

	"github.com/go-resty/resty/v2"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type Client struct {
	client *resty.Client
	config *Config
}

// NewClient would create a client with persisted cookie file.
func NewClient(config *Config) *Client {
	// Create cookiejar.
	cookieFile := path.Join(config.DownloadPath, config.CookieFile)
	cookieJar, err := NewCookieJar(cookieFile)
	if err != nil {
		log.Fatal(err)
	}

	client := resty.New()
	client.SetCookieJar(cookieJar)
	client.SetTimeout(config.Timeout)
	client.SetDisableWarn(true)
	client.SetRetryCount(config.Retry)
	client.SetHeader("User-Agent", config.UserAgent)
	client.SetDebug(config.Debug)
	client.SetDebugBodyLimit(1024 * 10)
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))
	return &Client{config: config, client: client}
}

func (c *Client) R() *resty.Request {
	return c.client.R()
}

func (c *Client) Download(link string, save func(filename string, contentLength int64, data io.ReadCloser) error) error {
	resp, err := c.client.GetClient().Get(link)
	if err != nil {
		return err
	}
	filename := Filename(resp)
	contentLength := resp.ContentLength
	err = save(filename, contentLength, resp.Body)
	return err
}
