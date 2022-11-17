package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/bookstairs/bookhunter/internal/log"
)

var (
	gitVersion = ""
	gitCommit  = "" // sha1 from git, output of $(git rev-parse HEAD)
	buildDate  = "" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	goVersion  = runtime.Version()
	platform   = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the bookhunter version info",
	Run: func(cmd *cobra.Command, args []string) {
		log.NewPrinter().
			Title("bookhunter version info").
			Row("Version", gitVersion).
			Row("Commit", gitCommit).
			Row("Build Date", buildDate).
			Row("Go Version", goVersion).
			Row("Platform", platform).
			Print()
	},
}
