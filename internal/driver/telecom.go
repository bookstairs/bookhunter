package driver

import (
	"io"
	"strconv"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/telecom"
)

func newTelecomDriver(config *client.Config, properties map[string]string) (Driver, error) {
	// Create the pan client.
	t, err := telecom.New(config, properties["username"], properties["password"])
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

func (t *telecomDriver) Resolve(shareLink string, passcode string) ([]Share, error) {
	code, err := t.client.ShareCode(shareLink)
	if err != nil {
		return nil, err
	}

	info, files, err := t.client.ShareFiles(shareLink, passcode)
	if err != nil {
		return nil, err
	}

	// Convert this into a share entity.
	shares := make([]Share, 0, len(files))
	for _, file := range files {
		shares = append(shares, Share{
			FileName: file.Name,
			Properties: map[string]string{
				"fileID":    strconv.FormatInt(file.ID, 10),
				"shareID":   strconv.FormatInt(info.ShareID, 10),
				"shareCode": code,
				"fileSize":  strconv.FormatInt(file.Size, 10),
			},
		})
	}

	return shares, nil
}

func (t *telecomDriver) Download(share Share) (io.ReadCloser, int64, error) {
	shareCode := share.Properties["shareCode"]
	shareID := share.Properties["shareID"]
	fileID := share.Properties["fileID"]
	size, _ := strconv.ParseInt(share.Properties["fileSize"], 10, 64)

	// Resolve the link.
	url, err := t.client.DownloadURL(shareCode, shareID, fileID)
	if err != nil {
		return nil, 0, err
	}

	// Download the file.
	file, err := t.client.DownloadFile(url)

	return file, size, err
}
