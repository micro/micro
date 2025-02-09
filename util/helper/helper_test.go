package helper

import (
	"net/http"
	"testing"

	metadata "github.com/micro/micro/v5/service/context"
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
		ctx := RequestToContext(d.request)
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
