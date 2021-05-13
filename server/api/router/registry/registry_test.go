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
// Original source: github.com/micro/go-micro/v3/api/router/registry/registry_test.go

package registry

import (
	"testing"

	"github.com/micro/micro/v3/service/registry"
	"github.com/stretchr/testify/assert"
)

func TestStoreRegex(t *testing.T) {
	router := newRouter()
	router.store("micro", []*registry.Service{
		{
			Name:    "Foobar",
			Version: "latest",
			Endpoints: []*registry.Endpoint{
				{
					Name: "foo",
					Metadata: map[string]string{
						"endpoint":    "FooEndpoint",
						"description": "Some description",
						"method":      "POST",
						"path":        "^/foo/$",
						"handler":     "rpc",
					},
				},
			},
			Metadata: map[string]string{},
		},
	},
	)

	assert.Len(t, router.namespaces["micro"].ceps["Foobar.foo"].pcreregs, 1)
}
