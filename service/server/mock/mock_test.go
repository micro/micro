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
// Original source: github.com/micro/go-micro/v3/server/mock/mock_test.go

package mock

import (
	"testing"

	"github.com/micro/micro/v3/service/server"
)

func TestMockServer(t *testing.T) {
	srv := NewServer(
		server.Name("mock"),
		server.Version("latest"),
	)

	if srv.Options().Name != "mock" {
		t.Fatalf("Expected name mock, got %s", srv.Options().Name)
	}

	if srv.Options().Version != "latest" {
		t.Fatalf("Expected version latest, got %s", srv.Options().Version)
	}

	srv.Init(server.Version("test"))
	if srv.Options().Version != "test" {
		t.Fatalf("Expected version test, got %s", srv.Options().Version)
	}

	h := srv.NewHandler(func() string { return "foo" })
	if err := srv.Handle(h); err != nil {
		t.Fatal(err)
	}

	sub := srv.NewSubscriber("test", func() string { return "foo" })
	if err := srv.Subscribe(sub); err != nil {
		t.Fatal(err)
	}

	if sub.Topic() != "test" {
		t.Fatalf("Expected topic test got %s", sub.Topic())
	}

	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	if err := srv.Register(); err != nil {
		t.Fatal(err)
	}

	if err := srv.Deregister(); err != nil {
		t.Fatal(err)
	}

	if err := srv.Stop(); err != nil {
		t.Fatal(err)
	}
}
