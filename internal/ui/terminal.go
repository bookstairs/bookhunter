package ui

import (
	"os"

	"github.com/mattn/go-isatty"
)

// TerminalCapabilities 终端能力检测
type TerminalCapabilities struct {
	SupportsColor   bool
	SupportsUnicode bool
	Width           int
	IsTTY           bool
}

var termCaps *TerminalCapabilities

// DetectTerminalCapabilities 检测终端能力
func DetectTerminalCapabilities() *TerminalCapabilities {
	if termCaps != nil {
		return termCaps
	}

	caps := &TerminalCapabilities{
		SupportsColor:   true,
		SupportsUnicode: true,
		Width:           80,
		IsTTY:           false,
	}

	// 检测是否为 TTY
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		caps.IsTTY = true
	}

	// 检测颜色支持
	if !caps.IsTTY {
		caps.SupportsColor = false
	} else {
		term := os.Getenv("TERM")
		if term == "dumb" || term == "" {
			caps.SupportsColor = false
		}

		// 检查 NO_COLOR 环境变量
		if os.Getenv("NO_COLOR") != "" {
			caps.SupportsColor = false
		}
	}

	// 检测 Unicode 支持
	// 简单检测：如果 LANG 包含 UTF-8，则支持
	lang := os.Getenv("LANG")
	if lang == "" || !containsUTF8(lang) {
		caps.SupportsUnicode = false
	}

	termCaps = caps
	return caps
}

// containsUTF8 检查字符串是否包含 UTF-8
func containsUTF8(s string) bool {
	return len(s) > 0 && (s[len(s)-1:] == "8" ||
		len(s) >= 5 && s[len(s)-5:] == "UTF-8" ||
		len(s) >= 5 && s[len(s)-5:] == "utf-8")
}

// IsColorSupported 是否支持颜色
func IsColorSupported() bool {
	caps := DetectTerminalCapabilities()
	return caps.SupportsColor
}

// IsUnicodeSupported 是否支持 Unicode
func IsUnicodeSupported() bool {
	caps := DetectTerminalCapabilities()
	return caps.SupportsUnicode
}

// IsTTY 是否为终端
func IsTTY() bool {
	caps := DetectTerminalCapabilities()
	return caps.IsTTY
}

// GetFallbackIcons 获取降级图标
func GetFallbackIcons() map[string]string {
	if IsUnicodeSupported() {
		return map[string]string{
			"success": "✓",
			"error":   "✗",
			"info":    "→",
			"warning": "⚠",
			"skip":    "⋯",
		}
	}

	// ASCII 降级
	return map[string]string{
		"success": "[OK]",
		"error":   "[X]",
		"info":    "[>]",
		"warning": "[!]",
		"skip":    "[-]",
	}
}

// SafeIcon 安全获取图标（自动降级）
func SafeIcon(iconType string) string {
	icons := GetFallbackIcons()
	if icon, ok := icons[iconType]; ok {
		return icon
	}
	return iconType
}

// SafeRenderSuccess 安全渲染成功消息（自动降级）
func SafeRenderSuccess(msg string) string {
	icon := SafeIcon("success")
	if IsColorSupported() {
		return IconSuccess + " " + StyleSuccess.Render(msg)
	}
	return icon + " " + msg
}

// SafeRenderError 安全渲染错误消息（自动降级）
func SafeRenderError(msg string) string {
	icon := SafeIcon("error")
	if IsColorSupported() {
		return IconError + " " + StyleError.Render(msg)
	}
	return icon + " " + msg
}

// SafeRenderWarning 安全渲染警告消息（自动降级）
func SafeRenderWarning(msg string) string {
	icon := SafeIcon("warning")
	if IsColorSupported() {
		return IconWarning + " " + StyleWarning.Render(msg)
	}
	return icon + " " + msg
}

// SafeRenderInfo 安全渲染信息消息（自动降级）
func SafeRenderInfo(msg string) string {
	icon := SafeIcon("info")
	if IsColorSupported() {
		return IconInfo + " " + StyleInfo.Render(msg)
	}
	return icon + " " + msg
}

// SafeRenderSkip 安全渲染跳过消息（自动降级）
func SafeRenderSkip(msg string) string {
	icon := SafeIcon("skip")
	if IsColorSupported() {
		return IconSkip + " " + StyleMuted.Render(msg)
	}
	return icon + " " + msg
}

