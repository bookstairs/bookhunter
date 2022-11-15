package sanqiu

import (
	"encoding/json"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bookstairs/bookhunter/internal/log"
)

// ParseAPIResponse will remove the unneeded error str in JSON response and try to parse it.
func ParseAPIResponse(resp *resty.Response, result any) error {
	// Remove error messages.
	str := resp.String()
	str = str[strings.LastIndex(str, "\n")+1:]

	log.Debug("Response: ", str)

	decoder := json.NewDecoder(strings.NewReader(str))

	return decoder.Decode(result)
}

// BookResp is the response for /wp-json/wp/v2/posts
type BookResp struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	DateGmt string `json:"date_gmt"`
	GUID    struct {
		Rendered string `json:"rendered"`
	} `json:"guid"`
	Modified    string `json:"modified"`
	ModifiedGmt string `json:"modified_gmt"`
	Slug        string `json:"slug"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Link        string `json:"link"`
	Title       struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered  string `json:"rendered"`
		Protected bool   `json:"protected"`
	} `json:"content"`
	Excerpt struct {
		Rendered  string `json:"rendered"`
		Protected bool   `json:"protected"`
	} `json:"excerpt"`
	Author        int           `json:"author"`
	FeaturedMedia int           `json:"featured_media"`
	CommentStatus string        `json:"comment_status"`
	PingStatus    string        `json:"ping_status"`
	Sticky        bool          `json:"sticky"`
	Template      string        `json:"template"`
	Format        string        `json:"format"`
	Meta          []interface{} `json:"meta"`
	Categories    []int         `json:"categories"`
	Tags          []int         `json:"tags"`
	Links         struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
		Collection []struct {
			Href string `json:"href"`
		} `json:"collection"`
		About []struct {
			Href string `json:"href"`
		} `json:"about"`
		Author []struct {
			Embeddable bool   `json:"embeddable"`
			Href       string `json:"href"`
		} `json:"author"`
		Replies []struct {
			Embeddable bool   `json:"embeddable"`
			Href       string `json:"href"`
		} `json:"replies"`
		VersionHistory []struct {
			Count int    `json:"count"`
			Href  string `json:"href"`
		} `json:"version-history"`
		WpAttachment []struct {
			Href string `json:"href"`
		} `json:"wp:attachment"`
		WpTerm []struct {
			Taxonomy   string `json:"taxonomy"`
			Embeddable bool   `json:"embeddable"`
			Href       string `json:"href"`
		} `json:"wp:term"`
		Curies []struct {
			Name      string `json:"name"`
			Href      string `json:"href"`
			Templated bool   `json:"templated"`
		} `json:"curies"`
	} `json:"_links"`
}
