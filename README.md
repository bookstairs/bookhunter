# bookhunter

Downloading books from [talebook](https://github.com/talebook/talebook), [www.sanqiu.mobi](https://www.sanqiu.mobi/)
This is totally rewrite fork compare to its [original version](https://github.com/hellojukay/dl-talebook).

## Development

1. [Go Releaser](https://github.com/goreleaser/goreleaser) is used for releasing and local building
2. [golangci-lint](https://github.com/golangci/golangci-lint) is used for code style.
3. [pre-commit](https://pre-commit.com/) is used for checking code before committing.

## Feature

### Download books from Talebook

1. Download from previous progress.
2. Register account on website.
3. Bypass the ratelimit from cloudflare.

```shell
Usage:
  bookhunter talebook [command]

Available Commands:
  download    Download the books from talebook.
  register    Register account on talebook.

Flags:
  -h, --help   help for talebook

Global Flags:
  -c, --config string       The config path for bookhunter.
      --proxy string        The request proxy.
  -a, --user-agent string   The request user-agent. (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging.
```

```shell
Usage:
  bookhunter talebook download [flags]

Flags:
  -d, --download string   The book directory you want to use, default would be current working directory. (default ".")
  -f, --format strings    The file formats you want to download. (default [epub,azw3,mobi,pdf,zip])
  -h, --help              help for download
  -i, --initial int       The book id you want to start download. It should exceed 0. (default 1)
  -p, --password string   The account password.
      --ratelimit int     The request per minutes. (default 30)
  -r, --rename            Rename the book file by book ID.
  -t, --thread int        The number of concurrent download thead. (default 1)
  -u, --username string   The account login name.
  -w, --website string    The talebook website.

Global Flags:
  -c, --config string       The config path for bookhunter.
      --proxy string        The request proxy.
  -a, --user-agent string   The request user-agent. (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging.
```

### Download books from Sanqiu

1. Find all the books update from [www.sanqiu.mobi](https://www.sanqiu.mobi/).
2. Download required formats from 189 cloud drive.
3. Download required formats from aliyun drive.
4. Record the download progress for crontab jobs.

```shell
Usage:
  bookhunter sanqiu [flags]

Flags:
  -d, --download string          The book directory you want to use, default would be current working directory. (default ".")
  -e, --extract                  Extract the archive file for filtering.
  -f, --format strings           The file formats you want to download. (default [epub,azw3,mobi,pdf,zip])
  -h, --help                     help for sanqiu
  -i, --initial int              The book id you want to start download. It should exceed 0. (default 1)
      --ratelimit int            The request per minutes. (default 30)
      --refreshToken string      We would try to download from the aliyun drive if you provide this token.
  -r, --rename                   Rename the book file by book ID.
      --telecomPassword string   Used to download file from telecom drive
      --telecomUsername string   Used to download file from telecom drive
  -t, --thread int               The number of concurrent download thead. (default 1)

Global Flags:
  -c, --config string       The config path for bookhunter.
      --proxy string        The request proxy.
  -a, --user-agent string   The request user-agent. (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging.
```

### Download books from Telegram groups.

1. Download all the valid book formats from a telegram channel.
2. Record the download progress for crontab jobs.
3. Support proxy from terminal environments.

```shell
Usage:
  bookhunter telegram [flags]

Flags:
      --appHash string     The appHash for telegram. Refer to https://core.telegram.org/api/obtaining_api_id to create your own appHash
      --appID int          The appID for telegram. Refer https://core.telegram.org/api/obtaining_api_id to create your own appID
  -k, --channelID string   The channelId for telegram.
  -d, --download string    The book directory you want to use, default would be current working directory. (default ".")
  -e, --extract            Extract the archive file for filtering.
  -f, --format strings     The file formats you want to download. (default [epub,azw3,mobi,pdf,zip])
  -h, --help               help for telegram
  -i, --initial int        The book id you want to start download. It should exceed 0. (default 1)
      --ratelimit int      The request per minutes. (default 30)
      --refresh            Refresh the login session.
  -r, --rename             Rename the book file by book ID.
  -t, --thread int         The number of concurrent download thead. (default 1)

Global Flags:
  -c, --config string       The config path for bookhunter.
      --proxy string        The request proxy.
  -a, --user-agent string   The request user-agent. (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging.
```

Example command: `bookhunter telegram --appID ****** --appHash ****** -k https://t.me/MothLib`

Please refer [Creating your Telegram Application](https://core.telegram.org/api/obtaining_api_id) to obtain your `appID` and `appHash`.

### Download books from Sobooks.

TODO

### Download books from Tianlang Books.

TODO

## Install

### Homebrew (for macOS, Linux)

```shell
brew tap bookstairs/tap
brew install bookhunter
```

### Scope (for Windows)

```shell
scoop bucket add bookstairs https://github.com/bookstairs/scoop-bucket.git
scoop install bookstairs/bookhunter
```

### Manually

Download the latest release in [release page](https://github.com/bookstairs/bookhunter/releases). Choose related
tarball by your running environment.

## Usage

Execute `bookhunter -h` to see how to use this download tools.
