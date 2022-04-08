package talebook

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
	"github.com/bibliolater/bookhunter/talebook"
)

// Used for register account on talebook website.
type registerConfig struct {
	website   string
	username  string
	password  string
	email     string
	userAgent string
}

// Arguments instance.
var regConf = registerConfig{}

// RegisterCmd represents the register command.
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account on talebook.",
	Long: `Some talebook website need a user account for downloading books.
You can use this register command for creating account.`,
	Run: func(cmd *cobra.Command, args []string) {
		register()
	},
}

func init() {
	// Add flags for use info.
	RegisterCmd.Flags().StringVarP(&regConf.website, "website", "w", "", "The talebook website.")
	RegisterCmd.Flags().StringVarP(&regConf.username, "username", "u", "", "The account login name.")
	RegisterCmd.Flags().StringVarP(&regConf.password, "password", "p", "", "The account password.")
	RegisterCmd.Flags().StringVarP(&regConf.email, "email", "e", "", "The account email.")
	RegisterCmd.Flags().StringVarP(&regConf.userAgent, "user-agent", "a", spider.DefaultUserAgent, "The account email.")

	_ = RegisterCmd.MarkFlagRequired("website")
	_ = RegisterCmd.MarkFlagRequired("username")
	_ = RegisterCmd.MarkFlagRequired("password")
	_ = RegisterCmd.MarkFlagRequired("email")
}

// register will create account on given website
func register() {
	// Print download configuration.
	log.PrintTable("Register Config Info", table.Row{"Config Key", "Config Value"}, &regConf)

	// Create register request.
	website := spider.GenerateUrl(regConf.website, "/api/user/sign_up")
	referer := spider.GenerateUrl(regConf.website, "/signup")
	values := url.Values{
		"username": {regConf.username},
		"password": {regConf.password},
		"nickname": {regConf.username},
		"email":    {regConf.email},
	}

	req, err := http.NewRequest(http.MethodPost, website, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatal("Illegal login request: %w", err)
	}
	req.Header.Set("User-Agent", regConf.userAgent)
	req.Header.Set("referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = form.Body.Close() }()
	if form.StatusCode != http.StatusOK {
		log.Fatalf("Error in register user, message: %s", form.Status)
	}

	result := &talebook.CommonResponse{}
	if err = spider.DecodeResponse(form, result); err != nil {
		log.Fatal(err)
	}

	if result.Err == talebook.SuccessStatus {
		log.Info("Register success.")
	} else {
		log.Fatalf("Register failed, reason: %s", result.Err)
	}
}
