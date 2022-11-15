package fetcher

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/yi-ge/unzip"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/naming"
	"github.com/bookstairs/bookhunter/internal/progress"
)

const (
	progressFile = "progress.db"
)

// Fetcher exposes the download method to the command line.
type Fetcher interface {
	// Download the books from the given service.
	Download() error
}

// commonFetcher is the basic common download service the multiple thread support.
type commonFetcher struct {
	*Config
	service  service
	progress progress.Progress
	wait     sync.WaitGroup
	errs     chan error
}

// Download the books from the given service.
func (f *commonFetcher) Download() error {
	// Create the config path.
	configPath, err := f.ConfigPath()
	if err != nil {
		return err
	}

	// Query the total download amount from the given service.
	size, err := f.service.size()
	if err != nil {
		return err
	}

	// Create download progress with ratelimit.
	f.progress, err = progress.NewProgress(f.InitialBookID, size, f.RateLimit, filepath.Join(configPath, progressFile))
	if err != nil {
		return err
	}

	// Create the download directory if it's not existed.
	err = os.MkdirAll(f.DownloadPath, 0755)
	if err != nil {
		return err
	}

	// Create the download thread and save the files.
	f.errs = make(chan error, f.Thread)
	defer close(f.errs)

	for i := 0; i < f.Thread; i++ {
		f.wait.Add(1)
		go f.startDownload()
	}

	f.wait.Wait()

	// Acquire the download error.
	for e := range f.errs {
		if e != nil {
			return e
		}
	}

	return nil
}

// startDownload will start a download thread.
func (f *commonFetcher) startDownload() {
	for {
		bookID := f.progress.AcquireBookID()
		if bookID == progress.NoBookToDownload {
			// Finish this thread.
			f.finishDownload(nil)
			break
		}

		// Start download the given book ID.
		// The error will be sent to the channel.

		// Acquire the available file formats
		formats, err := f.service.formats(bookID)
		if err != nil {
			f.finishDownload(err)
			break
		}

		// Filter the formats.
		formats = f.filterFormats(formats)

		// Download the file by formats one by one.
		for format, share := range formats {
			err := f.downloadFile(bookID, format, share)
			if err != nil {
				f.finishDownload(err)
				break
			}
		}

		// Save the download progress
		err = f.progress.SaveBookID(bookID)
		if err != nil {
			f.finishDownload(err)
			break
		}
	}
}

// downloadFile in a thread.
func (f *commonFetcher) downloadFile(bookID int64, format Format, share driver.Share) error {
	file, err := f.service.fetch(bookID, format, share)
	if err != nil {
		return err
	}
	defer func() { _ = file.content.Close() }()

	// Rename if it was required.
	prefix := strconv.FormatInt(bookID, 10)
	if f.Rename {
		share.FileName = prefix + "." + string(format)
	} else {
		share.FileName = prefix + "_" + share.FileName
	}

	// Escape the file name for avoiding the illegal characters.
	// Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	share.FileName = naming.EscapeFilename(share.FileName)

	// Generate the file path.
	path := filepath.Join(f.DownloadPath, share.FileName)

	// Remove the exist file.
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	// Create the file writer.
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = writer.Close() }()

	// Add download progress.
	bar := log.NewProgressBar(bookID, f.progress.Size(), share.FileName, file.size)
	defer func() { _ = bar.Close() }()

	// Write file content
	_, err = io.Copy(io.MultiWriter(writer, bar), file.content)
	if err != nil {
		return err
	}

	// Extract the archives. We only support the zip file now.
	if format.Archive() && f.Extract {
		u := unzip.New(path, f.DownloadPath)
		return u.Extract()
	}

	return nil
}

// filterFormats will find the valid formats by user configure.
func (f *commonFetcher) filterFormats(formats map[Format]driver.Share) map[Format]driver.Share {
	fs := make(map[Format]driver.Share)
	for format, share := range formats {
		for _, vf := range f.Formats {
			if format == vf {
				fs[format] = share
				break
			}
		}
	}
	return fs
}

// finishDownload will exist the download thread.
func (f *commonFetcher) finishDownload(err error) {
	if err != nil {
		f.errs <- err
	}
	f.wait.Done()
}