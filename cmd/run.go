package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/cmd/flags"
	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/task"
	taskui "github.com/bookstairs/bookhunter/internal/task/ui"
	"github.com/bookstairs/bookhunter/internal/ui"
)

var (
	runTags       []string
	runParallel   bool
	runDryRun     bool
	runForce      bool
	runConfigPath string
)

var runCmd = &cobra.Command{
	Use:   "run [task names...]",
	Short: "执行配置的任务",
	Long: `执行一个或多个配置的任务。

如果不指定任务名称，默认执行所有启用的任务。

示例:
  # 执行所有启用的任务
  bookhunter run

  # 执行指定任务
  bookhunter run mothlib-telegram sobooks-download

  # 按标签执行
  bookhunter run --tag daily --tag telegram

  # 并行执行所有任务
  bookhunter run --parallel

  # 预览模式（不实际执行）
  bookhunter run --dry-run

  # 强制执行禁用的任务
  bookhunter run disabled-task --force`,
	Run: runTask,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringSliceVarP(&runTags, "tag", "t", []string{}, "按标签筛选任务（可多次指定）")
	runCmd.Flags().BoolVarP(&runParallel, "parallel", "p", false, "并行执行所有任务")
	runCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "预览模式，显示将要执行的任务但不实际执行")
	runCmd.Flags().BoolVarP(&runForce, "force", "f", false, "强制执行禁用的任务")
	runCmd.Flags().StringVarP(&runConfigPath, "config-file", "", "", "任务配置文件路径（默认: ~/.config/bookhunter/tasks.yaml）")
}

func runTask(cmd *cobra.Command, args []string) {
	// 确定配置文件路径
	configPath := runConfigPath
	if configPath == "" {
		configPath = task.GetDefaultConfigPath()
	}

	// 检查配置文件是否存在
	if !task.ConfigExists(configPath) {
		fmt.Println(ui.RenderError("配置文件不存在: " + configPath))
		fmt.Println(ui.RenderInfo("运行 'bookhunter task create-example' 创建示例配置"))
		os.Exit(1)
	}

	// 创建任务管理器
	manager, err := task.NewManager(configPath)
	log.Exit(err)

	// 获取要执行的任务
	var tasksToRun []task.Task

	if len(runTags) > 0 {
		// 按标签筛选
		tasksToRun = manager.GetTasksByTag(runTags...)
		if len(tasksToRun) == 0 {
			log.Exit(fmt.Errorf("没有找到包含标签 %v 的任务", runTags))
		}
	} else if len(args) > 0 {
		// 按名称指定
		tasksToRun, err = manager.GetTasksByName(args...)
		if err != nil {
			log.Exit(err)
		}
	} else {
		// 执行所有启用的任务
		tasksToRun = manager.GetEnabledTasks()
		if len(tasksToRun) == 0 {
			log.Exit(fmt.Errorf("没有启用的任务"))
		}
	}

	// 检查禁用的任务
	if !runForce {
		var disabledTasks []string
		for _, t := range tasksToRun {
			if !t.Enabled {
				disabledTasks = append(disabledTasks, t.Name)
			}
		}

		if len(disabledTasks) > 0 {
			fmt.Println(ui.RenderWarning(fmt.Sprintf("以下任务已禁用: %v", disabledTasks)))
			fmt.Println(ui.RenderInfo("使用 --force 参数强制执行禁用的任务"))

			// 过滤掉禁用的任务
			var enabledTasks []task.Task
			for _, t := range tasksToRun {
				if t.Enabled {
					enabledTasks = append(enabledTasks, t)
				}
			}
			tasksToRun = enabledTasks

			if len(tasksToRun) == 0 {
				log.Exit(fmt.Errorf("没有可执行的启用任务"))
			}
		}
	}

	// 显示任务预览
	fmt.Println(taskui.RenderTaskPreview(tasksToRun))
	fmt.Println()

	// Dry-run 模式
	if runDryRun {
		fmt.Println(ui.RenderInfo("预览模式 - 不会实际执行任务"))
		os.Exit(0)
	}

	// 创建执行器
	globalArgs := &task.GlobalArgs{
		ConfigRoot: flags.ConfigRoot,
		Proxy:      flags.Proxy,
		Retry:      flags.Retry,
		SkipError:  flags.SkipError,
		DryRun:     runDryRun,
		Parallel:   runParallel,
	}

	executor := task.NewExecutor(manager, globalArgs)

	// 执行任务
	var results []task.ExecutionResult
	var execErr error

	if len(args) > 0 {
		results, execErr = executor.Run(args)
	} else if len(runTags) > 0 {
		results, execErr = executor.RunByTags(runTags)
	} else {
		results, execErr = executor.RunAll()
	}

	// 显示结果
	task.PrintResults(results)

	// 处理错误
	if execErr != nil {
		log.Exit(execErr)
	}

	// 检查是否有失败的任务
	hasFailure := false
	for _, result := range results {
		if !result.Success {
			hasFailure = true
			break
		}
	}

	if hasFailure {
		os.Exit(1)
	}
}
