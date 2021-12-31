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
// Original source: github.com/micro/go-micro/v3/api/api.go

package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"

	"github.com/micro/micro/v3/service/api/resolver"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/util/acme"
	"github.com/micro/micro/v3/util/codec"
	"github.com/micro/micro/v3/util/codec/bytes"
	"github.com/micro/micro/v3/util/codec/jsonrpc"
	"github.com/micro/micro/v3/util/codec/protorpc"
	"github.com/micro/micro/v3/util/qson"
	"github.com/oxtoacart/bpool"
)

var (
	bufferPool = bpool.NewSizedBufferPool(1024, 8)
)

type buffer struct {
	io.ReadCloser
}

func (b *buffer) Write(_ []byte) (int, error) {
	return 0, nil
}

type API interface {
	// Initialise options
	Init(...Option) error
	// Get the options
	Options() Options
	// Register a http handler
	Register(*Endpoint) error
	// Register a route
	Deregister(*Endpoint) error
	// Implementation of api
	String() string
}

// Server serves api requests
type Server interface {
	Address() string
	Init(opts ...Option) error
	Handle(path string, handler http.Handler)
	Start() error
	Stop() error
}

type Options struct {
	EnableACME   bool
	EnableCORS   bool
	ACMEProvider acme.Provider
	EnableTLS    bool
	ACMEHosts    []string
	TLSConfig    *tls.Config
	Resolver     resolver.Resolver
	Wrappers     []Wrapper
}

type Option func(*Options)

type Wrapper func(h http.Handler) http.Handler

func WrapHandler(w ...Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w...)
	}
}

func EnableCORS(b bool) Option {
	return func(o *Options) {
		o.EnableCORS = b
	}
}

func EnableACME(b bool) Option {
	return func(o *Options) {
		o.EnableACME = b
	}
}

func ACMEHosts(hosts ...string) Option {
	return func(o *Options) {
		o.ACMEHosts = hosts
	}
}

func ACMEProvider(p acme.Provider) Option {
	return func(o *Options) {
		o.ACMEProvider = p
	}
}

