# bookhunter

Downloading books from [talebook](https://github.com/talebook/talebook), [www.sanqiu.cc](https://www.sanqiu.cc/)
This is totally rewrite fork compare to its [original version](https://github.com/hellojukay/bookhunter).

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

TODO

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
