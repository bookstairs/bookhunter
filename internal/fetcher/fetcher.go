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
	defaultProgressFile = "progress.db"
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
	log.Infof("Successfully query the download content counts: %d", size)

	// Create download progress with ratelimit.
	if f.precessFile == "" {
		f.precessFile = defaultProgressFile
	}
	f.progress, err = progress.NewProgress(f.InitialBookID, size, f.RateLimit, filepath.Join(configPath, f.precessFile))
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

	var wait sync.WaitGroup
	for i := 0; i < f.Thread; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			f.startDownload()
		}()
	}
	wait.Wait()

	// Acquire the download errors.
	for e := range f.errs {
		return e
	}

	return nil
}

// startDownload will start a download thread.
func (f *commonFetcher) startDownload() {
thread:
	for {
		bookID := f.progress.AcquireBookID()
		if bookID == progress.NoBookToDownload {
			// Finish this thread.
			break thread
		}

		// Start download the given book ID.
		// The error will be sent to the channel.

		// Acquire the available file formats
		formats, err := f.service.formats(bookID)
		if err != nil {
			f.errs <- err
			break thread
		}
		log.Debugf("Book id %d formats: %v", bookID, formats)

		// Filter the formats.
		formats = f.filterFormats(formats)
		if len(formats) == 0 {
			log.Warnf("[%d/%d] No downloadable files found.", bookID, f.progress.Size())
		}

		// Download the file by formats one by one.
		for format, share := range formats {
			err := f.downloadFile(bookID, format, share)
			if err != nil {
				f.errs <- err
				break thread
			}
		}

		// Save the download progress
		err = f.progress.SaveBookID(bookID)
		if err != nil {
			f.errs <- err
			break thread
		}
	}
}

// downloadFile in a thread.
func (f *commonFetcher) downloadFile(bookID int64, format Format, share driver.Share) error {
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
	bar := log.NewProgressBar(bookID, f.progress.Size(), share.FileName, share.Size)
	defer func() { _ = bar.Close() }()

	// Write file content.
	if err := f.service.fetch(bookID, format, share, io.MultiWriter(writer, bar)); err != nil {
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
