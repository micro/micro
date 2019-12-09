// Package stats provides a service that collects stats from all services in the registry.
package stats

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"github.com/micro/micro/debug/stats/handler"
	stats "github.com/micro/micro/debug/stats/proto"
)

// Run is the entrypoint for debug/stats
func Run(c *cli.Context) {
	statsService := micro.NewService(
		micro.Name("go.micro.debug.stats"),
	)

	// Create handler
	h, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register Handler
	stats.RegisterStatsHandler(statsService.Server(), h)

	// Run service
	if err := statsService.Run(); err != nil {
		log.Fatal(err)
	}
}
