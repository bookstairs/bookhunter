package file

const (
	maxLength = 100
	empty     = " "
)

// escape the filename in *nix like systems and limit the max name size.
func escape(filename string) string {
	filename = replacer.Replace(filename)

	if name := []rune(filename); len(name) > maxLength {
		return string(name[0:maxLength])
	} else {
		return filename
	}
}
