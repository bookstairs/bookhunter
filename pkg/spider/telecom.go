package spider

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

const telecomAPI = "https://api.noki.top/pan/cloud189/shareToDown"

var archiveFormats = []string{
	"ZIP",
	"TAR",
	"GZ",
	"RAR",
}

type TelecomResponse struct {
	ShareID   string `json:"shareId"`
	Directory struct {
		Count int `json:"count"`
		Files []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"fileList"`
	} `json:"fileListAO"`
}

// ResolveTelecom reclusive translate the telecom link to a direct download link.
func ResolveTelecom(client *Client, url, passcode string, formats ...string) ([]string, error) {
	results := map[string][]string{}
	if err := resolveTelegram(client, url, passcode, "", "", results); err != nil {
		return nil, err
	}

	// Add while list formats,
	formats = append(formats, archiveFormats...)

	var links []string
	for _, format := range formats {
		if s, ok := results[format]; ok {
			links = append(links, s...)
		}
	}

	return links, nil
}

// resolveTelegram will find all the downloadable links
func resolveTelegram(client *Client, url, passcode, shareId, fileId string, results map[string][]string) error {
	// Build queries
	queries := make([]*Query, 0, 4)
	queries = append(queries, &Query{Key: "url", Value: url})
	if passcode != "" {
		queries = append(queries, &Query{Key: "passCode", Value: passcode})
	}
	if shareId != "" {
		queries = append(queries, &Query{Key: "shareId", Value: shareId})
	}
	if passcode != "" {
		queries = append(queries, &Query{Key: "fileId", Value: fileId})
	}

	content, err := requestContent(client, queries...)
	if err != nil {
		return err
	}

	if strings.HasPrefix(content, "{") {
		// This is a directory based link
		response := &TelecomResponse{}
		decoder := json.NewDecoder(strings.NewReader(content))
		if err := decoder.Decode(response); err != nil {
			return err
		}

		shareID := response.ShareID
		for _, file := range response.Directory.Files {
			err := resolveTelegram(client, url, passcode, shareID, file.Id, results)
			if err != nil {
				return err
			}
		}

		return nil
	} else if strings.HasPrefix(content, "http") {
		// This is a download link. We won't filter the format.
		extension, _ := Extension(content)
		format := strings.ToUpper(extension)

		if links, ok := results[format]; ok {
			links = append(links, content)
			results[format] = links
		} else {
			results[format] = []string{content}
		}

		return nil
	} else if strings.Trim(content, " ") == "" {
		return nil
	} else {
		return errors.New(content)
	}
}

// requestContent will perform http request and return the response in string.
func requestContent(client *Client, queries ...*Query) (string, error) {
	// Perform the request.
	resp, err := client.Get(telecomAPI, "", queries...)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// Read the request content.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
