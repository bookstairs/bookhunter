package cmd

import (
	"fmt"
	"os"

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

可以列出和验证任务配置文件。
要添加/修改/删除任务，请直接编辑配置文件。

示例:
  # 列出所有任务
  bookhunter task list

  # 验证配置文件
  bookhunter task validate

  # 创建示例配置
  bookhunter task create-example`,
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
	taskCmd.AddCommand(taskValidateCmd)
	taskCmd.AddCommand(taskCreateExampleCmd)

	// 全局参数
	taskCmd.PersistentFlags().StringVarP(&taskConfigPath, "config-file", "", "", "任务配置文件路径")

	// list 子命令参数
	taskListCmd.Flags().StringSliceVarP(&listTags, "tag", "t", []string{}, "按标签筛选")
	taskListCmd.Flags().BoolVar(&listEnabled, "enabled", false, "仅显示启用的任务")
	taskListCmd.Flags().BoolVar(&listDisabled, "disabled", false, "仅显示禁用的任务")
	taskListCmd.Flags().StringSliceVar(&listCategories, "category", []string{}, "按类别筛选")
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
		fmt.Println(ui.RenderInfo("运行 'bookhunter task create-example' 创建示例配置"))
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

// CreateExample 子命令
var taskCreateExampleCmd = &cobra.Command{
	Use:   "create-example",
	Short: "创建示例配置文件",
	Long:  "在默认位置或指定位置创建示例任务配置文件。",
	Run:   createExampleConfig,
}

func createExampleConfig(cmd *cobra.Command, args []string) {
	configPath := getTaskConfigPath()

	// 检查文件是否已存在
	if task.ConfigExists(configPath) {
		fmt.Println(ui.RenderWarning("配置文件已存在: " + configPath))
		fmt.Println(ui.RenderInfo("如需重新创建，请先删除或移动现有配置文件"))
		os.Exit(1)
	}

	// 创建示例配置
	if err := task.CreateExampleConfig(configPath); err != nil {
		log.Exit(fmt.Errorf("创建示例配置失败: %w", err))
	}

	fmt.Println(ui.RenderSuccess("✓ 示例配置已创建: " + configPath))
	fmt.Println()
	fmt.Println(ui.RenderInfo("接下来的步骤:"))
	fmt.Println("  1. 编辑配置文件，填写您的实际参数")
	fmt.Println("  2. 运行 'bookhunter task validate' 验证配置")
	fmt.Println("  3. 运行 'bookhunter task list' 查看任务列表")
	fmt.Println("  4. 运行 'bookhunter run' 开始下载")
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

	fmt.Println(ui.RenderSuccess("✓ 配置文件格式正确"))
	fmt.Println(ui.RenderInfo(fmt.Sprintf("  版本: %s", config.Version)))
	fmt.Println(ui.RenderInfo(fmt.Sprintf("  任务数: %d", len(config.Tasks))))

	// 统计信息
	enabledCount := 0
	for _, t := range config.Tasks {
		if t.Enabled {
			enabledCount++
		}
	}

	fmt.Println()
	fmt.Println(ui.StyleMuted.Render(fmt.Sprintf("  ✓ 启用: %d", enabledCount)))
	fmt.Println(ui.StyleMuted.Render(fmt.Sprintf("  ⋯ 禁用: %d", len(config.Tasks)-enabledCount)))
}
