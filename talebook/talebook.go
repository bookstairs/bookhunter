package talebook

import (
	"errors"
	"net/http"
	"path"

	"github.com/bibliolater/bookhunter/pkg/log"
	"github.com/bibliolater/bookhunter/pkg/progress"
	"github.com/bibliolater/bookhunter/pkg/spider"
)

var ErrNeedSignin = errors.New("need user account to download books")

// NewDownloader will create the download instance.
func NewDownloader(c *spider.DownloadConfig) *downloadWorker {
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
	last, err := lastBookID(c.Website, client)
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
		website:      c.Website,
		progress:     p,
		client:       client,
		userAgent:    c.UserAgent,
		retry:        c.Retry,
		downloadPath: c.DownloadPath,
		formats:      c.Formats,
		rename:       c.Rename,
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

// lastBookID will return the last available book ID.
func lastBookID(website string, client *spider.Client) (int64, error) {
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
