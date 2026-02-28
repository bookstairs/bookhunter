package telegram

import (
	"context"
	"time"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/file"
)

// sessionRefreshInterval defines how often the session should be refreshed.
const sessionRefreshInterval = 3 * time.Hour

type (
	Telegram struct {
		channelID          string
		mobile             string
		appID              int64
		appHash            string
		client             *telegram.Client
		ctx                context.Context
		lastSessionRefresh time.Time
	}

	ChannelInfo struct {
		ID         int64
		AccessHash int64
		LastMsgID  int64
	}

	// File is the file info from the telegram channel.
	File struct {
		ID       int64
		Name     string
		Format   file.Format
		Size     int64
		Document *tg.InputDocumentFileLocation
	}
)

// New will create a telegram client.
func New(channelID, mobile string, appID int64, appHash string, sessionPath, proxy string) (*Telegram, error) {
	// Create the http proxy dial.
	dialFunc, err := createProxy(proxy)
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
			},
		},
	)

	ctx := context.Background()
	_, err = bg.Connect(client, bg.WithContext(ctx)) // No need to close this client.
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		ctx:       ctx,
		channelID: channelID,
		mobile:    mobile,
		appID:     appID,
		appHash:   appHash,
		client:    client,
	}

	if err := t.Authentication(); err != nil {
		return nil, err
	}

	return t, nil
}

// RefreshSession will refresh the telegram session using Auth().IfNecessary().
func (t *Telegram) RefreshSession() error {
	flow := auth.NewFlow(&terminalAuth{mobile: t.mobile}, auth.SendCodeOptions{})
	if err := t.client.Auth().IfNecessary(t.ctx, flow); err != nil {
		return err
	}

	t.lastSessionRefresh = time.Now()
	return nil
}

// EnsureSessionValid ensures the session is valid, refreshing if necessary (every 3 hours).
func (t *Telegram) EnsureSessionValid() error {
	if time.Since(t.lastSessionRefresh) > sessionRefreshInterval {
		return t.RefreshSession()
	}
	return nil
}

// checkSessionRefresh is an alias for EnsureSessionValid for compatibility
func (t *Telegram) checkSessionRefresh() error {
	return t.EnsureSessionValid()
}
