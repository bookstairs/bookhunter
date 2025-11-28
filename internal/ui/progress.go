package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// ProgressTracker 进度跟踪器
type ProgressTracker struct {
	total        int
	current      int
	startTime    time.Time
	lastMessage  string
	width        int
	showTimeInfo bool // 是否显示时间信息
	showMessage  bool // 是否显示消息
	compactMode  bool // 紧凑模式（单行）
	successCount int  // 成功计数
	failureCount int  // 失败计数
	skippedCount int  // 跳过计数
}

// NewProgressTracker 创建新的进度跟踪器
func NewProgressTracker(total int) *ProgressTracker {
	return &ProgressTracker{
		total:        total,
		current:      0,
		startTime:    time.Now(),
		width:        40,
		showTimeInfo: true,
		showMessage:  true,
		compactMode:  false,
	}
}

// NewCompactProgressTracker 创建紧凑模式的进度跟踪器
func NewCompactProgressTracker(total int) *ProgressTracker {
	p := NewProgressTracker(total)
	p.compactMode = true
	p.showTimeInfo = false
	return p
}

// SetCompact 设置紧凑模式
func (p *ProgressTracker) SetCompact(compact bool) {
	p.compactMode = compact
	if compact {
		p.showTimeInfo = false
	}
}

// SetShowTimeInfo 设置是否显示时间信息
func (p *ProgressTracker) SetShowTimeInfo(show bool) {
	p.showTimeInfo = show
}

// SetShowMessage 设置是否显示消息
func (p *ProgressTracker) SetShowMessage(show bool) {
	p.showMessage = show
}

// IncrementSuccess 增加成功计数
func (p *ProgressTracker) IncrementSuccess() {
	p.current++
	p.successCount++
}

// IncrementFailure 增加失败计数
func (p *ProgressTracker) IncrementFailure() {
	p.current++
	p.failureCount++
}

// IncrementSkipped 增加跳过计数
func (p *ProgressTracker) IncrementSkipped() {
	p.current++
	p.skippedCount++
}

// GetStats 获取统计信息
func (p *ProgressTracker) GetStats() (success, failure, skipped int) {
	return p.successCount, p.failureCount, p.skippedCount
}

// Increment 增加进度
func (p *ProgressTracker) Increment() {
	p.current++
}

// SetMessage 设置当前消息
func (p *ProgressTracker) SetMessage(msg string) {
	p.lastMessage = msg
}

// Render 渲染进度条
func (p *ProgressTracker) Render() string {
	if p.total == 0 {
		return ""
	}

	if p.compactMode {
		return p.RenderCompact()
	}

	percentage := float64(p.current) / float64(p.total)
	filled := int(percentage * float64(p.width))

	// 构建进度条
	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	percent := int(percentage * 100)

	// 计算已用时间
	elapsed := time.Since(p.startTime)
	elapsedStr := formatDuration(elapsed)

	// 估算剩余时间
	var etaStr string
	if p.current > 0 {
		avgTime := elapsed / time.Duration(p.current)
		remaining := time.Duration(p.total-p.current) * avgTime
		etaStr = formatDuration(remaining)
	}

	progressText := fmt.Sprintf("[%s] %d%% (%d/%d)",
		StyleSuccess.Render(bar),
		percent,
		p.current,
		p.total,
	)

	var result string
	result = progressText

	if p.showTimeInfo {
		timeInfo := StyleMuted.Render(fmt.Sprintf("已用 %s", elapsedStr))
		if etaStr != "" && p.current < p.total {
			timeInfo += StyleMuted.Render(fmt.Sprintf(" | 预计剩余 %s", etaStr))
		}
		result += "\n" + timeInfo
	}

	if p.showMessage && p.lastMessage != "" {
		result += "\n" + StyleInfo.Render(p.lastMessage)
	}

	return result
}

// RenderSimple 渲染简单版本的进度（单行）
func (p *ProgressTracker) RenderSimple() string {
	if p.total == 0 {
		return ""
	}

	percentage := float64(p.current) / float64(p.total)
	percent := int(percentage * 100)

	return StyleMuted.Render(fmt.Sprintf("进度: %d%% (%d/%d)", percent, p.current, p.total))
}

// RenderCompact 渲染紧凑版本（单行带统计）
func (p *ProgressTracker) RenderCompact() string {
	if p.total == 0 {
		return ""
	}

	percentage := float64(p.current) / float64(p.total)
	percent := int(percentage * 100)
	filled := int(percentage * float64(20)) // 更短的进度条

	bar := strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)

	stats := ""
	if p.successCount > 0 {
		stats += StyleSuccess.Render(fmt.Sprintf("✓%d", p.successCount))
	}
	if p.failureCount > 0 {
		if stats != "" {
			stats += " "
		}
		stats += StyleError.Render(fmt.Sprintf("✗%d", p.failureCount))
	}
	if p.skippedCount > 0 {
		if stats != "" {
			stats += " "
		}
		stats += StyleMuted.Render(fmt.Sprintf("⋯%d", p.skippedCount))
	}

	result := fmt.Sprintf("[%s] %d%% (%d/%d)",
		StyleSuccess.Render(bar),
		percent,
		p.current,
		p.total,
	)

	if stats != "" {
		result += " " + stats
	}

	if p.showMessage && p.lastMessage != "" {
		// 截断消息以适应单行
		msg := p.lastMessage
		if len(msg) > 40 {
			msg = msg[:37] + "..."
		}
		result += " " + StyleMuted.Render(msg)
	}

	return result
}

