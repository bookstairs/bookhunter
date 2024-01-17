package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const hsuWebsite = "https://book.hsu.life"

var hsuCmd = &cobra.Command{
	Use:   "hsu",
	Short: "A tool for downloading book from hsu.life",
	Run: func(cmd *cobra.Command, args []string) {
		log.NewPrinter().
			Title("hsu.life Download Information").
			Head(log.DefaultHead...).
			Row("Username", flags.Username).
			Row("Password", flags.HideSensitive(flags.Password)).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("Formats", flags.Formats).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Keywords", flags.Keywords).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Print()

		flags.Website = hsuWebsite

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.Hsu, map[string]string{
			"username": flags.Username,
			"password": flags.Password,
		})
		log.Exit(err)

		// Start downloading the books.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	// Add flags for use info.
	f := hsuCmd.Flags()

	// Talebook related flags.
	f.StringVarP(&flags.Username, "username", "u", flags.Username, "The hsu.life username")
	f.StringVarP(&flags.Password, "password", "p", flags.Password, "The hsu.life password")

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")

	// Mark some flags as required.
	_ = hsuCmd.MarkFlagRequired("username")
	_ = hsuCmd.MarkFlagRequired("password")
}
