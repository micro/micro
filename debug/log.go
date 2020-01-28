package debug

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/micro/cli/v2"
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

	if ctx.Args().Len() == 0 {
		fmt.Println("Require service name")
		return
	}

	name := ctx.Args().Get(0)

	// must specify service name
	if len(name) == 0 {
		fmt.Println(logUsage)
		return
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
		fmt.Println(err)
		return
	}

	output := ctx.String("output")
	for record := range logs.Chan() {
		switch output {
		case "json":
			b, _ := json.Marshal(record)
			fmt.Printf("%v\n", string(b))
		default:
			fmt.Printf("%v\n", record.Message)
		}
	}
}

// logFlags is shared flags so we don't have to continually re-add
func logFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "version",
			Usage: "Set the version of the service to debug",
		},
		&cli.StringFlag{
			Name:  "output, o",
			Usage: "Set the output format e.g json, text",
		},
		&cli.BoolFlag{
			Name:  "stream",
			Usage: "Set to stream logs continuously (default: true)",
			Value: true,
		},
		&cli.StringFlag{
			Name:  "since",
			Usage: "Set to the relative time from which to show the logs for e.g. 1h",
		},
		&cli.IntFlag{
			Name:  "count",
			Usage: "Set to query the last number of log events",
		},
	}
}
