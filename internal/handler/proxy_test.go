package handler

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/micro/go-api"
	"github.com/micro/go-api/router"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	rmock "github.com/micro/go-micro/registry/mock"
)

func testProxy(t *testing.T, path, service string) {
	r := rmock.NewRegistry()
	cmd.DefaultCmd = cmd.NewCmd(cmd.Registry(&r))

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	parts := strings.Split(l.Addr().String(), ":")

	var host string
	var port int

	host = parts[0]
	port, _ = strconv.Atoi(parts[1])

	s := &registry.Service{
		Name: service,
		Nodes: []*registry.Node{
			&registry.Node{
				Id:      service + "-1",
				Address: host,
				Port:    port,
			},
		},
	}

	r.Register(s)
	defer r.Deregister(s)

	// setup the test handler
	m := http.NewServeMux()
	m.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`you got served`))
	})

	// start http test serve
	go http.Serve(l, m)

	// create new request and writer
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	// initialise the handler
	rt := router.NewRouter(router.WithHandler(api.Proxy))

	p := Proxy(rt, nil, false)

	// execute the handler
	p.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 response got %d %s", w.Code, w.Body.String())
	}

	if w.Body.String() != "you got served" {
		t.Fatalf("Expected body: you got served. Got: %s", w.Body.String())
	}
}

func TestProxyHandler(t *testing.T) {
	testData := []struct {
		path    string
		service string
	}{
		{
			"/test/foo",
			"go.micro.api.test",
		},
		{
			"/test/foo/baz",
			"go.micro.api.test",
		},
		{
			"/v1/foo",
			"go.micro.api.v1.foo",
		},
		{
			"/v1/foo/bar",
			"go.micro.api.v1.foo",
		},
		{
			"/v2/baz",
			"go.micro.api.v2.baz",
		},
		{
			"/v2/baz/bar",
			"go.micro.api.v2.baz",
		},
	}

	for _, d := range testData {
		testProxy(t, d.path, d.service)
	}
}
