package talebook

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/pkg/log"
	"github.com/bookstairs/bookhunter/pkg/spider"
	"github.com/bookstairs/bookhunter/talebook"
)

// Used for downloading books from talebook website.
var downloadConfig = spider.NewConfig()

// DownloadCmd represents the download command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the book from talebook.",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config
		spider.ValidateDownloadConfig(downloadConfig)

		// Print download configuration.
		log.PrintTable("Download Config Info", table.Row{"Config Key", "Config Value"}, downloadConfig, true)

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
