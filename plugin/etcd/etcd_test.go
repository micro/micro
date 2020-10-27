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
// Original source: github.com/micro/go-plugins/v3/registry/etcd/etcd_test.go

package etcd

import (
	"testing"
)

// test whether the name matches
func TestEtcdHasName(t *testing.T) {
	testCases := []struct {
		key    string
		prefix string
		name   string
		domain string
		expect bool
	}{
		{
			"/micro/registry/micro/registry",
			"/micro/registry",
			"registry",
			"micro",
			true,
		},
		{
			"/micro/registry/micro",
			"/micro/registry",
			"store",
			"micro",
			false,
		},
		{
			"/prefix/baz/*/registry",
			"/prefix/baz",
			"registry",
			"*",
			true,
		},
		{
			"/prefix/baz",
			"/prefix/baz",
			"store",
			"micro",
			false,
		},
		{
			"/prefix/baz/foobar/registry",
			"/prefix/baz",
			"registry",
			"foobar",
			true,
		},
	}

	for _, c := range testCases {
		domain, service, ok := getName(c.key, c.prefix)
		if ok != c.expect {
			t.Fatalf("Expected %t for %v got: %t", c.expect, c, ok)
		}
		if !ok {
			continue
		}
		if service != c.name {
			t.Fatalf("Expected service %s got %s", c.name, service)
		}
		if domain != c.domain {
			t.Fatalf("Expected domain %s got %s", c.domain, domain)
		}
	}
}
