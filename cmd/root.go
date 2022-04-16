package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bookhunter",
	Short: "A command line base downloader for downloading books from internet.",
	Long: `You can use this command to download book from these websites.

1. Self-hosted talebook websites
2. https://www.sanqiu.cc
3. Telegram channel`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(talebookCmd)
	rootCmd.AddCommand(sanqiuCmd)
	rootCmd.AddCommand(telegramCmd)
	rootCmd.AddCommand(versionCmd)
}
