package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
)

var (
	ansiStdout = ansi.NewAnsiStdout()
	lock       = sync.RWMutex{}
)

// Infof would print the log with in info level.
func Infof(format string, v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[green] [INFO] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
}

// Info would print the log with in info level.
func Info(v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[green] [INFO] [reset]")
	fmt.Println(v...)
}

// Warnf would print the log with in warn level.
func Warnf(format string, v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[yellow] [WARN] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
}

// Warn would print the log with in warn level.
func Warn(v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[yellow] [WARN] [reset]")
	fmt.Println(v...)
}

// Fatalf would print the log with in fatal level. And exit the program.
func Fatalf(format string, v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[red] [Fatal] [reset]")
	fmt.Printf(format, v...)
	fmt.Println()
	os.Exit(-1)
}

// Fatal would print the log with in fatal level. And exit the program.
func Fatal(v ...any) {
	lock.Lock()
	defer lock.Unlock()

	printPrefix("[red] [Fatal] [reset]")
	last := len(v) - 1
	if err, ok := v[last].(error); ok {
		fmt.Print(v[0:last]...)
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Println(v...)
	}
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
