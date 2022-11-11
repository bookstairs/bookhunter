//go:build !windows

package naming

import "strings"

var replacer = strings.NewReplacer(
	`/`, empty,
	`\`, empty,
	`*`, empty,
	`:`, empty,
	`"`, empty,
)
