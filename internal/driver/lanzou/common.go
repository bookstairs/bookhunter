package lanzou

import (
	"fmt"
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

type Lanzou struct {
	*client.Client
}

var (
	// 蓝奏云的主域名有时会挂掉, 此时尝试切换到备用域名
	availableHostnames = []string{
		"lanzouw.com",
		"lanzoui.com",
		"lanzoux.com",
		"lanzouy.com",
		"lanzoup.com",
	}
)

func checkOrSwitchHostname(c *client.Client) error {
	checkHostname := func(hostname string) bool {
		c.SetDefaultHostname(hostname)
		head, err := c.R().Head("/")
		return err == nil && !head.IsError()
	}

	for _, hostnames := range availableHostnames {
		if available := checkHostname(hostnames); available {
			return nil
		}
	}

	return fmt.Errorf("no available lanzou hostname")
}

func New(config *client.Config) (*Lanzou, error) {
	cl, err := client.New(&client.Config{
		HTTPS:      true,
		Host:       availableHostnames[0],
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

	if err := checkOrSwitchHostname(cl); err != nil {
		return nil, err
	}

	return &Lanzou{Client: cl}, nil
}

func (l *Lanzou) DownloadFile(downloadURL string) (io.ReadCloser, int64, error) {
	log.Debugf("Start to download file from aliyun drive: %s", downloadURL)

	resp, err := l.R().
		SetDoNotParseResponse(true).
		Get(downloadURL)
	if err != nil {
		return nil, 0, err
	}

	return resp.RawBody(), resp.RawResponse.ContentLength, nil
}
