package log

import (
	"fmt"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// NewProgressBar is used to print a beautiful download progress.
func NewProgressBar(index, total int64, filename string, bytes int64) *progressbar.ProgressBar {
	title := fmt.Sprintf("%s %s [%d/%d] %s", logTime(), info, index, total, filename)

	return progressbar.NewOptions64(bytes,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\n")
		}),
		progressbar.OptionSetDescription(title),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
