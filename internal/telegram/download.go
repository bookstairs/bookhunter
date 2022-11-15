package telegram

import (
	"context"
	"io"
	"math"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/naming"
)

func (t *Telegram) DownloadFile(file *File, writer io.Writer) error {
	return t.execute(func(ctx context.Context, client *telegram.Client) error {
		tool := downloader.NewDownloader()
		thread := int(math.Ceil(float64(file.Size) / (512 * 1024)))
		_, err := tool.Download(client.API(), file.Document).WithThreads(thread).Stream(ctx, writer)

		return err
	})
}

// ParseMessage will parse the given message id.
func (t *Telegram) ParseMessage(info *ChannelInfo, msgID int64) ([]File, error) {
	var files []File
	err := t.execute(func(ctx context.Context, client *telegram.Client) error {
		// This API is translated from official C++ client.
		api := client.API()
		history, err := api.MessagesSearch(ctx, &tg.MessagesSearchRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  info.ID,
				AccessHash: info.AccessHash,
			},
			Filter:    &tg.InputMessagesFilterEmpty{},
			Q:         "",
			OffsetID:  int(msgID),
			Limit:     1,
			AddOffset: -1,
		})
		if err != nil {
			return err
		}

		messages := history.(*tg.MessagesChannelMessages)
		for i := len(messages.Messages) - 1; i >= 0; i-- {
			message := messages.Messages[i]

			if file, ok := parseFile(message); ok {
				files = append(files, *file)
			} else {
				log.Warnf("[%d/%d] No downloadable files found.", message.GetID(), info.LastMsgID)
				continue
			}
		}

		return nil
	})

	return files, err
}

func parseFile(message tg.MessageClass) (*File, bool) {
	if message == nil {
		return nil, false
	}
	msg, ok := message.(*tg.Message)
	if !ok {
		return nil, false
	}
	if msg.Media == nil {
		return nil, false
	}
	s, ok := msg.Media.(*tg.MessageMediaDocument)
	if !ok {
		return nil, false
	}
	document := s.Document.(*tg.Document)
	fileName := ""
	for _, attribute := range document.Attributes {
		x, ok := attribute.(*tg.DocumentAttributeFilename)
		if ok {
			fileName = x.FileName
		}
	}
	if fileName == "" {
		return nil, false
	}
	format, _ := naming.Extension(fileName)

	return &File{
		ID:       int64(msg.ID),
		Name:     fileName,
		Format:   format,
		Size:     document.Size,
		Document: document.AsInputDocumentFileLocation(),
	}, true
}
