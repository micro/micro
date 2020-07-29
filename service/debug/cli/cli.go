package cli

import (
	"fmt"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/debug/service"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/cmd"
)

func init() {
	cmd.Register(&cli.Command{
		Name:   "trace",
		Usage:  "Get tracing info from a service",
		Action: getTrace,
	})
}

const (
	// logUsage message for logs command
	traceUsage = "Required usage: micro trace example"
)

func getTrace(ctx *cli.Context) error {
	log.Trace("debug")

	// TODO look for trace id

	if ctx.Args().Len() == 0 {
		return fmt.Errorf("Require service name")
	}

	name := ctx.Args().Get(0)

	// must specify service name
	if len(name) == 0 {
		return fmt.Errorf(traceUsage)
	}

	// initialise a new service log
	// TODO: allow "--source" e.g. kubernetes
	srv := service.NewClient(name)

	spans, err := srv.Trace()
	if err != nil {
		return err
	}

	if len(spans) == 0 {
		return nil
	}

	fmt.Println("Id\tName\tTime\tDuration\tStatus")

	for _, span := range spans {
		fmt.Printf("%s\t%s\t%s\t%v\t%s\n",
			span.Trace,
			span.Name,
			time.Unix(0, int64(span.Started)).String(),
			time.Duration(span.Duration),
			"",
		)
	}

	return nil
}
