package file

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/schollz/progressbar/v3"

	"github.com/bookstairs/bookhunter/internal/log"
)

type Creator interface {
	NewWriter(id, total int64, name string, format Format, size int64) (Writer, error)
}

func NewCreator(rename bool, downloadPath string, formats []Format, extract bool) Creator {
	fs := make(map[Format]bool)
	for _, format := range formats {
		fs[format] = true
	}

	return &creator{rename: rename, downloadPath: downloadPath, formats: fs, extract: extract}
}

type creator struct {
	rename       bool
	extract      bool
	formats      map[Format]bool
	downloadPath string
}

func (c *creator) NewWriter(id, total int64, name string, format Format, size int64) (Writer, error) {
	// Rename if it was required.
	filename := strconv.FormatInt(id, 10)
	if c.rename {
		filename = filename + "." + string(format)
	} else {
		filename = filename + "_" + name
	}

	// Escape the file name for avoiding the illegal characters.
	// Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	filename = escape(filename)

	// Generate the file path.
	path := filepath.Join(c.downloadPath, filename)

	// Remove the exist file.
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}

	// Add download progress, no need to close.
	bar := log.NewProgressBar(id, total, name, size)

	// Create file io. and remember to close it manually.
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &writer{
		file:     file,
		name:     filename,
		download: c.downloadPath,
		extract:  c.extract && format.Archive(),
		formats:  c.formats,
		bar:      bar,
	}, nil
}

type Writer interface {
	io.Writer
	io.Closer
	SetSize(int64)
}

type writer struct {
	file     *os.File
	name     string
	download string
	formats  map[Format]bool
	extract  bool
	bar      *progressbar.ProgressBar
}

func (p *writer) Close() error {
	_ = p.bar.Close()
	err := p.file.Close()
	if err != nil {
		return err
	}

	// Extract the file if user enabled this.
	if p.extract {
		if err := p.decompress(); err != nil {
			log.Fatal(err)
			return nil
		}

		// Remove the compress files.
		_ = os.Remove(p.filePath())
	}

	return err
}

func (p *writer) filePath() string {
	return filepath.Join(p.download, p.name)
}

func (p *writer) Write(b []byte) (n int, err error) {
	_, _ = p.bar.Write(b)
	return p.file.Write(b)
}

func (p *writer) SetSize(i int64) {
	p.bar.ChangeMax64(i)
}
