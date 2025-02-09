// Package runtime is the micro runtime
package runtime

import (
	"github.com/micro/micro/v5/cmd"
	"github.com/urfave/cli/v2"
)

// flags is shared flags so we don't have to continually re-add
var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "name",
		Usage: "Set the name of the service. Otherwise defaults to directory name",
	},
	&cli.StringFlag{
		Name:  "source",
		Usage: "Set the source url of the service e.g github.com/micro/services",
	},
	&cli.StringFlag{
		Name:  "image",
		Usage: "Set the image to use for the container",
	},
	&cli.StringFlag{
		Name:  "command",
		Usage: "Command to exec",
	},
	&cli.StringFlag{
		Name:  "args",
		Usage: "Command args",
	},
	&cli.StringFlag{
		Name:  "type",
		Usage: "The type of service operate on",
	},
	&cli.StringSliceFlag{
		Name:  "env_vars",
		Usage: "Set the environment variables e.g. foo=bar",
	},
	&cli.IntFlag{
		Name:    "instances",
		Aliases: []string{"i"},
		Usage:   "Number of instances to run",
		Value:   1,
	},
	&cli.StringSliceFlag{
		Name:  "metadata",
		Usage: "Set any metadata on the service e.g. foo=bar",
	},
	&cli.BoolFlag{
		Name: "watch",
		Usage: `Enable live-reloading, watch the file changes of source directories, then rebuild and restart the service. 
		only watching *.go and *.proto files now`,
	},
	&cli.IntFlag{
		Name: "watch_delay",
		Usage: `Watching delay milliseconds for live-reloading, only valid when --watch=true. 
		e.g. watch_delay=500 means watching delay time is 500ms.`,
		Value: 1000,
	},
	&cli.BoolFlag{
		Name:  "force",
		Usage: "Force rebuild and restart the service even though the service is running.",
	},
}

func init() {
	cmd.Register(
		&cli.Command{
			// In future we'll also have `micro run [x]` hence `micro run service` requiring "service"
			Name:  "run",
			Usage: RunUsage,
			Description: `Examples:
			micro run github.com/micro/services/helloworld
			micro run .  # deploy local folder to your local micro server
			micro run ../path/to/folder # deploy local folder to your local micro server
			micro run helloworld # deploy latest version, translates to micro run github.com/micro/services/helloworld
			micro run helloworld@9342934e6180 # deploy certain version
			micro run helloworld@branchname	# deploy certain branch`,
			Flags:  flags,
			Action: runService,
		},
		&cli.Command{
			Name:  "update",
			Usage: UpdateUsage,
			Description: `Examples:
			micro update github.com/micro/services/helloworld
			micro update .  # deploy local folder to your local micro server
			micro update ../path/to/folder # deploy local folder to your local micro server
			micro update helloworld # deploy master branch, translates to micro update github.com/micro/services/helloworld
			micro update helloworld@branchname	# deploy certain branch`,
			Flags:  flags,
			Action: updateService,
		},
		&cli.Command{
			Name:  "kill",
			Usage: KillUsage,
			Flags: flags,
			Description: `Examples:
			micro kill github.com/micro/services/helloworld
			micro kill .  # kill service deployed from local folder
			micro kill ../path/to/folder # kill service deployed from local folder
			micro kill helloworld # kill serviced deployed from master branch, translates to micro kill github.com/micro/services/helloworld
			micro kill helloworld@branchname	# kill service deployed from certain branch`,
			Action: killService,
		},
		&cli.Command{
			Name:   "status",
			Usage:  GetUsage,
			Flags:  flags,
			Action: getService,
		},
		&cli.Command{
			Name:   "logs",
			Usage:  "Get logs for a service e.g. micro logs helloworld",
			Action: getLogs,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "version",
					Usage: "Set the version of the service to debug",
				},
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "Set the output format e.g json, text",
				},
				&cli.BoolFlag{
					Name:    "follow",
					Aliases: []string{"f"},
					Usage:   "Set to stream logs continuously",
				},
				&cli.StringFlag{
					Name:  "since",
					Usage: "Set to the relative time from which to show the logs for e.g. 1h",
				},
				&cli.IntFlag{
					Name:    "lines",
					Aliases: []string{"n"},
					Usage:   "Set to query the last number of log events",
				},
			},
		},
	)
}
