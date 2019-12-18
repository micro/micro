package debug

import (
	"fmt"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/debug/service"
	ulog "github.com/micro/go-micro/util/log"
)

const (
	// logUsage message for logs command
	logUsage = "Required usage: micro log example"
)

func getLog(ctx *cli.Context, srvOpts ...micro.Option) {
	ulog.Name("debug")

	// get the args
	since := ctx.String("since")
	count := ctx.Int("count")
	stream := ctx.Bool("stream")

	if len(ctx.Args()) == 0 {
		ulog.Fatal("Require service name")
	}

	name := ctx.Args()[0]

	// must specify service name
	if len(name) == 0 {
		ulog.Fatal(logUsage)
	}

	// initialise a new service log
	// TODO: allow "--source" e.g. kubernetes
	service := service.NewClient(name)

	var readSince time.Time
	d, err := time.ParseDuration(since)
	if err == nil {
		readSince = time.Now().Add(-d)
	}

	logs, err := service.Log(readSince, count, stream)
	if err != nil {
		ulog.Fatal(err)
	}

	var json bool
	switch ctx.String("output") {
	case "json":
		json = true
	}

	for record := range logs.Chan() {
		if json {
			fmt.Printf("%v\n", record)
		} else {
			fmt.Printf("%v\n", record.Value)
		}
	}
}

// logFlags is shared flags so we don't have to continually re-add
func logFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "Set the output format e.g json, text",
		},
		cli.BoolTFlag{
			Name:  "stream",
			Usage: "Set to stream logs continuously (default: true)",
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
