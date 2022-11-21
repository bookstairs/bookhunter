package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var encoding = simplifiedchinese.GB18030

// decompress - extract zip file.
func (p *writer) decompress() error {
	r, err := zip.OpenReader(p.filePath())
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	_ = os.MkdirAll(p.download, 0o755)

	for _, f := range r.File {
		if err := p.extractAndWriteFile(f); err != nil {
			return err
		}
	}

	return nil
}

// Closure to address file descriptors issue with all the deferred Close() methods.
func (p *writer) extractAndWriteFile(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	filename := encodingFilename(f.Name)
	ext, ok := Extension(filename)
	if !ok || !p.formats[ext] {
		// No need to extract this file. Skip.
		return nil
	}

	path, err := sanitizeArchivePath(p.download, filename)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(path, filepath.Clean(p.download)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", path)
	}

	if f.FileInfo().IsDir() {
		_ = os.MkdirAll(path, 0o755)
	} else {
		mode := f.FileHeader.Mode()
		if mode&os.ModeType == os.ModeSymlink {
			data, err := io.ReadAll(rc)
			if err != nil {
				return err
			}
			_ = writeSymbolicLink(path, string(data))
		} else {
			_ = os.MkdirAll(filepath.Dir(path), 0o755)
			_ = os.Remove(path)
			outFile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer func() {
				if err := outFile.Close(); err != nil {
					panic(err)
				}
			}()

			// G110: Potential DoS vulnerability via decompression bomb.
			for {
				_, err := io.CopyN(outFile, rc, bytes.MinRead)
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
			}
		}
	}

	return nil
}

func writeSymbolicLink(filePath string, targetPath string) error {
	err := os.MkdirAll(filepath.Dir(filePath), 0o755)
	if err != nil {
		return err
	}

	err = os.Symlink(targetPath, filePath)
	if err != nil {
		return err
	}

	return nil
}

// sanitizeArchivePath sanitize archive file pathing from "G305: Zip Slip vulnerability"
func sanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", t)
}

// encodingFilename will convert the GBK into UTF-8 and escape invalid characters.
func encodingFilename(name string) string {
	i := bytes.NewReader([]byte(name))
	decoder := transform.NewReader(i, encoding.NewDecoder())
	content, err := io.ReadAll(decoder)
	if err != nil {
		// Fallback to default UTF-8 encoding
		return escape(name)
	} else {
		return escape(string(content))
	}
}
