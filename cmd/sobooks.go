package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	lowestsobooksBookID = 18000
	sobooksWebsite      = "https://sobooks.net"
)

// sobooksCmd used for download books from sobooks.net
var sobooksCmd = &cobra.Command{
	Use:   "sobooks",
	Short: "A tool for downloading books from sobooks.net",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the default start index.
		if flags.InitialBookID < lowestsobooksBookID {
			flags.InitialBookID = lowestsobooksBookID
		}

		// Print download configuration.
		log.NewPrinter().
			Title("SoBooks Download Information").
			Head(log.DefaultHead...).
			Row("SoBooks Code", flags.SoBooksCode).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("UserAgent", flags.UserAgent).
			Row("Formats", flags.Formats).
			Row("Extract Archive", flags.Extract).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Request Per Minute", flags.RateLimit).
			Print()

		// Set the domain for using in the client.Client.
		flags.Website = sobooksWebsite
		flags.Driver = string(driver.LANZOU)

		// Create the fetcher.
		properties := flags.NewDriverProperties()
		properties["code"] = flags.SoBooksCode
		f, err := flags.NewFetcher(fetcher.SoBooks, properties)
		log.Exit(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	f := sobooksCmd.Flags()

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes")

	// SoBooks books flags.
	f.StringVar(&flags.SoBooksCode, "code", flags.SoBooksCode, "The secret code for SoBooks")

	_ = sobooksCmd.MarkFlagRequired("code")
}
