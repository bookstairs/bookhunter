package talebook

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/argument"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

// DownloadCmd represents the download command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the book from talebook.",
	Run: func(cmd *cobra.Command, args []string) {
		// Print download configuration.
		log.NewPrinter().
			Title("Talebook Download Information").
			Head(log.DefaultHead...).
			Row("Website", argument.Website).
			Row("Username", argument.Username).
			Row("Password", "******").
			Row("Config Path", argument.ConfigRoot).
			Row("Proxy", argument.Proxy).
			Row("UserAgent", argument.UserAgent).
			Row("Formats", argument.Formats).
			Row("Download Path", argument.DownloadPath).
			Row("Initial ID", argument.InitialBookID).
			Row("Rename File", argument.Rename).
			Row("Thread", argument.Thread).
			Print()

		// Create the fetcher.
		f, err := argument.NewFetcher(fetcher.Talebook, map[string]string{
			"username": argument.Username,
			"password": argument.Password,
		})
		log.Fatal(err)

		// Start download the books.
		err = f.Download()
		log.Fatal(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	// Add flags for use info.
	flags := DownloadCmd.Flags()

	// Talebook related flags.
	flags.StringVarP(&argument.Username, "username", "u", argument.Username, "The account login name.")
	flags.StringVarP(&argument.Password, "password", "p", argument.Password, "The account password.")
	flags.StringVarP(&argument.Website, "website", "w", argument.Website, "The talebook website.")

	// Common download flags.
	flags.StringSliceVarP(&argument.Formats, "format", "f", argument.Formats, "The file formats you want to download.")
	flags.StringVarP(&argument.DownloadPath, "download", "d", argument.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	flags.IntVarP(&argument.InitialBookID, "initial", "i", argument.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	flags.BoolVarP(&argument.Rename, "rename", "r", argument.Rename, "Rename the book file by book ID.")
	flags.IntVarP(&argument.Thread, "thread", "t", argument.Thread, "The number of concurrent download thead.")

	// Mark some flags as required.
	_ = DownloadCmd.MarkFlagRequired("website")
}
