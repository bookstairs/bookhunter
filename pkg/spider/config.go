package spider

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type DownloadConfig struct {
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
}

// NewDownloadConfig will return a default blank config.
func NewDownloadConfig() *DownloadConfig {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &DownloadConfig{
		DownloadPath:  dir,
		CookieFile:    "cookies",
		ProgressFile:  "progress",
		InitialBookID: 1,
		Formats:       []string{"EPUB", "MOBI", "PDF"},
		Timeout:       10 * time.Minute,
		Retry:         5,
		UserAgent:     DefaultUserAgent,
		Rename:        false,
	}
}

// BindDownloadArgs will bind the download arguments to cobra command.
func BindDownloadArgs(command *cobra.Command, config *DownloadConfig) {
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
}
