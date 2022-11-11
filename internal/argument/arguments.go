package argument

import "github.com/bookstairs/bookhunter/internal/fetcher"

var (
	// Flags for talebook registering.

	Username = ""
	Password = ""
	Email    = ""

	// Common flags.

	Website    = ""
	UserAgent  = ""
	Proxy      = ""
	ConfigRoot = ""

	// Common download flags.

	Formats       = fetcher.NormalizeFormats(fetcher.EPUB, fetcher.AZW3, fetcher.MOBI, fetcher.PDF, fetcher.ZIP)
	Extract       = false
	DownloadPath  = ""
	InitialBookID = 1
	Rename        = false
	Thread        = 1

	// Drive ISP configurations.

	RefreshToken = ""

	// Telegram configurations.

	ChannelID = ""
	Mobile    = ""
	ReLogin   = false
	AppID     = int64(0)
	AppHash   = ""
)
