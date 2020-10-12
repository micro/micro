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
// Original source: github.com/micro/go-micro/v3/api/server/acme/acme.go

// Package acme abstracts away various ACME libraries
package acme

import (
	"crypto/tls"
	"errors"
	"net"
)

var (
	// ErrProviderNotImplemented can be returned when attempting to
	// instantiate an unimplemented provider
	ErrProviderNotImplemented = errors.New("Provider not implemented")
)

// Provider is a ACME provider interface
type Provider interface {
	// Listen returns a new listener
	Listen(...string) (net.Listener, error)
	// TLSConfig returns a tls config
	TLSConfig(...string) (*tls.Config, error)
}

// The Let's Encrypt ACME endpoints
const (
	LetsEncryptStagingCA    = "https://acme-staging-v02.api.letsencrypt.org/directory"
	LetsEncryptProductionCA = "https://acme-v02.api.letsencrypt.org/directory"
)
