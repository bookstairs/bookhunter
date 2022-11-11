package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/log"
)

var (
	gitMajor     = "" // major version, always numeric
	gitMinor     = "" // minor version, numeric possibly followed by "+"
	gitVersion   = ""
	gitCommit    = "" // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = "" // state of git tree, either "clean" or "dirty"
	buildDate    = "" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion    = runtime.Version()
	compiler     = runtime.Compiler
	platform     = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the bookhunter version info",
	Run: func(cmd *cobra.Command, args []string) {
		log.NewPrinter().
			Title("bookhunter version info").
			Row("major", gitMajor).
			Row("minor", gitMinor).
			Row("gitVersion", gitVersion).
			Row("gitCommit", gitCommit).
			Row("gitTreeState", gitTreeState).
			Row("buildDate", buildDate).
			Row("goVersion", goVersion).
			Row("compiler", compiler).
			Row("platform", platform).
			Print()
	},
}
