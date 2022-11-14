package telecom

import (
	"strconv"
	"time"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

const (
	webPrefix  = "https://cloud.189.cn"
	authPrefix = "https://open.e.189.cn/api/logbox/oauth2"
	apiPrefix  = "https://api.cloud.189.cn"
)

type Telecom struct {
	client   *client.Client
	appToken *AppLoginToken
}

func New(c *client.Config, username, password string) (*Telecom, error) {
	cl, err := client.New(&client.Config{
		HTTPS:      false,
		Host:       "cloud.189.cn",
		UserAgent:  c.UserAgent,
		Proxy:      c.Proxy,
		ConfigRoot: c.ConfigRoot,
	})
	if err != nil {
		return nil, err
	}

	t := &Telecom{client: cl}

	if username == "" || password == "" {
		log.Warn("No username or password provide, we may not able to download from telecom disk.")
		t.appToken = &AppLoginToken{}
		return t, nil
	}

	// Start to sign in.
	if err := t.login(username, password); err != nil {
		return nil, err
	}

	return t, nil
}

// timeStamp is used to return the telecom required time str.
func timeStamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano()/1e6, 10)
}
