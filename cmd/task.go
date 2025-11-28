package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/log"
	"github.com/bookstairs/bookhunter/internal/task"
	taskui "github.com/bookstairs/bookhunter/internal/task/ui"
	"github.com/bookstairs/bookhunter/internal/ui"
)

var taskConfigPath string

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "管理下载任务",
	Long: `管理和配置下载任务。

可以列出、添加、删除、启用、禁用任务，以及验证配置文件。

示例:
  # 列出所有任务
  bookhunter task list

  # 显示任务详情
  bookhunter task show mothlib-telegram

  # 启用/禁用任务
  bookhunter task enable mothlib-telegram
  bookhunter task disable mothlib-telegram

  # 验证配置文件
  bookhunter task validate`,
}

// List 子命令
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有任务",
	Long: `列出配置文件中的所有任务。

可以使用选项筛选特定的任务:
  --tag: 按标签筛选
  --enabled: 仅显示启用的任务
  --disabled: 仅显示禁用的任务
  --category: 按类别筛选

示例:
  # 列出所有任务
  bookhunter task list

  # 列出启用的任务
  bookhunter task list --enabled

  # 按标签筛选
  bookhunter task list --tag telegram --tag daily`,
	Run: listTasks,
}

var (
	listTags       []string
	listEnabled    bool
	listDisabled   bool
	listCategories []string
)

func init() {
	rootCmd.AddCommand(taskCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskShowCmd)
	taskCmd.AddCommand(taskEnableCmd)
	taskCmd.AddCommand(taskDisableCmd)
	taskCmd.AddCommand(taskValidateCmd)
	taskCmd.AddCommand(taskRemoveCmd)

	// 全局参数
	taskCmd.PersistentFlags().StringVarP(&taskConfigPath, "config-file", "", "", "任务配置文件路径")

	// list 子命令参数
	taskListCmd.Flags().StringSliceVarP(&listTags, "tag", "t", []string{}, "按标签筛选")
	taskListCmd.Flags().BoolVar(&listEnabled, "enabled", false, "仅显示启用的任务")
	taskListCmd.Flags().BoolVar(&listDisabled, "disabled", false, "仅显示禁用的任务")
	taskListCmd.Flags().StringSliceVarP(&listCategories, "category", "c", []string{}, "按类别筛选")
}

func getTaskConfigPath() string {
	if taskConfigPath != "" {
		return taskConfigPath
	}
	return task.GetDefaultConfigPath()
}

func ensureTaskManagerOrExit() *task.Manager {
	configPath := getTaskConfigPath()

	if !task.ConfigExists(configPath) {
		fmt.Println(ui.RenderError("配置文件不存在: " + configPath))
		fmt.Println(ui.RenderInfo("运行 'bookhunter run' 创建示例配置"))
		os.Exit(1)
	}

	manager, err := task.NewManager(configPath)
	log.Exit(err)

	return manager
}

func listTasks(cmd *cobra.Command, args []string) {
	manager := ensureTaskManagerOrExit()

	filter := task.TaskFilter{
		Tags:       listTags,
		Categories: listCategories,
	}

	if listEnabled && !listDisabled {
		enabled := true
		filter.Enabled = &enabled
	} else if listDisabled && !listEnabled {
		disabled := false
		filter.Enabled = &disabled
	}

	tasks, err := manager.List(filter)
	log.Exit(err)

	if len(tasks) == 0 {
		fmt.Println(ui.RenderWarning("没有找到匹配的任务"))
		return
	}

	fmt.Println(taskui.RenderTaskList(tasks))
	fmt.Println()
	fmt.Println(taskui.RenderTaskSummary(tasks))
}

// Show 子命令
var taskShowCmd = &cobra.Command{
	Use:   "show <task-name>",
	Short: "显示任务详细信息",
	Long:  "显示指定任务的详细配置信息。",
	Args:  cobra.ExactArgs(1),
	Run:   showTask,
}

func showTask(cmd *cobra.Command, args []string) {
	manager := ensureTaskManagerOrExit()

	taskName := args[0]
	t, err := manager.Get(taskName)
	log.Exit(err)

	fmt.Println(taskui.RenderTaskDetail(t))
}

