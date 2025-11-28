package task

import (
	"fmt"

	"github.com/bookstairs/bookhunter/internal/ui"
)

// Manager 任务管理器
type Manager struct {
	configPath string
	config     *TaskConfig
}

// NewManager 创建任务管理器
func NewManager(configPath string) (*Manager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		configPath: configPath,
		config:     config,
	}, nil
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *TaskConfig {
	return m.config
}

// Save 保存配置到文件
func (m *Manager) Save() error {
	return SaveConfig(m.config, m.configPath)
}

// Reload 重新加载配置
func (m *Manager) Reload() error {
	config, err := LoadConfig(m.configPath)
	if err != nil {
		return err
	}
	m.config = config
	return nil
}

// List 列出任务
func (m *Manager) List(filter TaskFilter) ([]Task, error) {
	return m.config.FilterTasks(filter), nil
}

// Get 获取单个任务
func (m *Manager) Get(name string) (*Task, error) {
	return m.config.GetTaskByName(name)
}

// Add 添加任务
func (m *Manager) Add(task Task) error {
	if err := m.config.AddTask(task); err != nil {
		return err
	}
	return m.Save()
}

// Remove 删除任务
func (m *Manager) Remove(name string) error {
	if err := m.config.RemoveTask(name); err != nil {
		return err
	}
	return m.Save()
}

// Update 更新任务
func (m *Manager) Update(name string, task Task) error {
	if err := m.config.UpdateTask(name, task); err != nil {
		return err
	}
	return m.Save()
}

// Enable 启用任务
func (m *Manager) Enable(name string) error {
	if err := m.config.EnableTask(name); err != nil {
		return err
	}
	return m.Save()
}

// Disable 禁用任务
func (m *Manager) Disable(name string) error {
	if err := m.config.DisableTask(name); err != nil {
		return err
	}
	return m.Save()
}

// Validate 验证配置
func (m *Manager) Validate() error {
	return m.config.Validate()
}

// GetEnabledTasks 获取所有启用的任务
func (m *Manager) GetEnabledTasks() []Task {
	enabled := true
	return m.config.FilterTasks(TaskFilter{Enabled: &enabled})
}

// GetDisabledTasks 获取所有禁用的任务
func (m *Manager) GetDisabledTasks() []Task {
	disabled := false
	return m.config.FilterTasks(TaskFilter{Enabled: &disabled})
}

// GetTasksByTag 根据标签获取任务
func (m *Manager) GetTasksByTag(tags ...string) []Task {
	return m.config.FilterTasks(TaskFilter{Tags: tags})
}

// GetTasksByCategory 根据类别获取任务
func (m *Manager) GetTasksByCategory(categories ...string) []Task {
	return m.config.FilterTasks(TaskFilter{Categories: categories})
}

// GetTasksByName 根据名称获取任务（可多个）
func (m *Manager) GetTasksByName(names ...string) ([]Task, error) {
	tasks := m.config.FilterTasks(TaskFilter{Names: names})
	
	if len(tasks) != len(names) {
		// 有些任务不存在，找出来报告
		found := make(map[string]bool)
		for _, task := range tasks {
			found[task.Name] = true
		}
		
		var missing []string
		for _, name := range names {
			if !found[name] {
				missing = append(missing, name)
			}
		}
		
		if len(missing) > 0 {
			return tasks, fmt.Errorf("部分任务不存在: %v", missing)
		}
	}
	
	return tasks, nil
}

// ListTasksFormatted 格式化列出任务（使用 UI 组件）
func (m *Manager) ListTasksFormatted(filter TaskFilter) string {
	tasks := m.config.FilterTasks(filter)
	
	if len(tasks) == 0 {
		return ui.RenderWarning("没有找到匹配的任务")
	}

	config := ui.NewTableConfig()
	config.Headers = []string{" 任务名称 ", " 状态 ", " 标签 ", " 类别 ", " 描述 "}
	config.BorderStyle = "rounded"

	var rows [][]string
	for _, task := range tasks {
		status := ui.IconSuccess + " 启用"
		if !task.Enabled {
			status = ui.IconSkip + " 禁用"
		}

		tags := "-"
		if len(task.Tags) > 0 {
			tags = fmt.Sprintf("%v", task.Tags)
		}

		description := task.Description
		if description == "" {
			description = "-"
		}

		rows = append(rows, []string{
			" " + task.Name + " ",
			" " + status + " ",
			" " + tags + " ",
			" " + task.Command.Category + " ",
			" " + description + " ",
		})
	}

	config.Rows = rows
	table := ui.NewTable(config)
	return table.Render()
}

// ShowTaskDetail 显示任务详细信息
func (m *Manager) ShowTaskDetail(name string) (string, error) {
	task, err := m.Get(name)
	if err != nil {
		return "", err
	}

	var details string
	details += ui.RenderTitle(task.Name) + "\n"
	details += ui.RenderSeparator(50) + "\n\n"
	
	if task.Description != "" {
		details += ui.StyleMuted.Render("描述: ") + task.Description + "\n"
	}
	
	status := "启用"
	if !task.Enabled {
		status = "禁用"
	}
	details += ui.StyleMuted.Render("状态: ") + status + "\n"
	
	if len(task.Tags) > 0 {
		details += ui.StyleMuted.Render("标签: ") + fmt.Sprintf("%v", task.Tags) + "\n"
	}
	
	details += ui.StyleMuted.Render("类别: ") + task.Command.Category + "\n"
	details += "\n" + ui.StyleMuted.Render("参数:") + "\n"
	
	for key, value := range task.Command.Args {
		details += fmt.Sprintf("  %s: %v\n", key, value)
	}
	
	return details, nil
}

// GetStats 获取统计信息
func (m *Manager) GetStats() map[string]int {
	stats := make(map[string]int)
	stats["total"] = len(m.config.Tasks)
	
	for _, task := range m.config.Tasks {
		if task.Enabled {
			stats["enabled"]++
		} else {
			stats["disabled"]++
		}
		
		// 按类别统计
		stats["category_"+task.Command.Category]++
	}
	
	return stats
}

