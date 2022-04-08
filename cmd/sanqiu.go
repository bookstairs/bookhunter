package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/sanqiu"
)

// Used for downloading books from sanqiu website.
var c = spider.NewDownloadConfig()

// sanqiuCmd used for download books from sanqiu.com
var sanqiuCmd = &cobra.Command{
	Use:   "sanqiu",
	Short: "A tool for downloading books from sanqiu.com",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config
		spider.ValidateDownloadConfig(c)

		// Create the downloader
		downloader := sanqiu.NewDownloader(c)

		// Start download books.
		downloader.Download()

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	sanqiuCmd.Flags().StringVarP(&sanqiu.Website, "website", "w", sanqiu.Website,
		"The website for sanqiu. You don't need to override the default url.")

	// Set common download config arguments.
	spider.BindDownloadArgs(sanqiuCmd, c)

	sanqiuCmd.Flags().IntVarP(&c.Thread, "thread", "t", c.Thread, "The number of download threads.")
}
