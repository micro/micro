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
// Original source: github.com/micro/go-micro/v3/util/token/options.go

package token

import (
	"time"

	"github.com/micro/micro/v3/service/store"
)

type Options struct {
	// Store to persist the tokens
	Store store.Store
	// PublicKey base64 encoded, used by JWT
	PublicKey string
	// PrivateKey base64 encoded, used by JWT
	PrivateKey string
}

type Option func(o *Options)

// WithStore sets the token providers store
func WithStore(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// WithPublicKey sets the JWT public key
func WithPublicKey(key string) Option {
	return func(o *Options) {
		o.PublicKey = key
	}
}

// WithPrivateKey sets the JWT private key
func WithPrivateKey(key string) Option {
	return func(o *Options) {
		o.PrivateKey = key
	}
}

func NewOptions(opts ...Option) Options {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	//set default store
	if options.Store == nil {
		options.Store = store.DefaultStore
	}
	return options
}

type GenerateOptions struct {
	// Expiry for the token
	Expiry time.Duration
}

type GenerateOption func(o *GenerateOptions)

// WithExpiry for the generated account's token expires
func WithExpiry(d time.Duration) GenerateOption {
	return func(o *GenerateOptions) {
		o.Expiry = d
	}
}

// NewGenerateOptions from a slice of options
func NewGenerateOptions(opts ...GenerateOption) GenerateOptions {
	var options GenerateOptions
	for _, o := range opts {
		o(&options)
	}
	//set default Expiry of token
	if options.Expiry == 0 {
		options.Expiry = time.Minute * 15
	}
	return options
}
