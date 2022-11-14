package driver

import (
	"errors"
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
		// Resolve the given link and return the file name with the download link.
		Resolve(shareLink string, passcode string) []Share

		// Download the given link.
		Download(share Share) io.ReadCloser
	}
)

const (
	ALIYUN  Source = "aliyun"
	LANZOU  Source = "lanzou"
	TELECOM Source = "telecom"
)

// New will create the basic driver service.
func New(source Source, config *client.Config) (Driver, error) {
	switch source {
	case ALIYUN:
		return newAliyunDriver(config)
	case TELECOM:
		return nil, errors.New("we don't support telecom currently")
	case LANZOU:
		return nil, errors.New("we don't support lanzou currently")
	default:
		return nil, fmt.Errorf("invalid driver service %s", source)
	}
}
