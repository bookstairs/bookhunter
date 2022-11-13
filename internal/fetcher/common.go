package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/client"
)

var (
	ErrOverrideRedirectHandler = errors.New("couldn't override the existed redirect handler")
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
)

const (
	Talebook Category = "talebook"
	SanQiu   Category = "sanqiu"
	SoBooks  Category = "sobooks"
	TianLang Category = "tianlang"
	Telegram Category = "telegram"
)

// Archive will return if this format is an archive.
func (f Format) Archive() bool {
	return f == ZIP
}

// Config is used to define a common config for a specified fetcher service.
type Config struct {
	Category      Category // The identity of the fetcher service.
	Formats       []Format // The formats that the user wants.
	Extract       bool     // Extract the archives after download.
	DownloadPath  string   // The path for storing the file.
	InitialBookID int64    // The book id start to download.
	Rename        bool     // Rename the file by using book ID.
	Thread        int      // The number of download threads.
	RateLimit     int      // Request per minute.

	// The extra configuration for a custom fetcher services.
	Properties map[string]string

	*client.Config
}

// Property will require an existed property from the config.
func (c *Config) Property(name string) string {
	if v, ok := c.Properties[name]; ok {
		return v
	}
	return ""
}

func (c *Config) SetRedirect(redirect func(request *http.Request, requests []*http.Request) error) error {
	if c.Config.Redirect != nil {
		return ErrOverrideRedirectHandler
	}
	c.Config.Redirect = resty.RedirectPolicyFunc(redirect)

	return nil
}

// ParseFormats will create the format array from the string slice.
func ParseFormats(formats []string) ([]Format, error) {
	var fs []Format
	for _, format := range formats {
		f, err := ParseFormat(format)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}

// ParseFormat will create the format from the string.
func ParseFormat(format string) (Format, error) {
	f := Format(strings.ToLower(format))
	if !IsValidFormat(f) {
		return "", fmt.Errorf("invalid format %s", format)
	}
	return f, nil
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
	default:
		return false
	}
}

// New create a fetcher service for downloading books.
func New(c *Config) (Fetcher, error) {
	s, err := newService(c)
	if err != nil {
		return nil, err
	}

	return &commonFetcher{
		Config:  c,
		service: s,
	}, nil
}
