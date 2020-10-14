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
// Original source: github.com/micro/go-micro/v3/server/mucp/extractor_test.go

package mucp

import (
	"context"
	"reflect"
	"testing"

	"github.com/micro/micro/v3/service/registry"
)

type testHandler struct{}

type testRequest struct{}

type testResponse struct{}

func (t *testHandler) Test(ctx context.Context, req *testRequest, rsp *testResponse) error {
	return nil
}

func TestExtractEndpoint(t *testing.T) {
	handler := &testHandler{}
	typ := reflect.TypeOf(handler)

	var endpoints []*registry.Endpoint

	for m := 0; m < typ.NumMethod(); m++ {
		if e := extractEndpoint(typ.Method(m)); e != nil {
			endpoints = append(endpoints, e)
		}
	}

	if i := len(endpoints); i != 1 {
		t.Errorf("Expected 1 endpoint, have %d", i)
	}

	if endpoints[0].Name != "Test" {
		t.Errorf("Expected handler Test, got %s", endpoints[0].Name)
	}

	if endpoints[0].Request == nil {
		t.Error("Expected non nil request")
	}

	if endpoints[0].Response == nil {
		t.Error("Expected non nil request")
	}

	if endpoints[0].Request.Name != "testRequest" {
		t.Errorf("Expected testRequest got %s", endpoints[0].Request.Name)
	}

	if endpoints[0].Response.Name != "testResponse" {
		t.Errorf("Expected testResponse got %s", endpoints[0].Response.Name)
	}

	if endpoints[0].Request.Type != "testRequest" {
		t.Errorf("Expected testRequest type got %s", endpoints[0].Request.Type)
	}

	if endpoints[0].Response.Type != "testResponse" {
		t.Errorf("Expected testResponse type got %s", endpoints[0].Response.Type)
	}

}
