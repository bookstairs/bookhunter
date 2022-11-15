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

func (a *aliyunDriver) Source() Source {
	return ALIYUN
}

func (a *aliyunDriver) Resolve(shareLink string, passcode string) ([]Share, error) {
	c := a.client
	token, err := c.ShareToken(shareLink, passcode)
	if err != nil {
		return nil, err
	}
	shareFiles, err := c.Share(shareLink, token.ShareToken)
	if err != nil {
		return nil, err
	}
	var shares []Share
	for index := range shareFiles {
		item := shareFiles[index]
		shares = append(shares, Share{
			FileName: item.Name,
			URL:      item.FileID,
			Properties: map[string]string{
				"shareToken": token.ShareToken,
				"shareID":    shareLink,
			},
		})
	}
	return shares, nil
}

func (a *aliyunDriver) Download(share Share) (io.ReadCloser, int64, error) {
	c := a.client
	url, err := c.DownloadURL(share.Properties["shareToken"], share.Properties["shareID"], share.URL)
	if err != nil {
		return nil, 0, err
	}
	return c.DownloadFile(url)
}
