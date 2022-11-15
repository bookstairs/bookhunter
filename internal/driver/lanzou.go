package driver

import (
	"fmt"
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/lanzou"
)

func newLanzouDriver(c *client.Config, _ map[string]string) (Driver, error) {
	drive, err := lanzou.NewDrive(c)
	if err != nil {
		return nil, err
	}
	return &lanzouDriver{
		driver: drive,
	}, nil
}

type lanzouDriver struct {
	driver *lanzou.Drive
}

func (l *lanzouDriver) Source() Source {
	return LANZOU
}

func (l *lanzouDriver) Resolve(shareLink string, passcode string) ([]Share, error) {
	resp, err := l.driver.ResolveShareURL(shareLink, passcode)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("parsed faild: %v", resp.Msg)
	}
	return []Share{
		{FileName: resp.Data.Name, URL: resp.Data.URL, Properties: nil},
	}, err
}

func (l *lanzouDriver) Download(share Share) (io.ReadCloser, int64, error) {
	return l.driver.DownloadFile(share.URL)
}
