package cmd

import (
	"github.com/spf13/cobra"

	"github.com/syhily/bookhunter/cmd/talebook"
)

// talebookCmd used to download books from talebook
var talebookCmd = &cobra.Command{
	Use:   "talebook",
	Short: "A command line base downloader for downloading books from talebook server.",
	Long: `You can use this command to register account and download book.
The url for talebook should be provided, the formats is also
optional.`,
}

func init() {
	talebookCmd.AddCommand(talebook.DownloadCmd)
	talebookCmd.AddCommand(talebook.RegisterCmd)
}
