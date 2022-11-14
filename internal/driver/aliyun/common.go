package aliyun

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/client"
)

var (
	// Token will be refreshed before this time.
	acceleratedExpirationDuration = 10 * time.Minute
)

type Aliyun struct {
	client       *client.Client
	tokenCache   *TokenResp
	refreshToken string
}

// New will create a aliyun download service.
func New(c *client.Config, refreshToken string) (*Aliyun, error) {
	if refreshToken == "" {
		return nil, errors.New("refreshToken is required")
	}

	cl, err := client.New(&client.Config{
		HTTPS:      true,
		Host:       "api.aliyundrive.com",
		UserAgent:  c.UserAgent,
		Proxy:      c.Proxy,
		ConfigRoot: c.ConfigRoot,
	})
	if err != nil {
		return nil, err
	}

	// Set extra middleware for cleaning up the header.
	cl.SetPreRequestHook(removeContentType)

	return &Aliyun{client: cl, refreshToken: refreshToken}, nil
}

// removeContentType is used to remove the useless content type by setting.
func removeContentType(_ *resty.Client, req *http.Request) error {
	if req.Header.Get("x-empty-content-type") != "" {
		req.Header.Del("x-empty-content-type")
		req.Header.Set("content-type", "")
	}

	return nil
}

// Try to get the cached auth token.
func (ali *Aliyun) cachedToken() *TokenResp {
	tokenResp := ali.tokenCache
	if tokenResp != nil {
		if time.Now().Add(acceleratedExpirationDuration).Before(tokenResp.ExpireTime) {
			return tokenResp
		} else {
			// Expire the cached token.
			ali.tokenCache = nil
		}
	}
	return nil
}

// cacheToken will save the token into cache.
func (ali *Aliyun) cacheToken(resp *TokenResp) {
	ali.refreshToken = resp.RefreshToken
	ali.tokenCache = resp
}
