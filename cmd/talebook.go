package cmd

import (
	"fmt"

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
	Short: "A tool for downloading books from talebook server",
}

// talebookDownloadCmd represents the download command
var talebookDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the books from talebook",
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
			Row("Formats", flags.Formats).
			Row("Download Path", flags.DownloadPath).
			Row("Initial ID", flags.InitialBookID).
			Row("Rename File", flags.Rename).
			Row("Thread", flags.Thread).
			Row("Keywords", flags.Keywords).
			Row("Thread Limit (req/min)", flags.RateLimit).
			Print()

		// Create the fetcher.
		f, err := flags.NewFetcher(fetcher.Talebook, map[string]string{
			"username": flags.Username,
			"password": flags.Password,
		})
		log.Exit(err)

		// Start downloading the books.
		err = f.Download()
		log.Exit(err)

		// Finished all the tasks.
		log.Info("Successfully download all the books.")
	},
}

// talebookRegisterCmd represents the register command.
var talebookRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account on talebook",
	Long: `Some talebook website need a user account for downloading books
You can use this register command for creating account`,
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
			Row("Proxy", flags.Proxy).
			Print()

		// Create client config.
		config, err := client.NewConfig(flags.Website, flags.Proxy, flags.ConfigRoot)
		log.Exit(err)

		// Create http client.
		c, err := client.New(config)
		log.Exit(err)

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
		log.Exit(err)

		result := resp.Result().(*talebook.CommonResp)
		if result.Err == talebook.SuccessStatus {
			log.Info("Register success.")
		} else {
			log.Exit(fmt.Errorf("register failed, reason: %s", result.Err))
		}
	},
}

func init() {
	/// Add the download command.

	// Add flags for use info.
	f := talebookDownloadCmd.Flags()

	// Talebook related flags.
	f.StringVarP(&flags.Username, "username", "u", flags.Username, "The talebook username")
	f.StringVarP(&flags.Password, "password", "p", flags.Password, "The talebook password")
	f.StringVarP(&flags.Website, "website", "w", flags.Website, "The talebook link")

	// Common download flags.
	f.StringSliceVarP(&flags.Formats, "format", "f", flags.Formats, "The file formats you want to download")
	f.StringVarP(&flags.DownloadPath, "download", "d", flags.DownloadPath, "The book directory you want to use")
	f.Int64VarP(&flags.InitialBookID, "initial", "i", flags.InitialBookID, "The book id you want to start download")
	f.BoolVarP(&flags.Rename, "rename", "r", flags.Rename, "Rename the book file by book id")
	f.IntVarP(&flags.Thread, "thread", "t", flags.Thread, "The number of download thead")
	f.IntVar(&flags.RateLimit, "ratelimit", flags.RateLimit, "The allowed requests per minutes for every thread")

	// Mark some flags as required.
	_ = talebookDownloadCmd.MarkFlagRequired("website")

	talebookCmd.AddCommand(talebookDownloadCmd)

	/// Add the register command.

	f = talebookRegisterCmd.Flags()

	// Add flags for registering.
	f.StringVarP(&flags.Username, "username", "u", flags.Username, "The talebook username")
	f.StringVarP(&flags.Password, "password", "p", flags.Password, "The talebook password")
	f.StringVarP(&flags.Email, "email", "e", flags.Email, "The talebook email")
	f.StringVarP(&flags.Website, "website", "w", flags.Website, "The talebook link")

	// Mark some flags as required.
	_ = talebookRegisterCmd.MarkFlagRequired("website")
	_ = talebookRegisterCmd.MarkFlagRequired("username")
	_ = talebookRegisterCmd.MarkFlagRequired("password")
	_ = talebookRegisterCmd.MarkFlagRequired("email")

	talebookCmd.AddCommand(talebookRegisterCmd)
}
