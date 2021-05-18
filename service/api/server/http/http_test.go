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
// Original source: github.com/micro/go-micro/v3/api/server/http/http_test.go

package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHTTPServer(t *testing.T) {
	testResponse := "hello world"

	s := NewServer("localhost:0")

	s.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, testResponse)
	}))

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	rsp, err := http.Get(fmt.Sprintf("http://%s/", s.Address()))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != testResponse {
		t.Fatalf("Unexpected response, got %s, expected %s", string(b), testResponse)
	}

	if err := s.Stop(); err != nil {
		t.Fatal(err)
	}
}
