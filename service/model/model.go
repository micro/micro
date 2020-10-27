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
// Original source: github.com/micro/go-micro/v3/model/model.go

// Package model is an interface for data modelling
package model

import (
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/internal/sync"
	"github.com/micro/micro/v3/service/store"
)

var (
	// DefaultModel for the service
	DefaultModel Model
)

// Model provides an interface for data modelling
type Model interface {
	// Initialise options
	Init(...Option) error
	// NewEntity creates a new entity to store or access
	NewEntity(name string, value interface{}) Entity
	// Create a value
	Create(Entity) error
	// Read values
	Read(...ReadOption) ([]Entity, error)
	// Update the value
	Update(Entity) error
	// Delete an entity
	Delete(...DeleteOption) error
	// Implementation of the model
	String() string
}

type Entity interface {
	// Unique id of the entity
	Id() string
	// Name of the entity
	Name() string
	// The value associated with the entity
	Value() interface{}
	// Attributes of the entity
	Attributes() map[string]interface{}
}

type Options struct {
	// Database to write to
	Database string
	// for serialising
	Codec codec.Marshaler
	// for locking
	Sync sync.Sync
	// for storage
	Store store.Store
}

type Option func(o *Options)

type ReadOptions struct{}

type ReadOption func(o *ReadOptions)

type DeleteOptions struct{}

type DeleteOption func(o *DeleteOptions)
