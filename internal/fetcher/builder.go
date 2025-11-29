package fetcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/file"
)

// CategoryDefaults 定义每个 category 的默认配置
type CategoryDefaults struct {
	Website string
	Driver  string
}

// GlobalArgs 包含全局级别的参数（不特定于某个 category）
type GlobalArgs struct {
	ConfigRoot string
	Proxy      string
	Retry      int
	SkipError  bool
}

// FetcherBuilder 用于统一创建 Fetcher 的逻辑
// 无论是 CLI 命令还是任务系统，都应该使用 Builder 来创建 Fetcher
type FetcherBuilder struct {
	Category   Category
	Args       map[string]interface{}
	GlobalArgs *GlobalArgs
}

// GetCategoryDefaults 返回指定 category 的默认配置
func GetCategoryDefaults(category Category) CategoryDefaults {
	switch category {
	case SoBooks:
		return CategoryDefaults{
			Website: "https://sobooks.cc",
			Driver:  "lanzou",
		}
	case K12:
		return CategoryDefaults{
			Website: "https://bp.pep.com.cn",
			Driver:  "telecom",
		}
	case Hsu:
		return CategoryDefaults{
			Website: "https://hsu.life",
			Driver:  "aliyun",
		}
	case Talebook:
		// Talebook 需要用户指定私有服务器地址
		return CategoryDefaults{
			Website: "",
			Driver:  "",
		}
	case Telegram:
		// Telegram 不需要 website 和 driver
		return CategoryDefaults{
			Website: "",
			Driver:  "",
		}
	default:
		return CategoryDefaults{}
	}
}

// Build 构建 Fetcher 实例
// 这是统一的 fetcher 创建入口，处理所有参数合并和配置逻辑
func (b *FetcherBuilder) Build() (Fetcher, error) {
	// 1. 获取 category 默认配置
	defaults := GetCategoryDefaults(b.Category)

	// 2. 合并 website 配置（用户配置覆盖默认值）
	website := defaults.Website
	if w, ok := b.Args["website"].(string); ok && w != "" {
		website = w
	}

	// 3. 合并 driver 配置（用户配置覆盖默认值）
	driver := defaults.Driver
	if d, ok := b.Args["driver"].(string); ok && d != "" {
		driver = d
	}

	// 4. 创建 client.Config
	var clientConfig *client.Config
	var err error

	if b.Category == Telegram {
		// Telegram 不需要 website，但需要 ConfigRoot 来存储 session
		configRoot := b.GlobalArgs.ConfigRoot
		if configRoot == "" {
			configRoot, err = client.DefaultConfigRoot()
			if err != nil {
				return nil, fmt.Errorf("获取默认配置目录失败: %w", err)
			}
		}
		clientConfig = &client.Config{
			Host:       "telegram", // 用于创建配置子目录
			Proxy:      b.GlobalArgs.Proxy,
			ConfigRoot: configRoot,
		}
	} else {
		clientConfig, err = client.NewConfig(website, b.GlobalArgs.Proxy, b.GlobalArgs.ConfigRoot)
		if err != nil {
			return nil, fmt.Errorf("创建客户端配置失败: %w", err)
		}
	}

	// 5. 解析格式
	formats, err := b.parseFormats()
	if err != nil {
		return nil, fmt.Errorf("解析格式失败: %w", err)
	}

	// 6. 解析关键词
	keywords := b.parseKeywords()

	// 7. 创建 properties（包括 driver 和其他特定参数）
	properties := b.buildProperties(driver)

	// 8. 创建 fetcher.Config
	fetcherConfig := &Config{
		Config:        clientConfig,
		Category:      b.Category,
		Formats:       formats,
		Keywords:      keywords,
		Extract:       b.getBool("extract", false),
		DownloadPath:  b.getString("directory", "."),
		InitialBookID: b.getInt64("initialBookID", 1),
		Rename:        b.getBool("rename", false),
		Thread:        b.getInt("thread", 4),
		RateLimit:     b.getInt("ratelimit", 30),
		Properties:    properties,
		Retry:         b.GlobalArgs.Retry,
		SkipError:     b.GlobalArgs.SkipError,
	}

	// 9. 调用现有的 New() 函数创建 Fetcher
	return New(fetcherConfig)
}

