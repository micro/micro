package api

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"

	"golang.org/x/net/context"
)

// Translates /foo/bar/zool into api service go.micro.api.foo method Bar.Zool
// Translates /foo/bar into api service go.micro.api.foo method Foo.Bar
func pathToReceiver(p string) (string, string) {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")

	if len(parts) <= 2 {
		service := Namespace + "." + strings.Join(parts[:len(parts)-1], ".")
		method := strings.Title(strings.Join(parts, "."))
		return service, method
	}

	service := Namespace + "." + strings.Join(parts[:len(parts)-2], ".")
	method := strings.Title(strings.Join(parts[len(parts)-2:], "."))
	return service, method
}

func requestToProto(r *http.Request) (*api.Request, error) {
	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("Error parsing form: %v", err)
	}

	req := &api.Request{
		Path:   r.URL.Path,
		Method: r.Method,
		Header: make(map[string]*api.Pair),
		Get:    make(map[string]*api.Pair),
		Post:   make(map[string]*api.Pair),
	}

	ct, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		ct = "application/x-www-form-urlencoded"
		r.Header.Set("Content-Type", ct)
	}

	switch ct {
	case "application/x-www-form-urlencoded":
		// expect form vals
	default:
		data, _ := ioutil.ReadAll(r.Body)
		req.Body = string(data)
	}

	// Get data
	for key, vals := range r.URL.Query() {
		header, ok := req.Get[key]
		if !ok {
			header = &api.Pair{
				Key: key,
			}
			req.Get[key] = header
		}
		header.Values = vals
	}

	// Post data
	for key, vals := range r.PostForm {
		header, ok := req.Post[key]
		if !ok {
			header = &api.Pair{
				Key: key,
			}
			req.Post[key] = header
		}
		header.Values = vals
	}

	// Pass through custom headers
	for key, vals := range r.Header {
		if !strings.HasPrefix(key, HeaderPrefix) {
			continue
		}
		header, ok := req.Header[key]
		if !ok {
			header = &api.Pair{
				Key: key,
			}
			req.Post[key] = header
		}
		header.Values = vals
	}

	return req, nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	request, err := requestToProto(r)
	if err != nil {
		er := errors.InternalServerError("go.micro.api", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	service, method := pathToReceiver(r.URL.Path)
	req := (*cmd.DefaultOptions().Client).NewRequest(service, method, request)
	rsp := &api.Response{}
	if err := (*cmd.DefaultOptions().Client).Call(context.Background(), req, rsp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	for _, header := range rsp.GetHeader() {
		for _, val := range header.Values {
			w.Header().Add(header.Key, val)
		}
	}

	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(int(rsp.StatusCode))
	w.Write([]byte(rsp.Body))
}
