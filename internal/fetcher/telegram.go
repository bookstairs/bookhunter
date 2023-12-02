package fetcher

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotd/td/tg"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/telegram"
)

func newTelegramService(config *Config) (service, error) {
	// Create the session file.
	path, err := config.ConfigPath()
	if err != nil {
		return nil, err
	}
	sessionPath := filepath.Join(path, "session.db")
	if refresh, _ := strconv.ParseBool(config.Property("reLogin")); refresh {
		_ = os.Remove(sessionPath)
	}

	channelID := config.Property("channelID")
	mobile := config.Property("mobile")
	appID, _ := strconv.ParseInt(config.Property("appID"), 10, 64)
	appHash := config.Property("appHash")

	// Change the process file name.
	config.processFile = strings.ReplaceAll(channelID, "/", "_") + ".db"

	tel, err := telegram.New(channelID, mobile, appID, appHash, sessionPath, config.Proxy)
	if err != nil {
		return nil, err
	}

	return &telegramService{config: config, telegram: tel}, nil
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

func (s *telegramService) formats(id int64) (map[file.Format]driver.Share, error) {
	files, err := s.telegram.ParseMessage(s.info, id)
	if err != nil {
		return nil, err
	}

	res := make(map[file.Format]driver.Share)
	for _, f := range files {
		res[f.Format] = driver.Share{
			FileName: f.Name,
			Size:     f.Size,
			Properties: map[string]any{
				"fileID":   f.ID,
				"document": f.Document,
			},
		}
	}

	return res, nil
}

func (s *telegramService) fetch(_ int64, f file.Format, share driver.Share, writer file.Writer) error {
	o := &telegram.File{
		ID:       share.Properties["fileID"].(int64),
		Name:     share.FileName,
		Format:   f,
		Size:     share.Size,
		Document: share.Properties["document"].(*tg.InputDocumentFileLocation),
	}

	return s.telegram.DownloadFile(o, writer)
}
