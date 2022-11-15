package driver

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
)

func newLanzouDriver(_ *client.Config, _ map[string]string) (Driver, error) {
	// TODO Add this implementation.
	return &lanzouDriver{}, nil
}

type lanzouDriver struct{}

func (l *lanzouDriver) Source() Source {
	return LANZOU
}

func (l *lanzouDriver) Resolve(shareLink string, passcode string) []Share {
	panic("TODO implement me")
}

func (l *lanzouDriver) Download(share Share) io.ReadCloser {
	panic("TODO implement me")
}
