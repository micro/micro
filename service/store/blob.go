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
// Original source: github.com/micro/go-micro/v3/store/blob.go

package store

import (
	"errors"
	"io"
)

var (
	// ErrMissingKey is returned when no key is passed to blob store Read / Write
	ErrMissingKey = errors.New("missing key")
)

// BlobStore is an interface for reading / writing blobs
type BlobStore interface {
	Read(key string, opts ...BlobOption) (io.Reader, error)
	Write(key string, blob io.Reader, opts ...BlobOption) error
	Delete(key string, opts ...BlobOption) error
	// List returns any keys that match, or an empty list with no error if none matched.
	List(opts ...BlobListOption) ([]string, error)
}

// BlobOptions contains options to use when interacting with the store
type BlobOptions struct {
	// Namespace to  from
	Namespace   string
	Public      bool
	ContentType string
}

// BlobOption sets one or more BlobOptions
type BlobOption func(o *BlobOptions)

// BlobNamespace sets the Namespace option
func BlobNamespace(ns string) BlobOption {
	return func(o *BlobOptions) {
		o.Namespace = ns
	}
}

// BlobNamespace sets the Public option
func BlobPublic(p bool) BlobOption {
	return func(o *BlobOptions) {
		o.Public = p
	}
}

// BlobNamespace sets the Public option
func BlobContentType(contentType string) BlobOption {
	return func(o *BlobOptions) {
		o.ContentType = contentType
	}
}

type BlobListOptions struct {
	Namespace string
	Prefix    string
}

type BlobListOption func(o *BlobListOptions)

func BlobListNamespace(namespace string) BlobListOption {
	return func(o *BlobListOptions) {
		o.Namespace = namespace
	}
}

func BlobListPrefix(prefix string) BlobListOption {
	return func(o *BlobListOptions) {
		o.Prefix = prefix
	}
}