// parseFormats 解析格式列表
func (b *FetcherBuilder) parseFormats() ([]file.Format, error) {
	// 从 args 中获取 formats
	var formatStrs []string

	if v, ok := b.Args["formats"]; ok {
		switch val := v.(type) {
		case []string:
			formatStrs = val
		case []interface{}:
			for _, item := range val {
				if s, ok := item.(string); ok {
					formatStrs = append(formatStrs, s)
				}
			}
		case string:
			// 单个格式
			formatStrs = []string{val}
		}
	}

	// 如果没有指定格式，使用默认格式
	if len(formatStrs) == 0 {
		formatStrs = []string{"epub", "azw3", "mobi", "pdf", "zip"}
	}

	return ParseFormats(formatStrs)
}

// parseKeywords 解析关键词列表
func (b *FetcherBuilder) parseKeywords() []string {
	if v, ok := b.Args["keywords"]; ok {
		switch val := v.(type) {
		case []string:
			return val
		case []interface{}:
			var keywords []string
			for _, item := range val {
				if s, ok := item.(string); ok {
					keywords = append(keywords, s)
				}
			}
			return keywords
		case string:
			// 单个关键词
			return []string{val}
		}
	}
	return nil
}

// buildProperties 构建 properties map
func (b *FetcherBuilder) buildProperties(driver string) map[string]string {
	properties := make(map[string]string)

	// 设置 driver
	if driver != "" {
		properties["driver"] = driver
	}

	// 根据 category 设置特定参数
	switch b.Category {
	case SoBooks:
		// SoBooks 需要 code
		if code := b.getString("code", ""); code != "" {
			properties["code"] = code
		}

	case Telegram:
		// Telegram 需要 appID, appHash, channelID 等
		if appID := b.getInt64("appID", 0); appID != 0 {
			properties["appID"] = strconv.FormatInt(appID, 10)
		}
		if appHash := b.getString("appHash", ""); appHash != "" {
			properties["appHash"] = appHash
		}
		if channelID := b.getString("channelID", ""); channelID != "" {
			properties["channelID"] = channelID
		}
		if mobile := b.getString("mobile", ""); mobile != "" {
			properties["mobile"] = mobile
		}
		if relogin := b.getBool("relogin", false); relogin {
			properties["relogin"] = "true"
		}

	case Talebook, Hsu:
		// Talebook 和 Hsu 需要用户名密码
		if username := b.getString("username", ""); username != "" {
			properties["username"] = username
		}
		if password := b.getString("password", ""); password != "" {
			properties["password"] = password
		}
		if email := b.getString("email", ""); email != "" {
			properties["email"] = email
		}
	}

	// 添加驱动特定的参数
	if driver == "aliyun" {
		if refreshToken := b.getString("refreshToken", ""); refreshToken != "" {
			properties["refreshToken"] = refreshToken
		}
	} else if driver == "telecom" {
		if telecomUsername := b.getString("telecomUsername", ""); telecomUsername != "" {
			properties["telecomUsername"] = telecomUsername
		}
		if telecomPassword := b.getString("telecomPassword", ""); telecomPassword != "" {
			properties["telecomPassword"] = telecomPassword
		}
	}

	return properties
}

// 辅助方法：从 Args 中获取字符串值
func (b *FetcherBuilder) getString(key string, defaultValue string) string {
	if v, ok := b.Args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return defaultValue
}

// 辅助方法：从 Args 中获取 int 值
func (b *FetcherBuilder) getInt(key string, defaultValue int) int {
	if v, ok := b.Args[key]; ok {
		switch val := v.(type) {
		case int:
			return val
		case int64:
			return int(val)
		case float64:
			return int(val)
		case string:
			if i, err := strconv.Atoi(val); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// 辅助方法：从 Args 中获取 int64 值
func (b *FetcherBuilder) getInt64(key string, defaultValue int64) int64 {
	if v, ok := b.Args[key]; ok {
		switch val := v.(type) {
		case int64:
			return val
		case int:
			return int64(val)
		case float64:
			return int64(val)
		case string:
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// 辅助方法：从 Args 中获取 bool 值
func (b *FetcherBuilder) getBool(key string, defaultValue bool) bool {
	if v, ok := b.Args[key]; ok {
		switch val := v.(type) {
		case bool:
			return val
		case string:
			return strings.ToLower(val) == "true"
		}
	}
	return defaultValue
}
