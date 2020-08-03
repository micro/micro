// Package broker is the micro broker
package broker

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/micro/v3/internal/cmd"
	"github.com/pkg/errors"
)

func init() {
	cmd.Init(func(ctx *cli.Context) error {

		// Setup broker options.
		opts := []broker.Option{}
		if len(ctx.String("broker_address")) > 0 {
			opts = append(opts, broker.Addrs(ctx.String("broker_address")))
		}

		// Parse broker TLS certs
		if len(ctx.String("broker_tls_cert")) > 0 || len(ctx.String("broker_tls_key")) > 0 {
			cert, err := tls.LoadX509KeyPair(ctx.String("broker_tls_cert"), ctx.String("broker_tls_key"))
			if err != nil {
				errors.Wrap(err, "Error loading broker TLS cert")
			}

			// load custom certificate authority
			caCertPool := x509.NewCertPool()
			if len(ctx.String("broker_tls_ca")) > 0 {
				crt, err := ioutil.ReadFile(ctx.String("broker_tls_ca"))
				if err != nil {
					errors.Wrap(err, "Error loading broker TLS certificate authority")
				}
				caCertPool.AppendCertsFromPEM(crt)
			}

			cfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: caCertPool}
			opts = append(opts, broker.TLSConfig(cfg))
		}

		return DefaultBroker.Init(opts...)
	})
}

// DefaultBroker implementation
var DefaultBroker broker.Broker

// Publish a message to a topic
func Publish(topic string, m *broker.Message, opts ...broker.PublishOption) error {
	return DefaultBroker.Publish(topic, m, opts...)
}

// Subscribe to a topic
func Subscribe(topic string, h broker.Handler, opts ...broker.SubscribeOption) (broker.Subscriber, error) {
	return DefaultBroker.Subscribe(topic, h, opts...)
}
