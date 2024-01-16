package aliyun

import (
	"github.com/bookstairs/bookhunter/internal/client"
)

type Aliyun struct {
	*client.Client
	authentication *authentication
}

// New will create an aliyun download service.
func New(c *client.Config, refreshToken string) (*Aliyun, error) {
	c = &client.Config{
		HTTPS:      true,
		Host:       "api.aliyundrive.com",
		Proxy:      c.Proxy,
		ConfigRoot: c.ConfigRoot,
	}

	authentication, err := newAuthentication(c, refreshToken)
	if err != nil {
		return nil, err
	}
	if err := authentication.Auth(); err != nil {
		return nil, err
	}

	cl, err := client.New(c)
	if err != nil {
		return nil, err
	}

	// Set extra middleware for cleaning up the header and authentication.
	cl.SetPreRequestHook(authentication.authenticationHook())

	return &Aliyun{Client: cl, authentication: authentication}, nil
}
