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
// Original source: github.com/micro/go-micro/v3/api/handler/rpc/rpc_test.go

package rpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/protobuf/proto"
	go_api "github.com/micro/micro/v3/proto/api"
)

func TestRequestPayloadFromRequest(t *testing.T) {

	// our test event so that we can validate serialising / deserializing of true protos works
	protoEvent := go_api.Event{
		Name: "Test",
	}

	protoBytes, err := proto.Marshal(&protoEvent)
	if err != nil {
		t.Fatal("Failed to marshal proto", err)
	}

	jsonBytes, err := json.Marshal(protoEvent)
	if err != nil {
		t.Fatal("Failed to marshal proto to JSON ", err)
	}

	jsonUrlBytes := []byte(`{"key1":"val1","key2":"val2","name":"Test"}`)

	t.Run("extracting a json from a POST request with url params", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path?key1=val1&key2=val2", bytes.NewReader(jsonBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonUrlBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), jsonUrlBytes)
		}
	})

	t.Run("extracting a proto from a POST request", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path", bytes.NewReader(protoBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(protoBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(protoBytes))
		}
	})

	t.Run("extracting JSON from a POST request", func(t *testing.T) {
		r, err := http.NewRequest("POST", "http://localhost/my/path", bytes.NewReader(jsonBytes))
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(jsonBytes))
		}
	})

	t.Run("extracting params from a GET request", func(t *testing.T) {

		r, err := http.NewRequest("GET", "http://localhost/my/path", nil)
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		q := r.URL.Query()
		q.Add("name", "Test")
		r.URL.RawQuery = q.Encode()

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != string(jsonBytes) {
			t.Fatalf("Expected %v and %v to match", string(extByte), string(jsonBytes))
		}
	})

	t.Run("GET request with no params", func(t *testing.T) {

		r, err := http.NewRequest("GET", "http://localhost/my/path", nil)
		if err != nil {
			t.Fatalf("Failed to created http.Request: %v", err)
		}

		extByte, err := requestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != "" {
			t.Fatalf("Expected %v and %v to match", string(extByte), "")
		}
	})
}
