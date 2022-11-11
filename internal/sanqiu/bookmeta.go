package sanqiu

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/spider"
)

type LinkType string

const (
	ALIYUN  LinkType = "aliyun"
	LANZOU  LinkType = "lanzou"
	CTFILE  LinkType = "ctfile"
	TELECOM LinkType = "telecom"
)

var (
	// driveNamings is the chinese name mapping of the drive's provider.
	driveNamings = map[LinkType]string{
		ALIYUN:  "阿里",
		LANZOU:  "蓝奏",
		CTFILE:  "城通",
		TELECOM: "天翼",
	}

	passcodeRe = regexp.MustCompile(".*?([a-zA-Z0-9]+).*?")
)

type BookMeta struct {
	ID    int64                  // The book id in website.
	Title string                 // The name of the book.
	Links map[LinkType]*BookLink // All the available download links from drive.
}

type BookLink struct {
	URL  string // The url to access the download page.
	Code string // The passcode for querying the file content.
}

// bookMetadata will find all the available books.
func (d *downloader) bookMetadata(bookID int64) *BookMeta {
	page := spider.GenerateURL(d.config.Website, "/download.php?id="+strconv.FormatInt(bookID, 10))
	referer := spider.GenerateURL(d.config.Website, "/"+strconv.FormatInt(bookID, 10)+".html")

	resp, err := d.client.Get(page, referer)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find all the links.
	links := map[LinkType]*BookLink{}
	doc.Find(".downfile a").Each(func(i int, selection *goquery.Selection) {
		driveName := selection.Text()
		href, exists := selection.Attr("href")
		if exists {
			for linkType, name := range driveNamings {
				if strings.Contains(driveName, name) {
					links[linkType] = &BookLink{URL: href}
					break
				}
			}
		}
	})

	if len(links) == 0 {
		return nil
	}

	// Find all the passcodes.
	doc.Find(".plus_l li").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		for linkType, link := range links {
			name := driveNamings[linkType]
			if strings.Contains(text, name) {
				match := passcodeRe.FindStringSubmatch(text)
				if len(match) == 2 {
					links[linkType] = &BookLink{
						URL:  link.URL,
						Code: match[1],
					}
				}
			}
		}
	})

	// Query the book title
	title := ""
	te := doc.Find(".content h2")
	if te.Length() < 1 || te.Text() == "" {
		// Use title from book page.
		resp, err := d.client.Get(referer, "")
		if err != nil {
			log.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		title = doc.Find("title").Text()
	} else {
		title = te.Text()
	}

	return &BookMeta{
		ID:    bookID,
		Title: title,
		Links: links,
	}
}