// RenderWithStats 渲染带统计信息的进度
func (p *ProgressTracker) RenderWithStats() string {
	if p.total == 0 {
		return ""
	}

	percentage := float64(p.current) / float64(p.total)
	filled := int(percentage * float64(p.width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", p.width-filled)
	percent := int(percentage * 100)

	progressText := fmt.Sprintf("[%s] %d%% (%d/%d)",
		StyleSuccess.Render(bar),
		percent,
		p.current,
		p.total,
	)

	// 统计信息
	stats := ""
	if p.successCount > 0 {
		stats += IconSuccess + StyleSuccess.Render(fmt.Sprintf(" %d", p.successCount))
	}
	if p.failureCount > 0 {
		if stats != "" {
			stats += "  "
		}
		stats += IconError + StyleError.Render(fmt.Sprintf(" %d", p.failureCount))
	}
	if p.skippedCount > 0 {
		if stats != "" {
			stats += "  "
		}
		stats += IconSkip + StyleMuted.Render(fmt.Sprintf(" %d", p.skippedCount))
	}

	result := progressText
	if stats != "" {
		result += "\n" + stats
	}

	if p.showMessage && p.lastMessage != "" {
		result += "\n" + StyleInfo.Render(p.lastMessage)
	}

	return result
}

// formatDuration 格式化时间长度
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

// SimpleProgress 简单进度条（用于静态显示）
type SimpleProgress struct {
	prog progress.Model
}

// NewSimpleProgress 创建简单进度条
func NewSimpleProgress() *SimpleProgress {
	prog := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return &SimpleProgress{
		prog: prog,
	}
}

// Render 渲染进度条
func (sp *SimpleProgress) Render(percent float64) string {
	return sp.prog.ViewAs(percent)
}

// SpinnerFrames Spinner 动画帧
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Spinner 简单的 Spinner
type Spinner struct {
	frame     int
	message   string
	style     lipgloss.Style
	startTime time.Time
	showTime  bool
}

// NewSpinner 创建新的 Spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frame:     0,
		message:   message,
		style:     lipgloss.NewStyle().Foreground(ColorPrimary),
		startTime: time.Now(),
		showTime:  false,
	}
}

// NewSpinnerWithTime 创建带时间显示的 Spinner
func NewSpinnerWithTime(message string) *Spinner {
	s := NewSpinner(message)
	s.showTime = true
	return s
}

// Next 下一帧
func (s *Spinner) Next() {
	s.frame = (s.frame + 1) % len(SpinnerFrames)
}

// SetMessage 设置消息
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// SetStyle 设置样式
func (s *Spinner) SetStyle(style lipgloss.Style) {
	s.style = style
}

// Render 渲染 Spinner
func (s *Spinner) Render() string {
	result := s.style.Render(SpinnerFrames[s.frame]) + " " + s.message

	if s.showTime {
		elapsed := time.Since(s.startTime)
		result += " " + StyleMuted.Render(fmt.Sprintf("(%s)", formatDuration(elapsed)))
	}

	return result
}

// RenderInline 渲染内联版本（不带换行）
func (s *Spinner) RenderInline() string {
	return s.style.Render(SpinnerFrames[s.frame]) + " " + s.message
}

// MultiSpinner 多任务 Spinner
type MultiSpinner struct {
	tasks     []string
	completed []bool
	current   int
	spinner   *Spinner
}

// NewMultiSpinner 创建多任务 Spinner
func NewMultiSpinner(tasks []string) *MultiSpinner {
	return &MultiSpinner{
		tasks:     tasks,
		completed: make([]bool, len(tasks)),
		current:   0,
		spinner:   NewSpinner(""),
	}
}

// NextTask 移动到下一个任务
func (ms *MultiSpinner) NextTask() {
	if ms.current < len(ms.tasks) {
		ms.completed[ms.current] = true
		ms.current++
	}
}

// SetCurrentMessage 设置当前任务消息
func (ms *MultiSpinner) SetCurrentMessage(msg string) {
	if ms.current < len(ms.tasks) {
		ms.spinner.SetMessage(msg)
	}
}

// Tick 更新动画帧
func (ms *MultiSpinner) Tick() {
	ms.spinner.Next()
}

// Render 渲染多任务进度
func (ms *MultiSpinner) Render() string {
	var lines []string

	for i, task := range ms.tasks {
		if ms.completed[i] {
			lines = append(lines, IconSuccess+" "+StyleMuted.Render(task))
		} else if i == ms.current {
			lines = append(lines, ms.spinner.RenderInline()+" "+task)
		} else {
			lines = append(lines, StyleMuted.Render("  "+task))
		}
	}

	return strings.Join(lines, "\n")
}

