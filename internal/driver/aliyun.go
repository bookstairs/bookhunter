package driver

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/aliyun"
)

func newAliyunDriver(c *client.Config) (Driver, error) {
	return &aliyunDriver{}, nil
}

type aliyunDriver struct {
	Aliyun *aliyun.Aliyun
}

func (a *aliyunDriver) Resolve(shareLink string, passcode string) []Share {
	panic("TODO implement me")
}

func (a *aliyunDriver) Download(share Share) io.ReadCloser {
	panic("TODO implement me")
}
