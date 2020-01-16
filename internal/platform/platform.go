// Package platform manages the runtime services as a platform
package platform

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

var (
	// Date of the build
	// TODO: move elsewhere
	Version string

	// list of services managed
	services = []string{
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
	log.Name("init")

	if len(context.Args()) > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// create the combined list of services
	initServices := append(services, dashboards...)
	initServices = append(services, apis...)

	// get the service prefix
	if namespace := context.GlobalString("namespace"); len(namespace) > 0 {
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

// Run runs the entire platform
func Run(context *cli.Context) {
	log.Name("micro")

	if len(context.Args()) > 0 {
		cli.ShowSubcommandHelp(context)
		os.Exit(1)
	}

	// get the network flag
	network := context.GlobalString("network")
	local := context.GlobalBool("local")

	// pass through the environment
	// TODO: perhaps don't do this
	env := os.Environ()

	if network == "local" || local {
		// no op for now
		log.Info("Setting local network")
	} else {
		log.Info("Setting global network")

		if v := os.Getenv("MICRO_NETWORK_NODES"); len(v) == 0 {
			// set the resolver to use https://micro.mu/network
			env = append(env, "MICRO_NETWORK_NODES=network.micro.mu")
			log.Log("Setting default network micro.mu")
		}
		if v := os.Getenv("MICRO_NETWORK_TOKEN"); len(v) == 0 {
			// set the network token
			env = append(env, "MICRO_NETWORK_TOKEN=micro.mu")
			log.Log("Setting default network token")
		}
	}

	log.Info("Loading core services")

	// create new micro runtime
	muRuntime := cmd.DefaultCmd.Options().Runtime

	// Use default update notifier
	if context.GlobalBool("auto_update") {
		options := []gorun.Option{
			gorun.WithScheduler(update.NewScheduler(Version)),
		}
		(*muRuntime).Init(options...)
	}

	for _, service := range services {
		name := service

		if namespace := context.GlobalString("namespace"); len(namespace) > 0 {
			name = fmt.Sprintf("%s.%s", namespace, service)
		}

		log.Infof("Registering %s", name)

		// runtime based on environment we run the service in
		args := []gorun.CreateOption{
			gorun.WithCommand(os.Args[0], service),
			gorun.WithEnv(env),
			gorun.WithOutput(os.Stdout),
		}

		// NOTE: we use Version right now to check for the latest release
		muService := &gorun.Service{Name: name, Version: Version}
		if err := (*muRuntime).Create(muService, args...); err != nil {
			log.Errorf("Failed to create runtime enviroment: %v", err)
			return
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Info("Starting service runtime")

	// start the runtime
	if err := (*muRuntime).Start(); err != nil {
		log.Fatal(err)
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
	}

	log.Info("Service runtime shutdown")

	// exit success
	os.Exit(0)
}
