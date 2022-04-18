package talebook

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/progress"
	"github.com/bibliolater/bookhunter/pkg/rename"
	"github.com/bibliolater/bookhunter/pkg/spider"
)

var ErrNeedSignin = errors.New("need user account to download books")

// downloadWorker is the download instance.
type downloadWorker struct {
	progress *progress.Progress
	client   *spider.Client
	config   *spider.Config
}

// NewDownloader will create the download instance.
func NewDownloader(c *spider.Config) *downloadWorker {
	// Create common http client.
	client := spider.NewClient(c)

	// Disable login redirect.
	loginUrl := spider.GenerateUrl(c.Website, "/login")
	client.CheckRedirect(
		func(req *http.Request, via []*http.Request) error {
			if req.URL.String() == loginUrl {
				return ErrNeedSignin
			}

			// Allow 10 redirects by default.
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	)

	// Try to signin if required.
	if err := login(c.Username, c.Password, c.Website, client); err != nil {
		log.Fatal(err)
	}

	// Try to find last book ID.
	last, err := latestBookID(c.Website, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Find the last book ID: %d", last)

	// Create book storage.
	storageFile := path.Join(c.DownloadPath, c.ProgressFile)
	p, err := progress.NewProgress(int64(c.InitialBookID), last, storageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create download worker
	return &downloadWorker{
		progress: p,
		client:   client,
		config:   c,
	}
}

// login to the given website by username and password. We will save the cookie into file.
// Thus, you don't need to signin twice.
func login(username, password, website string, client *spider.Client) error {
	if username == "" || password == "" {
		// No need to login.
		return nil
	}

	log.Info("You have provided user information, start to login.")

	site := spider.GenerateUrl(website, "/api/user/sign_in")
	referer := spider.GenerateUrl(website, "/login")
	form := spider.Form{
		spider.Field{Key: "username", Value: username},
		spider.Field{Key: "password", Value: password},
	}

	resp, err := client.FormPost(site, referer, form)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	result := &LoginResponse{}
	if err := spider.DecodeResponse(resp, result); err != nil {
		return err
	}

	if result.Err != SuccessStatus {
		return errors.New(result.Msg)
	}

	log.Info("Login success. Save cookies into file.")
	return nil
}

// latestBookID will return the last available book ID.
func latestBookID(website string, client *spider.Client) (int64, error) {
	site := spider.GenerateUrl(website, "/api/recent")
	referer := spider.GenerateUrl(website, "/recent")

	resp, err := client.Get(site, referer)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	result := &BookListResponse{}
	if err := spider.DecodeResponse(resp, result); err != nil {
		return 0, err
	}

	if result.Err != SuccessStatus {
		return 0, errors.New(result.Msg)
	}

	bookID := int64(0)
	for _, book := range result.Books {
		if book.ID > bookID {
			bookID = book.ID
		}
	}

	if bookID == 0 {
		return 0, errors.New("couldn't find available books")
	}

	return bookID, nil
}

// Download would start download books from given website.
func (worker *downloadWorker) Download() {
	bookID := worker.progress.AcquireBookID()
	log.Infof("Start to download book from %d.", bookID)

	// Try to acquire book ID from storage.
	for ; bookID != progress.NoBookToDownload; bookID = worker.progress.AcquireBookID() {
		// Acquire book info.
		var info *BookResponse
		err := worker.client.Retry(func() error {
			var err error
			info, err = worker.queryBookInfo(bookID)
			return err
		})
		if err != nil {
			log.Fatal(err)
		}

		if info == nil {
			log.Warnf("[%d/%d] Book with ID %d is not exist on target website.", bookID, worker.progress.Size(), bookID)
			worker.downloadedBook(bookID)
			continue
		}

		// Find formats to download.
		for _, file := range info.Book.Files {
			err := worker.client.Retry(func() error {
				return worker.downloadBook(bookID, info.Book.Title, file.Format, file.Href)
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		worker.downloadedBook(bookID)
	}
}

// queryBookInfo will find the required book information.
func (worker *downloadWorker) queryBookInfo(bookID int64) (*BookResponse, error) {
	site := spider.GenerateUrl(worker.config.Website, "/api/book", strconv.FormatInt(bookID, 10))

	resp, err := worker.client.Get(site, "")
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	result := &BookResponse{}
	if err := spider.DecodeResponse(resp, result); err != nil {
		return nil, err
	}

	switch result.Err {
	case SuccessStatus:
		return result, nil
	case BookNotFoundStatus:
		return nil, nil
	default:
		return nil, errors.New(result.Msg)
	}
}

// downloadBook will download the book file from
func (worker *downloadWorker) downloadBook(bookID int64, title, format, href string) error {
	valid := false
	for _, f := range worker.config.Formats {
		if f == format {
			valid = true
			break
		}
	}

	if !valid {
		// Skip this format.
		return nil
	}

	// Start download.
	site := ""
	if strings.HasPrefix(href, "http") {
		// Backward API support.
		site = href
	} else {
		site = spider.GenerateUrl(worker.config.Website, href)
	}

	resp, err := worker.client.Get(site, "")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Generate file name.
	filename := strconv.FormatInt(bookID, 10) + "." + strings.ToLower(format)
	// Use readable name.
	if !worker.config.Rename {
		name := spider.Filename(resp)
		if name == "" {
			filename = title + "." + strings.ToLower(format)
		} else {
			filename = name
		}
	}

	// Remove illegal characters. Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	filename = rename.EscapeFilename(filename)

	// Generate the file path.
	file := filepath.Join(worker.config.DownloadPath, filename)

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
	bar := log.NewProgressBar(bookID, worker.progress.Size(), format+" "+title, resp.ContentLength)

	// Write file content
	_, err = io.Copy(io.MultiWriter(writer, bar), resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// downloadedBook would record the download statue into storage.
func (worker *downloadWorker) downloadedBook(bookID int64) {
	if err := worker.progress.SaveBookID(bookID); err != nil {
		log.Fatal(err)
	}
}
