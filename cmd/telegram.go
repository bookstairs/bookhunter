package cmd

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/argument"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const defaultZone = "+86"

// telegramCmd used for download books from the telegram channel
var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "A tool for downloading books from telegram channel",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate the mobile number. Add default zone if need.
		if argument.Mobile != "" {
			if !strings.HasPrefix(argument.Mobile, "+") {
				argument.Mobile = defaultZone + argument.Mobile
			}
			for i, c := range argument.Mobile {
				if i > 0 && !unicode.IsDigit(c) {
					log.Fatalf("Invalid mobile number: %s", argument.Mobile)
				}
			}
		}

		// Remove prefix for telegram.
		argument.Website = argument.ChannelID
		argument.ChannelID = strings.TrimPrefix(argument.ChannelID, "https://t.me/")

		// Print download configuration.
		log.NewPrinter().
			Title("Telegram Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", argument.ConfigRoot).
			Row("Proxy", argument.Proxy).
			Row("UserAgent", argument.UserAgent).
			Row("Channel ID", argument.ChannelID).
			Row("Mobile", "******").
			Row("AppID", argument.AppID).
			Row("AppHash", argument.AppHash).
			Row("Formats", argument.Formats).
			Row("Extract Archive", argument.Extract).
			Row("Download Path", argument.DownloadPath).
			Row("Initial ID", argument.InitialBookID).
			Row("Rename File", argument.Rename).
			Row("Thread", argument.Thread).
			Row("Request Per Minute", argument.RateLimit).
			Print()

		// Create the fetcher.
		f, err := argument.NewFetcher(fetcher.Telegram, map[string]string{
			"channelID": argument.ChannelID,
			"mobile":    argument.Mobile,
			"reLogin":   strconv.FormatBool(argument.ReLogin),
			"appID":     strconv.FormatInt(argument.AppID, 10),
			"appHash":   argument.AppHash,
		})
		log.Fatal(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Fatal(err)

		// Finished all the tasks.
		log.Info("Successfully download all the telegram books.")
	},
}

func init() {
	flags := telegramCmd.Flags()

	// Telegram download arguments.
	flags.StringVarP(&argument.ChannelID, "channelID", "k", argument.ChannelID, "The channelId for telegram.")
	flags.StringVarP(&argument.Mobile, "mobile", "b", argument.Mobile, "The mobile number, default (+86).")
	flags.BoolVar(&argument.ReLogin, "refresh", argument.ReLogin, "Refresh the login session.")
	flags.Int64VarP(&argument.AppID, "appID", "", argument.AppID,
		"The appID for telegram. Refer https://core.telegram.org/api/obtaining_api_id to create your own appID")
	flags.StringVarP(&argument.AppHash, "appHash", "", argument.AppHash,
		"The appHash for telegram. Refer to https://core.telegram.org/api/obtaining_api_id to create your own appHash")

	// Common download flags.
	flags.StringSliceVarP(&argument.Formats, "format", "f", argument.Formats, "The file formats you want to download.")
	flags.BoolVarP(&argument.Extract, "extract", "e", argument.Extract, "Extract the archive file for filtering.")
	flags.StringVarP(&argument.DownloadPath, "download", "d", argument.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	flags.IntVarP(&argument.InitialBookID, "initial", "i", argument.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	flags.BoolVarP(&argument.Rename, "rename", "r", argument.Rename, "Rename the book file by book ID.")
	flags.IntVarP(&argument.Thread, "thread", "t", argument.Thread, "The number of concurrent download thead.")
	flags.IntVarP(&argument.RateLimit, "ratelimit", "", argument.RateLimit, "The request per minutes.")

	// Bind the required arguments
	_ = telegramCmd.MarkFlagRequired("channelID")
	_ = telegramCmd.MarkFlagRequired("appID")
	_ = telegramCmd.MarkFlagRequired("appHash")
}
