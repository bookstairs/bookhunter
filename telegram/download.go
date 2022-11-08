package telegram

import (
	"context"
	"io"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/pkg/log"
	"github.com/bookstairs/bookhunter/pkg/progress"
	"github.com/bookstairs/bookhunter/pkg/rename"
	"github.com/bookstairs/bookhunter/pkg/spider"
)

// tgFile is the file info from telegram channel.
type tgFile struct {
	id       int
	name     string
	format   string
	size     int
	document *tg.InputDocumentFileLocation
	dest     string
}

// download would start download books from given telegram channel.
func (d *tgDownloader) download(ctx context.Context, client *telegram.Client) {
	defer d.wait.Done()

	channel := d.channel
	api := client.API()

	msgID := d.progress.AcquireBookID()
	log.Infof("Start to download book from %d.", msgID)

	for ; msgID != progress.NoBookToDownload; msgID = d.progress.AcquireBookID() {
		history, err := api.MessagesSearch(ctx, &tg.MessagesSearchRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  channel.id,
				AccessHash: channel.accessHash,
			},
			Filter:    &tg.InputMessagesFilterEmpty{},
			Q:         "",
			OffsetID:  int(msgID),
			Limit:     1,
			AddOffset: -1,
		})

		if err != nil {
			log.Fatal(err)
		}

		messages := history.(*tg.MessagesChannelMessages)
		for i := range messages.Messages {
			message := messages.Messages[len(messages.Messages)-i-1]
			file, ok := d.parseFile(message)

			if !ok {
				log.Warnf("[%d/%d] No downloadable files found.", message.GetID(), channel.lastMsgID)
				continue
			}
			if !d.formatMatcher(file.format) {
				log.Warnf("[%d/%d] No matched file format.", message.GetID(), channel.lastMsgID)
				continue
			}

			d.downloadFile(ctx, client, file)
		}

		if err := d.progress.SaveBookID(msgID); err != nil {
			log.Fatal(err)
		}
	}
}

// parseFile would acquire the file info from telegram message.
func (d *tgDownloader) parseFile(message tg.MessageClass) (*tgFile, bool) {
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
	format, _ := spider.Extension(fileName)

	file := strconv.Itoa(msg.GetID()) + "_" + fileName
	if d.config.Rename {
		file = strconv.FormatInt(int64(msg.ID), 10) + "." + strings.ToLower(format)
	}
	dest := path.Join(d.config.DownloadPath, rename.EscapeFilename(file))

	return &tgFile{
		id:       msg.ID,
		name:     fileName,
		format:   format,
		size:     int(document.Size),
		document: document.AsInputDocumentFileLocation(),
		dest:     dest,
	}, true
}

func (d *tgDownloader) downloadFile(ctx context.Context, client *telegram.Client, file *tgFile) {
	tool := downloader.NewDownloader()

	// Remove the exist file.
	if _, err := os.Stat(file.dest); err == nil {
		if err := os.Remove(file.dest); err != nil {
			log.Fatal(err)
		}
	}

	writer, err := os.Create(file.dest)
	if err != nil {
		log.Fatalf("Create file err [%s]  %s", file.dest, err)
	}
	defer func() { _ = writer.Close() }()

	thread := int(math.Ceil(float64(file.size) / (512 * 1024)))

	// Create download progress
	bar := log.NewProgressBar(int64(file.id), d.progress.Size(), file.format+" "+file.name, int64(file.size))

	_, err2 := tool.Download(client.API(), file.document).WithThreads(thread).Stream(ctx, io.MultiWriter(writer, bar))
	if err2 != nil {
		log.Fatal(err)
	}
}

func (d *tgDownloader) formatMatcher(fileName string) bool {
	for _, f := range d.config.Formats {
		if strings.HasSuffix(fileName, strings.ToLower(f)) {
			return true
		}
	}
	return false
}
