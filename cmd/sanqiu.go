package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/argument"
	"github.com/bookstairs/bookhunter/internal/client"
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
		if argument.InitialBookID < lowestBookID {
			argument.InitialBookID = lowestBookID
		}

		// Print download configuration.
		log.NewPrinter().
			Title("Sanqiu Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", argument.ConfigRoot).
			Row("Proxy", argument.Proxy).
			Row("UserAgent", argument.UserAgent).
			Row("Formats", argument.Formats).
			Row("Extract Archive", argument.Extract).
			Row("Download Path", argument.DownloadPath).
			Row("Initial ID", argument.InitialBookID).
			Row("Rename File", argument.Rename).
			Row("Thread", argument.Thread).
			Row("Aliyun RefreshToken", "******").
			Print()

		// Create the fetcher config.
		cc, err := client.NewConfig(sanqiuWebsite, argument.UserAgent, argument.Proxy, argument.ConfigRoot)
		log.Fatal(err)
		fs, err := fetcher.ParseFormats(argument.Formats)
		log.Fatal(err)

		// Create the fetcher.
		f, err := fetcher.New(&fetcher.Config{
			Config:        cc,
			Category:      fetcher.SanQiu,
			Formats:       fs,
			Extract:       argument.Extract,
			DownloadPath:  argument.DownloadPath,
			InitialBookID: argument.InitialBookID,
			Rename:        argument.Rename,
			Thread:        argument.Thread,
			Properties: map[string]string{
				"refreshToken": argument.RefreshToken,
			},
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
	flags := sanqiuCmd.Flags()

	// Common download flags.
	flags.StringSliceVarP(&argument.Formats, "format", "f", argument.Formats, "The file formats you want to download.")
	flags.BoolVarP(&argument.Extract, "extract", "e", argument.Extract, "Extract the archive file for filtering.")
	flags.StringVarP(&argument.DownloadPath, "download", "d", argument.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	flags.IntVarP(&argument.InitialBookID, "initial", "i", argument.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	flags.BoolVarP(&argument.Rename, "rename", "r", argument.Rename, "Rename the book file by book ID.")
	flags.IntVarP(&argument.Thread, "thread", "t", argument.Thread, "The number of concurrent download thead.")

	// Drive ISP flags.
	flags.StringVarP(&argument.RefreshToken, "refreshToken", "", argument.RefreshToken,
		"We would try to download from the aliyun drive if you provide this token.")
}
