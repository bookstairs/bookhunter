package sobooks

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestParseSobooksUrl(t *testing.T) {
	client := resty.New().SetBaseURL("https://sobooks.net")

	// unsupported 14320 11000
	// 20081 16899 16240
	id := int64(18021)
	resp, err := client.R().
		SetCookie(&http.Cookie{
			Name:   "mpcode",
			Value:  "844283",
			Path:   "/",
			Domain: "sobooks.net",
		}).
		SetPathParam("bookId", strconv.FormatInt(id, 10)).
		SetHeader("referer", client.BaseURL).
		Get("/books/{bookId}.html")
	if err != nil {
		t.Error(err)
	}
	_, links, _ := ParseLinks(resp.String(), id)
	fmt.Printf("%v", links)
}
