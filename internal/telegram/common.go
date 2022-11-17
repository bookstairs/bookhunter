package telegram

import (
	"context"

	"github.com/gotd/td/telegram"
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
func New(channelID string, mobile string, appID int64, appHash string, client *telegram.Client) *Telegram {
	return &Telegram{
		channelID: channelID,
		mobile:    mobile,
		appID:     appID,
		appHash:   appHash,
		client:    client,
	}
}

func (t *Telegram) Execute(f func() error) error {
	return t.client.Run(context.Background(), func(ctx context.Context) error {
		if err := t.Authentication(ctx); err != nil {
			return err
		}

		t.ctx = ctx

		return f()
	})
}
