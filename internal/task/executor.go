package task

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bookstairs/bookhunter/internal/client"
	"github.com/bookstairs/bookhunter/internal/fetcher"
	"github.com/bookstairs/bookhunter/internal/file"
	"github.com/bookstairs/bookhunter/internal/ui"
)

// Executor 任务执行器
type Executor struct {
	config     *TaskConfig
	manager    *Manager
	globalArgs *GlobalArgs
}

// GlobalArgs 全局参数
type GlobalArgs struct {
	ConfigRoot string // bookhunter 配置根目录
	Proxy      string // 代理设置
	Retry      int    // 重试次数
	SkipError  bool   // 跳过错误继续执行
	DryRun     bool   // 仅预览不执行
	Parallel   bool   // 并行执行
}

// NewExecutor 创建任务执行器
func NewExecutor(manager *Manager, globalArgs *GlobalArgs) *Executor {
	return &Executor{
		config:     manager.GetConfig(),
		manager:    manager,
		globalArgs: globalArgs,
	}
}

// ExecutionResult 执行结果
type ExecutionResult struct {
	TaskName string
	Success  bool
	Error    error
	Duration time.Duration
}

// Run 执行指定的任务
func (e *Executor) Run(taskNames []string) ([]ExecutionResult, error) {
	// 获取要执行的任务
	tasks, err := e.manager.GetTasksByName(taskNames...)
	if err != nil {
		return nil, err
	}

	return e.executeTasks(tasks)
}

// RunByTags 根据标签执行任务
func (e *Executor) RunByTags(tags []string) ([]ExecutionResult, error) {
	tasks := e.manager.GetTasksByTag(tags...)
	if len(tasks) == 0 {
		return nil, fmt.Errorf("没有找到包含标签 %v 的任务", tags)
	}

	return e.executeTasks(tasks)
}

// RunAll 执行所有启用的任务
func (e *Executor) RunAll() ([]ExecutionResult, error) {
	tasks := e.manager.GetEnabledTasks()
	if len(tasks) == 0 {
		return nil, fmt.Errorf("没有启用的任务")
	}

	return e.executeTasks(tasks)
}

// executeTasks 执行任务列表
func (e *Executor) executeTasks(tasks []Task) ([]ExecutionResult, error) {
	if len(tasks) == 0 {
		return nil, fmt.Errorf("没有要执行的任务")
	}

	// Dry-run 模式
	if e.globalArgs.DryRun {
		fmt.Println(ui.RenderInfo("预览模式 - 不会实际执行任务"))
		fmt.Println(e.manager.ListTasksFormatted(TaskFilter{Names: getTaskNames(tasks)}))
		return nil, nil
	}

	// 并行执行
	if e.globalArgs.Parallel {
		return e.executeParallel(tasks)
	}

	// 串行执行
	return e.executeSequential(tasks)
}

// executeSequential 串行执行任务
func (e *Executor) executeSequential(tasks []Task) ([]ExecutionResult, error) {
	results := make([]ExecutionResult, 0, len(tasks))

	for i, task := range tasks {
		fmt.Printf("\n%s\n", ui.RenderSeparator(60))
		fmt.Printf("%s (%d/%d)\n", ui.RenderTitle(task.Name), i+1, len(tasks))
		if task.Description != "" {
			fmt.Println(ui.StyleMuted.Render(task.Description))
		}
		fmt.Printf("%s\n\n", ui.RenderSeparator(60))

		result := e.executeTask(&task)
		results = append(results, result)

		if !result.Success {
			fmt.Println(ui.RenderError(fmt.Sprintf("任务失败: %v", result.Error)))
			if !e.globalArgs.SkipError {
				return results, fmt.Errorf("任务 %s 执行失败，停止后续任务", task.Name)
			}
			fmt.Println(ui.RenderWarning("跳过错误，继续执行下一个任务"))
		} else {
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("任务完成，耗时: %s", result.Duration)))
		}
	}

	return results, nil
}

