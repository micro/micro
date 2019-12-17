package debug

import (
	"fmt"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/log"
	"github.com/micro/go-micro/debug/service"
	ulog "github.com/micro/go-micro/util/log"
)

func getLogs(ctx *cli.Context, srvOpts ...micro.Option) {
	ulog.Name("debug")

	// get the args
	name := ctx.String("name")
	since := ctx.String("since")
	count := ctx.Int("count")
	stream := ctx.Bool("stream")

	// must specify service name
	if len(name) == 0 {
		ulog.Fatal(LogsUsage)
	}

	// initialise a new service log
	// TODO: allow "--log_source"
	service := service.NewClient(name)

	var options []log.ReadOption

	d, err := time.ParseDuration(since)
	if err == nil {
		readSince := time.Now().Add(-d)
		options = append(options, log.Since(readSince))
	}

	if count > 0 {
		options = append(options, log.Count(count))
	}

	if stream {
		options = append(options, log.Stream(stream))
	}

	logs, err := service.Log(options...)
	if err != nil {
		ulog.Fatal(err)
	}

	for record := range logs {
		fmt.Printf("%v\n", record)
	}
}

// logFlags is shared flags so we don't have to continually re-add
func logFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "Set the name of the service to debug",
		},
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
			Value: "latest",
		},
		cli.BoolFlag{
			Name:  "stream",
			Usage: "Set to stream logs continuously",
		},
		cli.StringFlag{
			Name:  "since",
			Usage: "Set to the relative time from which to show the logs for e.g. 1h",
		},
		cli.IntFlag{
			Name:  "count",
			Usage: "Set to query the last number of log events",
		},
	}
}
