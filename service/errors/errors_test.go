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
// Original source: github.com/micro/go-micro/v3/errors/errors_test.go
package errors

import (
	er "errors"
	"net/http"
	"testing"
)

func TestFromError(t *testing.T) {
	err := NotFound("go.micro.test", "%s", "example")
	merr := FromError(err)
	if merr.id != "go.micro.test" || merr.code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}
	err = er.New(err.Error())
	merr = FromError(err)
	if merr.id != "go.micro.test" || merr.code != 404 {
		t.Fatalf("invalid conversation %v != %v", err, merr)
	}

}

func TestEqual(t *testing.T) {
	err1 := NotFound("myid1", "msg1")
	err2 := NotFound("myid2", "msg2")

	if !Equal(err1, err2) {
		t.Fatal("errors must be equal")
	}

	err3 := er.New("my test err")
	if Equal(err1, err3) {
		t.Fatal("errors must be not equal")
	}

}

func TestErrors(t *testing.T) {
	testData := []*Error{
		{
			id:     "test",
			code:   500,
			detail: "Internal server error",
			status: http.StatusText(500),
		},
	}

	for _, e := range testData {
		ne := New(e.id, e.detail, e.code)

		if e.Error() != ne.Error() {
			t.Fatalf("Expected %s got %s", e.Error(), ne.Error())
		}

		pe := Parse(ne.Error())

		if pe == nil {
			t.Fatalf("Expected error got nil %v", pe)
		}

		if pe.id != e.id {
			t.Fatalf("Expected %s got %s", e.id, pe.id)
		}

		if pe.detail != e.detail {
			t.Fatalf("Expected %s got %s", e.detail, pe.detail)
		}

		if pe.code != e.code {
			t.Fatalf("Expected %d got %d", e.code, pe.code)
		}

		if pe.status != e.status {
			t.Fatalf("Expected %s got %s", e.status, pe.status)
		}
	}
}
