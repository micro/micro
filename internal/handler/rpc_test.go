package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/metadata"
	rmock "github.com/micro/go-micro/registry/mock"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"

	"golang.org/x/net/context"
)

type TestHandler struct {
	t      *testing.T
	expect metadata.Metadata
}

type TestRequest struct{}
type TestResponse struct{}

func (t *TestHandler) Exec(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return fmt.Errorf("Expected metadata got %t", ok)
	}

	for k, v := range t.expect {
		if val := md[k]; val != v {
			return fmt.Errorf("Expected %s for key %s got %s", v, k, val)
		}
	}

	t.t.Logf("Received request %+v", req)
	t.t.Logf("Received metadata %+v", md)

	return nil
}

func TestRPCHandler(t *testing.T) {
	r := rmock.NewRegistry()

	(*cmd.DefaultOptions().Client).Init(
		client.Registry(r),
		client.Selector(selector.NewSelector(selector.Registry(r))),
	)

	(*cmd.DefaultOptions().Server).Init(
		server.Name("test"),
		server.Registry(r),
	)

	(*cmd.DefaultOptions().Server).Handle(
		(*cmd.DefaultOptions().Server).NewHandler(&TestHandler{t, metadata.Metadata{"Foo": "Bar"}}),
	)

	if err := server.Start(); err != nil {
		t.Fatal(err)
	}

	if err := server.Register(); err != nil {
		t.Fatal(err)
	}

	rpc := NewRPCWithWhitelist("test")

	testCases := []struct {
		service string
		method  string
		request string
		status  int
	}{
		{
			service: "test",
			method:  "TestHandler.Exec",
			request: "{}",
			status:  200,
		},
		{
			service: "restricted",
			method:  "TestHandler.Exec",
			request: "{}",
			status:  403,
		},
	}

	for i, testCase := range testCases {
		w := httptest.NewRecorder()

		request := map[string]string{
			"service": testCase.service,
			"method":  testCase.method,
			"request": testCase.request,
		}

		rb, err := json.Marshal(request)
		if err != nil {
			t.Errorf("test case %d: %s", i+1, err)
			continue
		}

		b := bytes.NewBuffer(rb)

		req, err := http.NewRequest("POST", "/rpc", b)
		if err != nil {
			t.Errorf("test case %d: %s", i+1, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Foo", "Bar")

		rpc.ServeHTTP(w, req)

		if w.Code != testCase.status {
			t.Errorf("test case %d: expected %d response got %d %s",
				i+1, testCase.status, w.Code, w.Body.String())
			continue
		}
	}

	if err := server.Deregister(); err != nil {
		t.Fatal(err)
	}

	if err := server.Stop(); err != nil {
		t.Fatal(err)
	}
}
