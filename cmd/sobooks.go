package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	lowestSobooksBookID = 18000
)

// sobooksCmd used for download books from sobooks.cc
var sobooksCmd = &cobra.Command{
	Use:   "sobooks",
	Short: "A tool for downloading books from sobooks.cc",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the default start index.
		if flags.InitialBookID < lowestSobooksBookID {
			flags.InitialBookID = lowestSobooksBookID
		}

		// Print download configuration.
		log.NewPrinter().
			Title("SoBooks Download Information").
			Head(log.DefaultHead...).
			Row("SoBooks Code", flags.SoBooksCode).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("Formats", flags.Formats).
			Row("Extract Archive", flags.Extract).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Keywords", flags.Keywords).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Print()

		// 构建参数 map
		cmdArgs := map[string]interface{}{
			"code":          flags.SoBooksCode,
			"formats":       flags.Formats,
			"extract":       flags.Extract,
			"directory":     flags.DownloadPath,
			"initialBookID": flags.InitialBookID,
			"rename":        flags.Rename,
			"thread":        flags.Thread,
			"ratelimit":     flags.RateLimit,
			"keywords":      flags.Keywords,
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
			Category:   fetcher.SoBooks,
			Args:       cmdArgs,
			GlobalArgs: globalArgs,
		}

		f, err := builder.Build()
		log.Exit(err)

		// Wait all the threads have finished.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

func init() {
	f := sobooksCmd.Flags()

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.BoolVarP(&flags.Extract, "extract", "e", flags.Extract, "Extract the archive file for filtering")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")

	// SoBooks books flags.
	f.StringVar(&flags.SoBooksCode, "code", flags.SoBooksCode, "The secret code for SoBooks")

	_ = sobooksCmd.MarkFlagRequired("code")
}
