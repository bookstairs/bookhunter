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
			panic(err)
		}

		//for i := 0; i < c.Thread; i++ {
		//	// Create a thread.
		//	downloader.Fork()
		//	// Download books in this thread.
		//	go downloader.Download()
		//}
		//
		//// Wait all the thread have finished.
		//downloader.Join()

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
	telegramCmd.Flags().IntVar(&telegram.ChunkSize, "chunkSize", telegram.ChunkSize,
		"The ChunkSize for download telegram. 4096 < ChunkSize < 512 * 1024")
	telegramCmd.Flags().IntVar(&telegram.LoadMessageSize, "loadMessageSize", telegram.LoadMessageSize,
		"The loadMessageSize is used to set the size of the number of messages obtained by requesting telegram API. 0 < loadMessageSize < 100")

	// Set common download config arguments.
	spider.BindDownloadArgs(telegramCmd, d)
	telegramCmd.Flags().IntVarP(&c.Thread, "thread", "t", c.Thread, "The number of download threads.")
}