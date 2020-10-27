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
// Original source: github.com/micro/go-micro/v3/api/server/acme/autocert/autocert.go

// Package autocert is the ACME provider from golang.org/x/crypto/acme/autocert
// This provider does not take any config.
package autocert

import (
	"crypto/tls"
	"net"
	"os"

	"github.com/micro/micro/v3/internal/api/server/acme"
	"github.com/micro/micro/v3/service/logger"
	"golang.org/x/crypto/acme/autocert"
)

// autoCertACME is the ACME provider from golang.org/x/crypto/acme/autocert
type autocertProvider struct{}

// Listen implements acme.Provider
func (a *autocertProvider) Listen(hosts ...string) (net.Listener, error) {
	return autocert.NewListener(hosts...), nil
}

// TLSConfig returns a new tls config
func (a *autocertProvider) TLSConfig(hosts ...string) (*tls.Config, error) {
	// create a new manager
	m := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}
	if len(hosts) > 0 {
		m.HostPolicy = autocert.HostWhitelist(hosts...)
	}
	dir := cacheDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		if logger.V(logger.InfoLevel, logger.DefaultLogger) {
			logger.Infof("warning: autocert not using a cache: %v", err)
		}
	} else {
		m.Cache = autocert.DirCache(dir)
	}
	return m.TLSConfig(), nil
}

// New returns an autocert acme.Provider
func NewProvider() acme.Provider {
	return &autocertProvider{}
}
