package fetcher

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bookstairs/bookhunter/internal/client"
)

type (
	Format   string // The supported file extension.
	Category string // The fetcher service identity.
)

const (
	EPUB Format = "epub"
	MOBI Format = "mobi"
	AZW  Format = "azw"
	AZW3 Format = "azw3"
	PDF  Format = "pdf"
	ZIP  Format = "zip"
	RAR  Format = "rar"
	TAR  Format = "tar"
	GZIP Format = "gz"
)

const (
	Talebook Category = "talebook"
	SanQiu   Category = "sanqiu"
	SoBooks  Category = "sobooks"
	TianLang Category = "tianlang"
	Telegram Category = "telegram"
)

// Config is used to define a common config for a specified fetcher service.
type Config struct {
	Category      Category // The identity of the fetcher service.
	Formats       []Format // The formats that the user wants.
	Extract       bool     // Extract the archives after download.
	DownloadPath  string   // The path for storing the file.
	InitialBookID int      // The book id start to download.
	Rename        bool     // Rename the file by using book ID.
	Thread        int      // The number of download threads.
	RateLimit     int      // Request per minute.

	// The extra configuration for a custom fetcher services.
	Properties map[string]string

	*client.Config
}

func (c *Config) Property(name string) (string, error) {
	if v, ok := c.Properties[name]; ok {
		return v, nil
	}
	return "", fmt.Errorf("no such config key [%s] existed", name)
}

type Fetcher interface {
	Download() error
}

// ParseFormats will create the format array from the string slice.
func ParseFormats(formats []string) ([]Format, error) {
	var fs []Format
	for _, format := range formats {
		f := Format(strings.ToLower(format))
		if !IsValidFormat(f) {
			return nil, fmt.Errorf("invalid format %s", format)
		}
		fs = append(fs, f)
	}
	return fs, nil
}

// New create a fetcher service for downloading books.
func New(c *Config) (Fetcher, error) {
	switch c.Category {
	case Talebook:
		return nil, errors.New("we don't talebook talebook now")
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

// IsValidFormat judge if this format was supported.
func IsValidFormat(format Format) bool {
	switch format {
	case EPUB:
		return true
	case MOBI:
		return true
	case AZW:
		return true
	case AZW3:
		return true
	case PDF:
		return true
	case ZIP:
		return true
	case RAR:
		return true
	case TAR:
		return true
	case GZIP:
		return true
	default:
		return false
	}
}
