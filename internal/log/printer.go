package log

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

var DefaultHead = []any{"Config Key", "Config Value"}

type Printer interface {
	Title(title string) Printer // Title adds the table title.
	Head(heads ...any) Printer  // Head adds the table head.
	Row(fields ...any) Printer  // Row add a row to the table.
	AllowZeroValue() Printer    // AllowZeroValue The row will be printed if it contains zero value.
	Print()                     // Print would print a table-like message from the given config.
}

type tablePrinter struct {
	title     string
	heads     []any
	rows      [][]any
	allowZero bool
}

func (t *tablePrinter) Title(title string) Printer {
	t.title = title
	return t
}

func (t *tablePrinter) Head(heads ...any) Printer {
	t.heads = heads
	return t
}

func (t *tablePrinter) Row(fields ...any) Printer {
	if len(fields) > 0 {
		t.rows = append(t.rows, fields)
	}
	return t
}

func (t *tablePrinter) AllowZeroValue() Printer {
	t.allowZero = true
	return t
}

func (t *tablePrinter) Print() {
	w := table.NewWriter()
	w.SetOutputMirror(os.Stdout)
	w.SetTitle(t.title)
	if len(t.heads) > 0 {
		w.AppendHeader(t.heads)
	}

	for _, row := range t.rows {
		if len(row) == 1 {
			w.AppendRow([]any{row[0]})
		} else {
			zero := true
			for _, r := range row[1:] {
				v := reflect.ValueOf(r)
				if !v.IsZero() {
					zero = false
					break
				}
			}

			if !zero || t.allowZero {
				w.AppendRow(row)
			}
		}
	}

	w.Render()
}

// NewPrinter will return a printer for table-like logs.
func NewPrinter() Printer {
	return &tablePrinter{}
}
