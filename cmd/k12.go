package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

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
			Row("Download Path", flags.DownloadPath).
			Row("Thread", flags.Thread).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Row("Keywords", flags.Keywords).
			Print()

		// 构建参数 map
		cmdArgs := map[string]interface{}{
			"formats":   flags.Formats,
			"extract":   flags.Extract,
			"directory": flags.DownloadPath,
			"rename":    flags.Rename,
			"thread":    flags.Thread,
			"ratelimit": flags.RateLimit,
			"keywords":  flags.Keywords,
		}

		// 全局参数
		globalArgs := &fetcher.GlobalArgs{
			ConfigRoot: flags.ConfigRoot,
			Proxy:      flags.Proxy,
			Retry:      flags.Retry,
			SkipError:  flags.SkipError,
		}

		// 使用 FetcherBuilder 创建 fetcher
		builder := &fetcher.FetcherBuilder{
			Category:   fetcher.K12,
			Args:       cmdArgs,
			GlobalArgs: globalArgs,
		}

		f, err := builder.Build()
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
