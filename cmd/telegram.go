package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/telegram"
)

// Used for downloading books from telegram channel .
var d = spider.NewConfig()

// telegramCmd used for download books from telegram channel
var telegramCmd = &cobra.Command{
	Use:   "telegram",
	Short: "A tool for downloading books from telegram channel",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate config
		spider.ValidateDownloadConfig(d)

		// Create the downloader
		downloader := telegram.NewDownloader(d)

		err := downloader.Exec()
		if err != nil {
			log.Fatal(err)
		}
		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	telegramCmd.Flags().StringVarP(&telegram.ChannelId, "channelId", "k", telegram.ChannelId,
		"The channelId for telegram. You must set value.")
	telegramCmd.Flags().StringVarP(&telegram.SessionPath, "sessionPath", "s", telegram.SessionPath,
		"The session file for telegram.")
	telegramCmd.Flags().BoolVar(&telegram.ReLogin, "reLogin", telegram.ReLogin,
		"force re-login.")
	telegramCmd.Flags().IntVar(&telegram.AppID, "appId", telegram.AppID,
		"The appID for telegram.")
	telegramCmd.Flags().StringVar(&telegram.AppHash, "appHash", telegram.AppHash,
		"The appHash for telegram.")
	telegramCmd.Flags().IntVar(&telegram.LoadMessageSize, "loadMessageSize", telegram.LoadMessageSize,
		"The loadMessageSize is used to set the size of the number of messages obtained by requesting telegram API. 0 < loadMessageSize < 100")

	// Set common download config arguments.
	spider.BindDownloadArgs(telegramCmd, d)
	telegramCmd.Flags().IntVarP(&d.Thread, "thread", "t", d.Thread, "The number of download threads.")
}
