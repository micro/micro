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
// Original source: github.com/micro/go-micro/v3/client/grpc/request_test.go

package grpc

import (
	"testing"
)

func TestMethodToGRPC(t *testing.T) {
	testData := []struct {
		service string
		method  string
		expect  string
	}{
		{
			"helloworld",
			"Greeter.SayHello",
			"/helloworld.Greeter/SayHello",
		},
		{
			"helloworld",
			"/helloworld.Greeter/SayHello",
			"/helloworld.Greeter/SayHello",
		},
		{
			"",
			"/helloworld.Greeter/SayHello",
			"/helloworld.Greeter/SayHello",
		},
		{
			"",
			"Greeter.SayHello",
			"/Greeter/SayHello",
		},
	}

	for _, d := range testData {
		method := methodToGRPC(d.service, d.method)
		if method != d.expect {
			t.Fatalf("expected %s got %s", d.expect, method)
		}
	}
}
