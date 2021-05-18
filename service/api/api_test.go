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
// Original source: github.com/micro/go-micro/v3/api/api_test.go

package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	go_api "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/api"
)

func TestEncoding(t *testing.T) {
	testData := []*api.Endpoint{
		nil,
		{
			Name:        "Foo.Bar",
			Description: "A test endpoint",
			Handler:     "meta",
			Host:        []string{"foo.com"},
			Method:      []string{"GET"},
			Path:        []string{"/test"},
		},
	}

	compare := func(expect, got []string) bool {
		// no data to compare, return true
		if len(expect) == 0 && len(got) == 0 {
			return true
		}
		// no data expected but got some return false
		if len(expect) == 0 && len(got) > 0 {
			return false
		}

		// compare expected with what we got
		for _, e := range expect {
			var seen bool
			for _, g := range got {
				if e == g {
					seen = true
					break
				}
			}
			if !seen {
				return false
			}
		}

		// we're done, return true
		return true
	}

	for _, d := range testData {
		// encode
		e := api.Encode(d)
		// decode
		de := api.Decode(e)

		// nil endpoint returns nil
		if d == nil {
			if e != nil {
				t.Fatalf("expected nil got %v", e)
			}
			if de != nil {
				t.Fatalf("expected nil got %v", de)
			}

			continue
		}

		// check encoded map
		name := e["endpoint"]
		desc := e["description"]
		method := strings.Split(e["method"], ",")
		path := strings.Split(e["path"], ",")
		host := strings.Split(e["host"], ",")
		handler := e["handler"]

		if name != d.Name {
			t.Fatalf("expected %v got %v", d.Name, name)
		}
		if desc != d.Description {
			t.Fatalf("expected %v got %v", d.Description, desc)
		}
		if handler != d.Handler {
			t.Fatalf("expected %v got %v", d.Handler, handler)
		}
		if ok := compare(d.Method, method); !ok {
			t.Fatalf("expected %v got %v", d.Method, method)
		}
		if ok := compare(d.Path, path); !ok {
			t.Fatalf("expected %v got %v", d.Path, path)
		}
		if ok := compare(d.Host, host); !ok {
			t.Fatalf("expected %v got %v", d.Host, host)
		}

		if de.Name != d.Name {
			t.Fatalf("expected %v got %v", d.Name, de.Name)
		}
		if de.Description != d.Description {
			t.Fatalf("expected %v got %v", d.Description, de.Description)
		}
		if de.Handler != d.Handler {
			t.Fatalf("expected %v got %v", d.Handler, de.Handler)
		}
		if ok := compare(d.Method, de.Method); !ok {
			t.Fatalf("expected %v got %v", d.Method, de.Method)
		}
		if ok := compare(d.Path, de.Path); !ok {
			t.Fatalf("expected %v got %v", d.Path, de.Path)
		}
		if ok := compare(d.Host, de.Host); !ok {
			t.Fatalf("expected %v got %v", d.Host, de.Host)
		}
	}
}

func TestValidate(t *testing.T) {
	epPcre := &api.Endpoint{
		Name:        "Foo.Bar",
		Description: "A test endpoint",
		Handler:     "meta",
		Host:        []string{"foo.com"},
		Method:      []string{"GET"},
		Path:        []string{"^/test/?$"},
	}
	if err := api.Validate(epPcre); err != nil {
		t.Fatal(err)
	}

	epGpath := &api.Endpoint{
		Name:        "Foo.Bar",
		Description: "A test endpoint",
		Handler:     "meta",
		Host:        []string{"foo.com"},
		Method:      []string{"GET"},
		Path:        []string{"/test/{id}"},
	}
	if err := api.Validate(epGpath); err != nil {
		t.Fatal(err)
	}

	epPcreInvalid := &api.Endpoint{
		Name:        "Foo.Bar",
		Description: "A test endpoint",
		Handler:     "meta",
		Host:        []string{"foo.com"},
		Method:      []string{"GET"},
		Path:        []string{"/test/?$"},
	}
	if err := api.Validate(epPcreInvalid); err == nil {
		t.Fatalf("invalid pcre %v", epPcreInvalid.Path[0])
	}

}

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

		extByte, err := api.RequestPayload(r)
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

		extByte, err := api.RequestPayload(r)
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

		extByte, err := api.RequestPayload(r)
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

		extByte, err := api.RequestPayload(r)
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

		extByte, err := api.RequestPayload(r)
		if err != nil {
			t.Fatalf("Failed to extract payload from request: %v", err)
		}
		if string(extByte) != "" {
			t.Fatalf("Expected %v and %v to match", string(extByte), "")
		}
	})
}
