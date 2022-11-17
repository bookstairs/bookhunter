//go:build !windows

package file

import "strings"

var replacer = strings.NewReplacer(
	`/`, empty,
	`\`, empty,
	`*`, empty,
	`:`, empty,
	`"`, empty,
)
