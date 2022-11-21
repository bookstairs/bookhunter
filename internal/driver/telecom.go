package driver

import (
	"io"
	"strconv"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/telecom"
)

func newTelecomDriver(config *client.Config, properties map[string]string) (Driver, error) {
	// Create the pan client.
	t, err := telecom.New(config, properties["telecomUsername"], properties["telecomPassword"])
	if err != nil {
		return nil, err
	}

	return &telecomDriver{client: t}, nil
}

type telecomDriver struct {
	client *telecom.Telecom
}

func (t *telecomDriver) Source() Source {
	return TELECOM
}

func (t *telecomDriver) Resolve(link, passcode string) ([]Share, error) {
	code, err := t.client.ShareCode(link)
	if err != nil {
		return nil, err
	}

	info, files, err := t.client.ShareFiles(link, passcode)
	if err != nil {
		return nil, err
	}

	// Convert this into a share entity.
	shares := make([]Share, 0, len(files))
	for _, file := range files {
		shares = append(shares, Share{
			FileName: file.Name,
			Size:     file.Size,
			Properties: map[string]any{
				"shareCode": code,
				"shareID":   strconv.FormatInt(info.ShareID, 10),
				"fileID":    strconv.FormatInt(file.ID, 10),
			},
		})
	}

	return shares, nil
}

func (t *telecomDriver) Download(share Share) (io.ReadCloser, int64, error) {
	shareCode := share.Properties["shareCode"].(string)
	shareID := share.Properties["shareID"].(string)
	fileID := share.Properties["fileID"].(string)

	// Resolve the link.
	url, err := t.client.DownloadURL(shareCode, shareID, fileID)
	if err != nil {
		return nil, 0, err
	}

	// Download the file.
	file, err := t.client.DownloadFile(url)

	return file, 0, err
}
