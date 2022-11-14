package driver

import (
	"errors"
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/aliyun"
)

var (
	ErrNoRefreshToken = errors.New("")
)

// newAliyunDriver will create the aliyun driver.
func newAliyunDriver(c *client.Config, properties map[string]string) (Driver, error) {
	// Get the refreshToken.
	token, exist := properties["refreshToken"]
	if !exist || token == "" {
		return nil, ErrNoRefreshToken
	}

	// Create the aliyun client.
	a, err := aliyun.New(c, token)
	if err != nil {
		return nil, err
	}

	// Check if the refresh token is valid at the beginning.
	_, err = a.AuthToken()
	if err != nil {
		return nil, err
	}

	return &aliyunDriver{client: a}, nil
}

type aliyunDriver struct {
	client *aliyun.Aliyun
}

func (a *aliyunDriver) Resolve(shareLink string, passcode string) []Share {
	panic("TODO implement me")
}

func (a *aliyunDriver) Download(share Share) io.ReadCloser {
	panic("TODO implement me")
}
