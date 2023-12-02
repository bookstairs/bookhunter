package fetcher

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/file"
)

var (
	ErrOverrideRedirectHandler = errors.New("couldn't override the existed redirect handler")
	ErrFileNotExist            = errors.New("current file does not exist")
)

type Category string // The fetcher service identity.

const (
	Talebook Category = "talebook"
	SanQiu   Category = "sanqiu"
	SoBooks  Category = "sobooks"
	TianLang Category = "tianlang"
	Telegram Category = "telegram"
	K12      Category = "k12"
)

// Config is used to define a common config for a specified fetcher service.
type Config struct {
	Category      Category      // The identity of the fetcher service.
	Formats       []file.Format // The formats that the user wants.
	Keywords      []string      // The keywords that the user wants.
	Extract       bool          // Extract the archives after download.
	DownloadPath  string        // The path for storing the file.
	InitialBookID int64         // The book id start to download.
	Rename        bool          // Rename the file by using book ID.
	Thread        int           // The number of download threads.
	RateLimit     int           // Request per minute for a thread.
	Retry         int           // The retry times for a failed download.
	SkipError     bool          // Continue to download the next book if the current book download failed.
	processFile   string        // Define the download process.

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

func (c *Config) SetRedirect(redirect func(*http.Request, []*http.Request) error) error {
	if c.Config.Redirect != nil {
		return ErrOverrideRedirectHandler
	}
	c.Config.Redirect = resty.RedirectPolicyFunc(redirect)

	return nil
}

// ParseFormats will create the format array from the string slice.
func ParseFormats(formats []string) ([]file.Format, error) {
	var fs []file.Format
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
func ParseFormat(format string) (file.Format, error) {
	f := file.Format(strings.ToLower(format))
	if !IsValidFormat(f) {
		return "", fmt.Errorf("invalid format %s", format)
	}
	return f, nil
}

// IsValidFormat judge if this format was supported.
func IsValidFormat(format file.Format) bool {
	switch format {
	case file.EPUB:
		return true
	case file.MOBI:
		return true
	case file.AZW:
		return true
	case file.AZW3:
		return true
	case file.PDF:
		return true
	case file.ZIP:
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

	return &fetcher{Config: c, service: s}, nil
}
