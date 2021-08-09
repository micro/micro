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
// Original source: github.com/micro/go-micro/v3/api/handler/rpc/rpc.go

// Package rpc is a go-micro rpc handler.
package rpc

import (
	bts "bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/api/handler"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/codec/bytes"
	"github.com/micro/micro/v3/util/ctx"
	"github.com/micro/micro/v3/util/router"
)

const (
	Handler = "rpc"
)

var (
	// supported json codecs
	jsonCodecs = []string{
		"application/grpc+json",
		"application/json",
		"application/json-rpc",
	}

	// support proto codecs
	protoCodecs = []string{
		"application/grpc",
		"application/grpc+proto",
		"application/proto",
		"application/protobuf",
		"application/proto-rpc",
		"application/octet-stream",
	}
)

type rpcHandler struct {
	opts handler.Options
	s    *api.Service
}

type buffer struct {
	io.ReadCloser
}

func (b *buffer) Write(_ []byte) (int, error) {
	return 0, nil
}

// see https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and/28596225
func jsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bts.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return bts.TrimRight(buffer.Bytes(), "\n"), err
}

func (h *rpcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bsize := handler.DefaultMaxRecvSize
	if h.opts.MaxRecvSize > 0 {
		bsize = h.opts.MaxRecvSize
	}

	r.Body = http.MaxBytesReader(w, r.Body, bsize)

	defer r.Body.Close()
	var service *api.Service

	if h.s != nil {
		// we were given the service
		service = h.s
	} else if h.opts.Router != nil {
		// try get service from router
		s, err := h.opts.Router.Route(r)
		if err != nil {
			writeError(w, r, errors.InternalServerError("go.micro.api", err.Error()))
			return
		}
		service = s
	} else {
		// we have no way of routing the request
		writeError(w, r, errors.InternalServerError("go.micro.api", "no route found"))
		return
	}

	ct := r.Header.Get("Content-Type")

	// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}

	// micro client
	c := h.opts.Client

	// create context
	cx := ctx.FromRequest(r)

	// set merged context to request
	*r = *r.Clone(cx)
	// if stream we currently only support json
	if isStream(r, service) {
		serveStream(cx, w, r, service, c)
		return
	}

	// create custom router
	callOpt := client.WithRouter(router.New(service.Services))

	// walk the standard call path
	// get payload
	br, err := api.RequestPayload(r)
	if err != nil {
		writeError(w, r, err)
		return
	}

	var rsp []byte

	switch {
	// proto codecs
	case hasCodec(ct, protoCodecs):
		var request *bytes.Frame
		// if the extracted payload isn't empty lets use it
		if len(br) > 0 {
			request = &bytes.Frame{Data: br}
		}

		// create the request
		req := c.NewRequest(
			service.Name,
			service.Endpoint.Name,
			request,
			client.WithContentType(ct),
		)

		// make the call
		var response *bytes.Frame
		if err := c.Call(cx, req, response, callOpt); err != nil {
			writeError(w, r, err)
			return
		}
		rsp = response.Data
	default:
		// if json codec is not present set to json
		if !hasCodec(ct, jsonCodecs) {
			ct = "application/json"
		}

		// default to trying json
		var request json.RawMessage
		// if the extracted payload isn't empty lets use it
		if len(br) > 0 {
			request = json.RawMessage(br)
		}

		// create request/response
		var response interface{}

		req := c.NewRequest(
			service.Name,
			service.Endpoint.Name,
			&request,
			client.WithContentType(ct),
		)
		// make the call
		if err := c.Call(cx, req, &response, callOpt); err != nil {
			writeError(w, r, err)
			return
		}

		// marshall response
		// see https://play.golang.org/p/oBNxUjVTzus
		rsp, err = jsonMarshal(response)
		if err != nil {
			writeError(w, r, err)
			return
		}
	}

	// write the response
	writeResponse(w, r, rsp)
}

func (rh *rpcHandler) String() string {
	return "rpc"
}

func hasCodec(ct string, codecs []string) bool {
	for _, codec := range codecs {
		if ct == codec {
			return true
		}
	}
	return false
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	// response content type
	w.Header().Set("Content-Type", "application/json")

	// parse out the error code
	ce := errors.Parse(err.Error())

	switch ce.Code {
	case 0:
		// assuming it's totally screwed
		ce.Code = 500
		ce.Id = "go.micro.api"
		ce.Status = http.StatusText(500)
		ce.Detail = "error during request: " + ce.Detail
		w.WriteHeader(500)
	default:
		w.WriteHeader(int(ce.Code))
	}

	// Set trailers
	if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		w.Header().Set("Trailer", "grpc-status")
		w.Header().Set("Trailer", "grpc-message")
		w.Header().Set("grpc-status", "13")
		w.Header().Set("grpc-message", ce.Detail)
	}

	_, werr := w.Write([]byte(ce.Error()))
	if werr != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(werr)
		}
	}
}

func writeResponse(w http.ResponseWriter, r *http.Request, rsp []byte) {
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", strconv.Itoa(len(rsp)))

	// Set trailers
	if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		w.Header().Set("Trailer", "grpc-status")
		w.Header().Set("Trailer", "grpc-message")
		w.Header().Set("grpc-status", "0")
		w.Header().Set("grpc-message", "")
	}

	// write 204 status if rsp is nil
	if len(rsp) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	// write response
	_, err := w.Write(rsp)
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error(err)
		}
	}

}

func NewHandler(opts ...handler.Option) handler.Handler {
	options := handler.NewOptions(opts...)
	return &rpcHandler{
		opts: options,
	}
}

func WithService(s *api.Service, opts ...handler.Option) handler.Handler {
	options := handler.NewOptions(opts...)
	return &rpcHandler{
		opts: options,
		s:    s,
	}
}
