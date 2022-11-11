package telegram

import (
	"context"
	"os"
	"time"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"
	"golang.org/x/time/rate"

	"github.com/bookstairs/bookhunter/internal/log"
)

type executor struct {
	sessionPath string
	config      *Config
}

// NewExecutor will create a wrapper for executing client.
func NewExecutor(config *Config) *executor {
	// Remove session file for forcing login.
	if config.Refresh {
		err := os.Remove(config.CookieFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &executor{
		sessionPath: config.CookieFile,
		config:      config,
	}
}

func (e *executor) Execute(f func(context.Context, *telegram.Client) error) error {
	// The backend client.
	client := telegram.NewClient(
		e.config.AppID,
		e.config.AppHash,
		telegram.Options{
			Resolver:       dcs.Plain(dcs.PlainOptions{Dial: proxy.Dial}),
			SessionStorage: &session.FileStorage{Path: e.sessionPath},
			Middlewares: []telegram.Middleware{
				floodwait.NewSimpleWaiter().WithMaxRetries(uint(e.config.Retry)),
				ratelimit.New(rate.Every(100*time.Millisecond), 5),
			},
		},
	)

	ctx := context.Background()
	return client.Run(ctx, func(ctx context.Context) error {
		// Login the telegram account.
		err := login(ctx, client, e.config)
		if err != nil {
			return err
		}

		return f(ctx, client)
	})
}
