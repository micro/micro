package cmd

var (
	// populated by ldflags
	GitCommit string
	GitTag    string
	BuildDate string

	version = "latest"
)

func buildVersion() string {
	verStr := version

	// check for git tag via ldflags
	if len(GitTag) > 0 {
		verStr = GitTag
	}

	return verStr
}
