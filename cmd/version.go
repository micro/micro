package cmd

import (
	"fmt"

	ver "github.com/hashicorp/go-version"
)

var (
	// populated by ldflags
	GitCommit string
	GitTag    string
	BuildDate string

	version    = "v3.0.0"
	prerelease = "" // blank if full release
)

func buildVersion() string {
	verStr := version
	if prerelease != "" {
		verStr = fmt.Sprintf("%s-%s", version, prerelease)
	}

	// check for git tag via ldflags
	if len(GitTag) > 0 {
		verStr = GitTag
	}

	// make sure we fail fast (panic) if bad version - this will get caught in CI tests
	ver.Must(ver.NewVersion(verStr))
	return verStr
}
