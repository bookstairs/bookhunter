package fetcher

import (
	"errors"
	"net/http"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/talebook"
)

var ErrTalebookNeedSignin = errors.New("need user account to download books")

type talebookService struct {
	config *Config
	client *client.Client
}

func newTalebookService(config *Config) (service, error) {
	// Add login check in redirect handler.
	err := config.SetRedirect(func(request *http.Request, requests []*http.Request) error {
		if request.URL.Path == "/login" {
			return ErrTalebookNeedSignin
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Create the resty client for HTTP handing.
	c, err := client.New(config.Config)
	if err != nil {
		return nil, err
	}

	// Start to sign in if required.
	username := config.Property("username")
	password := config.Property("password")
	if username != "" && password != "" {
		log.Info("You have provided user information, start to login.")
		resp, err := c.R().
			SetFormData(map[string]string{
				"username": username,
				"password": password,
			}).
			SetResult(&talebook.LoginResp{}).
			ForceContentType("application/json").
			Post("/api/user/sign_in")

		if err != nil {
			return nil, err
		}

		result := resp.Result().(*talebook.LoginResp)
		if result.Err != talebook.SuccessStatus {
			return nil, errors.New(result.Msg)
		}

		log.Info("Login success. Save cookies into file.")
	}

	return &talebookService{
		config: config,
		client: c,
	}, nil
}

func (t *talebookService) size() (int64, error) {
	panic("implement me")
}

func (t *talebookService) formats() ([]Format, error) {
	panic("implement me")
}

func (t *talebookService) fetch(id int64, format Format) (*fetch, error) {
	f := createFetch(nil)
	return f, nil
}
