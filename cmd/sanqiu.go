package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/sanqiu"

	"github.com/bibliolater/bookhunter/pkg/spider"
)

// Used for downloading books from sanqiu website.
var c = spider.NewDownloadConfig()

// sanqiuCmd used for download books from sanqiu.com
var sanqiuCmd = &cobra.Command{
	Use:   "sanqiu",
	Short: "A tool for downloading books from sanqiu.com",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	sanqiuCmd.Flags().StringVarP(&sanqiu.Website, "website", "w", sanqiu.Website, "The website for sanqiu.")
	// Set common download config arguments.
	spider.BindDownloadArgs(sanqiuCmd, c)
}
