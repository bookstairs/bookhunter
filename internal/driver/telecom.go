package driver

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver/telecom"
)

func newTelecomDriver(config *client.Config, properties map[string]string) (Driver, error) {
	// Try to sign in to the telecom disk.
	username := properties["username"]
	password := properties["password"]

	// Create the pan client.
	t, err := telecom.New(config, username, password)
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
	panic("TODO implement me")
}

func (t *telecomDriver) Download(share Share) (io.ReadCloser, int64, error) {
	panic("TODO implement me")
}
