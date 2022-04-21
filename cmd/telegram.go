package cmd

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/telegram"
)

// Used for downloading books from telegram channel .
var tc = telegram.NewConfig()

// telegramCmd used for download books from telegram channel
var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "A tool for downloading books from telegram channel",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config
		spider.ValidateDownloadConfig(tc.Config)

		// Remove prefix.
		tc.ChannelID = strings.TrimPrefix(tc.ChannelID, "https://t.me/")

		// Add global zone in phone number.
		tc.Mobile = telegram.AddCountryCode(tc.Mobile)

		// Print download configuration.
		log.PrintTable("Download Config Info", table.Row{"Config Key", "Config Value"}, tc, false)

		// Create the downloader and download books.
		downloader := telegram.NewDownloader(tc)

		for i := 0; i < tc.Thread; i++ {
			// Create a thread and download books in this thread.
			downloader.Fork()
		}

		// Wait all the thread have finished.
		downloader.Join()

		// Finished all the tasks.
		log.Info("Successfully download all the telegram books.")
	},
}

func init() {
	telegramCmd.Flags().StringVarP(&tc.ChannelID, "channelID", "k", "", "The channelId for telegram.")
	telegramCmd.Flags().StringVarP(&tc.Mobile, "mobile", "b", "", "The mobile number for your telegram account, default (+86).")
	telegramCmd.Flags().BoolVar(&tc.Refresh, "refresh", tc.Refresh, "Refresh the login session.")
	telegramCmd.Flags().IntVar(&tc.AppID, "appID", 0,
		"The appID for telegram. How to get `appID` please refer to https://core.telegram.org/api/obtaining_api_id.")
	telegramCmd.Flags().StringVar(&tc.AppHash, "appHash", "",
		"The appHash for telegram. How to get `appHash` please refer to https://core.telegram.org/api/obtaining_api_id.")
	telegramCmd.Flags().StringVarP(&tc.CookieFile, "sessionPath", "s", tc.CookieFile, "The session file for telegram.")

	// Set common download config arguments.
	spider.BindDownloadArgs(telegramCmd, tc.Config)

	// Support multiple thread download.
	telegramCmd.Flags().IntVarP(&tc.Thread, "thread", "t", tc.Thread, "The number of download threads.")

	// Bind the required arguments
	_ = telegramCmd.MarkFlagRequired("channelID")
	_ = telegramCmd.MarkFlagRequired("appID")
	_ = telegramCmd.MarkFlagRequired("appHash")
}
