package spider

import (
	"net/http"
	"path"

	"github.com/go-resty/resty/v2"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type Client struct {
	client *resty.Client
	config *Config
}

// Query used for get
type Query struct {
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

func (c *Client) GetHttpClient() *http.Client {
	return c.client.GetClient()
}
