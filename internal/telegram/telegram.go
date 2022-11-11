package telegram

import (
	"context"
	"path"
	"sync"

	"github.com/gotd/td/telegram"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/progress"
	"github.com/bookstairs/bookhunter/internal/spider"
)

// Config extends the spider download config and add telegram related configuration.
type Config struct {
	*spider.Config        // Extends the default download config
	ChannelID      string // ChannelID for telegram.
	Mobile         string // Mobile is used for storing the login user.
	Refresh        bool   // Refresh will force sign in.
	AppID          int    // AppID is obtained from https://my.telegram.org/apps
	AppHash        string // AppHash is obtained from https://my.telegram.org/apps
}

// NewConfig will create a telegram download config.
func NewConfig() *Config {
	return &Config{
		Config:    spider.NewConfig(),
		ChannelID: "",
		Refresh:   false,
		AppID:     0,
		AppHash:   "",
	}
}

// tgDownloader is the downloader for downloading books from telegram.
type tgDownloader struct {
	config   *Config
	channel  *tgChannelInfo
	executor *executor
	progress *progress.Progress
	wait     *sync.WaitGroup
}

type tgChannelInfo struct {
	id         int64
	accessHash int64
	lastMsgID  int64
}

// NewDownloader create a telegram downloader.
func NewDownloader(config *Config) *tgDownloader {
	executor := NewExecutor(config)

	// Get last book ID
	var channel *tgChannelInfo
	err := executor.Execute(func(ctx context.Context, client *telegram.Client) error {
		var err error
		channel, err = channelInfo(ctx, client, config)
		return err
	})
	if err != nil {
		log.Fatalf("Couldn't find channel info. %v", err)
	}
	log.Infof("Find the last message ID: %d", channel.lastMsgID)

	// Create book storage.
	storageFile := path.Join(config.DownloadPath, config.ProgressFile)
	p, err := progress.NewProgress(int64(config.InitialBookID), channel.lastMsgID, storageFile)
	if err != nil {
		log.Fatal(err)
	}

	return &tgDownloader{
		config:   config,
		channel:  channel,
		executor: executor,
		progress: p,
		wait:     new(sync.WaitGroup),
	}
}

// Fork a running instance.
func (d *tgDownloader) Fork() {
	d.wait.Add(1)
	go func() {
		_ = d.executor.Execute(func(ctx context.Context, client *telegram.Client) error {
			d.download(ctx, client)
			return nil
		})
	}()
}

// Join will wait all the running instance be finished.
func (d *tgDownloader) Join() {
	d.wait.Wait()
}