// executeParallel 并行执行任务
func (e *Executor) executeParallel(tasks []Task) ([]ExecutionResult, error) {
	fmt.Println(ui.RenderInfo(fmt.Sprintf("并行执行 %d 个任务", len(tasks))))

	var wg sync.WaitGroup
	resultsChan := make(chan ExecutionResult, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func(t Task) {
			defer wg.Done()
			result := e.executeTask(&t)
			resultsChan <- result
		}(task)
	}

	wg.Wait()
	close(resultsChan)

	// 收集结果
	var results []ExecutionResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

// executeTask 执行单个任务
func (e *Executor) executeTask(task *Task) ExecutionResult {
	startTime := time.Now()

	result := ExecutionResult{
		TaskName: task.Name,
		Success:  false,
	}

	// 创建 fetcher
	f, err := e.createFetcher(task)
	if err != nil {
		result.Error = fmt.Errorf("创建下载器失败: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	// 执行下载
	if err := f.Download(); err != nil {
		result.Error = fmt.Errorf("下载失败: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}

	result.Success = true
	result.Duration = time.Since(startTime)
	return result
}

// createFetcher 从任务创建 fetcher
func (e *Executor) createFetcher(task *Task) (fetcher.Fetcher, error) {
	// 创建客户端配置
	cc, err := e.createClientConfig(task)
	if err != nil {
		return nil, err
	}

	// 解析 fetcher 配置
	fc, err := e.parseFetcherConfig(task, cc)
	if err != nil {
		return nil, err
	}

	// 确定 category
	var category fetcher.Category
	switch task.Command.Category {
	case "telegram":
		category = fetcher.Telegram
	case "talebook":
		category = fetcher.Talebook
	case "sobooks":
		category = fetcher.SoBooks
	case "k12":
		category = fetcher.K12
	case "hsu":
		category = fetcher.Hsu
	default:
		return nil, fmt.Errorf("不支持的任务类别: %s", task.Command.Category)
	}

	fc.Category = category

	// 创建 fetcher
	return fetcher.New(fc)
}

// createClientConfig 创建客户端配置
func (e *Executor) createClientConfig(task *Task) (*client.Config, error) {
	// 从任务参数获取网站地址（如果有）
	website := ""
	if w, ok := task.Command.Args["website"].(string); ok {
		website = w
	}

	// 使用全局代理或任务指定的代理
	proxy := e.globalArgs.Proxy
	if p, ok := task.Command.Args["proxy"].(string); ok && p != "" {
		proxy = p
	}

	return client.NewConfig(website, proxy, e.globalArgs.ConfigRoot)
}

// parseFetcherConfig 解析 fetcher 配置
func (e *Executor) parseFetcherConfig(task *Task, cc *client.Config) (*fetcher.Config, error) {
	fc := &fetcher.Config{
		Config:    cc,
		Retry:     e.globalArgs.Retry,
		SkipError: e.globalArgs.SkipError,
	}

	args := task.Command.Args

	// 解析文件格式
	if formats, ok := args["formats"]; ok {
		fmtStrs, err := toStringSlice(formats)
		if err != nil {
			return nil, fmt.Errorf("解析 formats 失败: %w", err)
		}
		
		parsedFormats, err := fetcher.ParseFormats(fmtStrs)
		if err != nil {
			return nil, fmt.Errorf("解析 formats 失败: %w", err)
		}
		fc.Formats = parsedFormats
	} else {
		// 默认格式
		fc.Formats = []file.Format{file.EPUB, file.AZW3, file.MOBI, file.PDF}
	}

	// 解析关键词
	if keywords, ok := args["keywords"]; ok {
		kws, err := toStringSlice(keywords)
		if err != nil {
			return nil, fmt.Errorf("解析 keywords 失败: %w", err)
		}
		fc.Keywords = kws
	}

	// 解析布尔值参数
	if extract, ok := args["extract"].(bool); ok {
		fc.Extract = extract
	}
	if rename, ok := args["rename"].(bool); ok {
		fc.Rename = rename
	}

	// 解析下载路径
	if dir, ok := args["directory"].(string); ok && dir != "" {
		fc.DownloadPath = dir
	} else {
		// 默认下载路径
		cwd, _ := os.Getwd()
		fc.DownloadPath = cwd
	}

	// 解析初始书籍 ID
	if initialBookID, ok := toInt64(args["initialBookID"]); ok {
		fc.InitialBookID = initialBookID
	} else {
		fc.InitialBookID = 1
	}

	// 解析线程数
	if thread, ok := toInt(args["thread"]); ok {
		fc.Thread = thread
	} else {
		fc.Thread = 4 // 默认 4 线程
	}

	// 解析速率限制
	if ratelimit, ok := toInt(args["ratelimit"]); ok {
		fc.RateLimit = ratelimit
	} else {
		fc.RateLimit = 30 // 默认 30 req/min
	}

	// 构建 Properties map
	fc.Properties = make(map[string]string)

	// Telegram 特定参数
	if task.Command.Category == "telegram" {
		if appID, ok := args["appID"]; ok {
			fc.Properties["appID"] = fmt.Sprint(appID)
		}
		if appHash, ok := args["appHash"].(string); ok {
			fc.Properties["appHash"] = appHash
		}
		if channelID, ok := args["channelID"].(string); ok {
			fc.Properties["channelID"] = channelID
		}
		if mobile, ok := args["mobile"].(string); ok {
			fc.Properties["mobile"] = mobile
		}
		if relogin, ok := args["relogin"].(bool); ok {
			fc.Properties["relogin"] = fmt.Sprint(relogin)
		}
	}

	// SoBooks 特定参数
	if task.Command.Category == "sobooks" {
		if code, ok := args["code"].(string); ok {
			fc.Properties["code"] = code
		}
	}

	// Talebook / Hsu 网站参数
	if task.Command.Category == "talebook" || task.Command.Category == "hsu" {
		if username, ok := args["username"].(string); ok {
			fc.Properties["username"] = username
		}
		if password, ok := args["password"].(string); ok {
			fc.Properties["password"] = password
		}
		if email, ok := args["email"].(string); ok {
			fc.Properties["email"] = email
		}
	}

	// 驱动相关参数（用于某些特殊源）
	if driver, ok := args["driver"].(string); ok {
		fc.Properties["driver"] = driver
	}
	if source, ok := args["source"].(string); ok {
		fc.Properties["source"] = source
	}

	return fc, nil
}

// 辅助函数：类型转换

func toStringSlice(v interface{}) ([]string, error) {
	switch val := v.(type) {
	case []string:
		return val, nil
	case []interface{}:
		result := make([]string, len(val))
		for i, item := range val {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("元素 %d 不是字符串类型", i)
			}
			result[i] = str
		}
		return result, nil
	default:
		return nil, fmt.Errorf("无法转换为字符串切片: %T", v)
	}
}

func toInt(v interface{}) (int, bool) {
	switch val := v.(type) {
	case int:
		return val, true
	case int64:
		return int(val), true
	case float64:
		return int(val), true
	default:
		return 0, false
	}
}

func toInt64(v interface{}) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	default:
		return 0, false
	}
}

func getTaskNames(tasks []Task) []string {
	names := make([]string, len(tasks))
	for i, task := range tasks {
		names[i] = task.Name
	}
	return names
}

// PrintResults 打印执行结果
func PrintResults(results []ExecutionResult) {
	if len(results) == 0 {
		return
	}

	fmt.Printf("\n%s\n", ui.RenderSeparator(60))
	fmt.Println(ui.RenderTitle("执行结果"))
	fmt.Printf("%s\n\n", ui.RenderSeparator(60))

	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Println(ui.RenderSuccess(fmt.Sprintf("%s - 完成 (%s)", result.TaskName, result.Duration)))
		} else {
			failureCount++
			fmt.Println(ui.RenderError(fmt.Sprintf("%s - 失败: %v", result.TaskName, result.Error)))
		}
	}

	fmt.Println()
	stats := map[string]int{
		"passed": successCount,
		"failed": failureCount,
		"total":  len(results),
	}
	fmt.Println(ui.RenderStatsSummary(stats))
}

