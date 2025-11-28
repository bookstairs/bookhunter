# ⏬ bookhunter

[![LICENSE](https://img.shields.io/github/license/bookstairs/bookhunter)](https://github.com/bookstairs/bookhunter/blob/main/LICENSE)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bookstairs/bookhunter)](https://goreportcard.com/report/github.com/bookstairs/bookhunter)
![](https://img.shields.io/github/stars/bookstairs/bookhunter.svg)
![](https://img.shields.io/github/forks/bookstairs/bookhunter.svg)
![Release](https://github.com/bookstairs/bookhunter/workflows/release/badge.svg)

Downloading books from [talebook](https://github.com/talebook/talebook),
[SoBooks](https://sobooks.cc)
[中小学教材](https://basic.smartedu.cn/tchMaterial)
and Telegram Channels. This is a totally
rewritten fork compared to its [original version](https://github.com/hellojukay/dl-talebook).

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

### ⚡ Quick Start: Task Management (Recommended)

**NEW**: Manage multiple download tasks with a single configuration file!

```bash
# Create example configuration
bookhunter run  # Will prompt to create config if not exists

# List all configured tasks
bookhunter task list

# Run all enabled tasks
bookhunter run

# Run specific tasks
bookhunter run mothlib-telegram sobooks-download

# Interactive task selection
bookhunter run --interactive

# Enable/disable tasks
bookhunter task enable mothlib-telegram
bookhunter task disable sobooks-download
```

See [Task Management Guide](#task-management) for details.

### 📖 Supported Sources

| Website                                          | Address                                | Direct Download | [Aliyun](https://www.aliyundrive.com/) | [Lanzou](https://www.lanzou.com/) | [Telecom](https://cloud.189.cn/) |
|--------------------------------------------------|----------------------------------------|-----------------|----------------------------------------|-----------------------------------|----------------------------------|
| [智慧教育平台](#download-textbooks-for-kids)           | <https://basic.smartedu.cn/tchMaterial>   | ✅               | ❌                                      | ❌                                 | ❌                                |
| [Talebook](#download-books-from-talebook)        | <https://github.com/talebook/talebook> | ✅               | ❌                                      | ❌                                 | ❌                                |
| [SoBooks](#download-books-from-sobooks)          | <https://sobooks.cc>                  | ✅               | ❌                                      | ✅                                 | ❌                                |
| [Telegram](#download-books-from-telegram-groups) | <https://t.me>                         | ✅               | ❌                                      | ❌                                 | ❌                                |
| [Hsu Life](#download-books-from-hsu-life)        | <https://book.hsu.life>                | ✅               | ❌                                      | ❌                                 | ❌                                |

### Login Aliyundrive to get the `refreshToken`

We would show a QR code at the first time. And cache the `refreshToken` after successfully login.

```shell
bookhunter aliyun
```

### Download textbooks for Kids

```text
Usage:
  bookhunter k12 [flags]

Flags:
  -d, --download string   The book directory you want to use (default ".")
  -h, --help              help for k12
      --ratelimit int     The allowed requests per minutes for every thread (default 30)
  -t, --thread int        The number of download thead (default 1)

Global Flags:
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
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
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
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
      --ratelimit int     The allowed requests per minutes for every thread (default 30)
  -r, --rename            Rename the book file by book id
  -t, --thread int        The number of download thead (default 1)
  -u, --username string   The talebook username
  -w, --website string    The talebook link

Global Flags:
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
```

### Download books from SoBooks

```text
Usage:
  bookhunter sobooks [flags]

Flags:
      --code string       The secret code for SoBooks (default "244152")
  -d, --download string   The book directory you want to use (default ".")
  -e, --extract           Extract the archive file for filtering
  -f, --format strings    The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help              help for sobooks
  -i, --initial int       The book id you want to start download (default 1)
      --ratelimit int     The allowed requests per minutes for every thread (default 30)
  -r, --rename            Rename the book file by book id
  -t, --thread int        The number of download thead (default 1)

Global Flags:
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
```

### Download books from Telegram groups

Example command: `bookhunter telegram --appID ****** --appHash ****** --thread 4 --proxy http://127.0.0.1:7890 --channelID https://t.me/sharebooks4you`

Please refer [Creating your Telegram Application](https://core.telegram.org/api/obtaining_api_id) to obtain your `appID`
and `appHash`.

```text
Usage:
  bookhunter telegram [flags]

Flags:
      --appHash string     The app hash for telegram
      --appID int          The app id for telegram
      --channelID string   The channel id for telegram
  -d, --download string    The book directory you want to use (default ".")
  -e, --extract            Extract the archive file for filtering
  -f, --format strings     The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help               help for telegram
  -i, --initial int        The book id you want to start download (default 1)
      --mobile string      The mobile number, we will add +86 as default zone code
      --ratelimit int      The allowed requests per minutes for every thread (default 30)
      --refresh            Refresh the login session
  -r, --rename             Rename the book file by book id
  -t, --thread int         The number of download thead (default 1)

Global Flags:
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
```

### Download books from Hsu Life

Example command: `bookhunter hsu --username ****** --password ******`

```text
Usage:
  bookhunter hsu [flags]

Flags:
  -d, --download string   The book directory you want to use (default ".")
  -f, --format strings    The file formats you want to download (default [epub,azw3,mobi,pdf,zip])
  -h, --help              help for hsu
  -i, --initial int       The book id you want to start download (default 1)
  -p, --password string   The hsu.life password
      --ratelimit int     The allowed requests per minutes for every thread (default 30)
  -r, --rename            Rename the book file by book id
  -t, --thread int        The number of download thead (default 1)
  -u, --username string   The hsu.life username

Global Flags:
  -c, --config string     The config path for bookhunter
  -k, --keyword strings   The keywords for books
      --proxy string      The request proxy
      --retry int         The retry times for a failed download (default 3)
  -s, --skip-error        Continue to download the next book if the current book download failed (default true)
      --verbose           Print all the logs for debugging
```

## 🎯 Task Management

### Overview

Task management allows you to configure multiple download tasks in a YAML file and execute them easily. This is much more convenient than running individual commands.

**Benefits:**
- ✅ Configure once, run anytime
- ✅ Manage multiple sources in one place
- ✅ Enable/disable tasks without deleting configuration
- ✅ Tag-based filtering and batch operations
- ✅ Interactive task selection with TUI
- ✅ Parallel or sequential execution

### Configuration File

The default configuration file is located at `~/.config/bookhunter/tasks.yaml`. You can also specify a custom location using `--config-file`.

Example configuration (see `tasks.example.yaml` for complete examples):

```yaml
version: "1.0"
tasks:
  - name: "mothlib-telegram"
    description: "Download books from MothLib Telegram channel"
    enabled: true
    tags: ["telegram", "daily"]
    command:
      category: "telegram"
      args:
        appID: "YOUR_APP_ID"
        appHash: "YOUR_APP_HASH"
        formats: ["epub"]
        thread: 4
        ratelimit: 120
        directory: "./downloads/MothLib"
        channelID: "https://t.me/MothLib"

  - name: "sobooks-download"
    description: "Download books from SoBooks"
    enabled: false
    tags: ["sobooks", "daily"]
    command:
      category: "sobooks"
      args:
        extract: true
        formats: ["epub", "zip"]
        code: "881120"
        thread: 2
        directory: "./downloads/sobooks"
```

### Task Commands

**List tasks:**
```bash
# List all tasks
bookhunter task list

# List enabled tasks only
bookhunter task list --enabled

# Filter by tag
bookhunter task list --tag telegram

# Filter by category
bookhunter task list --category sobooks
```

**Show task details:**
```bash
bookhunter task show mothlib-telegram
```

**Enable/disable tasks:**
```bash
bookhunter task enable mothlib-telegram
bookhunter task disable sobooks-download
```

**Remove tasks:**
```bash
bookhunter task remove old-task
bookhunter task remove old-task --force  # Skip confirmation
```

**Validate configuration:**
```bash
bookhunter task validate
```

### Running Tasks

**Execute all enabled tasks:**
```bash
bookhunter run
```

**Execute specific tasks:**
```bash
bookhunter run mothlib-telegram sobooks-download
```

**Execute by tags:**
```bash
bookhunter run --tag telegram
bookhunter run --tag daily --tag telegram  # Multiple tags
```

**Interactive selection:**
```bash
bookhunter run --interactive
```

**Preview mode (dry-run):**
```bash
bookhunter run --dry-run
```

**Parallel execution:**
```bash
bookhunter run --parallel
```

**Force execute disabled tasks:**
```bash
bookhunter run disabled-task --force
```

### Getting Started

1. Run any command to create example configuration:
   ```bash
   bookhunter run
   ```

2. Edit the configuration file:
   ```bash
   vi ~/.config/bookhunter/tasks.yaml
   ```

3. Fill in your actual parameters (API keys, passwords, etc.)

4. Enable tasks you want to use:
   ```bash
   bookhunter task enable mothlib-telegram
   ```

5. Run your tasks:
   ```bash
   bookhunter run
   ```
