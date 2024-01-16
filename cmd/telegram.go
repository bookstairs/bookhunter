package cmd

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

// telegramCmd used for download books from the telegram channel
var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "A tool for downloading books from telegram channel",
	Run: func(cmd *cobra.Command, args []string) {
		// Remove prefix for telegram.
		flags.Website = flags.ChannelID
		flags.ChannelID = strings.TrimPrefix(flags.ChannelID, "https://t.me/")

		// Print download configuration.
		log.NewPrinter().
			Title("Telegram Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("Channel ID", flags.ChannelID).
			Row("Mobile", flags.HideSensitive(flags.Mobile)).
			Row("AppID", flags.HideSensitive(strconv.FormatInt(flags.AppID, 10))).
			Row("AppHash", flags.HideSensitive(flags.AppHash)).
			Row("Formats", flags.Formats).
			Row("Extract Archive", flags.Extract).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Keywords", flags.Keywords).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Print()

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.Telegram, map[string]string{
			"channelID": flags.ChannelID,
			"mobile":    flags.Mobile,
			"reLogin":   strconv.FormatBool(flags.ReLogin),
			"appID":     strconv.FormatInt(flags.AppID, 10),
			"appHash":   flags.AppHash,
		})
		log.Exit(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the telegram books.")
	},
}

func init() {
	f := telegramCmd.Flags()

	// Telegram download arguments.
	f.StringVarP(&flags.ChannelID, "channelID", "k", flags.ChannelID, "The channel id for telegram")
	f.StringVarP(&flags.Mobile, "mobile", "b", flags.Mobile, "The mobile number, we will add +86 as default zone code")
	f.BoolVar(&flags.ReLogin, "refresh", flags.ReLogin, "Refresh the login session")
	f.Int64Var(&flags.AppID, "appID", flags.AppID, "The app id for telegram")
	f.StringVar(&flags.AppHash, "appHash", flags.AppHash, "The app hash for telegram")

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")

	// Bind the required arguments
	_ = telegramCmd.MarkFlagRequired("channelID")
	_ = telegramCmd.MarkFlagRequired("appID")
	_ = telegramCmd.MarkFlagRequired("appHash")
}
