package fetcher

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/yi-ge/unzip"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/naming"
	"github.com/bookstairs/bookhunter/internal/progress"
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
	f.progress, err = progress.NewProgress(f.InitialBookID, size, f.RateLimit, filepath.Join(configPath, "process.db"))
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
		formats, err := f.service.formats()
		if err != nil {
			f.finishDownload(err)
			break
		}

		// Filter the formats.
		formats = f.filterFormats(formats)

		// Download the file by formats one by one.
		for _, format := range formats {
			err := f.downloadFile(bookID, format)
			if err != nil {
				f.finishDownload(err)
				break
			}
		}
	}
}

// downloadFile in a thread.
func (f *commonFetcher) downloadFile(bookID int64, format Format) error {
	file, err := f.service.fetch(bookID, format)
	if err != nil {
		return err
	}
	defer func() { _ = file.content.Close() }()

	// Rename if it was required.
	prefix := strconv.FormatInt(bookID, 10)
	if f.Rename {
		file.name = prefix + "." + string(format)
	} else {
		file.name = prefix + "_" + file.name
	}

	// Escape the file name for avoiding the illegal characters.
	// Ref: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
	file.name = naming.EscapeFilename(file.name)

	// Generate the file path.
	path := filepath.Join(f.DownloadPath, file.name)

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
	bar := log.NewProgressBar(bookID, f.progress.Size(), file.name, file.size)
	defer func() { _ = bar.Close() }()

	// Write file content
	_, err = io.Copy(io.MultiWriter(writer, bar), file.content)
	if err != nil {
		return err
	}

	// Extract the archives.
	if format.Archive() && f.Extract {
		u := unzip.New(path, f.DownloadPath)
		return u.Extract()
	}

	return nil
}

// finishDownload will exist the download thread.
func (f *commonFetcher) finishDownload(err error) {
	if err != nil {
		f.errs <- err
	}
	f.wait.Done()
}

// filterFormats will find the valid formats by user configure.
func (f *commonFetcher) filterFormats(formats []Format) []Format {
	var fs []Format
	for _, format := range formats {
		for _, vf := range f.Formats {
			if format == vf {
				fs = append(fs, format)
				break
			}
		}
	}
	return fs
}
