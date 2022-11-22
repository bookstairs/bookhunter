package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	youyiduWebsite = "https://www.youyidu.xyz/"
)

// youyiduCmd used for download books from sanqiu.mobi
var youyiduCmd = &cobra.Command{
	Use:   "youyidu",
	Short: "A tool for downloading books from youyidu.xyz",
	Run: func(cmd *cobra.Command, args []string) {
		// Print download configuration.
		log.NewPrinter().
			Title("Youyidu Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("UserAgent", flags.UserAgent).
			Row("Formats", flags.Formats).
			Row("Extract Archive", flags.Extract).
			Row("Rename File", flags.Rename).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Thread", flags.Thread).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Print()

		// Set the domain for using in the client.Client.
		flags.Website = youyiduWebsite

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.SanQiu, flags.NewDriverProperties())
		log.Exit(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	f := youyiduCmd.Flags()

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")

	// Drive ISP flags.
	f.StringVar(&flags.Driver, "source", flags.Driver, "The source (aliyun, telecom, lanzou) to download book")
	f.StringVar(&flags.RefreshToken, "refreshToken", flags.RefreshToken, "Refresh token for aliyun drive")
	f.StringVar(&flags.TelecomUsername, "telecomUsername", flags.TelecomUsername, "Telecom drive username")
	f.StringVar(&flags.TelecomPassword, "telecomPassword", flags.TelecomPassword, "Telecom drive password")

	_ = youyiduCmd.MarkFlagRequired("source")
}
