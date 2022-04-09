//go:build !windows

package rename

import "strings"

var replacer = strings.NewReplacer(
	`/`, empty,
	`\`, empty,
	`*`, empty,
	`:`, empty,
	`"`, empty,
)
