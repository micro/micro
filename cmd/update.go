package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/micro/cli/v2"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var (
	// SelfUpdate is set by gorelease LDFLAGS
	// We still prompt for update unless its disabled by env var
	// In future we may remove it entirely and always update
	SelfUpdate string
)

// confirmAndSelfUpdate looks for a new release of micro and upgrades in place
// we only execute this for select CLI commands rather than everything
func confirmAndSelfUpdate(ctx *cli.Context) (bool, error) {
	if SelfUpdate != "true" {
		return false, nil
	}

	latest, found, err := selfupdate.DetectLatest("micro/micro")
	if err != nil {
		return false, fmt.Errorf("Error occurred while detecting version: %s", err)
	}

	v := semver.MustParse(buildVersion())
	if !found || latest.Version.LTE(v) {
		// current version is the latest
		return false, nil
	}

	// if its not enabled via the update prompt bail out
	if ctx.Bool("prompt_update") {
		fmt.Print("New version found. Do you want to update to ", latest.Version, "? (yes/no): ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil || (input != "yes\n" && input != "no\n") {
			return false, fmt.Errorf("Invalid response")
		}
		if input == "no\n" {
			return false, nil
		}
	} else {
		fmt.Println("New version detected. Updating now...")
	}

	exe, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("Could not locate executable path")
	}
	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		return false, fmt.Errorf("Error occurred while updating binary: %s", err)
	}

	fmt.Println("Successfully updated to version", latest.Version)
	return true, nil
}
