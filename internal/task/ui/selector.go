package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/bookstairs/bookhunter/internal/task"
	baseui "github.com/bookstairs/bookhunter/internal/ui"
)

// TaskSelector 任务选择器
type TaskSelector struct {
	tasks    []task.Task
	selected map[int]bool
	cursor   int
	done     bool
}

// NewTaskSelector 创建任务选择器
func NewTaskSelector(tasks []task.Task) *TaskSelector {
	return &TaskSelector{
		tasks:    tasks,
		selected: make(map[int]bool),
		cursor:   0,
		done:     false,
	}
}

// Init 初始化
func (m *TaskSelector) Init() tea.Cmd {
	return nil
}

// Update 更新
func (m *TaskSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.done = true
			return m, tea.Quit

		case "enter":
			m.done = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case " ":
			// 切换选中状态
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "a":
			// 全选
			for i := range m.tasks {
				m.selected[i] = true
			}

		case "n":
			// 取消全选
			m.selected = make(map[int]bool)
		}
	}

	return m, nil
}

// View 渲染视图
func (m *TaskSelector) View() string {
	if m.done {
		return ""
	}

	var b strings.Builder

	b.WriteString(baseui.RenderTitle("选择要执行的任务"))
	b.WriteString("\n\n")

	// 渲染任务列表
	for i, t := range m.tasks {
		cursor := "  "
		if m.cursor == i {
			cursor = lipgloss.NewStyle().Foreground(baseui.ColorPrimary).Render("▶ ")
		}

		checked := "[ ]"
		if m.selected[i] {
			checked = baseui.StyleSuccess.Render("[✓]")
		}

		// 任务名称
		name := t.Name
		if !t.Enabled {
			name = baseui.StyleMuted.Render(name + " (已禁用)")
		}

		// 标签
		tags := ""
		if len(t.Tags) > 0 {
			tags = " " + baseui.StyleMuted.Render(fmt.Sprintf("[%s]", strings.Join(t.Tags, ", ")))
		}

		line := fmt.Sprintf("%s%s %s%s", cursor, checked, name, tags)
		
		if m.cursor == i {
			line = lipgloss.NewStyle().Background(lipgloss.Color("#333333")).Render(line)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(baseui.StyleMuted.Render("使用 ↑/↓ 移动, 空格选择, Enter 确认, q 退出"))
	b.WriteString("\n")
	b.WriteString(baseui.StyleMuted.Render("快捷键: a=全选, n=取消全选"))
	b.WriteString("\n")

	return b.String()
}

// GetSelectedTasks 获取选中的任务
func (m *TaskSelector) GetSelectedTasks() []task.Task {
	var result []task.Task
	for i, selected := range m.selected {
		if selected && i < len(m.tasks) {
			result = append(result, m.tasks[i])
		}
	}
	return result
}

// SelectTasks 运行任务选择器
func SelectTasks(tasks []task.Task) ([]task.Task, error) {
	if len(tasks) == 0 {
		return nil, fmt.Errorf("没有可选择的任务")
	}

	selector := NewTaskSelector(tasks)
	p := tea.NewProgram(selector)
	
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(*TaskSelector); ok {
		return m.GetSelectedTasks(), nil
	}

	return nil, fmt.Errorf("获取选择结果失败")
}

