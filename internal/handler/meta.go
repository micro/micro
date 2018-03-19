package handler

import (
	"net/http"

	"github.com/micro/go-api"
	"github.com/micro/go-api/handler"
	"github.com/micro/go-api/handler/event"
	"github.com/micro/go-api/router"
	"github.com/micro/go-micro/errors"
)

type metaHandler struct {
	r router.Router
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError("go.micro.api", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	switch service.Endpoint.Handler {
	// proxy handler
	case api.Proxy, api.Http:
		Proxy(nil, service, false).ServeHTTP(w, r)
	// rpcx handler
	case api.Rpc:
		RPCX(nil, service).ServeHTTP(w, r)
	// event handler
	case api.Event:
		ev := event.NewHandler(
			handler.WithNamespace(m.r.Options().Namespace),
		)
		ev.ServeHTTP(w, r)
	// api handler
	case api.Api:
		API(nil, service).ServeHTTP(w, r)
	// default handler: api
	default:
		API(nil, service).ServeHTTP(w, r)
	}
}

// Meta is a http.Handler that routes based on endpoint metadata
func Meta(namespace string) http.Handler {
	return &metaHandler{
		r: router.NewRouter(router.WithNamespace(namespace)),
	}
}
