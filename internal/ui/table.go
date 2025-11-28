package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// TableConfig 表格配置
type TableConfig struct {
	Headers      []string
	Rows         [][]string
	ColumnWidths []int  // 自定义列宽，nil 则自动计算
	BorderStyle  string // "normal", "rounded", "double", "thick", "none"
	HeaderStyle  lipgloss.Style
	CellStyle    lipgloss.Style
	BorderColor  lipgloss.Color
	ShowHeader   bool
	CompactMode  bool  // 紧凑模式（减少 padding）
	AlignRight   []int // 右对齐的列索引
}

// NewTableConfig 创建默认表格配置
func NewTableConfig() *TableConfig {
	return &TableConfig{
		BorderStyle: "normal",
		HeaderStyle: lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary),
		CellStyle:   lipgloss.NewStyle(),
		BorderColor: ColorMuted,
		ShowHeader:  true,
		CompactMode: false,
	}
}

// Table 表格生成器
type Table struct {
	config *TableConfig
}

// NewTable 创建新的表格
func NewTable(config *TableConfig) *Table {
	if config == nil {
		config = NewTableConfig()
	}
	return &Table{config: config}
}

// Render 渲染表格
func (t *Table) Render() string {
	if len(t.config.Rows) == 0 && !t.config.ShowHeader {
		return ""
	}

	// 计算列宽
	colWidths := t.calculateColumnWidths()

	// 获取边框字符
	borders := t.getBorderChars()

	var lines []string

	// 顶部边框
	if t.config.BorderStyle != "none" {
		lines = append(lines, t.renderTopBorder(colWidths, borders))
	}

	// 表头
	if t.config.ShowHeader && len(t.config.Headers) > 0 {
		lines = append(lines, t.renderRow(t.config.Headers, colWidths, t.config.HeaderStyle, borders))

		// 表头分隔线
		if t.config.BorderStyle != "none" {
			lines = append(lines, t.renderMiddleBorder(colWidths, borders))
		}
	}

	// 数据行
	for _, row := range t.config.Rows {
		lines = append(lines, t.renderRow(row, colWidths, t.config.CellStyle, borders))
	}

	// 底部边框
	if t.config.BorderStyle != "none" {
		lines = append(lines, t.renderBottomBorder(colWidths, borders))
	}

	return strings.Join(lines, "\n")
}

