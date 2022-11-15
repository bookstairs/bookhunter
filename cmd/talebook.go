package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/talebook"
)

// talebookCmd used to download books from talebook
var talebookCmd = &cobra.Command{
	Use:   "talebook",
	Short: "A command line base downloader for downloading books from talebook server.",
	Long: `You can use this command to register account and download book.
The url for talebook should be provided, the formats is also
optional.`,
}

// talebookDownloadCmd represents the download command
var talebookDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the book from talebook.",
	Run: func(cmd *cobra.Command, args []string) {
		// Print download configuration.
		log.NewPrinter().
			Title("Talebook Download Information").
			Head(log.DefaultHead...).
			Row("Website", flags.Website).
			Row("Username", flags.HideSensitive(flags.Username)).
			Row("Password", flags.HideSensitive(flags.Password)).
			Row("Config Path", flags.ConfigRoot).
			Row("Proxy", flags.Proxy).
			Row("UserAgent", flags.UserAgent).
			Row("Formats", flags.Formats).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Request Per Minute", flags.RateLimit).
			Print()

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.Talebook, map[string]string{
			"username": flags.Username,
			"password": flags.Password,
		})
		log.Fatal(err)

		// Start download the books.
		err = f.Download()
		log.Fatal(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

// talebookRegisterCmd represents the register command.
var talebookRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account on talebook.",
	Long: `Some talebook website need a user account for downloading books.
You can use this register command for creating account.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Print register configuration.
		log.NewPrinter().
			Title("Talebook Register Information").
			Head(log.DefaultHead...).
			Row("Website", flags.Website).
			Row("Username", flags.Username).
			Row("Password", flags.Password).
			Row("Email", flags.Email).
			Row("Config Path", flags.ConfigRoot).
			Row("UserAgent", flags.UserAgent).
			Row("Proxy", flags.Proxy).
			Print()

		// Create client config.
		config, err := client.NewConfig(flags.Website, flags.UserAgent, flags.Proxy, flags.ConfigRoot)
		log.Fatal(err)

		// Create http client.
		c, err := client.New(config)
		log.Fatal(err)

		// Execute the register request.
		resp, err := c.R().
			SetFormData(map[string]string{
				"username": flags.Username,
				"password": flags.Password,
				"nickname": flags.Username,
				"email":    flags.Email,
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
	/// Add the download command.

	// Add flags for use info.
	f := talebookDownloadCmd.Flags()

	// Talebook related flags.
	f.StringVarP(&flags.Username, "username", "u", flags.Username, "The account login name.")
	f.StringVarP(&flags.Password, "password", "p", flags.Password, "The account password.")
	f.StringVarP(&flags.Website, "website", "w", flags.Website, "The talebook website.")

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download.")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book ID.")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of concurrent download thead.")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The request per minutes.")

	// Mark some flags as required.
	_ = talebookDownloadCmd.MarkFlagRequired("website")

	talebookCmd.AddCommand(talebookDownloadCmd)

	/// Add the register command.

	f = talebookRegisterCmd.Flags()

	// Add flags for registering.
	f.StringVarP(&flags.Username, "username", "u", flags.Username, "The account login name.")
	f.StringVarP(&flags.Password, "password", "p", flags.Password, "The account password.")
	f.StringVarP(&flags.Email, "email", "e", flags.Email, "The account email.")
	f.StringVarP(&flags.Website, "website", "w", flags.Website, "The talebook website.")

	// Mark some flags as required.
	_ = talebookRegisterCmd.MarkFlagRequired("website")
	_ = talebookRegisterCmd.MarkFlagRequired("username")
	_ = talebookRegisterCmd.MarkFlagRequired("password")
	_ = talebookRegisterCmd.MarkFlagRequired("email")

	talebookCmd.AddCommand(talebookRegisterCmd)
}
