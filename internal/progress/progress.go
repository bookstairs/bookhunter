package progress

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bits-and-blooms/bitset"
	"go.uber.org/ratelimit"
)

const NoBookToDownload = -1

var (
	ErrStartBookID       = errors.New("the start book id should start from 1")
	ErrStartAndEndBookID = errors.New("start book id should below the available book id")
	ErrStorageFile       = errors.New("couldn't create file for storing download process")
)

type Progress interface {
	// TakeRateLimit would wait until the rate limit is available.
	TakeRateLimit()

	// AcquireBookID would find the book id from the assign array.
	AcquireBookID() int64

	// SaveBookID would save the download progress.
	SaveBookID(bookID int64) error

	// Finished would tell the called whether all the books have downloaded.
	Finished() bool

	// Size would return the book size.
	Size() int64
}

// bitProgress is a bit-based implementation with file persistence.
type bitProgress struct {
	limit    ratelimit.Limiter // The ratelimit for acquiring a book ID.
	progress *bitset.BitSet    // progress is used for file Progress.
	assigned *bitset.BitSet    // the assign status, memory based.
	lock     *sync.Mutex       // lock is used for concurrent request.
	file     *os.File          // The Progress file path for download progress.
}

// NewProgress Create a storage for save the download progress.
func NewProgress(start, size int64, rate int, path string) (Progress, error) {
	if start < 1 {
		return nil, ErrStartBookID
	}
	if start > size {
		return nil, ErrStartAndEndBookID
	}

	var progress *bitset.BitSet
	var file *os.File

	startIndex := func(set *bitset.BitSet) {
		for i := uint(0); i < uint(start-1); i++ {
			set.Set(i)
		}
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		// Create Progress related file
		if file, err = os.Create(path); err != nil {
			return nil, ErrStorageFile
		} else {
			// Enrich file content.
			progress = bitset.New(uint(size))
			startIndex(progress)

			if err := saveStorage(file, progress); err != nil {
				return nil, err
			}
		}
	} else {
		// Load Progress from file.
		if file, err = os.OpenFile(path, os.O_RDWR, 0o644); err != nil {
			return nil, err
		}
		if progress, err = loadStorage(file); err != nil {
			return nil, err
		}

		// Recalculate start index.
		startIndex(progress)

		// Support book update, increase the progress size.
		if progress.Len() < uint(size) {
			p := bitset.New(uint(size))
			progress.Copy(p)
			progress = p
		}
	}

	assigned := bitset.New(progress.Len())
	progress.Copy(assigned)

	// Create ratelimit
	limit := ratelimit.New(rate, ratelimit.Per(time.Minute))

	return &bitProgress{
		limit:    limit,
		progress: progress,
		assigned: assigned,
		lock:     new(sync.Mutex),
		file:     file,
	}, nil
}

func saveStorage(file *os.File, progress *bitset.BitSet) error {
	bytes, err := progress.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = file.WriteAt(bytes, 0)
	return err
}

func loadStorage(file *os.File) (*bitset.BitSet, error) {
	set := new(bitset.BitSet)
	if _, err := set.ReadFrom(file); err != nil {
		return nil, err
	}

	return set, nil
}

// TakeRateLimit block until the rate meets the given config.
func (storage *bitProgress) TakeRateLimit() {
	storage.limit.Take()
}

// AcquireBookID would find the book id from the assign array.
func (storage *bitProgress) AcquireBookID() int64 {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for i := uint(0); i < storage.assigned.Len(); i++ {
		if !storage.assigned.Test(i) {
			storage.assigned.Set(i)
			return int64(i + 1)
		}
	}

	return NoBookToDownload
}

// SaveBookID would save the download progress.
func (storage *bitProgress) SaveBookID(bookID int64) error {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	if bookID > int64(storage.progress.Len()) {
		return fmt.Errorf("invalid book id: %d", bookID)
	}

	i := uint(bookID - 1)
	storage.assigned.Set(i)
	storage.progress.Set(i)

	_ = saveStorage(storage.file, storage.progress)

	return nil
}

// Finished would tell the called whether all the books have downloaded.
func (storage *bitProgress) Finished() bool {
	return storage.progress.Count() == storage.progress.Len()
}

// Size would return the book size.
func (storage *bitProgress) Size() int64 {
	return int64(storage.progress.Len())
}
