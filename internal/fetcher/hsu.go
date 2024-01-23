package fetcher

import (
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
)

// HsuBookMeta is used for representing book information.
// https://book.hsu.life/api/series/all-v2
type HsuBookMeta struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	OriginalName        string `json:"originalName"`
	LocalizedName       string `json:"localizedName"`
	SortName            string `json:"sortName"`
	Pages               int    `json:"pages"`
	CoverImageLocked    bool   `json:"coverImageLocked"`
	PagesRead           int    `json:"pagesRead"`
	LatestReadDate      string `json:"latestReadDate"`
	LastChapterAdded    string `json:"lastChapterAdded"`
	UserRating          int    `json:"userRating"`
	HasUserRated        bool   `json:"hasUserRated"`
	Format              int    `json:"format"`
	Created             string `json:"created"`
	NameLocked          bool   `json:"nameLocked"`
	SortNameLocked      bool   `json:"sortNameLocked"`
	LocalizedNameLocked bool   `json:"localizedNameLocked"`
	WordCount           int    `json:"wordCount"`
	LibraryID           int    `json:"libraryId"`
	LibraryName         string `json:"libraryName"`
	MinHoursToRead      int    `json:"minHoursToRead"`
	MaxHoursToRead      int    `json:"maxHoursToRead"`
	AvgHoursToRead      int    `json:"avgHoursToRead"`
	FolderPath          string `json:"folderPath"`
	LastFolderScanned   string `json:"lastFolderScanned"`
}

type HsuBookMetaReq struct {
	Statements []struct {
		Comparison int    `json:"comparison"`
		Value      string `json:"value"`
		Field      int    `json:"field"`
	} `json:"statements"`
	Combination int `json:"combination"`
	LimitTo     int `json:"limitTo"`
	SortOptions struct {
		IsAscending bool `json:"isAscending"`
		SortField   int  `json:"sortField"`
	} `json:"sortOptions"`
}

// HsuLoginReq is used in POST https://book.hsu.life/api/account/login
type HsuLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	APIKey   string `json:"apiKey"`
}

// HsuLoginResp is used in https://book.hsu.life/api/account/login response
type HsuLoginResp struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	APIKey       string `json:"apiKey"`
	Preferences  struct {
		ReadingDirection           int    `json:"readingDirection"`
		ScalingOption              int    `json:"scalingOption"`
		PageSplitOption            int    `json:"pageSplitOption"`
		ReaderMode                 int    `json:"readerMode"`
		LayoutMode                 int    `json:"layoutMode"`
		EmulateBook                bool   `json:"emulateBook"`
		BackgroundColor            string `json:"backgroundColor"`
		SwipeToPaginate            bool   `json:"swipeToPaginate"`
		AutoCloseMenu              bool   `json:"autoCloseMenu"`
		ShowScreenHints            bool   `json:"showScreenHints"`
		BookReaderMargin           int    `json:"bookReaderMargin"`
		BookReaderLineSpacing      int    `json:"bookReaderLineSpacing"`
		BookReaderFontSize         int    `json:"bookReaderFontSize"`
		BookReaderFontFamily       string `json:"bookReaderFontFamily"`
		BookReaderTapToPaginate    bool   `json:"bookReaderTapToPaginate"`
		BookReaderReadingDirection int    `json:"bookReaderReadingDirection"`
		BookReaderWritingStyle     int    `json:"bookReaderWritingStyle"`
		Theme                      struct {
			ID              int    `json:"id"`
			Name            string `json:"name"`
			NormalizedName  string `json:"normalizedName"`
			FileName        string `json:"fileName"`
			IsDefault       bool   `json:"isDefault"`
			Provider        int    `json:"provider"`
			Created         string `json:"created"`
			LastModified    string `json:"lastModified"`
			CreatedUtc      string `json:"createdUtc"`
			LastModifiedUtc string `json:"lastModifiedUtc"`
		} `json:"theme"`
		BookReaderThemeName         string `json:"bookReaderThemeName"`
		BookReaderLayoutMode        int    `json:"bookReaderLayoutMode"`
		BookReaderImmersiveMode     bool   `json:"bookReaderImmersiveMode"`
		GlobalPageLayoutMode        int    `json:"globalPageLayoutMode"`
		BlurUnreadSummaries         bool   `json:"blurUnreadSummaries"`
		PromptForDownloadSize       bool   `json:"promptForDownloadSize"`
		NoTransitions               bool   `json:"noTransitions"`
		CollapseSeriesRelationships bool   `json:"collapseSeriesRelationships"`
		ShareReviews                bool   `json:"shareReviews"`
		Locale                      string `json:"locale"`
	} `json:"preferences"`
	AgeRestriction struct {
		AgeRating       int  `json:"ageRating"`
		IncludeUnknowns bool `json:"includeUnknowns"`
	} `json:"ageRestriction"`
	KavitaVersion string `json:"kavitaVersion"`
}

