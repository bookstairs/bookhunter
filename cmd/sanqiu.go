package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	lowestSanqiuBookID = 163
	sanqiuWebsite      = "https://www.sanqiu.mobi"
)

// sanqiuCmd used for download books from sanqiu.mobi
var sanqiuCmd = &cobra.Command{
	Use:   "sanqiu",
	Short: "A tool for downloading books from sanqiu.mobi",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the default start index.
		if flags.InitialBookID < lowestSanqiuBookID {
			flags.InitialBookID = lowestSanqiuBookID
		}

		// Print download configuration.
		log.NewPrinter().
			Title("Sanqiu Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("UserAgent", flags.UserAgent).
			Row("Formats", flags.Formats).
			Row("Extract Archive", flags.Extract).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Row("Aliyun RefreshToken", flags.HideSensitive(flags.RefreshToken)).
			Row("Telecom Username", flags.HideSensitive(flags.Username)).
			Row("Telecom Password", flags.HideSensitive(flags.Password)).
			Print()

		// Set the domain for using in the client.Client.
		flags.Website = sanqiuWebsite

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
	f := sanqiuCmd.Flags()

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")

	// Drive ISP flags.
	f.StringVar(&flags.Driver, "source", flags.Driver, "The source (aliyun, telecom, lanzou) to download book")
	f.StringVar(&flags.RefreshToken, "refreshToken", flags.RefreshToken, "Refresh token for aliyun drive")
	f.StringVar(&flags.TelecomUsername, "telecomUsername", flags.TelecomUsername, "Telecom drive username")
	f.StringVar(&flags.TelecomPassword, "telecomPassword", flags.TelecomPassword, "Telecom drive password")

	_ = sanqiuCmd.MarkFlagRequired("source")
}
