// Package shutdown will issue a shutdown signal
package shutdown

import (
	"os/exec"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/cmd"
)

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "shutdown",
			Usage:  "Shutdown a micro server",
			Action: shutdown,
		},
	)
}

func shutdown(ctx *cli.Context) error {
	// send the shutdown signal
	// TODO: send Server.Shutdown
	_, err := exec.Command("killall", "micro").CombinedOutput()
	return err
}
