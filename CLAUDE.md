# CLAUDE.md - BookHunter 项目规范文档

## 📋 项目概述

**BookHunter** 是一个使用 Go 语言开发的多源图书下载工具，用于从多个在线图书平台批量下载电子书。

### 核心特性

- ✅ **多平台支持**: 支持 Talebook、SoBooks、中小学教材平台、Telegram、Hsu Life 等多个图书源
- ✅ **多网盘驱动**: 支持阿里云盘、蓝奏云、电信云盘等网盘下载
- ✅ **并发下载**: 支持多线程并发下载，可配置线程数量
- ✅ **进度管理**: 使用 BitSet 实现断点续传和下载进度持久化
- ✅ **限流控制**: 内置速率限制器，避免请求过于频繁
- ✅ **格式筛选**: 支持按文件格式（EPUB、MOBI、PDF 等）筛选下载
- ✅ **关键词搜索**: 支持按关键词过滤要下载的图书
- ✅ **自动解压**: 支持自动解压压缩包文件
- ✅ **代理支持**: 支持 HTTP/SOCKS 代理配置

### 技术栈

- **语言**: Go 1.23.3+
- **CLI 框架**: Cobra
- **HTTP 客户端**: Resty v2
- **进度条**: progressbar/v3
- **并发控制**: go.uber.org/ratelimit
- **Telegram API**: gotd/td
- **构建工具**: GoReleaser
- **代码检查**: golangci-lint
- **版本控制**: pre-commit hooks

---

## 🏗️ 项目架构

### 目录结构

```
bookhunter/
├── cmd/                      # 命令行入口和子命令定义
│   ├── flags/               # 全局配置和命令行参数
│   ├── root.go              # 根命令定义
│   ├── talebook.go          # Talebook 子命令
│   ├── sobooks.go           # SoBooks 子命令
│   ├── telegram.go          # Telegram 子命令
│   ├── k12.go               # 中小学教材子命令
│   ├── hsu.go               # Hsu Life 子命令
│   └── aliyun.go            # 阿里云盘工具命令
├── internal/                 # 内部实现代码
│   ├── client/              # HTTP 客户端封装
│   ├── driver/              # 网盘驱动实现
│   │   ├── aliyun/          # 阿里云盘驱动
│   │   ├── lanzou/          # 蓝奏云驱动
│   │   └── telecom/         # 电信云盘驱动
│   ├── fetcher/             # 下载器核心实现
│   ├── file/                # 文件处理（写入、解压、格式化）
│   ├── log/                 # 日志和进度显示
│   ├── progress/            # 下载进度管理
│   ├── sobooks/             # SoBooks 元数据处理
│   ├── talebook/            # Talebook API 响应
│   └── telegram/            # Telegram 相关功能
├── main.go                   # 程序入口
├── go.mod                    # Go 模块依赖
└── README.md                 # 项目说明文档
```

### 架构设计原则

#### 1. 分层架构

```
┌─────────────────────────────────────┐
│         Command Layer (cmd/)        │  ← CLI 接口层
├─────────────────────────────────────┤
│       Fetcher Layer (fetcher/)      │  ← 业务逻辑层
├─────────────────────────────────────┤
│    Driver/Service Layer (driver/)   │  ← 数据源适配层
├─────────────────────────────────────┤
│   Infrastructure Layer (client/)    │  ← 基础设施层
└─────────────────────────────────────┘
```

#### 2. 接口驱动设计

项目广泛使用接口来实现解耦和可扩展性：

- **`Fetcher` 接口**: 定义下载器的核心行为
- **`Driver` 接口**: 定义网盘驱动的统一接口
- **`service` 接口**: 定义图书源服务的统一接口
- **`Progress` 接口**: 定义进度管理的行为
- **`Creator/Writer` 接口**: 定义文件创建和写入行为

---

## 🔧 核心组件详解

### 1. Fetcher（下载器）

**位置**: `internal/fetcher/fetcher.go`

Fetcher 是整个下载流程的核心协调者，负责：

- 初始化配置和进度管理
- 创建多个下载线程
- 协调并发下载
- 错误处理和重试逻辑
- 格式和关键词过滤

**关键实现**:

```go
type Fetcher interface {
    Download() error
}

type fetcher struct {
    *Config
    service  service       // 图书源服务适配器
    progress progress.Progress  // 进度管理器
    creator  file.Creator       // 文件创建器
    errs     chan error         // 错误通道
}
```

**下载流程**:

1. 查询图书总数
2. 创建或恢复下载进度
3. 启动多个下载线程（goroutine）
4. 每个线程从进度管理器获取待下载的图书 ID
5. 查询图书可用格式
6. 应用格式和关键词过滤
7. 下载文件并写入磁盘
8. 更新下载进度

### 2. Driver（网盘驱动）

**位置**: `internal/driver/`

Driver 是对不同网盘服务的统一抽象，支持：

- **阿里云盘** (`aliyun/`): OAuth 2.0 认证，支持分享链接解析
- **蓝奏云** (`lanzou/`): 直链解析
- **电信云盘** (`telecom/`): 加密算法和 Cookie 认证

