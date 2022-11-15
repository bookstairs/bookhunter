package driver

import (
	"errors"
	"io"
	"strings"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/aliyun"
)

var (
	ErrNoRefreshToken = errors.New("aliyun drive need refreshToken")
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
	shareID := strings.TrimPrefix(shareLink, "https://www.aliyundrive.com/s/")
	sharePwd := strings.TrimSpace(passcode)

	token, err := a.client.ShareToken(shareID, sharePwd)
	if err != nil {
		return nil, err
	}

	files, err := a.client.Share(shareID, token.ShareToken)
	if err != nil {
		return nil, err
	}

	var shares []Share
	for index := range files {
		item := files[index]
		share := Share{
			FileName: item.Name,
			Size:     int64(item.Size),
			URL:      item.FileID,
			Properties: map[string]any{
				"shareToken": token.ShareToken,
				"shareID":    shareID,
				"fileID":     item.FileID,
			},
		}

		shares = append(shares, share)
	}

	return shares, nil
}

func (a *aliyunDriver) Download(share Share) (io.ReadCloser, error) {
	shareToken := share.Properties["shareToken"].(string)
	shareID := share.Properties["shareID"].(string)
	fileID := share.Properties["fileID"].(string)

	url, err := a.client.DownloadURL(shareToken, shareID, fileID)
	if err != nil {
		return nil, err
	}

	return a.client.DownloadFile(url)
}
