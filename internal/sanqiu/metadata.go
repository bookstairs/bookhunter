package sanqiu

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bookstairs/bookhunter/internal/driver"
)

type LinkType string

var (
	// driveNamings is the chinese name mapping of the drive's provider.
	driveNamings = map[driver.Source]string{
		driver.ALIYUN:  "阿里",
		driver.LANZOU:  "蓝奏",
		driver.TELECOM: "天翼",
	}
	passcodeRe = regexp.MustCompile(".*?([a-zA-Z0-9]+).*?")
)

type BookLink struct {
	URL  string // The url to access the download page.
	Code string // The passcode for querying the file content.
}

// DownloadLinks will find all the available link in the different driver.
func DownloadLinks(content string) (map[driver.Source]BookLink, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	// Find all the links.
	links := map[driver.Source]BookLink{}
	doc.Find(".downfile a").Each(func(i int, selection *goquery.Selection) {
		driveName := selection.Text()
		href, exists := selection.Attr("href")
		if exists {
			for linkType, name := range driveNamings {
				if strings.Contains(driveName, name) {
					links[linkType] = BookLink{URL: href}
					break
				}
			}
		}
	})

	if len(links) == 0 {
		return map[driver.Source]BookLink{}, nil
	}

	// Find all the passcodes.
	doc.Find(".plus_l li").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		for linkType, link := range links {
			name := driveNamings[linkType]
			if strings.Contains(text, name) {
				match := passcodeRe.FindStringSubmatch(text)
				if len(match) == 2 {
					links[linkType] = BookLink{
						URL:  link.URL,
						Code: match[1],
					}
				}
			}
		}
	})

	return links, nil
}
