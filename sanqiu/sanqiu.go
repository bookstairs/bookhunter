package sanqiu

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/bibliolater/bookhunter/pkg/rename"

	"github.com/PuerkitoBio/goquery"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/progress"
	"github.com/bibliolater/bookhunter/pkg/spider"
)

var (
	// The Website for sanqiu.
	Website  = "https://www.sanqiu.cc"
	bookIDRe = regexp.MustCompile(".*?/(\\d+?).html")
)

type downloader struct {
	config       *spider.Config
	progress     *progress.Progress
	client       *spider.Client
	retry        int
	downloadPath string
	formats      []string
	rename       bool
	wait         *sync.WaitGroup
}

func NewDownloader(config *spider.Config) *downloader {
	// Create common http client.
	client := spider.NewClient(config)
	client.CheckRedirect(func(req *http.Request, via []*http.Request) error {
		// Allow 10 redirects by default.
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	})

	// Get last book ID
	last, err := latestBookID(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Find the last book ID: %d", last)

	// Create book storage.
	storageFile := path.Join(config.DownloadPath, config.ProgressFile)
	p, err := progress.NewProgress(int64(config.InitialBookID), last, storageFile)
	if err != nil {
		log.Fatal(err)
	}

	return &downloader{
		config:       config,
		progress:     p,
		client:       client,
		retry:        config.Retry,
		downloadPath: config.DownloadPath,
		formats:      config.Formats,
		rename:       config.Rename,
		wait:         new(sync.WaitGroup),
	}
}

// latestBookID will return the last available book ID.
func latestBookID(client *spider.Client) (int64, error) {
	resp, err := client.Get(Website, "")
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}

	lastID := -1

	// Find all the links is case of the website master changed the theme.
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if exists {
			match := bookIDRe.FindStringSubmatch(link)
			// This is a book link.
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

func (d *downloader) Fork() {
	d.wait.Add(1)
}

// Download would start download books from given website.
func (d *downloader) Download() {
	bookID := d.progress.AcquireBookID()
	log.Infof("Start to download book from %d.", bookID)

	// Try to acquire book ID from storage.
	for ; bookID != progress.NoBookToDownload; bookID = d.progress.AcquireBookID() {
		// Acquire book metadata from website.
		metadata := d.bookMetadata(bookID)
		if metadata == nil {
			log.Warnf("[%d/%d] Book with ID %d is not exist on target website.", bookID, d.progress.Size(), bookID)
			d.downloadedBook(bookID)
			continue
		}

		// Download books from telecom
		if link, ok := metadata.Links[TELECOM]; ok {
			links, err := spider.ResolveTelecom(d.client, link.Url, link.Code, d.formats...)
			if err != nil {
				log.Fatal(err)
			}

			if len(links) == 0 {
				log.Warnf("[%d/%d] No downloadable links found, this resource could be banned.", bookID, d.progress.Size())
			}

			for _, l := range links {
				for i := 0; i < d.retry; i++ {
					err := d.downloadBook(metadata, l)
					if err == nil {
						break
					} else if spider.IsTimeOut(err) && i < d.retry {
						continue
					} else {
						log.Fatal(err)
					}
				}
			}
		} else {
			log.Warnf("[%d/%d] Book with ID %d don't have telecom link, skip.", bookID, d.progress.Size(), bookID)
		}

		// Finished the book download.
		d.downloadedBook(bookID)
	}

	d.wait.Done()
}

func (d *downloader) Join() {
	d.wait.Wait()
}

// downloadBook would download the book to saving path.
func (d *downloader) downloadBook(meta *BookMeta, link string) error {
	resp, err := d.client.Get(link, "")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Generate file name.
	format := spider.Extension(link)
	filename := strconv.FormatInt(meta.Id, 10) + "." + strings.ToLower(format)
	if !d.rename {
		name := spider.Filename(resp)
		if name != "" {
			filename = name
		}
	}

	// Remove illegal characters. Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	filename = rename.EscapeFilename(filename)

	// Generate the file path.
	file := filepath.Join(d.downloadPath, filename)

	// Remove the exist file.
	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	// Create file writer.
	writer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = writer.Close() }()

	// Add download progress
	bar := log.NewProgressBar(meta.Id, d.progress.Size(), format+" "+meta.Title, resp.ContentLength)

	// Write file content
	_, err = io.Copy(io.MultiWriter(writer, bar), resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// downloadedBook would record the download statue into storage.
func (d *downloader) downloadedBook(bookID int64) {
	if err := d.progress.SaveBookID(bookID); err != nil {
		log.Fatal(err)
	}
}
