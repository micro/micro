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
// Original source: github.com/micro/go-micro/v3/internal/ctx/ctx_test.go

package ctx

import (
	"net/http"
	"testing"

	"github.com/micro/micro/v3/service/context/metadata"
)

func TestRequestToContext(t *testing.T) {
	testData := []struct {
		request *http.Request
		expect  metadata.Metadata
	}{
		{
			&http.Request{
				Header: http.Header{
					"Foo1": []string{"bar"},
					"Foo2": []string{"bar", "baz"},
				},
			},
			metadata.Metadata{
				"Foo1": "bar",
				"Foo2": "bar,baz",
			},
		},
	}

	for _, d := range testData {
		ctx := FromRequest(d.request)
		md, ok := metadata.FromContext(ctx)
		if !ok {
			t.Fatalf("Expected metadata for request %+v", d.request)
		}
		for k, v := range d.expect {
			if val := md[k]; val != v {
				t.Fatalf("Expected %s for key %s for expected md %+v, got md %+v", v, k, d.expect, md)
			}
		}
	}
}
