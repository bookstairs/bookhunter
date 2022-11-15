package fetcher

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/td/session"
	client "github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/telegram"
)

type telegramFetcher struct {
	fetcher  *commonFetcher
	telegram *telegram.Telegram
}

func newTelegramFetcher(config *Config) (Fetcher, error) {
	// Create the session file.
	path, err := config.ConfigPath()
	if err != nil {
		return nil, err
	}
	sessionPath := filepath.Join(path, "session.json")
	if refresh, _ := strconv.ParseBool(config.Property("reLogin")); refresh {
		_ = os.Remove(sessionPath)
	}

	channelID := config.Property("channelID")
	mobile := config.Property("mobile")
	appID, _ := strconv.ParseInt(config.Property("appID"), 10, 64)
	appHash := config.Property("appHash")

	// Change the process file name.
	config.precessFile = strings.ReplaceAll(channelID, "/", "_") + ".db"

	// Create the http proxy dial.
	dialFunc, err := telegram.CreateProxy(config.Proxy)
	if err != nil {
		return nil, err
	}

	// Create the backend telegram client.
	cl := client.NewClient(
		int(appID),
		appHash,
		client.Options{
			Resolver:       dcs.Plain(dcs.PlainOptions{Dial: dialFunc}),
			SessionStorage: &session.FileStorage{Path: sessionPath},
			Middlewares: []client.Middleware{
				floodwait.NewSimpleWaiter().WithMaxRetries(uint(3)),
			},
		},
	)

	tel := telegram.New(channelID, mobile, appID, appHash, cl)

	return &telegramFetcher{
		fetcher: &commonFetcher{
			Config: config,
			service: &telegramService{
				config:   config,
				telegram: tel,
			},
		},
		telegram: tel,
	}, nil
}

func (t *telegramFetcher) Download() error {
	return t.telegram.Execute(func() error {
		return t.fetcher.Download()
	})
}

type telegramService struct {
	config   *Config
	telegram *telegram.Telegram
	info     *telegram.ChannelInfo
}

func (s *telegramService) size() (int64, error) {
	info, err := s.telegram.ChannelInfo()
	if err != nil {
		return 0, err
	}
	s.info = info

	return info.LastMsgID, nil
}

func (s *telegramService) formats(id int64) (map[Format]driver.Share, error) {
	files, err := s.telegram.ParseMessage(s.info, id)
	if err != nil {
		return nil, err
	}

	res := make(map[Format]driver.Share)
	for _, file := range files {
		res[Format(file.Format)] = driver.Share{
			FileName: file.Name,
			Size:     file.Size,
			Properties: map[string]any{
				"fileID":   file.ID,
				"document": file.Document,
			},
		}
	}

	return res, nil
}

func (s *telegramService) fetch(_ int64, f Format, share driver.Share, writer io.Writer) error {
	file := &telegram.File{
		ID:       share.Properties["fileID"].(int64),
		Name:     share.FileName,
		Format:   string(f),
		Size:     share.Size,
		Document: share.Properties["document"].(*tg.InputDocumentFileLocation),
	}

	return s.telegram.DownloadFile(file, writer)
}
