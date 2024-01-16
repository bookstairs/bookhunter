package telecom

import (
	"errors"

	"github.com/bookstairs/bookhunter/internal/client"
)

const (
	webPrefix  = "https://cloud.189.cn"
	authPrefix = "https://open.e.189.cn/api/logbox/oauth2"
	apiPrefix  = "https://api.cloud.189.cn"
)

type Telecom struct {
	*client.Client
	appToken *AppLoginToken
}

func New(c *client.Config, username, password string) (*Telecom, error) {
	cl, err := client.New(&client.Config{
		HTTPS:      false,
		Host:       "cloud.189.cn",
		Proxy:      c.Proxy,
		ConfigRoot: c.ConfigRoot,
	})
	if err != nil {
		return nil, err
	}

	cl.SetHeader("Accept", "application/json;charset=UTF-8")
	t := &Telecom{Client: cl}

	if username == "" || password == "" {
		return nil, errors.New("no username or password provide, we may not able to download from telecom disk")
	}

	// Start to sign in.
	if err := t.login(username, password); err != nil {
		return nil, err
	}

	return t, nil
}
