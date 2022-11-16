package lanzou

import (
	"fmt"
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

type Drive struct {
	client *client.Client
}

var (
	// 蓝奏云的主域名有时会挂掉, 此时尝试切换到备用域名
	availableDomains = []string{
		"lanzouw.com",
		"lanzoui.com",
		"lanzoux.com",
		"lanzouy.com",
		"lanzoup.com",
	}
)

// hostname is used to return a lanzou host.
func hostname() string {
	return availableDomains[0]
}

func checkOrSwitchDomain(c *client.Client) (err error) {
	head, err := c.R().Head("/")
	if head.IsError() || err != nil {
		if len(availableDomains) == 1 {
			return fmt.Errorf("no lanzou domains available")
		}
		availableDomains = availableDomains[1:]
		c.SetHost(hostname())
		err = checkOrSwitchDomain(c)
	}
	return err
}

func NewDrive(config *client.Config) (*Drive, error) {
	cl, err := client.New(&client.Config{
		HTTPS:      true,
		Host:       hostname(),
		UserAgent:  config.UserAgent,
		Proxy:      config.Proxy,
		ConfigRoot: config.ConfigRoot,
	})

	if err != nil {
		return nil, err
	}

	cl.Client.
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Referer", cl.BaseURL)

	err = checkOrSwitchDomain(cl)
	if err != nil {
		return nil, err
	}
	return &Drive{client: cl}, nil
}

func (l *Drive) DownloadFile(downloadURL string) (io.ReadCloser, error) {
	log.Debugf("Start to download file from aliyun drive: %s", downloadURL)

	resp, err := l.client.R().
		SetDoNotParseResponse(true).
		Get(downloadURL)
	if err != nil {
		return nil, err
	}

	return resp.RawBody(), err
}
