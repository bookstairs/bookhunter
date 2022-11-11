package log

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
)

var (
	ansiStdout  = ansi.NewAnsiStdout()
	EnableDebug = false // EnableDebug will enable the disabled debug log level.
)

const (
	debug = "[light_gray][DEBUG][reset]"
	info  = "[green][INFO][reset]"
	warn  = "[yellow][WARN][reset]"
	fatal = "[red][FATAL][reset]"
)

// Debugf would print the log with in debug level. The debug was disabled by default.
// You should use EnableDebug to enable it.
func Debugf(format string, v ...any) {
	if EnableDebug {
		printLog(debug, fmt.Sprintf(format, v...))
	}
}

// Debug would print the log with in debug level. The debug was disabled by default.
// // You should use EnableDebug to enable it.
func Debug(v ...any) {
	if EnableDebug {
		printLog(debug, formatArgs(v...))
	}
}

// Infof would print the log with in info level.
func Infof(format string, v ...any) {
	printLog(info, fmt.Sprintf(format, v...))
}

// Info would print the log with in info level.
func Info(v ...any) {
	printLog(info, formatArgs(v...))
}

// Warnf would print the log with in warn level.
func Warnf(format string, v ...any) {
	printLog(warn, fmt.Sprintf(format, v...))
}

// Warn would print the log with in warn level.
func Warn(v ...any) {
	printLog(warn, formatArgs(v...))
}

// Fatalf would print the log with in fatal level. And exit the program.
func Fatalf(format string, v ...any) {
	if l := len(v); l > 0 {
		if a := v[l-1]; a == nil {
			return
		}
	}
	printLog(fatal, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

// Fatal would print the log with in fatal level. And exit the program.
func Fatal(v ...any) {
	if l := len(v); l > 0 {
		if a := v[l-1]; a == nil {
			return
		}
	}
	printLog(fatal, formatArgs(v...))
	os.Exit(-1)
}

// formatArgs will format all the arguments.
func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	} else {
		return fmt.Sprint(args...)
	}
}

// printLog would print a colorful log level and log time.
func printLog(level, log string) {
	_, _ = fmt.Fprintln(ansiStdout, logTime(), colorstring.Color(level), log)
}

// logTime will print the current time
func logTime() string {
	return time.Now().Format("06/01/02 15:04:05")
}
