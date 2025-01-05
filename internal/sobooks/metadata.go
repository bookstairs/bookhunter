package sobooks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/log"
)

type LinkType string

var (
	// driveNamings is the chinese name mapping of the drive's provider.
	driveNamings = map[driver.Source]string{
		driver.ALIYUN:  "阿里",
		driver.CTFILE:  "城通",
		driver.BAIDU:   "百度",
		driver.QUARK:   "夸克",
		driver.LANZOU:  "蓝奏",
		driver.TELECOM: "天翼",
		driver.DIRECT:  "备份",
	}
	dateRe = regexp.MustCompile(`(?m)(\d{4})-(\d{2})-(\d{2})`)
	linkRe = regexp.MustCompile(`(?m)https://sobooks\.cc/go\.html\?url=(.*?)"(.*?[：:]\s?(\w+))?`)
)

type BookLink struct {
	URL  string // The url to access the download page.
	Code string // The passcode for querying the file content.
}

// ParseLinks will find all the available link in the different driver.
func ParseLinks(content string, id int64) (title string, links map[driver.Source]BookLink, err error) {
	// Find all the links.
	links = map[driver.Source]BookLink{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", links, err
	}

	dateText := doc.Find("div.bookinfo > ul > li:nth-child(5)").Text()

	if !dateRe.MatchString(dateText) {
		log.Fatal("not found book date ", id)
		return "", map[driver.Source]BookLink{}, fmt.Errorf("not found book date %v", id)
	}
	submatch := dateRe.FindStringSubmatch(dateText)
	year := submatch[1]
	month := submatch[2]
	titleDom := doc.Find(".article-title>a")
	title = titleDom.Text()

	// default link
	links[driver.DIRECT] = BookLink{URL: fmt.Sprintf("https://sobooks.cloud/%s/%s/%d.epub", year, month, id)}

	html, err := doc.Find(".e-secret").Html()
	if err != nil {
		return "", links, err
	}
	split := strings.Split(html, "<br/>")

	for _, s := range split {
		for linkType, name := range driveNamings {
			if strings.Contains(s, name) {
				match := linkRe.FindStringSubmatch(s)
				if len(match) > 2 {
					links[linkType] = BookLink{
						URL:  match[1],
						Code: match[3],
					}
				}
				break
			}
		}
	}
	return title, links, nil
}
