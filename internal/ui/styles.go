package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// 颜色定义
var (
	// 主题颜色
	ColorSuccess = lipgloss.Color("#00ff00") // 绿色 - 成功
	ColorError   = lipgloss.Color("#ff0000") // 红色 - 错误
	ColorWarning = lipgloss.Color("#ffff00") // 黄色 - 警告
	ColorInfo    = lipgloss.Color("#00ffff") // 青色 - 信息
	ColorPrimary = lipgloss.Color("#00aaff") // 蓝色 - 主要
	ColorMuted   = lipgloss.Color("#888888") // 灰色 - 次要

	// 细节颜色
	ColorPath      = lipgloss.Color("#aaaaaa") // 文件路径
	ColorHighlight = lipgloss.Color("#ff00ff") // 高亮
)

// 文本样式
var (
	// 基础样式
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	StyleWarning = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	StyleInfo    = lipgloss.NewStyle().Foreground(ColorInfo)
	StyleMuted   = lipgloss.NewStyle().Foreground(ColorMuted)

	// 文件路径样式
	StylePath = lipgloss.NewStyle().
			Foreground(ColorPath).
			Italic(true)

	// 标题样式
	StyleTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Underline(true)

	// 高亮样式
	StyleHighlight = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Bold(true)

	// 删除线样式（用于旧内容）
	StyleStrikethrough = lipgloss.NewStyle().
				Foreground(ColorWarning).
				Strikethrough(true)
)

// 布局样式
var (
	// 边框样式
	StyleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	// 警告框样式
	StyleWarningBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorWarning).
			Padding(1, 2)

	// 错误框样式
	StyleErrorBox = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorError).
			Padding(1, 2)

	// 表格边框
	StyleTableBorder = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorMuted)
)

// 状态指示符样式
var (
	// ✓ 成功图标
	IconSuccess = StyleSuccess.Render("✓")
	// ✗ 失败图标
	IconError = StyleError.Render("✗")
	// → 信息图标
	IconInfo = StyleInfo.Render("→")
	// ⚠ 警告图标
	IconWarning = StyleWarning.Render("⚠")
	// ⋯ 跳过图标
	IconSkip = StyleMuted.Render("⋯")
)

// RenderSuccess 渲染成功消息
func RenderSuccess(msg string) string {
	return IconSuccess + " " + StyleSuccess.Render(msg)
}

// RenderError 渲染错误消息
func RenderError(msg string) string {
	return IconError + " " + StyleError.Render(msg)
}

// RenderWarning 渲染警告消息
func RenderWarning(msg string) string {
	return IconWarning + " " + StyleWarning.Render(msg)
}

// RenderInfo 渲染信息消息
func RenderInfo(msg string) string {
	return IconInfo + " " + StyleInfo.Render(msg)
}

// RenderSkip 渲染跳过消息
func RenderSkip(msg string) string {
	return IconSkip + " " + StyleMuted.Render(msg)
}

// RenderPath 渲染文件路径
func RenderPath(path string) string {
	return StylePath.Render(path)
}

// RenderTitle 渲染标题
func RenderTitle(title string) string {
	return StyleTitle.Render(title)
}

// RenderOldValue 渲染旧值（带删除线）
func RenderOldValue(value string) string {
	return StyleStrikethrough.Render(value)
}

// RenderNewValue 渲染新值（高亮）
func RenderNewValue(value string) string {
	return StyleSuccess.Render(value)
}

