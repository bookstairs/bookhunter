package driver

import (
	"io"
	"strings"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/aliyun"
)

// newAliyunDriver will create the aliyun driver.
func newAliyunDriver(c *client.Config, properties map[string]string) (Driver, error) {
	// Create the aliyun client.
	a, err := aliyun.New(c, properties["refreshToken"])
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

func (a *aliyunDriver) Resolve(link, passcode string) ([]Share, error) {
	shareID := strings.TrimPrefix(link, "https://www.aliyundrive.com/s/")
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

func (a *aliyunDriver) Download(share Share) (io.ReadCloser, int64, error) {
	shareToken := share.Properties["shareToken"].(string)
	shareID := share.Properties["shareID"].(string)
	fileID := share.Properties["fileID"].(string)

	url, err := a.client.DownloadURL(shareToken, shareID, fileID)
	if err != nil {
		return nil, 0, err
	}

	file, err := a.client.DownloadFile(url)

	return file, 0, err
}
