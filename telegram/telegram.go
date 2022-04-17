package telegram

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	"golang.org/x/time/rate"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/progress"
	"github.com/bibliolater/bookhunter/pkg/rename"
	"github.com/bibliolater/bookhunter/pkg/spider"
)

var (
	// ChannelId for telegram.
	ChannelId       = "https://t.me/haoshufenxiang"
	SessionPath     = ".tg-session"
	ReLogin         = false
	AppID           = 0
	AppHash         = ""
	LoadMessageSize = 20
	LastId          = 0
)

type tgDownloader struct {
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
	filePath     string
}

func NewDownloader(config *spider.Config) *tgDownloader {
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
		Middlewares: []telegram.Middleware{
			floodwait.NewSimpleWaiter().WithMaxRetries(uint(config.Retry)),
			ratelimit.New(rate.Every(100*time.Millisecond), 5),
		},
	})
	ChannelId = strings.TrimPrefix(ChannelId, "https://t.me/")
	return &tgDownloader{
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
func (d *tgDownloader) latestBookID(info *tg.Channel) (int, error) {
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
		lastID = tmp.GetID()
	}
	return lastID, nil
}

func (d *tgDownloader) Exec() error {
	f := func(ctx context.Context) error {
		err := d.login()
		if err != nil {
			return err
		}
		ch := make(chan tgFile, LoadMessageSize)

		d.Fork()
		go d.startDownloads(ch)

		for i := 0; i < d.config.Thread; i++ {
			d.Fork()
			go func() {
				defer d.Done()
				for item := range ch {
					d.DownloadFile(&item)
				}
			}()
		}
		d.Join()
		return nil
	}
	if err := d.client.Run(d.context, f); err != nil {
		return err
	}
	return nil
}

func (d *tgDownloader) login() error {
	flow := auth.NewFlow(
		&TermAuth{},
		auth.SendCodeOptions{},
	)
	if err := d.client.Auth().IfNecessary(d.context, flow); err != nil {
		return err
	}
	status, _ := d.client.Auth().Status(d.context)
	if !status.Authorized {
		return errors.New("login fail, Please use '--reLogin' login again")
	}
	log.Info("Login Success")
	return nil
}

func (d *tgDownloader) startDownloads(ch chan tgFile) {
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
			entity, ok := d.toFile(message, saveDir)

			if !ok {
				log.Warnf("[%d/%d] No downloadable files found, this resource could be banned.", message.GetID(), last)
				continue
			}
			if !d.formatMatcher(entity.format) {
				log.Warnf("[%d/%d] No match file format, this resource could be banned.", message.GetID(), last)
				// Skip this format.
				continue
			}
			d.saveCurrentBookId(saveDir, entity.id, last)
			ch <- *entity
		}
	}
}

func (d *tgDownloader) saveCurrentBookId(saveDir string, current, last int) {
	// Create book storage.
	storageFile := path.Join(saveDir, d.config.ProgressFile)
	_, err := progress.NewProgress(int64(current), int64(last), storageFile)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *tgDownloader) DownloadFile(entity *tgFile) {
	tool := downloader.NewDownloader()

	// Remove the exist file.
	if _, err := os.Stat(entity.filePath); err == nil {
		if err := os.Remove(entity.filePath); err != nil {
			log.Fatal(err)
		}
	}
	writer, err := os.Create(entity.filePath)
	if err != nil {
		log.Fatalf("create file err [%s]  %s", entity.filePath, err)
	}
	defer func() { _ = writer.Close() }()

	thread := 1
	if entity.fileSize/(512*1024) > d.config.Thread*4 {
		thread = d.config.Thread
	}
	// Add startDownloads progress
	bar := log.NewProgressBar(int64(entity.id), int64(LastId), entity.format+" "+entity.filename, int64(entity.fileSize))

	_, err2 := tool.Download(d.client.API(), entity.documentFile).WithThreads(thread).Stream(d.context, io.MultiWriter(writer, bar))
	if err2 != nil {
		log.Fatal(err)
	}
}

func (d *tgDownloader) toFile(message tg.MessageClass, dir string) (*tgFile, bool) {
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
	format := spider.Extension(fileName)
	outFilename := strconv.FormatInt(int64(msg.ID), 10) + "." + strings.ToLower(format)
	if !d.rename {
		outFilename = strconv.Itoa(msg.GetID()) + "_" + fileName
	}
	saveFilePath := path.Join(dir, rename.EscapeFilename(outFilename))
	return &tgFile{
		id:           msg.ID,
		filename:     fileName,
		format:       format,
		fileSize:     document.Size,
		documentFile: document.AsInputDocumentFileLocation(),
		filePath:     saveFilePath,
	}, true
}

func (d *tgDownloader) formatMatcher(fileName string) bool {
	for _, f := range d.config.Formats {
		if strings.HasSuffix(fileName, strings.ToLower(f)) {
			return true
		}
	}
	return false
}

func (d *tgDownloader) Fork() {
	d.wait.Add(1)
}

func (d *tgDownloader) Done() {
	d.wait.Done()
}

func (d *tgDownloader) Join() {
	d.wait.Wait()
}

func generatePart(start, length, step int) []filePart {
	eachSize := length / step
	if eachSize == 0 {
		eachSize = 1
	}
	if length%step > 0 {
		eachSize++
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
	}
	return jobs
}

func IsDir(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return s.IsDir()
}
