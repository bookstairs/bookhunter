package fetcher

import (
	"fmt"
	"io"
	"regexp"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/wordpress"
)

var (
	// driveNamings is the chinese name mapping of the drive's provider.
	driveNamings = map[driver.Source]string{
		driver.ALIYUN:  "阿里",
		driver.LANZOU:  "蓝奏",
		driver.TELECOM: "天翼",
	}

	sanqiuPasscodeRe = regexp.MustCompile(".*?([a-zA-Z0-9]+).*?")

	tianLangLinkRe     = regexp.MustCompile(`(?m)location\.href\s?=\s?"(http.*?)(\n?)";`)
	tianlangPasscodeRe = regexp.MustCompile("密码.*?([a-zA-Z0-9]+).*?")
)

type (
	shareLinkResolver func(*client.Client, int64) (map[driver.Source]wordpress.ShareLink, error)

	wordpressService struct {
		config   *Config
		client   *client.Client
		driver   driver.Driver
		resolver shareLinkResolver
	}
)

func newWordpressService(config *Config, resolver shareLinkResolver) (service, error) {
	// Create the resty client for HTTP handing.
	c, err := client.New(config.Config)
	if err != nil {
		return nil, err
	}

	// Create the net disk driver.
	d, err := driver.New(config.Config, config.Properties)
	if err != nil {
		return nil, err
	}

	return &wordpressService{
		config:   config,
		client:   c,
		driver:   d,
		resolver: resolver,
	}, nil
}

func (w *wordpressService) size() (int64, error) {
	resp, err := w.client.R().
		SetQueryParams(map[string]string{
			"orderby":  "id",
			"order":    "desc",
			"per_page": "1",
		}).
		Get("/wp-json/wp/v2/posts")
	if err != nil {
		return 0, err
	}

	posts, err := wordpress.ParsePosts(resp)
	if err != nil {
		return 0, err
	}
	if len(posts) < 1 {
		return 0, fmt.Errorf("couldn't find available books in %s", w.config.Category)
	}

	return posts[0].ID, nil
}

func (w *wordpressService) formats(id int64) (map[file.Format]driver.Share, error) {
	links, err := w.resolver(w.client, id)
	if err != nil {
		log.Fatalf("Error in find downloadable links %v", err)
		return map[file.Format]driver.Share{}, nil
	}

	log.Debugf("Available download links: %v", links)

	for source, link := range links {
		if source != w.driver.Source() {
			continue
		}

		shares, err := w.driver.Resolve(link.URL, link.Code)
		if err != nil {
			log.Warnf("Error in resolve the share link for book %d, %v", id, err)
			break
		}

		res := make(map[file.Format]driver.Share)
		for _, share := range shares {
			if ext, has := file.LinkExtension(share.FileName); has {
				if IsValidFormat(ext) {
					res[ext] = share
				} else {
					log.Debugf("The file name %s don't have valid extension %s", share.FileName, ext)
				}
			} else {
				log.Debugf("The file name %s don't have the extension", share.FileName)
			}
		}

		return res, nil
	}

	log.Debugf("No downloadable files found in this book id %d.", id)
	return map[file.Format]driver.Share{}, nil
}

func (w *wordpressService) fetch(_ int64, _ file.Format, share driver.Share, writer file.Writer) error {
	content, size, err := w.driver.Download(share)
	if err != nil {
		return err
	}
	defer func() { _ = content.Close() }()

	// Set download progress if required.
	if size > 0 {
		writer.SetSize(size)
	}

	// Save the download content info files.
	_, err = io.Copy(writer, content)
	return err
}
