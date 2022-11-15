package telegram

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/time/rate"

	"github.com/bookstairs/bookhunter/internal/fetcher"
)

type (
	Telegram struct {
		channelID string
		mobile    string
		appID     int64
		appHash   string
		client    *telegram.Client
	}

	ChannelInfo struct {
		ID         int64
		AccessHash int64
		LastMsgID  int64
	}
)

// New will create a telegram client.
func New(config *fetcher.Config) (*Telegram, error) {
	// Create the session file.
	path, err := config.ConfigPath()
	if err != nil {
		return nil, err
	}
	sessionPath := filepath.Join(path, "session.json")
	if refresh, _ := strconv.ParseBool(config.Property("reLogin")); refresh {
		_ = os.Remove(sessionPath)
	}

	channelID := config.Property("channelID")
	mobile := config.Property("mobile")
	appID, _ := strconv.ParseInt(config.Property("appID"), 10, 64)
	appHash := config.Property("appHash")

	// Create the http proxy dial.
	dialFunc, err := createProxy(config.Proxy)
	if err != nil {
		return nil, err
	}

	// Create the backend telegram client.
	client := telegram.NewClient(
		int(appID),
		appHash,
		telegram.Options{
			Resolver:       dcs.Plain(dcs.PlainOptions{Dial: dialFunc}),
			SessionStorage: &session.FileStorage{Path: sessionPath},
			Middlewares: []telegram.Middleware{
				floodwait.NewSimpleWaiter().WithMaxRetries(uint(3)),
				ratelimit.New(rate.Every(time.Minute), config.RateLimit),
			},
		},
	)

	t := &Telegram{
		channelID: channelID,
		mobile:    mobile,
		appID:     appID,
		appHash:   appHash,
		client:    client,
	}

	// Try to sign in with the mobile.
	if err := t.login(); err != nil {
		return nil, err
	}

	return t, nil
}

// Every telegram execution should be wrapped in a client Run session.
// We have to expose this method for internal usage.
func (t *Telegram) execute(f func(context.Context, *telegram.Client) error) error {
	return t.client.Run(context.Background(), func(ctx context.Context) error {
		if err := t.authentication(ctx); err != nil {
			return err
		}

		return f(ctx, t.client)
	})
}
