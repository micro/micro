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
// Original source: github.com/micro/go-micro/v3/util/pool/default_test.go

package pool

import (
	"testing"
	"time"

	"github.com/micro/micro/v3/internal/network/transport"
	"github.com/micro/micro/v3/internal/network/transport/memory"
)

func testPool(t *testing.T, size int, ttl time.Duration) {
	// mock transport
	tr := memory.NewTransport()

	options := Options{
		TTL:       ttl,
		Size:      size,
		Transport: tr,
	}
	// zero pool
	p := newPool(options)

	// listen
	l, err := tr.Listen(":0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	// accept loop
	go func() {
		for {
			if err := l.Accept(func(s transport.Socket) {
				for {
					var msg transport.Message
					if err := s.Recv(&msg); err != nil {
						return
					}
					if err := s.Send(&msg); err != nil {
						return
					}
				}
			}); err != nil {
				return
			}
		}
	}()

	for i := 0; i < 10; i++ {
		// get a conn
		c, err := p.Get(l.Addr())
		if err != nil {
			t.Fatal(err)
		}

		msg := &transport.Message{
			Body: []byte(`hello world`),
		}

		if err := c.Send(msg); err != nil {
			t.Fatal(err)
		}

		var rcv transport.Message

		if err := c.Recv(&rcv); err != nil {
			t.Fatal(err)
		}

		if string(rcv.Body) != string(msg.Body) {
			t.Fatalf("got %v, expected %v", rcv.Body, msg.Body)
		}

		// release the conn
		p.Release(c, nil)

		p.Lock()
		if i := len(p.conns[l.Addr()]); i > size {
			p.Unlock()
			t.Fatalf("pool size %d is greater than expected %d", i, size)
		}
		p.Unlock()
	}
}

func TestClientPool(t *testing.T) {
	testPool(t, 0, time.Minute)
	testPool(t, 2, time.Minute)
}
