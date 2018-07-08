package version

import (
	"fmt"
	"runtime"
)

var version string

// Info contains versioning information.
type Info struct {
	Version      string `json:"Version"`
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.GitTag
}

// SetVersion for setup version string.
func SetVersion(v string) {
	version = v
}

// GetVersion for get current version.
func GetVersion() string {
	return version
}

func Get() Info {
	return Info{
		Version:			version,
		GitTag:       gitTag,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
