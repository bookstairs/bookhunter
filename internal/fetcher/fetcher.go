package fetcher

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/bookstairs/bookhunter/internal/driver"
	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/progress"
	"github.com/bookstairs/bookhunter/internal/unzip"
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
	creator  file.Creator
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

	// Create the file creator.
	f.creator = file.NewCreator(f.Rename, f.DownloadPath)

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
	select {
	case err := <-f.errs:
		return err
	default:
		log.Debug("All the fetch thread have been finished.")
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
			log.Debugf("No book to download in [%s] service.", f.Category)
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
		log.Debugf("Book id %d formats: %v.", bookID, formats)

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
	// Create the file writer.
	writer, err := f.creator.NewWriter(bookID, f.progress.Size(), share.FileName, string(format), share.Size)
	if err != nil {
		return err
	}

	// Remove the archive file.
	if format.Archive() && f.Extract {
		defer func() { _ = os.Remove(writer.Path()) }()
	}
	defer func() { _ = writer.Close() }()

	// Write file content.
	if err := f.service.fetch(bookID, format, share, writer); err != nil {
		return err
	}

	// Extract the archives. We only support the zip file now.
	if format.Archive() && f.Extract {
		path := writer.Path()
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
