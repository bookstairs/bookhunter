package fetcher

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/wordpress"
)

func newSanqiuService(config *Config) (service, error) {
	return newWordpressService(config, func(c *client.Client, id int64) (map[driver.Source]wordpress.ShareLink, error) {
		resp, err := c.R().
			SetQueryParam("id", strconv.FormatInt(id, 10)).
			Get("/download.php")
		if err != nil {
			return nil, err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))
		if err != nil {
			return nil, err
		}

		// Find all the links.
		links := map[driver.Source]wordpress.ShareLink{}
		doc.Find(".downfile a").Each(func(i int, selection *goquery.Selection) {
			driveName := selection.Text()
			href, exists := selection.Attr("href")
			if exists {
				for linkType, name := range driveNamings {
					if strings.Contains(driveName, name) {
						links[linkType] = wordpress.ShareLink{URL: href}
						break
					}
				}
			}
		})

		if len(links) == 0 {
			return map[driver.Source]wordpress.ShareLink{}, nil
		}

		// Find all the passcodes.
		doc.Find(".plus_l li").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			for linkType, link := range links {
				name := driveNamings[linkType]
				if strings.Contains(text, name) {
					match := sanqiuPasscodeRe.FindStringSubmatch(text)
					if len(match) == 2 {
						links[linkType] = wordpress.ShareLink{
							URL:  link.URL,
							Code: match[1],
						}
					}
				}
			}
		})

		return links, nil
	})
}
