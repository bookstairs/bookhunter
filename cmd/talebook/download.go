package talebook

import (
	"errors"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/talebook"
)

// Used for downloading books from talebook website.
var downloadConfig = spider.NewDownloadConfig()

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
		log.PrintTable("Download Config Info", table.Row{"Config Key", "Config Value"}, downloadConfig)

		// Create the downloader
		downloader := talebook.NewDownloader(downloadConfig)

		// Start download books.
		downloader.Download()

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	// Add flags for use info.
	DownloadCmd.Flags().StringVarP(&downloadConfig.Website, "website", "w", "", "The talebook website.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.Username, "username", "u", "", "The account login name.")
	DownloadCmd.Flags().StringVarP(&downloadConfig.Password, "password", "p", "", "The account password.")

	spider.BindDownloadArgs(DownloadCmd, downloadConfig)

	_ = DownloadCmd.MarkFlagRequired("website")
}