func EnableTLS(b bool) Option {
	return func(o *Options) {
		o.EnableTLS = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func Resolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}

// Endpoint is a mapping between an RPC method and HTTP endpoint
type Endpoint struct {
	// RPC Method e.g. Greeter.Hello
	Name string
	// Description e.g what's this endpoint for
	Description string
	// API Handler e.g rpc, proxy
	Handler string
	// HTTP Host e.g example.com
	Host []string
	// HTTP Methods e.g GET, POST
	Method []string
	// HTTP Path e.g /greeter. Expect POSIX regex
	Path []string
	// Body destination
	// "*" or "" - top level message value
	// "string" - inner message value
	Body string
	// Stream flag
	Stream bool
}

// Service represents an API service
type Service struct {
	// Name of service
	Name string
	// The endpoint for this service
	Endpoint *Endpoint
	// Versions of this service
	Services []*registry.Service
}

// Encode encodes an endpoint to endpoint metadata
func Encode(e *Endpoint) map[string]string {
	if e == nil {
		return nil
	}

	// endpoint map
	ep := make(map[string]string)

	// set vals only if they exist
	set := func(k, v string) {
		if len(v) == 0 {
			return
		}
		ep[k] = v
	}

	set("endpoint", e.Name)
	set("description", e.Description)
	set("handler", e.Handler)
	set("method", strings.Join(e.Method, ","))
	set("path", strings.Join(e.Path, ","))
	set("host", strings.Join(e.Host, ","))

	return ep
}

// Decode decodes endpoint metadata into an endpoint
func Decode(e map[string]string) *Endpoint {
	if e == nil {
		return nil
	}

	return &Endpoint{
		Name:        e["endpoint"],
		Description: e["description"],
		Method:      slice(e["method"]),
		Path:        slice(e["path"]),
		Host:        slice(e["host"]),
		Handler:     e["handler"],
	}
}

// Validate validates an endpoint to guarantee it won't blow up when being served
func Validate(e *Endpoint) error {
	if e == nil {
		return errors.New("endpoint is nil")
	}

	if len(e.Name) == 0 {
		return errors.New("name required")
	}

	for _, p := range e.Path {
		ps := p[0]
		pe := p[len(p)-1]

		if ps == '^' && pe == '$' {
			_, err := regexp.CompilePOSIX(p)
			if err != nil {
				return err
			}
		} else if ps == '^' && pe != '$' {
			return errors.New("invalid path")
		} else if ps != '^' && pe == '$' {
			return errors.New("invalid path")
		}
	}

	if len(e.Handler) == 0 {
		return errors.New("invalid handler")
	}

	return nil
}

func WithEndpoint(e *Endpoint) server.HandlerOption {
	return server.EndpointMetadata(e.Name, Encode(e))
}

func slice(s string) []string {
	var sl []string

	for _, p := range strings.Split(s, ",") {
		if str := strings.TrimSpace(p); len(str) > 0 {
			sl = append(sl, strings.TrimSpace(p))
		}
	}

	return sl
}

// RequestPayload takes a *http.Request and returns the payload
// If the request is a GET the query string parameters are extracted and marshaled to JSON and the raw bytes are returned.
// If the request method is a POST the request body is read and returned
func RequestPayload(r *http.Request) ([]byte, error) {
	var err error

	// we have to decode json-rpc and proto-rpc because we suck
	// well actually because there's no proxy codec right now

	ct := r.Header.Get("Content-Type")
	switch {
	case strings.Contains(ct, "application/json-rpc"):
		msg := codec.Message{
			Type:   codec.Request,
			Header: make(map[string]string),
		}
		c := jsonrpc.NewCodec(&buffer{r.Body})
		if err = c.ReadHeader(&msg, codec.Request); err != nil {
			return nil, err
		}
		var raw json.RawMessage
		if err = c.ReadBody(&raw); err != nil {
			return nil, err
		}
		return ([]byte)(raw), nil
	case strings.Contains(ct, "application/proto-rpc"), strings.Contains(ct, "application/octet-stream"):
		msg := codec.Message{
			Type:   codec.Request,
			Header: make(map[string]string),
		}
		c := protorpc.NewCodec(&buffer{r.Body})
		if err = c.ReadHeader(&msg, codec.Request); err != nil {
			return nil, err
		}
		var raw *bytes.Frame
		if err = c.ReadBody(&raw); err != nil {
			return nil, err
		}
		return raw.Data, nil
	case strings.Contains(ct, "application/x-www-form-urlencoded"):
		r.ParseForm()

		// generate a new set of values from the form
		vals := make(map[string]string)
		for k, v := range r.Form {
			vals[k] = strings.Join(v, ",")
		}

		// marshal
		return json.Marshal(vals)
	case strings.Contains(ct, "multipart/form-data"):
		// 10MB buffer
		if err := r.ParseMultipartForm(int64(10 << 20)); err != nil {
			return nil, err
		}
		vals := make(map[string]interface{})
		for k, v := range r.MultipartForm.Value {
			vals[k] = strings.Join(v, ",")
		}
		for k, _ := range r.MultipartForm.File {
			f, _, err := r.FormFile(k)
			if err != nil {
				return nil, err
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return nil, err
			}
			vals[k] = b
		}
		return json.Marshal(vals)
		// TODO: application/grpc
	}

	// otherwise as per usual
	ctx := r.Context()
	// dont user metadata.FromContext as it mangles names
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// allocate maximum
	matches := make(map[string]interface{}, len(md))
	bodydst := ""

	// get fields from url path
	for k, v := range md {
		k = strings.ToLower(k)
		// filter own keys
		if strings.HasPrefix(k, "x-api-field-") {
			matches[strings.TrimPrefix(k, "x-api-field-")] = v
			delete(md, k)
		} else if k == "x-api-body" {
			bodydst = v
			delete(md, k)
		}
	}

	// map of all fields
	req := make(map[string]interface{}, len(md))

	// get fields from url values
	if len(r.URL.RawQuery) > 0 {
		umd := make(map[string]interface{})
		err = qson.Unmarshal(&umd, r.URL.RawQuery)
		if err != nil {
			return nil, err
		}
		for k, v := range umd {
			matches[k] = v
		}
	}

	// restore context without fields
	*r = *r.Clone(metadata.NewContext(ctx, md))

	for k, v := range matches {
		ps := strings.Split(k, ".")
		if len(ps) == 1 {
			req[k] = v
			continue
		}
		em := make(map[string]interface{})
		em[ps[len(ps)-1]] = v
		for i := len(ps) - 2; i > 0; i-- {
			nm := make(map[string]interface{})
			nm[ps[i]] = em
			em = nm
		}
		if vm, ok := req[ps[0]]; ok {
			// nested map
			nm := vm.(map[string]interface{})
			for vk, vv := range em {
				nm[vk] = vv
			}
			req[ps[0]] = nm
		} else {
			req[ps[0]] = em
		}
	}
	pathbuf := []byte("{}")
	if len(req) > 0 {
		pathbuf, err = json.Marshal(req)
		if err != nil {
			return nil, err
		}
	}

	urlbuf := []byte("{}")
	out, err := jsonpatch.MergeMergePatches(urlbuf, pathbuf)
	if err != nil {
		return nil, err
	}

	switch r.Method {
	case "GET":
		// empty response
		if strings.Contains(ct, "application/json") && string(out) == "{}" {
			return out, nil
		} else if string(out) == "{}" && !strings.Contains(ct, "application/json") {
			return []byte{}, nil
		}
		return out, nil
	case "PATCH", "POST", "PUT", "DELETE":
		bodybuf := []byte("{}")
		buf := bufferPool.Get()
		defer bufferPool.Put(buf)
		if _, err := buf.ReadFrom(r.Body); err != nil {
			return nil, err
		}
		if b := buf.Bytes(); len(b) > 0 {
			bodybuf = b
		}
		if bodydst == "" || bodydst == "*" {
			// jsonpatch resequences the json object so we avoid it if possible (some usecases such as
			// validating signatures require the request body to be unchangedd). We're keeping support
			// for the custom paramaters for backwards compatability reasons.
			if string(out) == "{}" {
				return bodybuf, nil
			}

			if out, err = jsonpatch.MergeMergePatches(out, bodybuf); err == nil {
				return out, nil
			}
		}
		var jsonbody map[string]interface{}
		if json.Valid(bodybuf) {
			if err = json.Unmarshal(bodybuf, &jsonbody); err != nil {
				return nil, err
			}
		}
		dstmap := make(map[string]interface{})
		ps := strings.Split(bodydst, ".")
		if len(ps) == 1 {
			if jsonbody != nil {
				dstmap[ps[0]] = jsonbody
			} else {
				// old unexpected behaviour
				dstmap[ps[0]] = bodybuf
			}
		} else {
			em := make(map[string]interface{})
			if jsonbody != nil {
				em[ps[len(ps)-1]] = jsonbody
			} else {
				// old unexpected behaviour
				em[ps[len(ps)-1]] = bodybuf
			}
			for i := len(ps) - 2; i > 0; i-- {
				nm := make(map[string]interface{})
				nm[ps[i]] = em
				em = nm
			}
			dstmap[ps[0]] = em
		}

		bodyout, err := json.Marshal(dstmap)
		if err != nil {
			return nil, err
		}

		if out, err = jsonpatch.MergeMergePatches(out, bodyout); err == nil {
			return out, nil
		}

		//fallback to previous unknown behaviour
		return bodybuf, nil
	}

	return []byte{}, nil
}
