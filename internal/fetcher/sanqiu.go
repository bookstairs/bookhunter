package fetcher

import (
	"errors"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
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
	var d driver.Driver
	refreshToken := config.Property("refreshToken")
	if refreshToken != "" {
		d, err = driver.New(driver.ALIYUN, config.Config, map[string]string{"refreshToken": refreshToken})
	} else {
		d, err = driver.New(driver.TELECOM, config.Config, map[string]string{
			"username": config.Property("telecomUsername"),
			"password": config.Property("telecomPassword"),
		})
	}
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
	err = sanqiu.ParseAPIResponse(resp, books)
	if err != nil {
		return 0, err
	}

	if len(books) < 1 {
		return 0, ErrEmptySanqiu
	}

	return books[0].ID, nil
}

func (s *sanqiuService) formats(id int64) (map[Format]driver.Share, error) {
	// TODO implement me
	panic("implement me")
}

func (s *sanqiuService) fetch(id int64, format Format, share driver.Share) (*fetch, error) {
	// TODO implement me
	panic("implement me")
}
