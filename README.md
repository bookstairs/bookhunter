# bookhunter

Downloading books from [talebook](https://github.com/talebook/talebook), [www.sanqiu.cc](https://www.sanqiu.cc/)
This is totally rewrite fork compare to its [original version](https://github.com/hellojukay/dl-talebook).

## Development

1. [Go Releaser](https://github.com/goreleaser/goreleaser) is used for releasing and local building
2. [golangci-lint](https://github.com/golangci/golangci-lint) is used for code style.
3. [goimports-reviser](https://github.com/incu6us/goimports-reviser) is used for sorting imports.

## Feature

### Download books from Talebook

1. Download from previous progress.
2. Register account on website.
3. Bypass the ratelimit from cloudflare.

### Download books from Sanqiu

1. Find all the books update from [www.sanqiu.cc](https://www.sanqiu.cc/).
2. Download required formats from 189 cloud drive.
3. Record the download progress for crontab jobs.

### Download books from Sobooks.

TODO

### Download books from my Telegram groups.

```
Usage:
  bookhunter telegram [flags]

Flags:
      --appHash string        The appHash for telegram.
      --appId int             The appID for telegram.
  -k, --channelId string      The channelId for telegram. You must set value. (default "https://t.me/haoshufenxiang")
  -d, --download string       The book directory you want to use, default would be current working directory. (default "/Users/zhaojianyun/Developer/project/github/bookhunter")
  -f, --format strings        The file formats you want to download. (default [EPUB,MOBI,PDF])
  -h, --help                  help for telegram
  -i, --initial int           The book id you want to start download. It should exceed 0. (default 1)
      --loadMessageSize int   The loadMessageSize is used to set the size of the number of messages obtained by requesting telegram API. 0 < loadMessageSize < 100 (default 20)
  -g, --progress string       The download progress file name you want to use, it would be saved under the download directory. (default "progress")
      --reLogin               force re-login.
  -n, --rename                Rename the book file by book ID.
  -r, --retry int             The max retry times for timeout download request. (default 5)
  -s, --sessionPath string    The session file for telegram. (default ".tg-session")
  -t, --thread int            The number of download threads. (default 1)
  -o, --timeout duration      The max pending time for download request. (default 10m0s)
```

Example command :
`bookhunter telegram --appId 12345 --appHash xxxxx -k https://t.me/MothLib`

How to get `appId` and `appHash` please refer to  [Creating your Telegram Application](https://core.telegram.org/api/obtaining_api_id)

## Install

### Homebrew (for macOS, Linux)

```shell
brew tap bibliolater/tap
brew install bookhunter
```

### Scope (for Windows)

```shell
scoop bucket add bibliolater https://github.com/bibliolater/scoop-bucket.git
scoop install bibliolater/bookhunter
```

### Manually

Download the latest release in [release page](https://github.com/bibliolater/bookhunter/releases). Choose related
tarball by your running environment.

## Usage

Execute `bookhunter -h` to see how to use this download tools.
