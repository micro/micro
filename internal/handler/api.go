package handler

import (
	"net/http"

	goapi "github.com/micro/go-api"
	"github.com/micro/go-api/router"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/selector"
	api "github.com/micro/micro/api/proto"
	"github.com/micro/micro/internal/helper"
)

type apiHandler struct {
	r router.Router
	s *goapi.Service
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

	var service *goapi.Service

	if a.r != nil {
		// try get service from router
		s, err := a.r.Route(r)
		if err != nil {
			er := errors.InternalServerError("go.micro.api", err.Error())
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			w.Write([]byte(er.Error()))
			return
		}
		service = s
	} else if a.s != nil {
		// we were given the service
		service = a.s
	} else {
		// we have no way of routing the request
		er := errors.InternalServerError("go.micro.api", "no route found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// create request and response
	req := (*cmd.DefaultOptions().Client).NewRequest(service.Name, service.Endpoint.Name, request)
	rsp := &api.Response{}

	// create the context from headers
	ctx := helper.RequestToContext(r)
	// create strategy
	so := selector.WithStrategy(strategy(service.Services))

	if err := (*cmd.DefaultOptions().Client).Call(ctx, req, rsp, client.WithSelectOption(so)); err != nil {
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

func API(r router.Router, s *goapi.Service) http.Handler {
	return &apiHandler{
		r: r,
		s: s,
	}
}
