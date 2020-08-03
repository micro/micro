// Package registry is the micro registry
package registry

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/micro/v3/internal/cmd"
	"github.com/micro/micro/v3/service/registry/client"
	"github.com/pkg/errors"
)

func init() {
	cmd.Init(func(ctx *cli.Context) error {
		var opts []registry.Option

		// Parse registry TLS certs
		if len(ctx.String("registry_tls_cert")) > 0 || len(ctx.String("registry_tls_key")) > 0 {
			cert, err := tls.LoadX509KeyPair(ctx.String("registry_tls_cert"), ctx.String("registry_tls_key"))
			if err != nil {
				return errors.Wrap(err, "Error loading registry tls cert")
			}

			// load custom certificate authority
			caCertPool := x509.NewCertPool()
			if len(ctx.String("registry_tls_ca")) > 0 {
				crt, err := ioutil.ReadFile(ctx.String("registry_tls_ca"))
				if err != nil {
					return errors.Wrap(err, "Error loading registry tls certificate authority")
				}
				caCertPool.AppendCertsFromPEM(crt)
			}

			cfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool}
			opts = append(opts, registry.TLSConfig(cfg))
		}

		if len(ctx.String("registry_address")) > 0 {
			addresses := strings.Split(ctx.String("registry_address"), ",")
			opts = append(opts, registry.Addrs(addresses...))
		}

		return DefaultRegistry.Init(opts...)
	})
}

var (
	// DefaultRegistry implementation
	DefaultRegistry registry.Registry = client.NewRegistry()
)

// Register a service
func Register(service *registry.Service, opts ...registry.RegisterOption) error {
	return DefaultRegistry.Register(service, opts...)
}

// Deregister a service
func Deregister(service *registry.Service, opts ...registry.DeregisterOption) error {
	return DefaultRegistry.Deregister(service, opts...)
}

// GetService from the registry
func GetService(service string, opts ...registry.GetOption) ([]*registry.Service, error) {
	return DefaultRegistry.GetService(service, opts...)
}

// ListServices in the registry
func ListServices(opts ...registry.ListOption) ([]*registry.Service, error) {
	return DefaultRegistry.ListServices(opts...)
}

// Watch the registry for updates
func Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}
