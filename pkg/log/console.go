package log

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
)

var (
	ansiStdout = ansi.NewAnsiStdout()
)

// Infof would print the log with in info level.
func Infof(format string, v ...any) {
	printLog("[green][INFO][reset]", fmt.Sprintf(format, v...))
}

// Info would print the log with in info level.
func Info(v ...any) {
	printLog("[green][INFO][reset]", formatArgs(v...))
}

// Warnf would print the log with in warn level.
func Warnf(format string, v ...any) {
	printLog("[yellow][WARN][reset]", fmt.Sprintf(format, v...))
}

// Warn would print the log with in warn level.
func Warn(v ...any) {
	printLog("[green][INFO][reset]", formatArgs(v...))
}

// Fatalf would print the log with in fatal level. And exit the program.
func Fatalf(format string, v ...any) {
	printLog("[red][Fatal][reset]", fmt.Sprintf(format, v...))
	os.Exit(-1)
}

// Fatal would print the log with in fatal level. And exit the program.
func Fatal(v ...any) {
	printLog("[green][INFO][reset]", formatArgs(v...))
	os.Exit(-1)
}

// printLog would print a colorful log level and log time.
func printLog(level, log string) {
	_, _ = fmt.Fprintln(ansiStdout, logTime(), colorstring.Color(level), log)
}

func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	} else {
		return fmt.Sprint(args...)
	}
}

// logTime will print the current time
func logTime() string {
	return time.Now().Format("06/01/02 15:04:05")
}
