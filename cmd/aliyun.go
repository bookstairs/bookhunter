package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/driver/aliyun"
	"github.com/bookstairs/bookhunter/internal/log"
)

var aliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "A command line tool for acquiring the refresh token from aliyundrive with QR code login",
	Run: func(cmd *cobra.Command, args []string) {
		// Create the client config.
		flags.Website = "https://api.aliyundrive.com"
		c, err := flags.NewClientConfig()
		if err != nil {
			log.Exit(err)
		}

		// Perform login.
		aliyun, err := aliyun.New(c, "")
		if err != nil {
			log.Exit(err)
		}
		log.Info("RefreshToken: ", aliyun.RefreshToken())
	},
}
