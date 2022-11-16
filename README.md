# ⏬ bookhunter

Downloading books from [talebook](https://github.com/talebook/talebook), [www.sanqiu.mobi](https://www.sanqiu.mobi/)
, [www.tianlangbooks.com](www.tianlangbooks.com) and Telegram Channels. This is a totally rewritten fork compared to
its [original version](https://github.com/hellojukay/dl-talebook).

## 🚧 Development

1. [Go Releaser](https://github.com/goreleaser/goreleaser) is used for releasing and local building
2. [golangci-lint](https://github.com/golangci/golangci-lint) is used for code style.
3. [pre-commit](https://pre-commit.com/) is used for checking code before committing.

## 💾 Install

### 🍎 Homebrew (for macOS, Linux)

```shell
brew tap bookstairs/tap
brew install bookhunter
```

### 💻 Scope (for Windows)

```shell
scoop bucket add bookstairs https://github.com/bookstairs/scoop-bucket.git
scoop install bookstairs/bookhunter
```

### 🛠 Manually

Download the latest release in [release page](https://github.com/bookstairs/bookhunter/releases). Choose related tarball
by your running environment.

## 📚 Usage

<!--ts-->

* [Login Aliyundrive to get the refreshToken](#login-aliyundrive-to-get-the-refreshtoken)
* [Register account in Talebook](#register-account-in-talebook)
* [Download books from Talebook](#download-books-from-talebook)
* [Download books from Tianlang](#download-books-from-tianlang)
* [Download books from Sanqiu](#download-books-from-sanqiu)
* [Download books from Telegram groups.](#download-books-from-telegram-groups)

<!--te-->

### Login Aliyundrive to get the `refreshToken`

We would show a QR code at the first time. And cache the `refreshToken` after successfully login.

```shell
bookhunter aliyun
```

### Register account in Talebook

```text
Usage:
  bookhunter talebook register [flags]

Flags:
  -e, --email string      The talebook email
  -h, --help              help for register
  -p, --password string   The talebook password
  -u, --username string   The talebook username
  -w, --website string    The talebook link

Global Flags:
  -c, --config string       The config path for bookhunter
      --proxy string        The request proxy
  -a, --user-agent string   The request user-agent (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging
```

### Download books from Talebook

```text
Usage:
  bookhunter talebook download [flags]

Flags:
  -d, --download string   The book directory you want to use (default ".")
  -f, --format strings    The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help              help for download
  -i, --initial int       The book id you want to start download (default 1)
  -p, --password string   The talebook password
      --ratelimit int     The allowed requests per minutes (default 30)
  -r, --rename            Rename the book file by book id
  -t, --thread int        The number of download thead (default 1)
  -u, --username string   The talebook username
  -w, --website string    The talebook link

Global Flags:
  -c, --config string       The config path for bookhunter
      --proxy string        The request proxy
  -a, --user-agent string   The request user-agent (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging
```

### Download books from Tianlang

```text
Usage:
  bookhunter tianlang [flags]

Flags:
  -d, --download string          The book directory you want to use (default "/Users/Yufan/Developer/Go/bookstairs/bookhunter")
  -e, --extract                  Extract the archive file for filtering
  -f, --format strings           The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help                     help for tianlang
  -i, --initial int              The book id you want to start download (default 1)
      --ratelimit int            The allowed requests per minutes (default 30)
      --refreshToken string      Refresh token for aliyun drive
  -r, --rename                   Rename the book file by book id
      --secretKey string         The secret key for tianlang (default "359198")
      --source string            The source (aliyun, telecom, lanzou) to download book (default "telecom")
      --telecomPassword string   Telecom drive password
      --telecomUsername string   Telecom drive username
  -t, --thread int               The number of download thead (default 1)

Global Flags:
  -c, --config string       The config path for bookhunter
      --proxy string        The request proxy
  -a, --user-agent string   The request user-agent (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging
```

### Download books from Sanqiu

```text
Usage:
  bookhunter sanqiu [flags]

Flags:
  -d, --download string          The book directory you want to use (default ".")
  -e, --extract                  Extract the archive file for filtering
  -f, --format strings           The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help                     help for sanqiu
  -i, --initial int              The book id you want to start download (default 1)
      --ratelimit int            The allowed requests per minutes (default 30)
      --refreshToken string      Refresh token for aliyun drive
  -r, --rename                   Rename the book file by book id
      --source string            The source (aliyun, telecom, lanzou) to download book (default "telecom")
      --telecomPassword string   Telecom drive password
      --telecomUsername string   Telecom drive username
  -t, --thread int               The number of download thead (default 1)

Global Flags:
  -c, --config string       The config path for bookhunter
      --proxy string        The request proxy
  -a, --user-agent string   The request user-agent (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging
```

### Download books from Telegram groups.

Example command: `bookhunter telegram --appID ****** --appHash ****** -k https://t.me/MothLib`

Please refer [Creating your Telegram Application](https://core.telegram.org/api/obtaining_api_id) to obtain your `appID`
and `appHash`.

```text
Usage:
  bookhunter telegram [flags]

Flags:
      --appHash string     The app hash for telegram
      --appID int          The app id for telegram
  -k, --channelID string   The channel id for telegram
  -d, --download string    The book directory you want to use (default ".")
  -e, --extract            Extract the archive file for filtering
  -f, --format strings     The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help               help for telegram
  -i, --initial int        The book id you want to start download (default 1)
  -b, --mobile string      The mobile number, we will add +86 as default zone code
      --ratelimit int      The allowed requests per minutes (default 30)
      --refresh            Refresh the login session
  -r, --rename             Rename the book file by book id
  -t, --thread int         The number of download thead (default 1)

Global Flags:
  -c, --config string       The config path for bookhunter
      --proxy string        The request proxy
  -a, --user-agent string   The request user-agent (default "Mozilla/5.0 (X11; Linux i686; rv:13.0) Gecko/13.0 Firefox/13.0")
      --verbose             Print all the logs for debugging
```
