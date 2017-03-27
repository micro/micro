package handler

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"path"
	"regexp"
	"strings"

	api "github.com/micro/go-api/proto"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/micro/internal/helper"
)

var (
	versionRe = regexp.MustCompilePOSIX("^v[0-9]+$")
)

type apiHandler struct {
	Namespace string
}

// Translates /foo/bar/zool into api service go.micro.api.foo method Bar.Zool
// Translates /foo/bar into api service go.micro.api.foo method Foo.Bar
func pathToReceiver(ns, p string) (string, string) {
	p = path.Clean(p)
	p = strings.TrimPrefix(p, "/")
	parts := strings.Split(p, "/")

	// If we've got two or less parts
	// Use first part as service
	// Use all parts as method
	if len(parts) <= 2 {
		service := ns + "." + strings.Join(parts[:len(parts)-1], ".")
		method := strings.Title(strings.Join(parts, "."))
		return service, method
	}

	// Treat /v[0-9]+ as versioning where we have 3 parts
	// /v1/foo/bar => service: v1.foo method: Foo.bar
	if len(parts) == 3 && versionRe.Match([]byte(parts[0])) {
		service := ns + "." + strings.Join(parts[:len(parts)-1], ".")
		method := strings.Title(strings.Join(parts[len(parts)-2:], "."))
		return service, method
	}

	// Service is everything minus last two parts
	// Method is the last two parts
	service := ns + "." + strings.Join(parts[:len(parts)-2], ".")
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

	// Set X-Forwarded-For if it does not exist
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if prior, ok := r.Header["X-Forwarded-For"]; ok {
			ip = strings.Join(prior, ", ") + ", " + ip
		}

		// Set the header
		req.Header["X-Forwarded-For"] = &api.Pair{
			Key:    "X-Forwarded-For",
			Values: []string{ip},
		}
	}

	// Host is stripped from net/http Headers so let's add it
	req.Header["Host"] = &api.Pair{
		Key:    "Host",
		Values: []string{r.Host},
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

	for key, vals := range r.Header {
		header, ok := req.Header[key]
		if !ok {
			header = &api.Pair{
				Key: key,
			}
			req.Header[key] = header
		}
		header.Values = vals
	}

	return req, nil
}

// API handler is the default handler which takes api.Request and returns api.Response
func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := requestToProto(r)
	if err != nil {
		er := errors.InternalServerError("go.micro.api", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// get service and method
	service, method := pathToReceiver(a.Namespace, r.URL.Path)

	// create request and response
	req := (*cmd.DefaultOptions().Client).NewRequest(service, method, request)
	rsp := &api.Response{}

	// create the context from headers
	ctx := helper.RequestToContext(r)

	if err := (*cmd.DefaultOptions().Client).Call(ctx, req, rsp); err != nil {
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

func API(namespace string) http.Handler {
	return &apiHandler{
		Namespace: namespace,
	}
}