**接口定义**:

```go
type Driver interface {
    Source() Source
    Resolve(link, passcode string) ([]Share, error)
    Download(share Share) (io.ReadCloser, int64, error)
}
```

**设计特点**:
- 使用工厂模式创建不同的 Driver 实例
- Share 结构体统一描述可下载的文件
- Properties 字段支持扩展元数据

### 3. Progress（进度管理）

**位置**: `internal/progress/progress.go`

基于 BitSet 实现高效的进度管理：

- **并发安全**: 使用互斥锁保护并发访问
- **持久化**: 进度存储在本地文件中，支持断点续传
- **速率限制**: 集成 uber ratelimit，控制请求频率
- **内存高效**: BitSet 只需 1 bit 表示一本书的下载状态

**核心方法**:

```go
type Progress interface {
    TakeRateLimit()              // 获取速率限制许可
    AcquireBookID() int64        // 获取待下载的图书 ID
    SaveBookID(bookID int64) error  // 保存下载进度
    Finished() bool              // 是否全部下载完成
    Size() int64                 // 图书总数
}
```

**BitSet 优势**:
- 100 万本书仅需约 125KB 内存
- O(1) 时间复杂度查询和更新
- 天然支持序列化到文件

### 4. Client（HTTP 客户端）

**位置**: `internal/client/client.go`

对 Resty 的封装，提供统一的 HTTP 客户端：

- 自动重试机制（3 次，间隔 3-10 秒）
- 随机 User-Agent
- Cookie 管理
- 代理支持（HTTP/SOCKS）
- 请求超时控制（5 分钟）
- 重定向策略配置

**配置管理**:

```go
type Config struct {
    HTTPS      bool
    Host       string
    Proxy      string
    ConfigRoot string
    Redirect   resty.RedirectPolicy
}
```

### 5. File（文件处理）

**位置**: `internal/file/`

文件处理模块包含：

- **Writer**: 文件写入和进度显示
- **Formats**: 文件格式定义（EPUB、MOBI、PDF、AZW3 等）
- **Decompress**: 自动解压 ZIP、RAR、7Z 等格式
- **Replacer**: 文件名非法字符替换

**文件命名策略**:
- 按 ID 重命名: `123.epub`
- 保留原名: `书名.epub`
- 限制文件名长度（60 字符）
- 替换非法字符（`/`, `\`, `?`, `*` 等）

---

## 🔌 扩展指南

### 添加新的图书源

1. **在 `internal/fetcher/` 中创建新的服务文件**:

```go
// internal/fetcher/newsource.go
package fetcher

type newSourceService struct {
    *Config
    client *client.Client
}

func (s *newSourceService) size() (int64, error) {
    // 返回图书总数
}

func (s *newSourceService) formats(bookID int64) (map[file.Format]driver.Share, error) {
    // 返回指定图书的可用格式
}

func (s *newSourceService) fetch(bookID int64, format file.Format, 
    share driver.Share, writer file.Writer) error {
    // 下载文件内容到 writer
}
```

2. **在 `internal/fetcher/common.go` 中注册服务**:

```go
const (
    NewSource Category = "newsource"
)

func New(config *Config) (Fetcher, error) {
    switch config.Category {
    case NewSource:
        return newNewSource(config)
    // ...
    }
}
```

3. **在 `cmd/` 中添加新的子命令**:

```go
// cmd/newsource.go
var newsourceCmd = &cobra.Command{
    Use:   "newsource",
    Short: "Download books from NewSource",
    Run: func(cmd *cobra.Command, args []string) {
        f, err := flags.NewFetcher(fetcher.NewSource, properties)
        log.Exit(err)
        err = f.Download()
        log.Exit(err)
    },
}

func init() {
    rootCmd.AddCommand(newsourceCmd)
}
```

### 添加新的网盘驱动

1. **在 `internal/driver/` 中创建驱动目录和文件**:

```go
// internal/driver/newdrive/newdrive.go
package newdrive

type newDriver struct {
    config *client.Config
    client *client.Client
}

func (d *newDriver) Source() driver.Source {
    return driver.NewDrive
}

func (d *newDriver) Resolve(link, passcode string) ([]driver.Share, error) {
    // 解析分享链接，返回文件列表
}

func (d *newDriver) Download(share driver.Share) (io.ReadCloser, int64, error) {
    // 返回文件的可读流和大小
}
```

2. **在 `internal/driver/common.go` 中注册驱动**:

```go
const (
    NEWDRIVE Source = "newdrive"
)