// Enable 子命令
var taskEnableCmd = &cobra.Command{
	Use:   "enable <task-name>",
	Short: "启用任务",
	Long:  "启用指定的任务，使其可以在 'bookhunter run' 中执行。",
	Args:  cobra.ExactArgs(1),
	Run:   enableTask,
}

func enableTask(cmd *cobra.Command, args []string) {
	manager := ensureTaskManagerOrExit()

	taskName := args[0]
	if err := manager.Enable(taskName); err != nil {
		log.Exit(err)
	}

	fmt.Println(ui.RenderSuccess(fmt.Sprintf("已启用任务: %s", taskName)))
}

// Disable 子命令
var taskDisableCmd = &cobra.Command{
	Use:   "disable <task-name>",
	Short: "禁用任务",
	Long:  "禁用指定的任务，使其不会在 'bookhunter run' 中自动执行。",
	Args:  cobra.ExactArgs(1),
	Run:   disableTask,
}

func disableTask(cmd *cobra.Command, args []string) {
	manager := ensureTaskManagerOrExit()

	taskName := args[0]
	if err := manager.Disable(taskName); err != nil {
		log.Exit(err)
	}

	fmt.Println(ui.RenderSuccess(fmt.Sprintf("已禁用任务: %s", taskName)))
}

// Validate 子命令
var taskValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证配置文件",
	Long:  "验证任务配置文件的格式和内容是否正确。",
	Run:   validateConfig,
}

func validateConfig(cmd *cobra.Command, args []string) {
	configPath := getTaskConfigPath()

	if !task.ConfigExists(configPath) {
		log.Exit(fmt.Errorf("配置文件不存在: %s", configPath))
	}

	config, err := task.LoadConfig(configPath)
	if err != nil {
		fmt.Println(ui.RenderError("配置验证失败"))
		log.Exit(err)
	}

	fmt.Println(ui.RenderSuccess("配置文件格式正确"))
	fmt.Println(ui.RenderInfo(fmt.Sprintf("版本: %s", config.Version)))
	fmt.Println(ui.RenderInfo(fmt.Sprintf("任务数: %d", len(config.Tasks))))

	// 统计信息
	enabledCount := 0
	for _, t := range config.Tasks {
		if t.Enabled {
			enabledCount++
		}
	}

	fmt.Println()
	fmt.Println(ui.StyleMuted.Render(fmt.Sprintf("✓ 启用: %d", enabledCount)))
	fmt.Println(ui.StyleMuted.Render(fmt.Sprintf("⋯ 禁用: %d", len(config.Tasks)-enabledCount)))
}

// Remove 子命令
var taskRemoveCmd = &cobra.Command{
	Use:   "remove <task-name>",
	Short: "删除任务",
	Long:  "从配置文件中删除指定的任务。",
	Args:  cobra.ExactArgs(1),
	Run:   removeTask,
}

var removeForce bool

func init() {
	taskRemoveCmd.Flags().BoolVarP(&removeForce, "force", "f", false, "强制删除不需要确认")
}

func removeTask(cmd *cobra.Command, args []string) {
	manager := ensureTaskManagerOrExit()

	taskName := args[0]

	// 确认删除
	if !removeForce {
		t, err := manager.Get(taskName)
		log.Exit(err)

		fmt.Println(ui.RenderWarning(fmt.Sprintf("即将删除任务: %s", taskName)))
		if t.Description != "" {
			fmt.Println(ui.StyleMuted.Render("  " + t.Description))
		}
		fmt.Println(ui.RenderInfo("确认删除？[y/N]"))

		var answer string
		fmt.Scanln(&answer)
		answer = strings.ToLower(strings.TrimSpace(answer))

		if answer != "y" && answer != "yes" {
			fmt.Println(ui.RenderWarning("已取消"))
			os.Exit(0)
		}
	}

	if err := manager.Remove(taskName); err != nil {
		log.Exit(err)
	}

	fmt.Println(ui.RenderSuccess(fmt.Sprintf("已删除任务: %s", taskName)))
}

