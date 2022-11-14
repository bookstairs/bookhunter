package driver

import (
	"io"

	"github.com/tickstep/cloudpan189-api/cloudpan"

	"github.com/bookstairs/bookhunter/internal/client"
)

func newTelecomDriver(_ *client.Config, properties map[string]string) (Driver, error) {
	webToken := &cloudpan.WebLoginToken{}
	appToken := &cloudpan.AppLoginToken{}

	// Try to sign in to the telecom disk.
	username := properties["username"]
	password := properties["password"]
	if username != "" && password != "" {
		token, err := cloudpan.AppLogin(username, password)
		if err != nil {
			return nil, err
		}

		webTokenStr := cloudpan.RefreshCookieToken(token.SessionKey)
		if webTokenStr != "" {
			webToken.CookieLoginUser = webTokenStr
		}
		appToken = token
	}

	// Create the pan client.
	c := cloudpan.NewPanClient(*webToken, *appToken)

	return &telecomDriver{client: c}, nil
}

type telecomDriver struct {
	client *cloudpan.PanClient
}

func (t *telecomDriver) Resolve(shareLink string, passcode string) []Share {
	panic("TODO implement me")
}

func (t *telecomDriver) Download(share Share) io.ReadCloser {
	panic("TODO implement me")
}
