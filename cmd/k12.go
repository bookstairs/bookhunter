package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const k12Website = "https://www.zxx.edu.cn"

var k12Cmd = &cobra.Command{
	Use:   "k12",
	Short: "A tool for downloading textbooks from www.zxx.edu.cn",
	Run: func(cmd *cobra.Command, args []string) {
		// Print download configuration.
		log.NewPrinter().
			Title("Textbook Download Information").
			Head(log.DefaultHead...).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("UserAgent", flags.UserAgent).
			Row("Download Path", flags.DownloadPath).
			Row("Thread", flags.Thread).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Row("Keywords", flags.Keywords).
			Print()

		flags.Website = k12Website
		f, err := flags.NewFetcher(fetcher.K12, map[string]string{})
		log.Exit(err)

		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the textbooks.")
	},
}

func init() {
	f := k12Cmd.Flags()

	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")
}
