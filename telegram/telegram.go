package telegram

import (
	"bytes"
	"context"
	"github.com/bibliolater/bookhunter/pkg/progress"
	"github.com/bibliolater/bookhunter/pkg/rename"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/spider"
)

var (
	// ChannelId for telegram.
	ChannelId       = "https://t.me/haoshufenxiang"
	SessionPath     = ".tg-session"
	ReLogin         bool
	AppID           int
	AppHash         string
	ChunkSize       = 512 * 1024
	LoadMessageSize = 20
	LastId          int
)

type downloader struct {
	config       *spider.Config
	client       *telegram.Client
	context      context.Context
	channelId    string
	retry        int
	downloadPath string
	formats      []string
	rename       bool
	wait         *sync.WaitGroup
}

type filePart struct {
	Index  int
	Limit  int
	Offset int
}

type tgFile struct {
	id           int
	filename     string
	format       string
	fileSize     int
	documentFile *tg.InputDocumentFileLocation
	savePath     string
}

func NewDownloader(config *spider.Config) *downloader {

	if ReLogin {
		err := os.Remove(path.Join(config.DownloadPath, SessionPath))
		if err != nil {
			panic(nil)
		}
	}

	client := telegram.NewClient(AppID, AppHash, telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: path.Join(config.DownloadPath, SessionPath),
		},
	})
	ChannelId = strings.TrimPrefix(ChannelId, "https://t.me/")

	return &downloader{
		channelId:    ChannelId,
		config:       config,
		client:       client,
		context:      context.Background(),
		retry:        config.Retry,
		downloadPath: config.DownloadPath,
		formats:      config.Formats,
		rename:       config.Rename,
		wait:         new(sync.WaitGroup),
	}
}

// latestBookID will return the last available book ID.
func (d *downloader) latestBookID(info *tg.Channel) (int, error) {
	a := make([]tg.InputChannelClass, 1)
	a[0] = &tg.InputChannel{
		ChannelID:  info.ID,
		AccessHash: info.AccessHash,
	}
	searchData, err := d.client.API().MessagesSearch(d.context, &tg.MessagesSearchRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  info.ID,
			AccessHash: info.AccessHash,
		},
		Filter:   &tg.InputMessagesFilterEmpty{},
		Q:        "",
		OffsetID: -1,
		Limit:    1,
	})
	if err != nil {
		return 0, err
	}
	channelInfo, ok := searchData.(*tg.MessagesChannelMessages)
	if !ok {
		return 0, err
	}
	lastID := -1
	for _, tmp := range channelInfo.Messages {
		if tmp == nil {
			continue
		}
		msg, ok := tmp.(*tg.Message)
		if !ok {
			continue
		}
		lastID = msg.ID
	}
	return lastID, nil
}

func (d *downloader) Exec() error {
	f := func(ctx context.Context) error {
		err := d.login()
		if err != nil {
			return err
		}
		ch := make(chan tgFile, d.config.Thread)

		d.Fork()
		go d.startDownloads(ch)

		d.Fork()
		go func() {
			defer d.Done()
			for entity := range ch {
				d.Fork()
				go d.DownloadFile(&entity)
			}
		}()

		d.Join()
		return nil
	}
	if err := d.client.Run(d.context, f); err != nil {
		return err
	}
	return nil
}

func (d *downloader) login() error {
	flow := auth.NewFlow(
		&TermAuth{},
		auth.SendCodeOptions{},
	)
	if err := d.client.Auth().IfNecessary(d.context, flow); err != nil {
		return err
	}

	log.Info("Login Success")
	return nil
}

