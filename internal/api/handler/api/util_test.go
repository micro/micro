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
// Original source: github.com/micro/go-micro/v3/api/handler/util_test.go

package api

import (
	"net/http"
	"net/url"
	"testing"
)

func TestRequestToProto(t *testing.T) {
	testData := []*http.Request{
		{
			Method: "GET",
			Header: http.Header{
				"Header": []string{"test"},
			},
			URL: &url.URL{
				Scheme:   "http",
				Host:     "localhost",
				Path:     "/foo/bar",
				RawQuery: "param1=value1",
			},
		},
	}

	for _, d := range testData {
		p, err := requestToProto(d)
		if err != nil {
			t.Fatal(err)
		}
		if p.Path != d.URL.Path {
			t.Fatalf("Expected path %s got %s", d.URL.Path, p.Path)
		}
		if p.Method != d.Method {
			t.Fatalf("Expected method %s got %s", d.Method, p.Method)
		}
		for k, v := range d.Header {
			if val, ok := p.Header[k]; !ok {
				t.Fatalf("Expected header %s", k)
			} else {
				if val.Values[0] != v[0] {
					t.Fatalf("Expected val %s, got %s", val.Values[0], v[0])
				}
			}
		}
	}
}
