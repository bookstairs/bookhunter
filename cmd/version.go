package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/bibliolater/bookhunter/pkg/log"
)

var (
	gitMajor string // major version, always numeric
	gitMinor string // minor version, numeric possibly followed by "+"

	gitVersion   = ""
	gitCommit    = "" // sha1 from git, output of $(git rev-parse HEAD)
	gitTreeState = "" // state of git tree, either "clean" or "dirty"

	buildDate = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

var versionInfo = version{
	Major:        gitMajor,
	Minor:        gitMinor,
	GitVersion:   gitVersion,
	GitCommit:    gitCommit,
	GitTreeState: gitTreeState,
	BuildDate:    buildDate,
	GoVersion:    runtime.Version(),
	Compiler:     runtime.Compiler,
	Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
}

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return the bookhunter version info",
	Run: func(cmd *cobra.Command, args []string) {
		log.PrintTable("bookhunter version info", nil, &versionInfo, true)
	},
}

// version contains versioning information.
// how we'll want to distribute that information.
type version struct {
	Major        string `json:"major"`
	Minor        string `json:"minor"`
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}
