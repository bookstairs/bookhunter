package fetcher

import (
	"errors"
	"fmt"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
)

// service is a real implementation for fetcher.
type service interface {
	// size is the total download amount of this given service.
	size() (int64, error)

	// formats will query the available downloadable file formats.
	formats(int64) (map[Format]driver.Share, error)

	// fetch the given book ID.
	fetch(int64, Format, driver.Share, file.Writer) error
}

// newService is the endpoint for creating all the supported download service.
func newService(c *Config) (service, error) {
	switch c.Category {
	case Talebook:
		return newTalebookService(c)
	case SanQiu:
		return newSanqiuService(c)
	case TianLang:
		return newTianlangService(c)
	case SoBooks:
		return nil, errors.New("TODO we don't support sobooks now")
	default:
		return nil, fmt.Errorf("no such fetcher service [%s] supported", c.Category)
	}
}
