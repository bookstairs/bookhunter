package driver

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/lanzou"
)

func newLanzouDriver(c *client.Config, _ map[string]string) (Driver, error) {
	cl, _ := client.New(&client.Config{
		HTTPS:      true,
		Host:       "lanzoux.com",
		UserAgent:  c.UserAgent,
		Proxy:      c.Proxy,
		ConfigRoot: c.ConfigRoot,
	})
	driver := &lanzou.Drive{
		Client:  cl.Client,
		BaseURL: "https://lanzoux.com",
	}
	return &lanzouDriver{
		driver: driver,
	}, nil
}

type lanzouDriver struct {
	driver *lanzou.Drive
}

func (l *lanzouDriver) Source() Source {
	return LANZOU
}

func (l *lanzouDriver) Resolve(shareLink string, passcode string) ([]Share, error) {
	_, _ = l.driver.ResolveShareURL(shareLink, passcode)
	panic("TODO implement me")
}

func (l *lanzouDriver) Download(share Share) (io.ReadCloser, int64, error) {
	panic("TODO implement me")
}
