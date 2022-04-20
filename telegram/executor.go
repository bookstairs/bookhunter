package telegram

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"
	"golang.org/x/time/rate"

	"github.com/bibliolater/bookhunter/pkg/log"
)

type executor struct {
	sessionPath string
	config      *Config
	context     context.Context
	cancel      context.CancelFunc
	taskCh      chan func(context.Context, *telegram.Client) error
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	return &executor{
		sessionPath: config.CookieFile,
		config:      config,
		context:     ctx,
		cancel:      cancel,
		taskCh:      make(chan func(context.Context, *telegram.Client) error, 10),
	}
}

func (e *executor) Execute() error {
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

	return client.Run(e.context, func(ctx context.Context) error {
		// Login the telegram account.
		err := login(ctx, client, e.config)
		if err != nil {
			return err
		}

		for task := range e.taskCh {
			err := task(ctx, client)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
}
