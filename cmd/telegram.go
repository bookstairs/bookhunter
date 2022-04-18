package cmd

import (
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

		// Create the downloader and download books.
		downloader := telegram.NewDownloader(tc)

		for i := 0; i < c.Thread; i++ {
			// Create a thread.
			downloader.Fork()
			// Download books in this thread.
			go downloader.Download()
		}

		// Wait all the thread have finished.
		downloader.Join()

		// Finished all the tasks.
		log.Info("Successfully download all the telegram books.")
	},
}

func init() {
	telegramCmd.Flags().StringVarP(&tc.ChannelID, "channelId", "k", "", "The channelId for telegram.")
	telegramCmd.Flags().StringVarP(&tc.Mobile, "mobile", "b", "", "The mobile number for your telegram account, default (+86)")
	telegramCmd.Flags().BoolVar(&tc.Refresh, "reLogin", tc.Refresh, "Refresh the login session.")
	telegramCmd.Flags().IntVar(&tc.AppID, "appId", 0,
		"The appID for telegram. How to get `appId` please refer to https://core.telegram.org/api/obtaining_api_id")
	telegramCmd.Flags().StringVar(&tc.AppHash, "appHash", "",
		"The appHash for telegram. How to get `appHash` please refer to https://core.telegram.org/api/obtaining_api_id")
	telegramCmd.Flags().StringVarP(&tc.CookieFile, "sessionPath", "s", tc.CookieFile, "The session file for telegram.")

	// Set common download config arguments.
	spider.BindDownloadArgs(telegramCmd, tc.Config)

	// Support multiple thread download.
	telegramCmd.Flags().IntVarP(&tc.Thread, "thread", "t", tc.Thread, "The number of download threads.")

	// Bind the required arguments
	_ = telegramCmd.MarkFlagRequired("channelId")
	_ = telegramCmd.MarkFlagRequired("appId")
	_ = telegramCmd.MarkFlagRequired("appHash")
}