func New(config *client.Config, properties map[string]string) (Driver, error) {
    switch Source(properties["driver"]) {
    case NEWDRIVE:
        return newNewDriveDriver(config, properties)
    // ...
    }
}
```

---

## 💡 最佳实践

### 1. 错误处理

- 使用 `errors.New()` 定义包级别的错误常量
- 错误应该包含足够的上下文信息
- 使用 `fmt.Errorf()` 包装错误链
- 区分可重试错误和不可重试错误（如 `ErrFileNotExist`）

### 2. 日志记录

项目使用自定义的日志系统 (`internal/log/`):

```go
log.Info("正常信息")
log.Warn("警告信息")
log.Error("错误信息")
log.Debug("调试信息")  // 需要 --verbose 参数
log.Exit(err)          // 如果 err != nil 则打印并退出
```

### 3. 并发控制

- 使用 `sync.WaitGroup` 等待所有 goroutine 完成
- 使用 `sync.Mutex` 保护共享状态
- 使用 channel 传递错误信息
- 使用 `uber.go/ratelimit` 控制请求频率

### 4. 配置管理

配置存储在 `~/.config/bookhunter/` 目录下：

```
~/.config/bookhunter/
├── example.com/          # 按域名分类
│   ├── progress.db       # 下载进度
│   └── config.json       # 站点配置
└── aliyun/
    └── token.json        # 阿里云盘 Token
```

### 5. 测试

- 为核心逻辑编写单元测试
- 使用 `testify` 进行断言
- Mock 外部依赖（HTTP 请求、文件系统）
- 测试文件命名: `*_test.go`

---

## 🐛 常见问题

### Q1: 下载进度丢失怎么办？

下载进度存储在配置目录的 `progress.db` 文件中。如果需要重新下载，删除此文件即可。

### Q2: 如何处理速率限制？

使用 `--ratelimit` 参数调整请求频率（每线程每分钟的请求数）。实际总速率 = `ratelimit × thread`。

### Q3: 代理配置

支持以下格式：
- HTTP: `http://127.0.0.1:7890`
- SOCKS5: `socks5://127.0.0.1:1080`
- 认证代理: `http://user:pass@proxy.com:8080`

### Q4: 文件名包含非法字符

项目会自动替换非法字符并限制文件名长度。可以使用 `--rename` 选项使用图书 ID 作为文件名。

### Q5: Telegram 登录问题

首次登录会缓存 session，存储在配置目录。使用 `--refresh` 强制重新登录。

---

## 🔒 安全考虑

1. **凭证存储**: 所有敏感信息（Token、密码）存储在本地配置文件，不应提交到版本控制
2. **命令行参数**: 避免在命令历史中暴露密码，优先使用配置文件
3. **代理使用**: 注意代理服务器的可信度
4. **网盘认证**: 阿里云盘使用 OAuth 2.0，Token 定期刷新

---

## 🚀 性能优化

### 并发策略

- 默认线程数为 CPU 核心数
- 每个线程独立下载，互不阻塞
- 使用 channel 进行线程间通信
- BitSet 提供 O(1) 的进度查询

### 内存管理

- 流式下载，不将整个文件加载到内存
- Progress 使用 BitSet，内存占用极小
- 及时关闭文件句柄和 HTTP 连接

### 网络优化

- 自动重试失败的请求
- 支持 HTTP Keep-Alive
- 使用连接池复用 TCP 连接

---

## 📚 代码规范

### 命名约定

- 包名: 小写单词，如 `fetcher`, `driver`
- 接口: 大写开头，动词或名词，如 `Fetcher`, `Driver`
- 结构体: 大写开头（导出）或小写开头（私有）
- 方法: 驼峰命名，大写开头（导出）
- 常量: 全大写或大写开头

### 代码组织

- 每个包应有清晰的职责
- 接口定义在使用方包中
- 实现细节对外隐藏
- 优先使用组合而非继承

### 注释规范

- 每个导出的类型、函数、常量都应有注释
- 包级别的注释在 `doc.go` 或包的主文件开头
- 注释应该说明"为什么"而不仅是"是什么"

---

## 🛠️ 开发工具

### 本地构建

```bash
# 使用 GoReleaser 构建所有平台
goreleaser build --snapshot --clean

# 仅构建当前平台
go build -o bookhunter main.go
```

### 代码检查

```bash
# 运行所有 linters
golangci-lint run

# 自动修复
golangci-lint run --fix
```

### Pre-commit Hooks

项目使用 pre-commit 进行代码质量检查：

```bash
# 安装 hooks
pre-commit install

# 手动运行
pre-commit run --all-files
```
---

## 📖 参考资料

### 官方文档

- [Go 官方文档](https://golang.org/doc/)
- [Cobra CLI 框架](https://github.com/spf13/cobra)
- [Resty HTTP 客户端](https://github.com/go-resty/resty)

### 相关项目

- [Talebook](https://github.com/talebook/talebook) - 个人图书管理系统
- [GoTD](https://github.com/gotd/td) - Telegram API Go 实现

---

## 🤝 贡献指南

1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

### 代码审查要点

- 代码符合 Go 惯用法
- 通过所有 linter 检查
- 有足够的单元测试
- 更新相关文档
- 提交信息清晰明确

---

## 📄 许可证

本项目采用 Apache 2.0 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 📮 联系方式

- **项目地址**: https://github.com/bookstairs/bookhunter
- **Issues**: https://github.com/bookstairs/bookhunter/issues
- **Discussions**: https://github.com/bookstairs/bookhunter/discussions

---

*最后更新: 2025-11-28*
