package file

import (
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

type Writer interface {
	io.Writer
	io.Closer
	SetSize(int64)
}

type progressWriter struct {
	bar  *progressbar.ProgressBar
	file *os.File
}

func (p *progressWriter) Close() error {
	_ = p.bar.Close()
	return p.file.Close()
}

func (p *progressWriter) Write(b []byte) (n int, err error) {
	_, _ = p.bar.Write(b)
	return p.file.Write(b)
}

func (p *progressWriter) SetSize(i int64) {
	p.bar.ChangeMax64(i)
}

func NewWriter(path string, bar *progressbar.ProgressBar) (Writer, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &progressWriter{file: file, bar: bar}, nil
}
