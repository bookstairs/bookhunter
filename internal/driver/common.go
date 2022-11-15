package driver

import (
	"fmt"
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
)

type (
	// Source is a net drive disk provider.
	Source string

	// Share is an atomic downloadable file.
	Share struct {
		// FileName is a file name with the file extension.
		FileName string
		// URL is the downloadable url for this file.
		URL string
		// Properties could be some metadata, such as the token for this downloadable share.
		Properties map[string]string
	}

	// Driver is used to resolve the links from a Source.
	Driver interface {
		// Source will return the driver identity.
		Source() Source

		// Resolve the given link and return the file name with the download link.
		Resolve(shareLink string, passcode string) ([]Share, error)

		// Download the given link.
		Download(share Share) (io.ReadCloser, int64, error)
	}
)

const (
	ALIYUN  Source = "aliyun"
	LANZOU  Source = "lanzou"
	TELECOM Source = "telecom"
)

// New will create the basic driver service.
func New(config *client.Config, properties map[string]string) (Driver, error) {
	source := Source(properties["driver"])
	switch source {
	case ALIYUN:
		return newAliyunDriver(config, properties)
	case TELECOM:
		return newTelecomDriver(config, properties)
	case LANZOU:
		return newLanzouDriver(config, properties)
	default:
		return nil, fmt.Errorf("invalid driver service %s", source)
	}
}
