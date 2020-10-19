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
// Original source: github.com/micro/go-micro/v3/util/grpc/grpc_test.go

package grpc

import (
	"testing"
)

func TestServiceMethod(t *testing.T) {
	type testCase struct {
		input   string
		service string
		method  string
		err     bool
	}

	methods := []testCase{
		{"Foo.Bar", "Foo", "Bar", false},
		{"/Foo/Bar", "Foo", "Bar", false},
		{"/package.Foo/Bar", "Foo", "Bar", false},
		{"/a.package.Foo/Bar", "Foo", "Bar", false},
		{"a.package.Foo/Bar", "", "", true},
		{"/Foo/Bar/Baz", "", "", true},
		{"Foo.Bar.Baz", "", "", true},
	}
	for _, test := range methods {
		service, method, err := ServiceMethod(test.input)
		if err != nil && test.err == true {
			continue
		}
		// unexpected error
		if err != nil && test.err == false {
			t.Fatalf("unexpected err %v for %+v", err, test)
		}
		// expecter error
		if test.err == true && err == nil {
			t.Fatalf("expected error for %+v: got service: %s method: %s", test, service, method)
		}

		if service != test.service {
			t.Fatalf("wrong service for %+v: got service: %s method: %s", test, service, method)
		}

		if method != test.method {
			t.Fatalf("wrong method for %+v: got service: %s method: %s", test, service, method)
		}
	}
}
