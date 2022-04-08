package sanqiu

import "github.com/bibliolater/bookhunter/pkg/spider"

// The Website for sanqiu.
var Website = "https://www.sanqiu.cc"

type downloader struct {
}

func NewDownloader(config *spider.DownloadConfig) *downloader {
	return nil
}

func (d *downloader) Download() {

}
