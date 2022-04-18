package log

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

// PrintTable would print a table-like message from the given struct.
func PrintTable(title string, head table.Row, data any, allowZeroValue bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	if len(head) > 0 {
		t.AppendHeader(head)
	}

	// Reflect from ptr to instance.
	cv := reflect.ValueOf(data)
	for cv.Kind() == reflect.Ptr || cv.Kind() == reflect.Interface {
		cv = cv.Elem()
	}

	// Print the struct.
	for i := 0; i < cv.NumField(); i++ {
		field := cv.Type().Field(i)
		value := cv.Field(i).Interface()

		// Print the field which doesn't have the zero value.
		if allowZeroValue || reflect.Zero(field.Type) != value {
			t.AppendRow(table.Row{field.Name, value})
		}
	}

	t.Render()
}
