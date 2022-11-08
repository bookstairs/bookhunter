package spider

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/pkg/log"
)

var (
	ErrInitialBookID = errors.New("illegal book id, it should exceed 0")
	ErrRetryTimes    = errors.New("illegal retry times, it should exceed 0")
	ErrThreadCounts  = errors.New("illegal download thread counts, it should exceed 0")
)

const DefaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"

type Config struct {
	Website       string        // The website for talebook.
	Username      string        // The login user.
	Password      string        // The password for login user.
	DownloadPath  string        // Use the executed directory as the default download path.
	CookieFile    string        // The cookie file to use in this download progress.
	ProgressFile  string        // The progress file serving the remaining book id.
	InitialBookID int           // The book id start to download.
	Formats       []string      // The file formats you want to download
	Timeout       time.Duration // The request timeout for a single request.
	Retry         int           // The maximum retry times for a timeout request.
	UserAgent     string        // The user agent for the download request.
	Rename        bool          // Rename the file by using book ID.
	Thread        int           // The number of download thread.
	Debug         bool          // Enable debug log.
}

// NewConfig will return a default blank config.
func NewConfig() *Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		DownloadPath:  dir,
		CookieFile:    "cookies",
		ProgressFile:  "progress",
		InitialBookID: 1,
		Formats:       []string{"EPUB", "MOBI", "PDF"},
		Timeout:       10 * time.Minute,
		Retry:         5,
		UserAgent:     DefaultUserAgent,
		Rename:        false,
		Thread:        1,
	}
}

// BindDownloadArgs will bind the download arguments to cobra command.
func BindDownloadArgs(command *cobra.Command, config *Config) {
	command.Flags().StringVarP(&config.DownloadPath, "download", "d", config.DownloadPath,
		"The book directory you want to use, default would be current working directory.")

	command.Flags().StringVarP(&config.CookieFile, "cookie", "c", config.CookieFile,
		"The cookie file name you want to use, it would be saved under the download directory.")

	command.Flags().StringVarP(&config.ProgressFile, "progress", "g", config.ProgressFile,
		"The download progress file name you want to use, it would be saved under the download directory.")

	command.Flags().IntVarP(&config.InitialBookID, "initial", "i", config.InitialBookID,
		"The book id you want to start download. It should exceed 0.")

	command.Flags().StringSliceVarP(&config.Formats, "format", "f", config.Formats, "The file formats you want to download.")
	command.Flags().DurationVarP(&config.Timeout, "timeout", "o", config.Timeout, "The max pending time for download request.")
	command.Flags().IntVarP(&config.Retry, "retry", "r", config.Retry, "The max retry times for timeout download request.")
	command.Flags().StringVarP(&config.UserAgent, "user-agent", "a", config.UserAgent, "Set User-Agent for download request.")
	command.Flags().BoolVarP(&config.Rename, "rename", "n", config.Rename, "Rename the book file by book ID.")
	command.Flags().BoolVar(&config.Debug, "debug", config.Debug, "Enable debug mode")
}

// ValidateDownloadConfig would print the final download config table.
func ValidateDownloadConfig(config *Config) {
	if config.InitialBookID < 1 {
		log.Fatal(ErrInitialBookID)
	}
	if config.Retry < 1 {
		log.Fatal(ErrRetryTimes)
	}
	for i, format := range config.Formats {
		// Make sure all the format should be upper case.
		config.Formats[i] = strings.ToUpper(format)
	}
	if config.Thread < 1 {
		log.Fatal(ErrThreadCounts)
	}
}
