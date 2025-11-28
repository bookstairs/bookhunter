package ui

import (
	"fmt"
	"strings"

	"github.com/bookstairs/bookhunter/internal/task"
	baseui "github.com/bookstairs/bookhunter/internal/ui"
)

// RenderTaskList 渲染任务列表
func RenderTaskList(tasks []task.Task) string {
	if len(tasks) == 0 {
		return baseui.RenderWarning("没有找到任务")
	}

	config := baseui.NewTableConfig()
	config.Headers = []string{" 任务名称 ", " 状态 ", " 标签 ", " 类别 ", " 描述 "}
	config.BorderStyle = "rounded"

	var rows [][]string
	for _, t := range tasks {
		status := baseui.IconSuccess + " 启用"
		if !t.Enabled {
			status = baseui.IconSkip + " 禁用"
		}

		tags := "-"
		if len(t.Tags) > 0 {
			tags = strings.Join(t.Tags, ", ")
		}

		description := t.Description
		if description == "" {
			description = "-"
		}
		
		// 限制描述长度
		if len(description) > 40 {
			description = description[:37] + "..."
		}

		rows = append(rows, []string{
			" " + t.Name + " ",
			" " + status + " ",
			" " + tags + " ",
			" " + t.Command.Category + " ",
			" " + description + " ",
		})
	}

	config.Rows = rows
	table := baseui.NewTable(config)
	return table.Render()
}

// RenderTaskDetail 渲染任务详细信息
func RenderTaskDetail(t *task.Task) string {
	var details string

	// 标题
	details += baseui.RenderTitle(t.Name) + "\n"
	details += baseui.RenderSeparator(60) + "\n\n"

	// 描述
	if t.Description != "" {
		details += baseui.StyleMuted.Render("描述: ") + t.Description + "\n"
	}

	// 状态
	status := baseui.IconSuccess + " 启用"
	if !t.Enabled {
		status = baseui.IconSkip + " 禁用"
	}
	details += baseui.StyleMuted.Render("状态: ") + status + "\n"

	// 标签
	if len(t.Tags) > 0 {
		details += baseui.StyleMuted.Render("标签: ") + strings.Join(t.Tags, ", ") + "\n"
	}

	// 类别
	details += baseui.StyleMuted.Render("类别: ") + t.Command.Category + "\n\n"

	// 参数
	details += baseui.StyleMuted.Render("参数:") + "\n"
	if len(t.Command.Args) == 0 {
		details += "  " + baseui.StyleMuted.Render("(无)") + "\n"
	} else {
		for key, value := range t.Command.Args {
			// 敏感信息脱敏
			strValue := fmt.Sprint(value)
			if key == "appHash" || key == "password" {
				if len(strValue) > 8 {
					strValue = strValue[:4] + "****" + strValue[len(strValue)-4:]
				} else {
					strValue = "****"
				}
			}
			details += fmt.Sprintf("  %s: %v\n", baseui.StyleInfo.Render(key), strValue)
		}
	}

	return details
}

// RenderTaskPreview 渲染任务预览（简短版本）
func RenderTaskPreview(tasks []task.Task) string {
	if len(tasks) == 0 {
		return baseui.RenderWarning("没有任务")
	}

	var preview string
	preview += baseui.RenderInfo(fmt.Sprintf("将执行 %d 个任务:", len(tasks))) + "\n\n"

	for i, t := range tasks {
		status := ""
		if !t.Enabled {
			status = baseui.StyleWarning.Render(" (已禁用)")
		}

		preview += fmt.Sprintf("  %d. %s%s - %s\n", 
			i+1, 
			baseui.StyleSuccess.Render(t.Name),
			status,
			baseui.StyleMuted.Render(t.Command.Category),
		)

		if t.Description != "" {
			desc := t.Description
			if len(desc) > 50 {
				desc = desc[:47] + "..."
			}
			preview += fmt.Sprintf("     %s\n", baseui.StyleMuted.Render(desc))
		}
	}

	return preview
}

// RenderTaskSummary 渲染任务摘要统计
func RenderTaskSummary(tasks []task.Task) string {
	if len(tasks) == 0 {
		return ""
	}

	stats := make(map[string]int)
	enabledCount := 0
	disabledCount := 0
	categoryCount := make(map[string]int)
	allTags := make(map[string]int)

	for _, t := range tasks {
		if t.Enabled {
			enabledCount++
		} else {
			disabledCount++
		}

		categoryCount[t.Command.Category]++

		for _, tag := range t.Tags {
			allTags[tag]++
		}
	}

	stats["total"] = len(tasks)
	
	if enabledCount > 0 {
		stats["passed"] = enabledCount // 使用 "passed" 显示为成功图标
	}
	if disabledCount > 0 {
		stats["skipped"] = disabledCount // 使用 "skipped" 显示为跳过图标
	}

	summary := baseui.RenderStatsSummary(stats)
	summary += "\n\n"

	// 类别统计
	if len(categoryCount) > 0 {
		summary += baseui.StyleMuted.Render("按类别统计:") + "\n"
		for cat, count := range categoryCount {
			summary += fmt.Sprintf("  %s: %d\n", cat, count)
		}
		summary += "\n"
	}

	// 标签统计
	if len(allTags) > 0 {
		summary += baseui.StyleMuted.Render("常用标签:") + "\n"
		for tag, count := range allTags {
			summary += fmt.Sprintf("  %s: %d\n", tag, count)
		}
	}

	return summary
}

