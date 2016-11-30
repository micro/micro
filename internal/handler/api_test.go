package handler

import (
	"net/http"
	"net/url"
	"testing"
)

func TestPathToReceiver(t *testing.T) {
	namespace := "go.micro.api"

	testData := []struct {
		path    string
		service string
		method  string
	}{
		{
			"/foo/bar",
			namespace + ".foo",
			"Foo.Bar",
		},
		{
			"/foo/foo/bar",
			namespace + ".foo",
			"Foo.Bar",
		},
		{
			"/foo/bar/baz",
			namespace + ".foo",
			"Bar.Baz",
		},
		{
			"/foo/bar/baz/cat",
			namespace + ".foo.bar",
			"Baz.Cat",
		},
		{
			"/foo/bar/baz/cat/car",
			namespace + ".foo.bar.baz",
			"Cat.Car",
		},
		{
			"/v1/foo/bar",
			namespace + ".v1.foo",
			"Foo.Bar",
		},
		{
			"/v1/foo/bar/baz",
			namespace + ".v1.foo",
			"Bar.Baz",
		},
		{
			"/v1/foo/bar/baz/cat",
			namespace + ".v1.foo.bar",
			"Baz.Cat",
		},
	}

	for _, d := range testData {
		s, m := pathToReceiver(namespace, d.path)
		if d.service != s {
			t.Fatalf("Expected service: %s for path: %s got: %s", d.service, d.path, s)
		}
		if d.method != m {
			t.Fatalf("Expected service: %s for path: %s got: %s", d.method, d.path, m)
		}
	}
}

func TestRequestToProto(t *testing.T) {
	testData := []*http.Request{
		&http.Request{
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
					t.Fatal("Expected val %s, got %s", val.Values[0], v[0])
				}
			}
		}
	}
}
