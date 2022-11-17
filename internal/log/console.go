package log

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
)

var (
	EnableDebug = false // EnableDebug will enable the disabled debug log level.

	ansiStdout = ansi.NewAnsiStdout()

	debug = colorstring.Color("[dark_gray][DEBUG][reset]")
	info  = colorstring.Color("[green][INFO] [reset]")
	warn  = colorstring.Color("[yellow][WARN] [reset]")
	fatal = colorstring.Color("[red][FATAL][reset]")
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
		printLog(debug, v...)
	}
}

// Infof would print the log with info level.
func Infof(format string, v ...any) {
	printLog(info, fmt.Sprintf(format, v...))
}

// Info would print the log with info level.
func Info(v ...any) {
	printLog(info, v...)
}

// Warnf would print the log with warn level.
func Warnf(format string, v ...any) {
	printLog(warn, fmt.Sprintf(format, v...))
}

// Warn would print the log with warn level.
func Warn(v ...any) {
	printLog(warn, v...)
}

// Fatalf would print the log with fatal level. And exit the program.
func Fatalf(format string, v ...any) {
	printLog(fatal, fmt.Sprintf(format, v...))
}

// Fatal would print the log with fatal level. And exit the program.
func Fatal(v ...any) {
	printLog(fatal, v...)
}

// Exit will print the error with fatal level and os.Exit if the error isn't nil.
func Exit(err error) {
	if err != nil {
		Fatal(err.Error())
		os.Exit(-1)
	}
}

// printLog would print a colorful log level and log time.
func printLog(level string, args ...any) {
	_, _ = fmt.Fprintln(ansiStdout, logTime(), level, fmt.Sprint(args...))
}

// logTime will print the current time
func logTime() string {
	return time.Now().Format("06/01/02 15:04:05")
}
