// Package init provides the micro init command for initialising plugins and profiles
package init

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/micro/micro/v3/cmd"
	"github.com/urfave/cli/v2"
)

var (
	// The import path we use for profiles
	Import = "github.com/micro/micro/profile"
	// Vesion of micro
	Version = "v3"
)

func Run(ctx *cli.Context) error {
	var profiles []string

	for _, val := range ctx.StringSlice("profile") {
		for _, profile := range strings.Split(val, ",") {
			p := strings.TrimSpace(profile)
			if len(p) == 0 {
				continue
			}
			profiles = append(profiles, p)
		}
	}

	if len(profiles) == 0 {
		return nil
	}

	f := os.Stdout

	if v := ctx.String("output"); v != "stdout" {
		var err error
		f, err = os.Create(v)
		if err != nil {
			return err
		}
	}

	fmt.Fprint(f, "package main\n\n")
	fmt.Fprint(f, "import (\n")

	// write the profiles
	for _, profile := range profiles {
		path := path.Join(Import, profile, Version)
		line := fmt.Sprintf("\t_ \"%s\"\n", path)
		fmt.Fprint(f, line)
	}

	fmt.Fprint(f, ")\n")
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:        "init",
		Usage:       "Generate a profile for micro plugins",
		Description: `'micro init' generates a profile.go file defining plugins and profiles`,
		Action:      Run,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "profile",
				Usage: "A comma separated list of profiles to load",
				Value: cli.NewStringSlice(),
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Where to output the file, by default stdout",
				Value: "stdout",
			},
		},
	})
}
