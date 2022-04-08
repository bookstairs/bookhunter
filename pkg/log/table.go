package log

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

// PrintTable would print a table-like message from the given struct.
func PrintTable(title string, head table.Row, data interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	if len(head) > 0 {
		t.AppendHeader(head)
	}
	cv := reflect.ValueOf(data)
	for cv.Kind() == reflect.Ptr || cv.Kind() == reflect.Interface {
		cv = cv.Elem()
	}
	for i := 0; i < cv.NumField(); i++ {
		name := cv.Type().Field(i).Name
		value := cv.Field(i).Interface()
		t.AppendRow(table.Row{name, value})
	}
	t.Render()
}
