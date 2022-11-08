# bookhunter

Downloading books from [talebook](https://github.com/talebook/talebook), [www.sanqiu.cc](https://www.sanqiu.cc/)
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
  download    Download the book from talebook.
  register    Register account on talebook.

Flags:
  -h, --help   help for talebook

Use "bookhunter talebook [command] --help" for more information about a command.
```

### Download books from Sanqiu

1. Find all the books update from [www.sanqiu.cc](https://www.sanqiu.cc/).
2. Download required formats from 189 cloud drive.
2. Download required formats from aliyun drive.
3. Record the download progress for crontab jobs.

```shell
Usage:
  bookhunter sanqiu [flags]

Flags:
  -c, --cookie string         The cookie file name you want to use, it would be saved under the download directory. (default "cookies")
      --debug                 Enable debug mode
  -d, --download string       The book directory you want to use, default would be current working directory. (default "/Users/jianyun/GolandProjects/bookhunter")
  -f, --format strings        The file formats you want to download. (default [EPUB,MOBI,PDF])
  -h, --help                  help for sanqiu
  -i, --initial int           The book id you want to start download. It should exceed 0. (default 1)
  -g, --progress string       The download progress file name you want to use, it would be saved under the download directory. (default "progress")
      --refreshToken string   The refreshToken for AliYun Drive.
  -n, --rename                Rename the book file by book ID.
  -r, --retry int             The max retry times for timeout download request. (default 5)
  -t, --thread int            The number of download threads. (default 1)
  -o, --timeout duration      The max pending time for download request. (default 10m0s)
  -a, --user-agent string     Set User-Agent for download request. (default "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
  -w, --website string        The website for sanqiu. You don't need to override the default url. (default "https://www.sanqiu.cc")
```

### Download books from my Telegram groups.

1. Download all the valid book formats from a telegram channel.
2. Record the download progress for crontab jobs.
3. Support proxy from terminal environments.

```shell
Usage:
  bookhunter telegram [flags]

Flags:
      --appHash appHash      The appHash for telegram. How to get appHash please refer to https://core.telegram.org/api/obtaining_api_id.
      --appID appID          The appID for telegram. How to get appID please refer to https://core.telegram.org/api/obtaining_api_id.
  -k, --channelID string     The channelId for telegram.
  -c, --cookie string        The cookie file name you want to use, it would be saved under the download directory. (default "cookies")
  -d, --download string      The book directory you want to use, default would be current working directory. (default "/Users/Yufan")
  -f, --format strings       The file formats you want to download. (default [EPUB,MOBI,PDF])
  -h, --help                 help for telegram
  -i, --initial int          The book id you want to start download. It should exceed 0. (default 1)
  -b, --mobile string        The mobile number for your telegram account, default (+86).
  -g, --progress string      The download progress file name you want to use, it would be saved under the download directory. (default "progress")
      --refresh              Refresh the login session.
  -n, --rename               Rename the book file by book ID.
  -r, --retry int            The max retry times for timeout download request. (default 5)
  -s, --sessionPath string   The session file for telegram. (default "cookies")
  -t, --thread int           The number of download threads. (default 1)
  -o, --timeout duration     The max pending time for download request. (default 10m0s)
  -a, --user-agent string    Set User-Agent for download request. (default "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36")
```

Example command: `bookhunter telegram --appID 12345 --appHash xxxxx -k https://t.me/MothLib`

Please refer [Creating your Telegram Application](https://core.telegram.org/api/obtaining_api_id) to obtain your `appID` and `appHash`.

### Download books from Sobooks.

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
