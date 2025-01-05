package fetcher

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/sobooks"
)

var (
	sobooksIDRe = regexp.MustCompile(`.*?/(\d+?).html`)
)

type sobooksService struct {
	config *Config
	*client.Client
	driver driver.Driver
}

func newSobooksService(config *Config) (service, error) {
	// Create the resty client for HTTP handing.
	c, err := client.New(config.Config)
	// Set code for viewing hidden content
	c.SetCookie(&http.Cookie{
		Name:   "mpcode",
		Value:  config.Property("code"),
		Path:   "/",
		Domain: config.Host,
	}).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec

	if err != nil {
		return nil, err
	}

	// Create the net disk driver.
	d, err := driver.New(config.Config, config.Properties)
	if err != nil {
		return nil, err
	}

	return &sobooksService{config: config, Client: c, driver: d}, nil
}

func (s *sobooksService) size() (int64, error) {
	resp, err := s.R().
		Get("/")
	if err != nil {
		return 0, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
	if err != nil {
		return 0, err
	}

	lastID := -1

	// Find all the links is case of the website primary changed the theme.
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		// This is a book link.
		link, exists := selection.Attr("href")
		if exists {
			// Extract bookId.
			match := sobooksIDRe.FindStringSubmatch(link)
			if len(match) > 0 {
				id, _ := strconv.Atoi(match[1])
				if id > lastID {
					lastID = id
				}
			}
		}
	})
	return int64(lastID), nil
}

func (s *sobooksService) formats(id int64) (map[file.Format]driver.Share, error) {
	resp, err := s.R().
		SetPathParam("bookId", strconv.FormatInt(id, 10)).
		SetHeader("referer", s.BaseURL).
		Get("/books/{bookId}.html")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		log.Debugf("The current book [%v] content does not exist ", id)
		return map[file.Format]driver.Share{}, nil
	}

	title, links, err := sobooks.ParseLinks(resp.String(), id)
	if err != nil {
		return map[file.Format]driver.Share{}, nil
	}

	res := make(map[file.Format]driver.Share)
	for source, link := range links {
		if source != s.driver.Source() {
			continue
		}

		shares, err := s.driver.Resolve(link.URL, link.Code)
		if err != nil {
			return nil, err
		}

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
	}

	if link, ok := links[driver.DIRECT]; ok && len(res) == 0 {
		res[file.EPUB] = driver.Share{
			FileName: title + ".epub",
			Size:     0,
			URL:      link.URL,
		}
	}
	return res, nil
}

func (s *sobooksService) fetch(_ int64, _ file.Format, share driver.Share, writer file.Writer) error {
	u, err := url.Parse(share.URL)
	if err != nil {
		return err
	}
	resp, err := s.R().
		SetDoNotParseResponse(true).
		Get(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() == 404 {
		return ErrFileNotExist
	}
	body := resp.RawBody()
	defer func() { _ = body.Close() }()

	// Save the download content info files.
	writer.SetSize(resp.RawResponse.ContentLength)
	_, err = io.Copy(writer, body)
	return err
}
