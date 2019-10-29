// Package dns provides a DNS registration service for autodiscovery of core network nodes.
package dns

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	"github.com/micro/micro/network/dns/handler"
	dns "github.com/micro/micro/network/dns/proto/dns"
)

// Run is the entrypoint for network/dns
func Run(c *cli.Context) {
	dnsService := micro.NewService(
		micro.Name("go.micro.network.dns"),
	)

	// Create handler
	h := handler.New()

	// Register Handler
	dns.RegisterDnsHandler(dnsService.Server(), h)

	// Run service
	if err := dnsService.Run(); err != nil {
		log.Fatal(err)
	}

}
