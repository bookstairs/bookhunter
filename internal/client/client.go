package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/log"
)

var (
	ErrInvalidRequestURL = errors.New("invalid request url, we only support https:// or http://")
)

const (
	cookieFile       = "cookies.json"
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko)" +
		" Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.42"
)

// Client is the wrapper for resty.Client we may provide extra method on this wrapper.
type Client struct {
	*resty.Client
	*Config
}

// Config is the basic configuration for creating the client.
type Config struct {
	HTTPS      bool   // If the request was under the https of http.
	Host       string // The request host name.
	UserAgent  string // Custom user agent for mocking as the browser client.
	Proxy      string // The proxy address, such as the http://127.0.0.1:7890, socks://127.0.0.1:7890
	ConfigRoot string // The root config path for whole bookhunter download service.

	// The custom redirect function.
	Redirect resty.RedirectPolicy `json:"-"`
}

// ConfigPath will return a unique path for this download service.
func (c *Config) ConfigPath() (string, error) {
	if c.ConfigRoot == "" {
		var err error
		c.ConfigRoot, err = DefaultConfigRoot()
		if err != nil {
			return "", err
		}
	}

	return mkdir(filepath.Join(c.ConfigRoot, c.Host))
}

func (c *Config) newCookieJar() (http.CookieJar, error) {
	configPath, err := c.ConfigPath()
	if err != nil {
		return nil, err
	}

	return newCookieJar(filepath.Join(configPath, cookieFile))
}

func (c *Config) redirectPolicy() []any {
	policies := []any{
		resty.FlexibleRedirectPolicy(5),
	}
	if c.Redirect != nil {
		policies = append(policies, c.Redirect)
	}

	return policies
}

func (c *Config) userAgent() string {
	if c.UserAgent == "" {
		return DefaultUserAgent
	}

	return c.UserAgent
}

func (c *Config) baseURL() string {
	if c.HTTPS {
		return "https://" + c.Host
	}

	return "http://" + c.Host
}

func (c *Client) SetDefaultHostname(host string) {
	c.Host = host
	c.Client.SetBaseURL(c.baseURL())
}

// DefaultConfigRoot will generate the default config path based on the user and his running environment.
func DefaultConfigRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return mkdir(filepath.Join(home, ".config", "bookhunter"))
}

func mkdir(path string) (string, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}

	return path, nil
}

// NewConfig will create a config instance by using the request url.
func NewConfig(rawURL, userAgent, proxy, configRoot string) (*Config, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf(rawURL, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalidRequestURL
	}

	if configRoot == "" {
		configRoot, err = DefaultConfigRoot()
		if err != nil {
			return nil, err
		}
	} else {
		if err := os.MkdirAll(configRoot, 0755); err != nil {
			return nil, err
		}
	}

	return &Config{
		HTTPS:      u.Scheme == "https",
		Host:       u.Host,
		UserAgent:  userAgent,
		Proxy:      proxy,
		ConfigRoot: configRoot,
	}, nil
}

// New will create a resty client with a lot of predefined settings.
func New(c *Config) (*Client, error) {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(3*time.Second).
		SetRetryMaxWaitTime(10*time.Second).
		SetAllowGetMethodPayload(true).
		SetTimeout(5*time.Minute).
		SetContentLength(true).
		SetDebug(log.EnableDebug).
		SetDisableWarn(true).
		SetHeader("User-Agent", c.userAgent())

	if c.Host != "" {
		client.SetBaseURL(c.baseURL())
	}

	if len(c.redirectPolicy()) > 0 {
		client.SetRedirectPolicy(c.redirectPolicy()...)
	}

	// Setting the cookiejar
	cookieJar, err := c.newCookieJar()
	if err != nil {
		return nil, err
	}
	client.SetCookieJar(cookieJar)

	// Setting the proxy for the resty client.
	// The proxy environment is also supported.
	if c.Proxy != "" {
		client.SetProxy(c.Proxy)
	}

	return &Client{Client: client, Config: c}, nil
}
