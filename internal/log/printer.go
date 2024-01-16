package log

import (
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

var DefaultHead = []any{"Config Key", "Config Value"}

type Printer interface {
	Title(title string) Printer   // Title adds the table title.
	Head(heads ...any) Printer    // Head adds the table head.
	MaxWidth(width uint8) Printer // MaxColWidth the large column will be trimmed.
	Row(fields ...any) Printer    // Row add a row to the table.
	AllowZeroValue() Printer      // AllowZeroValue The row will be printed if it contains zero value.
	Print()                       // Print would print a table-like message from the given config.
}

type tablePrinter struct {
	title     string
	heads     []any
	width     uint8
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

func (t *tablePrinter) MaxWidth(width uint8) Printer {
	t.width = width
	return t
}

func (t *tablePrinter) Row(fields ...any) Printer {
	if len(fields) > 0 {
		// Trim the fields into a small length.
		for i, field := range fields {
			if f, ok := field.(string); ok && len(f) > int(t.width) {
				fields[i] = f[:t.width] + "..."
			}
		}
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
		appendRow(w, row, t.allowZero)
	}

	w.Render()
}

func appendRow(writer table.Writer, row []any, allowZero bool) {
	if len(row) == 1 {
		writer.AppendRow(row)
	} else {
		zero := true
		for _, r := range row[1:] {
			v := reflect.ValueOf(r)
			if !v.IsZero() {
				zero = false
				break
			}
		}
		if !zero || allowZero {
			writer.AppendRow(row)
		}
	}
}

// NewPrinter will return a printer for table-like logs.
func NewPrinter() Printer {
	return &tablePrinter{width: 30}
}
