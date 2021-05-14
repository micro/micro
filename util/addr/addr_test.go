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
// Original source: github.com/micro/go-micro/v3/util/addr/addr_test.go

package addr

import (
	"net"
	"testing"
)

func TestIsLocal(t *testing.T) {
	testData := []struct {
		addr   string
		expect bool
	}{
		{"localhost", true},
		{"localhost:8080", true},
		{"127.0.0.1", true},
		{"127.0.0.1:1001", true},
		{"80.1.1.1", false},
	}

	for _, d := range testData {
		res := IsLocal(d.addr)
		if res != d.expect {
			t.Fatalf("expected %t got %t", d.expect, res)
		}
	}
}

func TestExtractor(t *testing.T) {
	testData := []struct {
		addr   string
		expect string
		parse  bool
	}{
		{"127.0.0.1", "127.0.0.1", false},
		{"10.0.0.1", "10.0.0.1", false},
		{"", "", true},
		{"0.0.0.0", "", true},
		{"[::]", "", true},
	}

	for _, d := range testData {
		addr, err := Extract(d.addr)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if d.parse {
			ip := net.ParseIP(addr)
			if ip == nil {
				t.Error("Unexpected nil IP")
			}

		} else if addr != d.expect {
			t.Errorf("Expected %s got %s", d.expect, addr)
		}
	}

}

func TestAppendPrivateBlocks(t *testing.T) {
	tests := []struct {
		addr   string
		expect bool
	}{
		{addr: "9.134.71.34", expect: true},
		{addr: "8.10.110.34", expect: false}, // not in private blocks
	}

	AppendPrivateBlocks("9.134.0.0/16")

	for _, test := range tests {
		t.Run(test.addr, func(t *testing.T) {
			res := isPrivateIP(test.addr)
			if res != test.expect {
				t.Fatalf("expected %t got %t", test.expect, res)
			}
		})
	}
}
