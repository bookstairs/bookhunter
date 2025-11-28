package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// TaskConfig 任务配置文件结构
type TaskConfig struct {
	Version string `yaml:"version"`
	Tasks   []Task `yaml:"tasks"`
}

// Task 单个任务配置
type Task struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Enabled     bool              `yaml:"enabled"`
	Tags        []string          `yaml:"tags"`
	Command     Command           `yaml:"command"`
}

// Command 任务执行命令
type Command struct {
	Category string                 `yaml:"category"`
	Args     map[string]interface{} `yaml:"args"`
}

// LoadConfig 从文件加载任务配置
func LoadConfig(configPath string) (*TaskConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件不存在: %s", configPath)
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config TaskConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存任务配置到文件
func SaveConfig(config *TaskConfig, configPath string) error {
	// 确保目录存在
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Validate 验证任务配置
func (tc *TaskConfig) Validate() error {
	if tc.Version == "" {
		return errors.New("配置版本不能为空")
	}

	if len(tc.Tasks) == 0 {
		return errors.New("任务列表不能为空")
	}

	// 检查任务名称重复
	names := make(map[string]bool)
	for i, task := range tc.Tasks {
		if err := task.Validate(); err != nil {
			return fmt.Errorf("任务 #%d 验证失败: %w", i+1, err)
		}

		if names[task.Name] {
			return fmt.Errorf("任务名称重复: %s", task.Name)
		}
		names[task.Name] = true
	}

	return nil
}

// Validate 验证单个任务
func (t *Task) Validate() error {
	if t.Name == "" {
		return errors.New("任务名称不能为空")
	}

	if strings.ContainsAny(t.Name, " \t\n\r/\\") {
		return errors.New("任务名称不能包含空格或特殊字符")
	}

	if t.Command.Category == "" {
		return errors.New("任务类别不能为空")
	}

	// 验证类别是否有效
	validCategories := map[string]bool{
		"telegram": true,
		"talebook": true,
		"sobooks":  true,
		"k12":      true,
		"hsu":      true,
	}

	if !validCategories[t.Command.Category] {
		return fmt.Errorf("无效的任务类别: %s (支持: telegram, talebook, sobooks, k12, hsu)", t.Command.Category)
	}

	return nil
}

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() string {
	// 优先级:
	// 1. 环境变量 BOOKHUNTER_TASKS_CONFIG
	// 2. ~/.config/bookhunter/tasks.yaml

	if envPath := os.Getenv("BOOKHUNTER_TASKS_CONFIG"); envPath != "" {
		return envPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "tasks.yaml" // 回退到当前目录
	}

	return filepath.Join(homeDir, ".config", "bookhunter", "tasks.yaml")
}

// CreateExampleConfig 创建示例配置文件
func CreateExampleConfig(configPath string) error {
	example := &TaskConfig{
		Version: "1.0",
		Tasks: []Task{
			{
				Name:        "mothlib-telegram",
				Description: "从 MothLib Telegram 频道下载书籍",
				Enabled:     true,
				Tags:        []string{"telegram", "daily"},
				Command: Command{
					Category: "telegram",
					Args: map[string]interface{}{
						"appID":     "YOUR_APP_ID",
						"appHash":   "YOUR_APP_HASH",
						"formats":   []string{"epub"},
						"rename":    true,
						"thread":    4,
						"ratelimit": 120,
						"directory": "./downloads/MothLib",
						"channelID": "https://t.me/MothLib",
					},
				},
			},
			{
				Name:        "sobooks-download",
				Description: "从 SoBooks 下载书籍",
				Enabled:     false,
				Tags:        []string{"sobooks", "daily"},
				Command: Command{
					Category: "sobooks",
					Args: map[string]interface{}{
						"extract":       true,
						"formats":       []string{"epub", "zip"},
						"code":          "881120",
						"rename":        true,
						"thread":        2,
						"directory":     "./downloads/sobooks",
						"initialBookID": 22000,
						"ratelimit":     120,
					},
				},
			},
		},
	}

	return SaveConfig(example, configPath)
}

// ConfigExists 检查配置文件是否存在
func ConfigExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return err == nil
}

// GetTaskByName 根据名称获取任务
func (tc *TaskConfig) GetTaskByName(name string) (*Task, error) {
	for i := range tc.Tasks {
		if tc.Tasks[i].Name == name {
			return &tc.Tasks[i], nil
		}
	}
	return nil, fmt.Errorf("任务不存在: %s", name)
}

// FilterTasks 根据条件筛选任务
func (tc *TaskConfig) FilterTasks(filter TaskFilter) []Task {
	var result []Task

	for _, task := range tc.Tasks {
		if filter.Matches(&task) {
			result = append(result, task)
		}
	}

	return result
}

// TaskFilter 任务筛选器
type TaskFilter struct {
	Names      []string // 任务名称（精确匹配）
	Tags       []string // 标签（包含任意一个即匹配）
	Enabled    *bool    // 是否启用（nil 表示不筛选）
	Categories []string // 类别（包含任意一个即匹配）
}

// Matches 检查任务是否匹配筛选条件
func (f *TaskFilter) Matches(task *Task) bool {
	// 检查名称
	if len(f.Names) > 0 {
		found := false
		for _, name := range f.Names {
			if task.Name == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查标签
	if len(f.Tags) > 0 {
		found := false
		for _, filterTag := range f.Tags {
			for _, taskTag := range task.Tags {
				if taskTag == filterTag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// 检查启用状态
	if f.Enabled != nil && task.Enabled != *f.Enabled {
		return false
	}

	// 检查类别
	if len(f.Categories) > 0 {
		found := false
		for _, cat := range f.Categories {
			if task.Command.Category == cat {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// AddTask 添加新任务
func (tc *TaskConfig) AddTask(task Task) error {
	// 验证任务
	if err := task.Validate(); err != nil {
		return err
	}

	// 检查名称是否重复
	for _, t := range tc.Tasks {
		if t.Name == task.Name {
			return fmt.Errorf("任务名称已存在: %s", task.Name)
		}
	}

	tc.Tasks = append(tc.Tasks, task)
	return nil
}

// RemoveTask 删除任务
func (tc *TaskConfig) RemoveTask(name string) error {
	for i, task := range tc.Tasks {
		if task.Name == name {
			tc.Tasks = append(tc.Tasks[:i], tc.Tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("任务不存在: %s", name)
}

// UpdateTask 更新任务
func (tc *TaskConfig) UpdateTask(name string, newTask Task) error {
	// 验证新任务
	if err := newTask.Validate(); err != nil {
		return err
	}

	// 如果改名，检查新名称是否重复
	if name != newTask.Name {
		for _, t := range tc.Tasks {
			if t.Name == newTask.Name {
				return fmt.Errorf("任务名称已存在: %s", newTask.Name)
			}
		}
	}

	// 查找并更新
	for i, task := range tc.Tasks {
		if task.Name == name {
			tc.Tasks[i] = newTask
			return nil
		}
	}

	return fmt.Errorf("任务不存在: %s", name)
}

// EnableTask 启用任务
func (tc *TaskConfig) EnableTask(name string) error {
	task, err := tc.GetTaskByName(name)
	if err != nil {
		return err
	}
	task.Enabled = true
	return nil
}

// DisableTask 禁用任务
func (tc *TaskConfig) DisableTask(name string) error {
	task, err := tc.GetTaskByName(name)
	if err != nil {
		return err
	}
	task.Enabled = false
	return nil
}

