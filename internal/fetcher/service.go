package fetcher

import (
	"errors"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/naming"
)

// service is a real implementation for fetcher.
type service interface {
	// size is the total download amount of this given service.
	size() (int64, error)

	// formats will query the available downloadable file formats.
	formats(id int64) (map[Format]string, error)

	// fetch the given book ID.
	fetch(id int64, format Format, url string) (*fetch, error)
}

// fetch is the result which can be created from a resty.Response
type fetch struct {
	name    string
	content io.ReadCloser
	size    int64
}

// createFetch will try to create a fetch. Remember to create the response with the `SetDoNotParseResponse` option.
func createFetch(resp *resty.Response) *fetch {
	// Try to parse a filename if it existed.
	name := naming.Filename(resp)
	return &fetch{
		name:    name,
		content: resp.RawBody(),
		size:    resp.RawResponse.ContentLength,
	}
}

// newService is the endpoint for creating all the supported download service.
func newService(c *Config) (service, error) {
	switch c.Category {
	case Talebook:
		return newTalebookService(c)
	case SanQiu:
		return nil, errors.New("we don't support sanqiu now")
	case Telegram:
		return nil, errors.New("we don't support telegram now")
	case SoBooks:
		return nil, errors.New("we don't support sobooks now")
	case TianLang:
		return nil, errors.New("we don't support tianlang now")
	default:
		return nil, fmt.Errorf("no such fetcher service [%s] supported", c.Category)
	}
}
