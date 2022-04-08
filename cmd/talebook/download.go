package talebook

import (
	"errors"
	"strings"

	"github.com/syhily/bookhunter/pkg/log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/syhily/bookhunter/talebook"
)

// Used for downloading books from talebook website.
var downloadConfig = *talebook.NewDownloadConfig()

var (
	ErrInitialBookID = errors.New("illegal book id, it should exceed 0")
	ErrRetryTimes    = errors.New("illegal retry times, it should exceed 0")
)

// DownloadCmd represents the download command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the book from talebook.",
	Run: func(cmd *cobra.Command, args []string) {
		if downloadConfig.InitialBookID < 1 {
			log.Fatal(ErrInitialBookID)
		}
		if downloadConfig.Retry < 1 {
			log.Fatal(ErrRetryTimes)
		}
		for i, format := range downloadConfig.Formats {
			// Make sure all the format should be upper case.
			downloadConfig.Formats[i] = strings.ToUpper(format)
		}

		// Print download configuration.
		log.PrintTable("Download Config Info", table.Row{"Config Key", "Config Value"}, &downloadConfig)

		// Create the downloader
		talebook := talebook.NewTalebook(&downloadConfig)

		// Start download books.
		talebook.Start()

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	// Add flags for use info.
	DownloadCmd.Flags().StringVarP(&downloadConfig.Website, "website", "w", "", "The talebook website.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.Username, "username", "u", "", "The account login name.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.Password, "password", "p", "", "The account password.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.DownloadPath, "download", "d", downloadConfig.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.CookieFile, "cookie", "c", downloadConfig.CookieFile,
		"The cookie file name you want to use, it would be saved under the download directory.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.ProgressFile, "progress", "g", downloadConfig.ProgressFile,
		"The download progress file name you want to use, it would be saved under the download directory.")
	DownloadCmd.Flags().IntVarP(&downloadConfig.InitialBookID, "initial", "i", downloadConfig.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	DownloadCmd.Flags().StringSliceVarP(&downloadConfig.Formats, "format", "f", downloadConfig.Formats,
		"The file formats you want to download.")
	DownloadCmd.Flags().DurationVarP(&downloadConfig.Timeout, "timeout", "o", downloadConfig.Timeout,
		"The max pending time for download request.")
	DownloadCmd.Flags().IntVarP(&downloadConfig.Retry, "retry", "r", downloadConfig.Retry, "The max retry times for timeout download request.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.UserAgent, "user-agent", "a", downloadConfig.UserAgent,
		"Set User-Agent for download request.")
	DownloadCmd.Flags().BoolVarP(&downloadConfig.Rename, "rename", "n", downloadConfig.Rename, "Rename the book file by book ID.")

	_ = DownloadCmd.MarkFlagRequired("website")
}
