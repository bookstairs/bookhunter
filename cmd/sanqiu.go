package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	lowestBookID  = 163
	sanqiuWebsite = "https://www.sanqiu.mobi"
)

// sanqiuCmd used for download books from sanqiu.mobi
var sanqiuCmd = &cobra.Command{
	Use:   "sanqiu",
	Short: "A tool for downloading books from sanqiu.mobi",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the default start index.
		if flags.InitialBookID < lowestBookID {
			flags.InitialBookID = lowestBookID
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
			Row("Request Per Minute", flags.RateLimit).
			Row("Aliyun RefreshToken", flags.HideSensitive(flags.RefreshToken)).
			Row("Telecom Username", flags.HideSensitive(flags.TelecomUsername)).
			Row("Telecom Password", flags.HideSensitive(flags.TelecomPassword)).
			Print()

		// Set the domain for using in the client.Client.
		flags.Website = sanqiuWebsite

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.SanQiu, map[string]string{
			"refreshToken":    flags.RefreshToken,
			"telecomUsername": flags.TelecomUsername,
			"telecomPassword": flags.TelecomPassword,
		})
		log.Fatal(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Fatal(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	f := sanqiuCmd.Flags()

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download.")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering.")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book ID.")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of concurrent download thead.")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The request per minutes.")

	// Drive ISP flags.
	f.StringVar(&flags.RefreshToken, "refreshToken", flags.RefreshToken,
		"We would try to download from the aliyun drive if you provide this token.")
	f.StringVar(&flags.TelecomUsername, "telecomUsername", flags.TelecomUsername,
		"Used to download file from telecom drive")
	f.StringVar(&flags.TelecomPassword, "telecomPassword", flags.TelecomPassword,
		"Used to download file from telecom drive")
}
