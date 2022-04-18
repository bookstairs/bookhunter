package cmd

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/sanqiu"
)

const lowestBookID = 163

// Used for downloading books from sanqiu website.
var c = spider.NewConfig()

// sanqiuCmd used for download books from sanqiu.cc
var sanqiuCmd = &cobra.Command{
	Use:   "sanqiu",
	Short: "A tool for downloading books from sanqiu.cc",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config.
		spider.ValidateDownloadConfig(c)

		// Set the default start index.
		if c.InitialBookID < lowestBookID {
			c.InitialBookID = lowestBookID
		}

		// Print download configuration.
		log.PrintTable("Download Config Info", table.Row{"Config Key", "Config Value"}, c, false)

		// Create the downloader.
		downloader := sanqiu.NewDownloader(c)

		for i := 0; i < c.Thread; i++ {
			// Create a thread and download books in this thread.
			downloader.Fork()
		}

		// Wait all the thread have finished.
		downloader.Join()

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	sanqiuCmd.Flags().StringVarP(&c.Website, "website", "w", sanqiu.DefaultWebsite,
		"The website for sanqiu. You don't need to override the default url.")

	// Set common download config arguments.
	spider.BindDownloadArgs(sanqiuCmd, c)

	sanqiuCmd.Flags().IntVarP(&c.Thread, "thread", "t", c.Thread, "The number of download threads.")
}
