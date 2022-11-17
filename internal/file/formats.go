package file

import (
	"net/url"
	"strings"
	"unicode"
)

type Format string // The supported file extension.

const (
	EPUB Format = "epub"
	MOBI Format = "mobi"
	AZW  Format = "azw"
	AZW3 Format = "azw3"
	PDF  Format = "pdf"
	ZIP  Format = "zip"
)

// The Archive will return if this format is an archive.
func (f Format) Archive() bool {
	return f == ZIP
}

func inArchive(filename string) bool {
	ext, ok := Extension(filename)
	if !ok {
		return false
	}
	return ext.Archive()
}

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// LinkExtension the file extension from the link.
func LinkExtension(link string) (Format, bool) {
	u, err := url.Parse(link)
	if err != nil {
		return "", false
	}
	return Extension(u.Path)
}

func Extension(filename string) (Format, bool) {
	start := strings.LastIndex(filename, ".") + 1
	ext := filename[start:]

	if isLetter(ext) {
		return Format(strings.ToLower(ext)), true
	}
	return "", false
}
