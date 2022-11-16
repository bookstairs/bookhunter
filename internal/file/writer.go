package file

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"unicode/utf8"

	"github.com/schollz/progressbar/v3"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/naming"
)

type Creator interface {
	NewWriter(id, total int64, name string, format string, size int64) (Writer, error)
}

func NewCreator(rename bool, downloadPath string) Creator {
	return &creator{rename: rename, downloadPath: downloadPath}
}

type creator struct {
	rename       bool
	downloadPath string
}

func (c *creator) NewWriter(id, total int64, name string, format string, size int64) (Writer, error) {
	// Rename if it was required.
	filename := strconv.FormatInt(id, 10)
	if c.rename {
		filename = filename + "." + format
	} else {
		filename = filename + "_" + name
	}

	// Escape the file name for avoiding the illegal characters.
	// Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	filename = naming.EscapeFilename(filename)

	// Generate the file path.
	path := filepath.Join(c.downloadPath, filename)

	// Remove the exist file.
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}

	// Add download progress, no need to close.
	display := name
	if utf8.RuneCountInString(display) > 30 {
		// Trim the display size for better printing.
		display = string([]rune(display)[:30]) + "..."
	}
	bar := log.NewProgressBar(id, total, display, size)

	// Create file io. and remember to close it manually.
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &writer{file: file, bar: bar, path: path}, nil
}

type Writer interface {
	io.Writer
	io.Closer
	Path() string
	SetSize(int64)
}

type writer struct {
	bar  *progressbar.ProgressBar
	file *os.File
	path string
}

func (p *writer) Path() string {
	return p.path
}

func (p *writer) Close() error {
	_ = p.bar.Close()
	return p.file.Close()
}

func (p *writer) Write(b []byte) (n int, err error) {
	_, _ = p.bar.Write(b)
	return p.file.Write(b)
}

func (p *writer) SetSize(i int64) {
	p.bar.ChangeMax64(i)
}