func (d *downloader) startDownloads(ch chan tgFile) {
	defer d.Done()
	defer close(ch)

	client := d.client
	api := client.API()
	ctx := d.context

	resolveUsername, err := api.ContactsResolveUsername(ctx, ChannelId)
	if err != nil {
		panic(err)
	}
	channelInfo := resolveUsername.Chats[0].(*tg.Channel)

	title := channelInfo.Title
	log.Infof("Start Download channel: %s", title)
	last, err := d.latestBookID(channelInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Find the last book ID: %d", last)

	LastId = last
	// Create the channel dir.
	saveDir := path.Join(d.config.DownloadPath, rename.EscapeFilename(title))
	if !IsDir(saveDir) {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	// Shards that generate query messages
	idParts := generatePart(d.config.InitialBookID, last-d.config.InitialBookID, 20)

	for _, part := range idParts {
		history, err := api.MessagesSearch(ctx, &tg.MessagesSearchRequest{
			Peer: &tg.InputPeerChannel{
				ChannelID:  channelInfo.ID,
				AccessHash: channelInfo.AccessHash,
			},
			Filter:    &tg.InputMessagesFilterEmpty{},
			Q:         "",
			OffsetID:  part.Offset,
			Limit:     part.Limit,
			AddOffset: -part.Limit,
		})
		if err != nil {
			panic(err)
		}

		messages := history.(*tg.MessagesChannelMessages)
		for i := range messages.Messages {
			message := messages.Messages[len(messages.Messages)-i-1]
			entity, ok := toFile(message, saveDir)

			if !ok {
				log.Warnf("[%d/%d] No downloadable files found, this resource could be banned.", message.GetID(), last)
				continue
			}
			if !d.formatMatcher(entity.format) {
				log.Warnf("[%d/%d] No match file format, this resource could be banned.", message.GetID(), last)
				// Skip this format.
				continue
			}

			//err := d.downloadFile(entity)
			//if err != nil {
			//	return err
			//}
			d.saveCurrentBookId(entity.id, last)
			ch <- *entity
		}
	}

	// Return to close client connection and free up resources.
}

func (d *downloader) saveCurrentBookId(current int, last int) {
	// Create book storage.
	storageFile := path.Join(d.config.DownloadPath, d.config.ProgressFile)
	_, err := progress.NewProgress(int64(current), int64(last), storageFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *downloader) DownloadFile(entity *tgFile) {
	defer d.Done()

	writer, err := os.Create(entity.savePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = writer.Close() }()

	fileParts := generatePart(0, entity.fileSize, ChunkSize)
	// Add startDownloads progress
	bar := log.NewProgressBar(int64(entity.id), int64(LastId), entity.format+" "+entity.filename, int64(entity.fileSize))
	for _, filePart := range fileParts {
		getFile, err := d.client.API().UploadGetFile(d.context, &tg.UploadGetFileRequest{
			Location:     entity.documentFile,
			Limit:        filePart.Limit,
			Offset:       filePart.Offset,
			CDNSupported: false,
			Precise:      false,
		})
		if err != nil {
			log.Fatal(err)
		}
		resp := getFile.(*tg.UploadFile)
		// Write file content
		_, err = io.Copy(io.MultiWriter(writer, bar), bytes.NewReader(resp.Bytes))
		if err != nil {
			log.Fatal(err)
		}
	}

	//for entity := range ch {
	//
	//}
}

func toFile(message tg.MessageClass, dir string) (*tgFile, bool) {
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
	if len(fileName) == 0 {
		return nil, false
	}
	format := spider.Extension(fileName)
	saveFilePath := path.Join(dir, rename.EscapeFilename(strconv.Itoa(msg.ID)+"_"+fileName))
	return &tgFile{
		id:           msg.ID,
		filename:     fileName,
		format:       format,
		fileSize:     document.Size,
		documentFile: document.AsInputDocumentFileLocation(),
		savePath:     saveFilePath,
	}, true
}

func (d *downloader) formatMatcher(fileName string) bool {
	for _, f := range d.config.Formats {
		if strings.HasSuffix(fileName, strings.ToLower(f)) {
			return true
		}
	}
	return false
}

func (d *downloader) Fork() {
	d.wait.Add(1)
}

func (d *downloader) Done() {
	d.wait.Done()
}

func (d *downloader) Join() {
	d.wait.Wait()
}

func generatePart(start int, length int, step int) []filePart {
	eachSize := length / step
	if eachSize == 0 {
		eachSize = 1
	}
	if length%step > 0 {
		eachSize = eachSize + 1
	}
	jobs := make([]filePart, eachSize)

	for i := range jobs {
		jobs[i].Index = i
		if i == 0 {
			jobs[i].Offset = start
		} else {
			jobs[i].Offset = jobs[i-1].Offset + jobs[i-1].Limit
		}
		jobs[i].Limit = step
		//log.Infof("part %d --> %d", jobs[i].Offset, jobs[i].Offset+jobs[i].Limit)
	}
	return jobs
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
