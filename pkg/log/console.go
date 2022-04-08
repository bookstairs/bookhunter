package log

import (
	"fmt"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
)

var ansiStdout = ansi.NewAnsiStdout()

// Infof would print the log with in info level.
func Infof(format string, v ...any) {
	printPrefix("[green] [INFO] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
}

// Info would print the log with in info level.
func Info(v ...any) {
	printPrefix("[green] [INFO] [reset]")
	fmt.Println(v...)
}

// Warnf would print the log with in warn level.
func Warnf(format string, v ...any) {
	printPrefix("[yellow] [WARN] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
}

// Warn would print the log with in warn level.
func Warn(v ...any) {
	printPrefix("[yellow] [WARN] [reset]")
	fmt.Println(v...)
}

// Fatalf would print the log with in fatal level. And exit the program.
func Fatalf(format string, v ...any) {
	printPrefix("[red] [Fatal] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
	os.Exit(-1)
}

// Fatal would print the log with in fatal level. And exit the program.
func Fatal(v ...any) {
	printPrefix("[red] [Fatal] [reset]")
	fmt.Println(v...)
	os.Exit(-1)
}

// printPrefix would print a colorful log level and log time.
func printPrefix(level string) {
	fmt.Print(logTime())
	_, _ = fmt.Fprint(ansiStdout, colorstring.Color(level))
}

// logTime will print the current time
func logTime() string {
	return time.Now().Format("06/01/02 15:04:05")
}
