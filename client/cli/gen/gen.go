// Package gen provides the micro gen command which simply runs go generate
package init

import (
	"fmt"
	"os/exec"

	"github.com/micro/micro/v3/cmd"
	"github.com/urfave/cli/v2"
)

var (
	Command = "go generate"
)

func Run(ctx *cli.Context) error {
	cmd := exec.Command("go", "generate")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Print(string(b))
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:        "gen",
		Usage:       "Generate a micro related dependencies e.g protobuf",
		Description: `'micro gen' will generate any micro related dependencies such as proto files`,
		Action:      Run,
		Flags:       []cli.Flag{},
	})
}
