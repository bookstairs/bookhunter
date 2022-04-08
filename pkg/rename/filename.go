//go:build !windows

package rename

import "strings"

var empty = " "
var replacer = strings.NewReplacer(
	`/`, empty,
	`\`, empty,
	`*`, empty,
	`:`, empty,
	`"`, empty,
	`.`, empty,
)

// EscapeFilename escape the filename in *nix like systems.
func EscapeFilename(filename string) string {
	return replacer.Replace(filename)
}
