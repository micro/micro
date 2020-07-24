package server

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/debug/log"
	"github.com/micro/go-micro/v2/debug/log/kubernetes"
	dservice "github.com/micro/go-micro/v2/debug/service"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/micro/v2/service"
	logHandler "github.com/micro/micro/v2/service/debug/log/handler"
	pblog "github.com/micro/micro/v2/service/debug/log/proto"
	statshandler "github.com/micro/micro/v2/service/debug/stats/handler"
	pbstats "github.com/micro/micro/v2/service/debug/stats/proto"
	tracehandler "github.com/micro/micro/v2/service/debug/trace/handler"
	pbtrace "github.com/micro/micro/v2/service/debug/trace/proto"
)

const (
	name    = "go.micro.debug"
	address = ":8089"
)

var (
	// Flags specific to the debug service
	Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "window",
			Usage:   "Specifies how many seconds of stats snapshots to retain in memory",
			EnvVars: []string{"MICRO_DEBUG_WINDOW"},
			Value:   60,
		},
	}
)

// Run micro debug
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name(name),
		service.Address(address),
	)

	// default log initialiser
	newLog := func(service string) log.Log {
		// service log calls the actual service for the log
		return dservice.NewLog(
			// log with service name
			log.Name(service),
		)
	}

	source := ctx.String("log")
	switch source {
	case "service":
		newLog = func(service string) log.Log {
			return dservice.NewLog(
				log.Name(service),
			)
		}
	case "kubernetes":
		newLog = func(service string) log.Log {
			return kubernetes.NewLog(
				log.Name(service),
			)
		}
	}

	done := make(chan bool)
	defer func() {
		close(done)
	}()

	// create a service cache
	c := newCache(done)

	// log handler
	lgHandler := &logHandler.Log{
		// create the log map
		Logs: make(map[string]log.Log),
		// Create the new func
		New: newLog,
	}

	// Register the logs handler
	pblog.RegisterLogHandler(srv.Server(), lgHandler)

	// stats handler
	statsHandler, err := statshandler.New(done, ctx.Int("window"), c.services)
	if err != nil {
		logger.Fatal(err)
	}

	// stats handler
	traceHandler, err := tracehandler.New(done, ctx.Int("window"), c.services)
	if err != nil {
		logger.Fatal(err)
	}

	// Register the stats handler
	pbstats.RegisterStatsHandler(srv.Server(), statsHandler)
	// register trace handler
	pbtrace.RegisterTraceHandler(srv.Server(), traceHandler)

	// TODO: implement debug service for k8s cruft

	// start debug service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
	return nil
}
