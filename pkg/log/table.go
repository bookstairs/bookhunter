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
	printField(t, cv, allowZeroValue)

	t.Render()
}

func printField(t table.Writer, cv reflect.Value, allowZeroValue bool) {
	for cv.Kind() == reflect.Ptr || cv.Kind() == reflect.Interface {
		cv = cv.Elem()
	}

	// Print the struct.
	for i := 0; i < cv.NumField(); i++ {
		v := cv.Field(i)

		for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			printField(t, v, allowZeroValue)
		} else if allowZeroValue || !v.IsZero() {
			// Print the field which doesn't have the zero value.
			field := cv.Type().Field(i)
			value := v.Interface()
			t.AppendRow(table.Row{field.Name, value})
		}
	}
}
