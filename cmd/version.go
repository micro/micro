package cmd

import (
	"fmt"

	ver "github.com/hashicorp/go-version"
)

var (
	version    = "v3.0.0"
	prerelease = "develop" // blank if full release
)

func buildVersion() string {
	verStr := version
	if prerelease != "" {
		verStr = fmt.Sprintf("%s-%s", version, prerelease)
	}
	// make sure we fail fast (panic) if bad version - this will get caught in CI tests
	ver.Must(ver.NewVersion(verStr))
	return verStr
}
