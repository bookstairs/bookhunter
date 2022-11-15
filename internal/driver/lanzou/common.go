package lanzou

import (
	"io"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
)

type Drive struct {
	client *client.Client
}

// hostname is used to return a valid lanzou host.
func hostname() string {
	// We may add a roundrobin list in the future.
	return "lanzoux.com"
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

	return &Drive{client: cl}, nil
}

func (l Drive) DownloadFile(downloadURL string) (io.ReadCloser, int64, error) {
	log.Debugf("Start to download file from aliyun drive: %s", downloadURL)

	resp, err := l.client.R().
		SetDoNotParseResponse(true).
		Get(downloadURL)
	response := resp.RawResponse

	return response.Body, response.ContentLength, err
}