// calculateColumnWidths 计算列宽
func (t *Table) calculateColumnWidths() []int {
	if t.config.ColumnWidths != nil {
		return t.config.ColumnWidths
	}

	numCols := len(t.config.Headers)
	if numCols == 0 && len(t.config.Rows) > 0 {
		numCols = len(t.config.Rows[0])
	}

	colWidths := make([]int, numCols)

	// 检查表头宽度
	for i, header := range t.config.Headers {
		if i < numCols {
			colWidths[i] = len(header)
		}
	}

	// 检查数据行宽度
	for _, row := range t.config.Rows {
		for i, cell := range row {
			if i < numCols && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 添加 padding
	padding := 2
	if t.config.CompactMode {
		padding = 1
	}

	for i := range colWidths {
		colWidths[i] += padding * 2
	}

	return colWidths
}

// borderChars 边框字符集
type borderChars struct {
	topLeft      string
	topRight     string
	bottomLeft   string
	bottomRight  string
	horizontal   string
	vertical     string
	middleLeft   string
	middleRight  string
	middleMiddle string
	topMiddle    string
	bottomMiddle string
}

// getBorderChars 获取边框字符
func (t *Table) getBorderChars() borderChars {
	switch t.config.BorderStyle {
	case "rounded":
		return borderChars{
			topLeft: "╭", topRight: "╮",
			bottomLeft: "╰", bottomRight: "╯",
			horizontal: "─", vertical: "│",
			middleLeft: "├", middleRight: "┤", middleMiddle: "┼",
			topMiddle: "┬", bottomMiddle: "┴",
		}
	case "double":
		return borderChars{
			topLeft: "╔", topRight: "╗",
			bottomLeft: "╚", bottomRight: "╝",
			horizontal: "═", vertical: "║",
			middleLeft: "╠", middleRight: "╣", middleMiddle: "╬",
			topMiddle: "╦", bottomMiddle: "╩",
		}
	case "thick":
		return borderChars{
			topLeft: "┏", topRight: "┓",
			bottomLeft: "┗", bottomRight: "┛",
			horizontal: "━", vertical: "┃",
			middleLeft: "┣", middleRight: "┫", middleMiddle: "╋",
			topMiddle: "┳", bottomMiddle: "┻",
		}
	default: // "normal"
		return borderChars{
			topLeft: "┌", topRight: "┐",
			bottomLeft: "└", bottomRight: "┘",
			horizontal: "─", vertical: "│",
			middleLeft: "├", middleRight: "┤", middleMiddle: "┼",
			topMiddle: "┬", bottomMiddle: "┴",
		}
	}
}

// renderTopBorder 渲染顶部边框
func (t *Table) renderTopBorder(colWidths []int, borders borderChars) string {
	var parts []string
	parts = append(parts, borders.topLeft)

	for i, width := range colWidths {
		parts = append(parts, strings.Repeat(borders.horizontal, width))
		if i < len(colWidths)-1 {
			parts = append(parts, borders.topMiddle)
		}
	}

	parts = append(parts, borders.topRight)
	return lipgloss.NewStyle().Foreground(t.config.BorderColor).Render(strings.Join(parts, ""))
}

// renderMiddleBorder 渲染中间边框
func (t *Table) renderMiddleBorder(colWidths []int, borders borderChars) string {
	var parts []string
	parts = append(parts, borders.middleLeft)

	for i, width := range colWidths {
		parts = append(parts, strings.Repeat(borders.horizontal, width))
		if i < len(colWidths)-1 {
			parts = append(parts, borders.middleMiddle)
		}
	}

	parts = append(parts, borders.middleRight)
	return lipgloss.NewStyle().Foreground(t.config.BorderColor).Render(strings.Join(parts, ""))
}

// renderBottomBorder 渲染底部边框
func (t *Table) renderBottomBorder(colWidths []int, borders borderChars) string {
	var parts []string
	parts = append(parts, borders.bottomLeft)

	for i, width := range colWidths {
		parts = append(parts, strings.Repeat(borders.horizontal, width))
		if i < len(colWidths)-1 {
			parts = append(parts, borders.bottomMiddle)
		}
	}

	parts = append(parts, borders.bottomRight)
	return lipgloss.NewStyle().Foreground(t.config.BorderColor).Render(strings.Join(parts, ""))
}

// renderRow 渲染一行
func (t *Table) renderRow(cells []string, colWidths []int, style lipgloss.Style, borders borderChars) string {
	var parts []string

	if t.config.BorderStyle != "none" {
		parts = append(parts, lipgloss.NewStyle().Foreground(t.config.BorderColor).Render(borders.vertical))
	}

	padding := 1
	if !t.config.CompactMode {
		padding = 1
	}

	for i, cell := range cells {
		if i >= len(colWidths) {
			break
		}

		width := colWidths[i] - padding*2
		cellContent := cell

		// 处理对齐
		isRightAlign := false
		for _, idx := range t.config.AlignRight {
			if idx == i {
				isRightAlign = true
				break
			}
		}

		if isRightAlign {
			cellContent = fmt.Sprintf("%*s", width, cell)
		} else {
			cellContent = fmt.Sprintf("%-*s", width, cell)
		}

		paddedCell := strings.Repeat(" ", padding) + style.Render(cellContent) + strings.Repeat(" ", padding)
		parts = append(parts, paddedCell)

		if t.config.BorderStyle != "none" {
			parts = append(parts, lipgloss.NewStyle().Foreground(t.config.BorderColor).Render(borders.vertical))
		}
	}

	return strings.Join(parts, "")
}

// QuickTable 快速创建简单表格
func QuickTable(headers []string, rows [][]string) string {
	config := NewTableConfig()
	config.Headers = headers
	config.Rows = rows

	table := NewTable(config)
	return table.Render()
}

// QuickStatsTable 快速创建统计表格
func QuickStatsTable(stats map[string]interface{}) string {
	config := NewTableConfig()
	config.Headers = []string{"  项目  ", " 值 "}
	config.AlignRight = []int{1} // 数值右对齐
	config.BorderStyle = "rounded"

	var rows [][]string
	for key, value := range stats {
		rows = append(rows, []string{
			fmt.Sprintf(" %s ", key),
			fmt.Sprintf(" %v ", value),
		})
	}

	config.Rows = rows
	table := NewTable(config)
	return table.Render()
}