type hsuService struct {
	config *Config
	books  []*HsuBookMeta
	*client.Client
}

// formatMapping is defined in https://github.com/Kareadita/Kavita/blob/develop/UI/Web/src/app/_models/manga-format.ts
var formatMapping = map[int]file.Format{
	1: file.ZIP,
	3: file.EPUB,
	4: file.PDF,
}

func newHsuService(config *Config) (service, error) {
	// Create the resty client for HTTP handing.
	c, err := client.New(config.Config)
	if err != nil {
		return nil, err
	}

	resp, err := c.R().
		SetBody(&HsuLoginReq{
			Username: config.Property("username"),
			Password: config.Property("password"),
			APIKey:   "",
		}).
		SetResult(&HsuLoginResp{}).
		ForceContentType("application/json").
		Post("/api/account/login")

	if err != nil {
		return nil, err
	}

	token := resp.Result().(*HsuLoginResp).Token
	if token == "" {
		return nil, fmt.Errorf("invalid login credential")
	}
	c.SetAuthToken(token)

	// Download books.
	resp, err = c.R().
		SetBody(&HsuBookMetaReq{
			Statements: []struct {
				Comparison int    `json:"comparison"`
				Value      string `json:"value"`
				Field      int    `json:"field"`
			}{
				{
					Comparison: 0,
					Value:      "",
					Field:      1,
				},
			},
			Combination: 1,
			LimitTo:     0,
			SortOptions: struct {
				IsAscending bool `json:"isAscending"`
				SortField   int  `json:"sortField"`
			}{
				IsAscending: true,
				SortField:   1,
			},
		}).
		SetResult(&[]HsuBookMeta{}).
		Post("/api/series/all-v2")

	if err != nil {
		return nil, err
	}

	metas := *resp.Result().(*[]HsuBookMeta)
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].ID < metas[j].ID
	})

	// Create a better slice for holding the books.
	books := make([]*HsuBookMeta, metas[len(metas)-1].ID)
	for i := range metas {
		meta := metas[i]
		books[meta.ID-1] = &meta
	}

	return &hsuService{config: config, Client: c, books: books}, nil
}

func (h *hsuService) size() (int64, error) {
	return int64(len(h.books)), nil
}

func (h *hsuService) formats(i int64) (map[file.Format]driver.Share, error) {
	book := h.books[i-1]

	if book != nil {
		if format, ok := formatMapping[book.Format]; ok {
			resp, err := h.R().
				SetQueryParam("seriesId", strconv.Itoa(int(i))).
				Get("/api/download/series-size")
			if err != nil {
				return nil, err
			}

			fileSize := resp.String()
			if fileSize == "" {
				return nil, fmt.Errorf("you are not allowed to download books")
			}
			size, err := strconv.ParseInt(fileSize, 10, 64)
			if err != nil {
				return nil, err
			}

			return map[file.Format]driver.Share{
				format: {
					FileName: book.Name,
					Size:     size,
				},
			}, nil
		}
	}

	return make(map[file.Format]driver.Share), nil
}

func (h *hsuService) fetch(i int64, _ file.Format, _ driver.Share, writer file.Writer) error {
	resp, err := h.R().
		SetDoNotParseResponse(true).
		SetQueryParam("seriesId", strconv.Itoa(int(i))).
		Get("/api/download/series")
	if err != nil {
		return err
	}

	body := resp.RawBody()
	defer func() { _ = body.Close() }()

	// Save the download content info files.
	_, err = io.Copy(writer, body)
	return err
}
