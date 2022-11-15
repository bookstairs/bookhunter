package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chyroc/go-aliyundrive"
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

var aliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "A command line tool for acquiring the refresh token from aliyundrive with QR code login.",
	Run: func(cmd *cobra.Command, args []string) {
		// Create the config path.
		config := &client.Config{Host: "api.aliyundrive.com"}
		path, err := config.ConfigPath()
		if err != nil {
			log.Fatal(err)
		}

		// Create the session file if it's not existed.
		file := filepath.Join(path, "session.json")
		open, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		_ = open.Close()

		// Create the aliyun client.
		store := aliyundrive.NewFileStore(file)
		ins := aliyundrive.New(aliyundrive.WithStore(store))
		ctx := context.Background()

		// Valid the token, we will sign in with QR code if this token is expired.
		user, err := ins.Auth.LoginByQrcode(ctx, &aliyundrive.LoginByQrcodeReq{})
		if err != nil {
			log.Fatal(err)
		}

		// Access the token from the storage.
		token, err := store.Get(ctx, "")
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Successfully sign into aliyun drive.")
		log.Infof("%s: %s", user.NickName, flags.HideSensitive(user.Phone))
		log.Infof("Refresh Token: %s", token.RefreshToken)
	},
}
