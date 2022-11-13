package fetcher

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/talebook"
)

var (
	ErrTalebookNeedSignin = errors.New("need user account to download books")
	ErrEmptyTalebook      = errors.New("couldn't find available books in talebook")

	redirectHandler = func(request *http.Request, requests []*http.Request) error {
		if request.URL.Path == "/login" {
			return ErrTalebookNeedSignin
		}
		return nil
	}
)

type talebookService struct {
	config *Config
	client *client.Client
}

func newTalebookService(config *Config) (service, error) {
	// Add login check in redirect handler.
	if err := config.SetRedirect(redirectHandler); err != nil {
		return nil, err
	}

	// Create the resty client for HTTP handing.
	c, err := client.New(config.Config)
	if err != nil {
		return nil, err
	}

	// Start to sign in if required.
	username := config.Property("username")
	password := config.Property("password")
	if username != "" && password != "" {
		log.Info("You have provided user information, start to login.")
		resp, err := c.R().
			SetFormData(map[string]string{
				"username": username,
				"password": password,
			}).
			SetResult(&talebook.LoginResp{}).
			ForceContentType("application/json").
			Post("/api/user/sign_in")

		if err != nil {
			return nil, err
		}

		result := resp.Result().(*talebook.LoginResp)
		if result.Err != talebook.SuccessStatus {
			return nil, errors.New(result.Msg)
		}

		log.Info("Login success. Save cookies into file.")
	}

	return &talebookService{
		config: config,
		client: c,
	}, nil
}

func (t *talebookService) size() (int64, error) {
	resp, err := t.client.R().
		SetResult(&talebook.BooksResp{}).
		Get("/api/recent")
	if err != nil {
		return 0, err
	}

	result := resp.Result().(*talebook.BooksResp)
	if result.Err != talebook.SuccessStatus {
		return 0, errors.New(result.Msg)
	}

	bookID := int64(0)
	for _, book := range result.Books {
		if book.ID > bookID {
			bookID = book.ID
		}
	}

	if bookID == 0 {
		return 0, ErrEmptyTalebook
	}

	return bookID, nil
}

func (t *talebookService) formats(id int64) (map[Format]string, error) {
	resp, err := t.client.R().
		SetResult(&talebook.BookResp{}).
		SetPathParam("bookID", strconv.FormatInt(id, 10)).
		Get("/api/book/{bookID}")
	if err != nil {
		return nil, err
	}

	result := resp.Result().(*talebook.BookResp)
	switch result.Err {
	case talebook.SuccessStatus:
		formats := make(map[Format]string)
		for _, file := range result.Book.Files {
			format, err := ParseFormat(file.Format)
			if err != nil {
				return nil, err
			}
			formats[format] = file.Href
		}
		return formats, nil
	case talebook.BookNotFoundStatus:
		return nil, nil
	default:
		return nil, errors.New(result.Msg)
	}
}

func (t *talebookService) fetch(_ int64, _ Format, url string) (*fetch, error) {
	resp, err := t.client.R().
		SetDoNotParseResponse(true).
		Get(url)
	if err != nil {
		return nil, err
	}

	return createFetch(resp), nil
}
