// Package platform manages the runtime services as a platform
package platform

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	gorun "github.com/micro/go-micro/v2/runtime"
	signalutil "github.com/micro/go-micro/v2/util/signal"
	"github.com/micro/micro/v2/cmd"

	// include usage

	"github.com/micro/micro/v2/internal/update"
	_ "github.com/micro/micro/v2/internal/usage"
)

var (
	// Date of the build
	// TODO: move elsewhere
	Version string

	// list of services managed
	Services = []string{
		// runtime services
		"config",   // ????
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"router",   // :8084
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"auth",     // :8010
		"web",      // :8082
		"bot",      // :????
		"init",     // no port, manage self
	}
)

type initScheduler struct {
	gorun.Scheduler
	services []string
}

func (i *initScheduler) Notify() (<-chan gorun.Event, error) {
	ch, err := i.Scheduler.Notify()
	if err != nil {
		return nil, err
	}

	// create new event channel
	evChan := make(chan gorun.Event, 32)

	go func() {
		for ev := range ch {
			// fire an event per service
			for _, service := range i.services {
				newEv := gorun.Event{
					Service: &gorun.Service{
						Name: service,
					},
					Timestamp: ev.Timestamp,
					Type:      ev.Type,
				}
				// Some updates don't come with version, e.g. filesystem watcher or the update notifier
				if ev.Service != nil {
					newEv.Service.Version = ev.Service.Version
				}
				evChan <- newEv
				// slow roll the change
				time.Sleep(time.Second)
			}
		}

		// we've reached the end
		close(evChan)
	}()

	return evChan, nil
}

func initNotify(n gorun.Scheduler, services []string) gorun.Scheduler {
	return &initScheduler{n, services}
}

// Init is the `micro init` command which manages the lifecycle
// of all the services. It does not start the services.
func Init(context *cli.Context) {
	log.Init(log.WithFields(map[string]interface{}{"service": "init"}))

	if context.Args().Len() > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// list of services to operate on
	initServices := Services

	// get the service prefix
	if namespace := context.String("namespace"); len(namespace) > 0 {
		for i, service := range initServices {
			initServices[i] = fmt.Sprintf("%s.%s", namespace, service)
		}
	}

	updateURL := context.String("update_url")
	if len(updateURL) == 0 {
		updateURL = update.DefaultURL
	}

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	notifier := update.NewScheduler(updateURL, Version)
	wrapped := initNotify(notifier, initServices)

	// specify with a notifier that fires
	// individual events for each service
	options := []gorun.Option{
		gorun.WithScheduler(wrapped),
		gorun.WithType("runtime"),
	}
	(*muRuntime).Init(options...)

	// used to signal when to shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, signalutil.Shutdown()...)

	log.Info("Starting service runtime")

	// start the runtime
	if err := (*muRuntime).Start(); err != nil {
		log.Fatal(err)
	}

	log.Info("Service runtime started")

	select {
	case <-shutdown:
		log.Info("Shutdown signal received")
		log.Info("Stopping service runtime")
	}

	// stop all the things
	if err := (*muRuntime).Stop(); err != nil {
		log.Fatal(err)
	}

	log.Info("Service runtime shutdown")

	// exit success
	os.Exit(0)
}
