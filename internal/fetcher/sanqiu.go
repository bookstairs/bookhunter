package fetcher

import (
	"errors"
	"io"
	"strconv"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/naming"
	"github.com/bookstairs/bookhunter/internal/sanqiu"
)

var (
	ErrEmptySanqiu = errors.New("couldn't find available books in sanqiu")
)

type sanqiuService struct {
	config *Config
	client *client.Client
	driver driver.Driver
}

func newSanqiuService(config *Config) (service, error) {
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

	return &sanqiuService{
		config: config,
		client: c,
		driver: d,
	}, nil
}

func (s *sanqiuService) size() (int64, error) {
	resp, err := s.client.R().
		SetQueryParams(map[string]string{
			"orderby":  "id",
			"order":    "desc",
			"per_page": "1",
		}).
		Get("/wp-json/wp/v2/posts")

	if err != nil {
		return 0, err
	}

	books := make([]sanqiu.BookResp, 0, 1)
	err = sanqiu.ParseAPIResponse(resp, &books)
	if err != nil {
		return 0, err
	}

	if len(books) < 1 {
		return 0, ErrEmptySanqiu
	}

	return books[0].ID, nil
}

func (s *sanqiuService) formats(id int64) (map[Format]driver.Share, error) {
	resp, err := s.client.R().
		SetQueryParam("id", strconv.FormatInt(id, 10)).
		Get("/download.php")
	if err != nil {
		return nil, err
	}

	links, err := sanqiu.DownloadLinks(resp.String())
	if err != nil {
		return nil, err
	}

	for source, link := range links {
		if source != s.driver.Source() {
			continue
		}

		shares, err := s.driver.Resolve(link.URL, link.URL)
		if err != nil {
			return nil, err
		}

		res := make(map[Format]driver.Share)
		for _, share := range shares {
			if ext, has := naming.Extension(share.FileName); has {
				if format, err := ParseFormat(ext); err == nil {
					res[format] = share
				} else {
					log.Debugf("The file name %s don't have valid extension %s", share.FileName, ext)
				}
			} else {
				log.Debugf("The file name %s don't have the extension", share.FileName)
			}
		}

		return res, nil
	}

	log.Debug("No downloadable files found in this link.")
	return map[Format]driver.Share{}, nil
}

func (s *sanqiuService) fetch(_ int64, _ Format, share driver.Share, writer io.Writer) error {
	content, err := s.driver.Download(share)
	if err != nil {
		return err
	}
	defer func() { _ = content.Close() }()

	// Save the download content info files.
	_, err = io.Copy(writer, content)
	return err
}
