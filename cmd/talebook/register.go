package talebook

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/argument"
	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/talebook"
)

// RegisterCmd represents the register command.
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account on talebook.",
	Long: `Some talebook website need a user account for downloading books.
You can use this register command for creating account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print register configuration.
		log.NewPrinter().
			Title("Talebook Register Information").
			Head(log.DefaultHead...).
			Row("Website", argument.Website).
			Row("Username", argument.Username).
			Row("Password", argument.Password).
			Row("Email", argument.Email).
			Row("Config Path", argument.ConfigRoot).
			Row("UserAgent", argument.UserAgent).
			Row("Proxy", argument.Proxy).
			Print()

		// Create client config.
		config, err := client.NewConfig(argument.Website, argument.UserAgent, argument.Proxy, argument.ConfigRoot)
		log.Fatal(err)

		// Create http client.
		c, err := client.New(config)
		log.Fatal(err)

		// Execute the register request.
		resp, err := c.R().
			SetFormData(map[string]string{
				"username": argument.Username,
				"password": argument.Password,
				"nickname": argument.Username,
				"email":    argument.Email,
			}).
			SetResult(&talebook.CommonResp{}).
			ForceContentType("application/json").
			Post("/api/user/sign_up")
		log.Fatal(err)

		result := resp.Result().(*talebook.CommonResp)
		if result.Err == talebook.SuccessStatus {
			log.Info("Register success.")
		} else {
			log.Fatalf("Register failed, reason: %s", result.Err)
		}
	},
}

func init() {
	flags := RegisterCmd.Flags()

	// Add flags for registering.
	flags.StringVarP(&argument.Username, "username", "u", argument.Username, "The account login name.")
	flags.StringVarP(&argument.Password, "password", "p", argument.Password, "The account password.")
	flags.StringVarP(&argument.Email, "email", "e", argument.Email, "The account email.")
	flags.StringVarP(&argument.Website, "website", "w", argument.Website, "The talebook website.")

	// Mark some flags as required.
	_ = RegisterCmd.MarkFlagRequired("website")
	_ = RegisterCmd.MarkFlagRequired("username")
	_ = RegisterCmd.MarkFlagRequired("password")
	_ = RegisterCmd.MarkFlagRequired("email")
}
