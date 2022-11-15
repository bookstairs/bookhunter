package telegram

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
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

	// File is the file info from the telegram channel.
	File struct {
		ID       int64
		Name     string
		Format   string
		Size     int64
		Document *tg.InputDocumentFileLocation
	}
)

// New will create a telegram client.
func New(channelID string, mobile string, appID int64, appHash string, client *telegram.Client) (*Telegram, error) {
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
