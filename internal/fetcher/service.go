package fetcher

import (
	"fmt"
	"strings"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
)

// service is a real implementation for fetcher.
type service interface {
	// size is the total download amount of this given service.
	size() (int64, error)

	// formats will query the available downloadable file formats.
	formats(int64) (map[file.Format]driver.Share, error)

	// fetch the given book ID.
	fetch(int64, file.Format, driver.Share, file.Writer) error
}

func matchKeywords(title string, keywords []string) bool {
	if len(keywords) == 0 {
		return true
	}

	for _, keyword := range keywords {
		if strings.Contains(title, keyword) {
			return true
		}
	}

	return false
}

// newService is the endpoint for creating all the supported download service.
func newService(c *Config) (service, error) {
	switch c.Category {
	case Talebook:
		return newTalebookService(c)
	case SanQiu:
		return newSanqiuService(c)
	case SoBooks:
		return newSobooksService(c)
	case TianLang:
		return newTianlangService(c)
	case Telegram:
		return newTelegramService(c)
	case K12:
		return newK12Service(c)
	default:
		return nil, fmt.Errorf("no such fetcher service [%s] supported", c.Category)
	}
}
