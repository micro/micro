package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/micro/cli"
	"github.com/micro/go-micro/config/cmd"
	gorun "github.com/micro/go-micro/runtime"
	"github.com/micro/go-micro/util/log"

	// include usage

	"github.com/micro/micro/internal/update"
	_ "github.com/micro/micro/internal/usage"
)

type initNotifier struct {
	gorun.Notifier
	services []string
}

func (i *initNotifier) Notify() (<-chan gorun.Event, error) {
	ch, err := i.Notifier.Notify()
	if err != nil {
		return nil, err
	}

	// create new event channel
	evChan := make(chan gorun.Event, 32)

	go func() {
		for ev := range ch {
			// fire an event per service
			for _, service := range i.services {
				evChan <- gorun.Event{
					Service:   service,
					Version:   ev.Version,
					Timestamp: ev.Timestamp,
					Type:      ev.Type,
				}
			}
		}

		// we've reached the end
		close(evChan)
	}()

	return evChan, nil
}

func initNotify(n gorun.Notifier, services []string) gorun.Notifier {
	return &initNotifier{n, services}
}

func initCommand(context *cli.Context) {
	log.Name("init")

	if len(context.Args()) > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// services to manage
	services := []string{
		// network services
		"network.api",
		"network.dns",
		"network.web",
		"debug.web",
		// runtime services
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"tunnel",   // :8083
		"router",   // :8084
		"monitor",  // :????
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"web",      // :8082
		"bot",      // :????
		"init",     // no port, manage self
	}

	// get the service prefix
	if namespace := context.GlobalString("namespace"); len(namespace) > 0 {
		for i, service := range services {
			services[i] = fmt.Sprintf("%s.%s", namespace, service)
		}
	}

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	notifier := update.NewNotifier(BuildDate)
	wrapped := initNotify(notifier, services)

	// specify with a notifier that fires
	// individual events for each service
	options := []gorun.Option{
		gorun.WithNotifier(wrapped),
		gorun.WithType("runtime"),
	}
	(*muRuntime).Init(options...)

	// used to signal when to shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

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
