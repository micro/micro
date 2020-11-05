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
// Original source: github.com/micro/go-micro/v3/util/kubernetes/client/util_test.go

package client

import (
	"bytes"
	"testing"

	"github.com/micro/micro/v3/service/runtime"
)

func TestTemplates(t *testing.T) {
	srv := &runtime.Service{Name: "foo", Version: "123"}
	opts := &runtime.CreateOptions{Type: "service", Namespace: "default"}

	// Render default service
	s := NewService(srv, opts)
	bs := new(bytes.Buffer)
	if err := renderTemplate(templates["service"], bs, s); err != nil {
		t.Errorf("Failed to render kubernetes service: %v", err)
	}

	// Render default deployment
	d := NewDeployment(srv, opts)
	bd := new(bytes.Buffer)
	if err := renderTemplate(templates["deployment"], bd, d); err != nil {
		t.Errorf("Failed to render kubernetes deployment: %v", err)
	}
}

func TestFormatName(t *testing.T) {
	testCases := []struct {
		name   string
		expect string
	}{
		{"foobar", "foobar"},
		{"foo-bar", "foo-bar"},
		{"foo.bar", "foo-bar"},
		{"Foo.Bar", "foo-bar"},
		{"go.micro.foo.bar", "go-micro-foo-bar"},
		{"go.micro.foo.bar", "go-micro-foo-bar"},
		{"foo/bar_baz", "foo-bar-baz"},
	}

	for _, test := range testCases {
		v := Format(test.name)
		if v != test.expect {
			t.Fatalf("Expected name %s for %s got: %s", test.expect, test.name, v)
		}
	}
}
