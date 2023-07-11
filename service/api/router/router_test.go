// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/api/router/router_test.go
package router_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"micro.dev/v4/service/client"
	gcli "micro.dev/v4/service/client/grpc"
	rmemory "micro.dev/v4/service/registry/memory"
	rt "micro.dev/v4/service/router"
	regRouter "micro.dev/v4/service/router/registry"
	"micro.dev/v4/service/server"
	gsrv "micro.dev/v4/service/server/grpc"
	pb "micro.dev/v4/service/server/grpc/proto"
)

// server is used to implement helloworld.GreeterServer.
type testServer struct {
	msgCount int
}

// TestHello implements helloworld.GreeterServer
func (s *testServer) Call(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	rsp.Msg = "Hello " + req.Uuid
	return nil
}

func initial(t *testing.T) (server.Server, client.Client) {
	r := rmemory.NewRegistry()

	// create a new client
	s := gsrv.NewServer(
		server.Name("foo"),
		server.Registry(r),
	)

	rtr := regRouter.NewRouter(
		rt.Registry(r),
	)

	// create a new server
	c := gcli.NewClient(
		client.Router(rtr),
	)

	h := &testServer{}
	pb.RegisterTestHandler(s, h)

	if err := s.Start(); err != nil {
		t.Fatalf("failed to start: %v", err)
	}

	return s, c
}

func check(t *testing.T, addr string, path string, expected string) {
	req, err := http.NewRequest("POST", fmt.Sprintf(path, addr), nil)
	if err != nil {
		t.Fatalf("Failed to created http.Request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatalf("Failed to created http.Request: %v", err)
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	jsonMsg := expected
	if string(buf) != jsonMsg {
		t.Fatalf("invalid message received, parsing error %s != %s", buf, jsonMsg)
	}
}
