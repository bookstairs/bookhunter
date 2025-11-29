package task

import (
	"fmt"
	"sync"
	"time"

	"github.com/bookstairs/bookhunter/internal/fetcher"
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

	// 使用 FetcherBuilder 创建 fetcher
	builder := &fetcher.FetcherBuilder{
		Category: category,
		Args:     task.Command.Args,
		GlobalArgs: &fetcher.GlobalArgs{
			ConfigRoot: e.globalArgs.ConfigRoot,
			Proxy:      e.globalArgs.Proxy,
			Retry:      e.globalArgs.Retry,
			SkipError:  e.globalArgs.SkipError,
		},
	}

	return builder.Build()
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
