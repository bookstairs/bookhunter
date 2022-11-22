package telegram

import (
	"context"

	"github.com/gotd/contrib/bg"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/file"
)

type (
	Telegram struct {
		channelID string
		mobile    string
		appID     int64
		appHash   string
		client    *telegram.Client
		ctx       context.Context
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

	if err := t.Authentication(t.ctx); err != nil {
		return nil, err
	}

	return t, nil
}
