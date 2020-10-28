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
// Original source: github.com/micro/go-micro/v3/util/net/net_test.go

package net

import (
	"net"
	"testing"
)

func TestListen(t *testing.T) {
	fn := func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	}

	// try to create a number of listeners
	for i := 0; i < 10; i++ {
		l, err := Listen("localhost:10000-11000", fn)
		if err != nil {
			t.Fatal(err)
		}
		defer l.Close()
	}

	// TODO nats case test
	// natsAddr := "_INBOX.bID2CMRvlNp0vt4tgNBHWf"
	// Expect addr DO NOT has extra ":" at the end!

}
