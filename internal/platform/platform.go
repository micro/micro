// Package platform manages the runtime services as a platform
package platform

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	gorun "github.com/micro/go-micro/v2/runtime"
	signalutil "github.com/micro/go-micro/v2/util/signal"

	// include usage

	"github.com/micro/micro/v2/internal/update"
	_ "github.com/micro/micro/v2/internal/usage"
)

var (
	// Date of the build
	// TODO: move elsewhere
	Version string

	// list of services managed
	services = []string{
		// runtime services
		"config",   // ????
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"tunnel",   // :8083
		"router",   // :8084
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"auth",     // :8010
		"web",      // :8082
		"bot",      // :????
		"init",     // no port, manage self
	}

	// list of web apps
	dashboards = []string{
		"network.web",
		"debug.web",
	}

	// list of apis
	apis = []string{
		"network.dns",
		"network.api",
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
				evChan <- gorun.Event{
					Service:   service,
					Version:   ev.Version,
					Timestamp: ev.Timestamp,
					Type:      ev.Type,
				}
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

	// create the combined list of services
	initServices := append(dashboards, apis...)
	initServices = append(initServices, services...)

	// get the service prefix
	if namespace := context.String("namespace"); len(namespace) > 0 {
		for i, service := range initServices {
			initServices[i] = fmt.Sprintf("%s.%s", namespace, service)
		}
	}

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	notifier := update.NewScheduler(Version)
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

// Run runs the entire platform
func Run(context *cli.Context) error {
	log.Init(log.WithFields(map[string]interface{}{"service": "micro"}))

	if context.Args().Len() > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// get the network flag
	local := context.Bool("local")
	peer := context.Bool("peer")

	// pass through the environment
	// TODO: perhaps don't do this
	env := os.Environ()

	// check either the peer or local flags are set
	// otherwise just return the hel
	if !peer && !local {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// connect to the network if specified
	if peer || !local {
		log.Info("Setting global network")

		if v := os.Getenv("MICRO_NETWORK_NODES"); len(v) == 0 {
			// set the resolver to use https://micro.mu/network
			env = append(env, "MICRO_NETWORK_RESOLVER=http")
			log.Info("Setting default network micro.mu")
		}
		if v := os.Getenv("MICRO_NETWORK_TOKEN"); len(v) == 0 {
			// set the network token
			env = append(env, "MICRO_NETWORK_TOKEN=micro.mu")
			log.Info("Setting default network token")
		}
	}

	log.Info("Loading core services")

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	if context.Bool("auto_update") {
		options := []gorun.Option{
			gorun.WithScheduler(update.NewScheduler(Version)),
		}
		(*muRuntime).Init(options...)
	}

	for _, service := range services {
		name := service

		if namespace := context.String("namespace"); len(namespace) > 0 {
			name = fmt.Sprintf("%s.%s", namespace, service)
		}

		log.Infof("Registering %s", name)

		// runtime based on environment we run the service in
		args := []gorun.CreateOption{
			gorun.WithCommand(os.Args[0]),
			gorun.WithArgs(service),
			gorun.WithEnv(env),
			gorun.WithOutput(os.Stdout),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &gorun.Service{Name: name, Version: Version}
		if err := (*muRuntime).Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime enviroment: %v", err)
			return err
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, signalutil.Shutdown()...)

	log.Info("Starting service runtime")

	// start the runtime
	if err := (*muRuntime).Start(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime started")

	// TODO: should we launch the console?
	// start the console
	// cli.Init(context)

	select {
	case <-shutdown:
		log.Info("Shutdown signal received")
	}

	log.Info("Stopping service runtime")

	// stop all the things
	if err := (*muRuntime).Stop(); err != nil {
		log.Fatal(err)
		return err
	}

	log.Info("Service runtime shutdown")

	// exit success
	os.Exit(0)
	return nil
}
