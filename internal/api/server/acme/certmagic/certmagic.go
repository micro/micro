// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/api/server/acme/certmagic/certmagic.go

// Package certmagic is the ACME provider from github.com/caddyserver/certmagic
package certmagic

import (
	"crypto/tls"
	"math/rand"
	"net"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/micro/micro/v3/internal/api/server/acme"
	"github.com/micro/micro/v3/service/logger"
)

type certmagicProvider struct {
	opts acme.Options
}

// TODO: set self-contained options
func (c *certmagicProvider) setup() {
	certmagic.DefaultACME.CA = c.opts.CA
	if c.opts.ChallengeProvider != nil {
		// Enabling DNS Challenge disables the other challenges
		certmagic.DefaultACME.DNSProvider = c.opts.ChallengeProvider
	}
	if c.opts.OnDemand {
		certmagic.Default.OnDemand = new(certmagic.OnDemandConfig)
	}
	if c.opts.Cache != nil {
		// already validated by new()
		certmagic.Default.Storage = c.opts.Cache.(certmagic.Storage)
	}
	// If multiple instances of the provider are running, inject some
	// randomness so they don't collide
	// RenewalWindowRatio [0.33 - 0.50)
	rand.Seed(time.Now().UnixNano())
	randomRatio := float64(rand.Intn(17)+33) * 0.01
	certmagic.Default.RenewalWindowRatio = randomRatio
}

func (c *certmagicProvider) Listen(hosts ...string) (net.Listener, error) {
	c.setup()
	return certmagic.Listen(hosts)
}

func (c *certmagicProvider) TLSConfig(hosts ...string) (*tls.Config, error) {
	c.setup()
	return certmagic.TLS(hosts)
}

// NewProvider returns a certmagic provider
func NewProvider(options ...acme.Option) acme.Provider {
	opts := acme.DefaultOptions()

	for _, o := range options {
		o(&opts)
	}

	if opts.Cache != nil {
		if _, ok := opts.Cache.(certmagic.Storage); !ok {
			logger.Fatal("ACME: cache provided doesn't implement certmagic's Storage interface")
		}
	}

	return &certmagicProvider{
		opts: opts,
	}
}
