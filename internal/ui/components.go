package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FormatFileOperation 格式化文件操作显示
// 用于显示 "old -> new" 这样的变更
func FormatFileOperation(label, oldValue, newValue string) string {
	var parts []string

	if label != "" {
		parts = append(parts, StyleMuted.Render(label+":"))
	}

	parts = append(parts, RenderOldValue(oldValue))
	parts = append(parts, StyleInfo.Render("→"))
	parts = append(parts, RenderNewValue(newValue))

	return strings.Join(parts, " ")
}

// FormatFilePath 格式化文件路径显示
func FormatFilePath(label, path string) string {
	if label != "" {
		return StyleMuted.Render(label+":") + " " + RenderPath(path)
	}
	return RenderPath(path)
}

// RenderProgressBar 渲染简单的进度条
func RenderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ""
	}

	percentage := float64(current) / float64(total)
	filled := int(percentage * float64(width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	percent := int(percentage * 100)

	style := lipgloss.NewStyle().Foreground(ColorSuccess)
	return fmt.Sprintf("[%s] %d%% (%d/%d)",
		style.Render(bar),
		percent,
		current,
		total,
	)
}

// RenderSimpleTable 渲染简单的两列表格
func RenderSimpleTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	// 计算列宽
	col1Width, col2Width := 0, 0
	for _, row := range rows {
		if len(row) >= 1 && len(row[0]) > col1Width {
			col1Width = len(row[0])
		}
		if len(row) >= 2 && len(row[1]) > col2Width {
			col2Width = len(row[1])
		}
	}

	// 添加一些 padding
	col1Width += 2
	col2Width += 2

	var lines []string

	// 顶部边框
	lines = append(lines, "┌"+strings.Repeat("─", col1Width)+"┬"+strings.Repeat("─", col2Width)+"┐")

	// 表格内容
	for i, row := range rows {
		col1 := row[0]
		col2 := ""
		if len(row) >= 2 {
			col2 = row[1]
		}

		// 填充空格
		col1 = col1 + strings.Repeat(" ", col1Width-len(col1))
		col2 = col2 + strings.Repeat(" ", col2Width-len(col2))

		lines = append(lines, "│"+col1+"│"+col2+"│")

		// 表头后添加分隔线
		if i == 0 {
			lines = append(lines, "├"+strings.Repeat("─", col1Width)+"┼"+strings.Repeat("─", col2Width)+"┤")
		}
	}

	// 底部边框
	lines = append(lines, "└"+strings.Repeat("─", col1Width)+"┴"+strings.Repeat("─", col2Width)+"┘")

	return strings.Join(lines, "\n")
}

// RenderStatsSummary 渲染统计摘要
func RenderStatsSummary(stats map[string]int) string {
	var rows [][]string

	// 表头
	rows = append(rows, []string{"  状态  ", " 数量 "})

	// 数据行
	if passed, ok := stats["passed"]; ok && passed > 0 {
		rows = append(rows, []string{
			" " + IconSuccess + " 通过 ",
			fmt.Sprintf(" %d ", passed),
		})
	}

	if failed, ok := stats["failed"]; ok && failed > 0 {
		rows = append(rows, []string{
			" " + IconError + " 失败 ",
			fmt.Sprintf(" %d ", failed),
		})
	}

	if handled, ok := stats["handled"]; ok && handled > 0 {
		rows = append(rows, []string{
			" " + IconInfo + " 已处理 ",
			fmt.Sprintf(" %d ", handled),
		})
	}

	if skipped, ok := stats["skipped"]; ok && skipped > 0 {
		rows = append(rows, []string{
			" " + IconSkip + " 跳过 ",
			fmt.Sprintf(" %d ", skipped),
		})
	}

	if updated, ok := stats["updated"]; ok && updated > 0 {
		rows = append(rows, []string{
			" " + IconSuccess + " 已更新 ",
			fmt.Sprintf(" %d ", updated),
		})
	}

	if total, ok := stats["total"]; ok {
		// 添加分隔线效果
		rows = append(rows, []string{"──────────", "──────"})
		rows = append(rows, []string{
			"  总计  ",
			fmt.Sprintf(" %d ", total),
		})
	}

	return RenderSimpleTable(rows)
}

// RenderMessageBox 渲染消息框
func RenderMessageBox(title, message string, boxType string) string {
	var style lipgloss.Style

	switch boxType {
	case "error":
		style = StyleErrorBox
		title = IconError + " " + title
	case "warning":
		style = StyleWarningBox
		title = IconWarning + " " + title
	case "info":
		style = StyleBorder
		title = IconInfo + " " + title
	default:
		style = StyleBorder
	}

	content := lipgloss.NewStyle().Bold(true).Render(title)
	if message != "" {
		content += "\n\n" + message
	}

	return style.Render(content)
}

// RenderSeparator 渲染分隔线
func RenderSeparator(width int) string {
	if width <= 0 {
		width = 50
	}
	return StyleMuted.Render(strings.Repeat("─", width))
}

// RenderHeader 渲染命令头部信息
func RenderHeader(commandName, description string) string {
	header := RenderTitle(commandName)
	if description != "" {
		header += "\n" + StyleMuted.Render(description)
	}
	return header + "\n" + RenderSeparator(50)
}

// FormatRenamePreview 格式化重命名预览
func FormatRenamePreview(oldName, newName string) string {
	maxLen := 40
	if len(oldName) > maxLen {
		oldName = "..." + oldName[len(oldName)-maxLen+3:]
	}
	if len(newName) > maxLen {
		newName = "..." + newName[len(newName)-maxLen+3:]
	}

	old := StylePath.Width(maxLen).Render(oldName)
	arrow := StyleInfo.Render(" → ")
	new := StyleSuccess.Width(maxLen).Render(newName)

	return old + arrow + new
}

