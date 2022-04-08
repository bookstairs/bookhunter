package talebook

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/syhily/bookhunter/pkg/progress"

	"github.com/syhily/bookhunter/pkg/spider"

	"github.com/syhily/bookhunter/pkg/log"
)

// The download config.
type config struct {
	Website       string        // The website for talebook.
	Username      string        // The login user.
	Password      string        // The password for login user.
	DownloadPath  string        // Use the executed directory as the default download path.
	CookieFile    string        // The cookie file to use in this download progress.
	ProgressFile  string        // The progress file serving the remaining book id.
	InitialBookID int           // The book id start to download.
	Formats       []string      // The file formats you want to download
	Timeout       time.Duration // The request timeout for a single request.
	Retry         int           // The maximum retry times for a timeout request.
	UserAgent     string        // The user agent for the download request.
	Rename        bool          // Rename the file by using book ID.
}

// The main instance for start downloading the book.
type talebook struct {
	wait       *sync.WaitGroup
	downloader *downloadWorker
}

// NewDownloadConfig will return a default blank config.
func NewDownloadConfig() *config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return &config{
		DownloadPath:  dir,
		CookieFile:    "cookies",
		ProgressFile:  "progress",
		InitialBookID: 1,
		Formats:       []string{"EPUB", "MOBI", "PDF"},
		Timeout:       10 * time.Minute,
		Retry:         5,
		UserAgent:     spider.DefaultUserAgent,
		Rename:        false,
	}
}

// NewTalebook will create the download instance.
func NewTalebook(c *config) *talebook {
	// Create cookiejar.
	cookieFile := path.Join(c.DownloadPath, c.CookieFile)
	cookieJar, err := spider.NewCookieJar(cookieFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create common http client.
	client := &http.Client{Jar: cookieJar, Timeout: c.Timeout}

	// Disable login redirect.
	loginUrl := spider.GenerateUrl(c.Website, "/login")
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if req.URL.String() == loginUrl {
			return spider.ErrNeedSignin
		}

		// Allow 10 redirects by default.
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}

	// Try to signin if required.
	if err := login(c.Username, c.Password, c.Website, c.UserAgent, client); err != nil {
		log.Fatal(err)
	}

	// Try to find last book ID.
	last, err := lastBookID(c.Website, c.UserAgent, client)
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
	downloader := &downloadWorker{
		website:      c.Website,
		progress:     p,
		client:       client,
		userAgent:    c.UserAgent,
		retry:        c.Retry,
		downloadPath: c.DownloadPath,
		formats:      c.Formats,
		rename:       c.Rename,
	}

	return &talebook{
		downloader: downloader,
	}
}

func (t *talebook) Start() {
	log.Infof("Start to download book.")
	t.downloader.Download()
}

// login to the given website by username and password. We will save the cookie into file.
// Thus, you don't need to signin twice.
func login(username, password, website, userAgent string, client *http.Client) error {
	if username == "" || password == "" {
		// No need to login.
		return nil
	}

	log.Info("You have provided user information, start to login.")

	site := spider.GenerateUrl(website, "/api/user/sign_in")
	referer := spider.GenerateUrl(website, "/login")
	values := url.Values{
		"username": {username},
		"password": {password},
	}

	req, err := http.NewRequest(http.MethodPost, site, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("illegal login request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = form.Body.Close() }()
	if form.StatusCode != http.StatusOK {
		return errors.New(form.Status)
	}

	result := &LoginResponse{}
	if err := spider.DecodeResponse(form, result); err != nil {
		return err
	}

	if result.Err != SuccessStatus {
		return errors.New(result.Msg)
	}

	log.Info("Login success. Save cookies into file.")
	return nil
}

// lastBookID will return the last available book ID.
func lastBookID(website, userAgent string, client *http.Client) (int64, error) {
	site := spider.GenerateUrl(website, "/api/recent")
	referer := spider.GenerateUrl(website, "/recent")

	req, err := http.NewRequest(http.MethodGet, site, http.NoBody)
	if err != nil {
		return 0, fmt.Errorf("illegal book id request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("referer", referer)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}

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
