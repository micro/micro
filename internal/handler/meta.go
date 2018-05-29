package handler

import (
	"net/http"

	"github.com/micro/go-api"
	"github.com/micro/go-api/handler"
	"github.com/micro/go-api/handler/event"
	"github.com/micro/go-api/router"
	"github.com/micro/go-micro/errors"

	// TODO: only import handler package
	aapi "github.com/micro/go-api/handler/api"
	ahttp "github.com/micro/go-api/handler/http"
	arpc "github.com/micro/go-api/handler/rpc"
	aweb "github.com/micro/go-api/handler/web"
)

type metaHandler struct {
	r router.Router
}

func (m *metaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service, err := m.r.Route(r)
	if err != nil {
		er := errors.InternalServerError(m.r.Options().Namespace, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(er.Error()))
		return
	}

	// TODO: don't do this ffs
	switch service.Endpoint.Handler {
	// web socket handler
	case api.Web:
		aweb.WithService(service).ServeHTTP(w, r)
	// proxy handler
	case api.Proxy, api.Http:
		ahttp.WithService(service).ServeHTTP(w, r)
	// rpcx handler
	case api.Rpc:
		arpc.WithService(service).ServeHTTP(w, r)
	// event handler
	case api.Event:
		ev := event.NewHandler(
			handler.WithNamespace(m.r.Options().Namespace),
		)
		ev.ServeHTTP(w, r)
	// api handler
	case api.Api:
		aapi.WithService(service).ServeHTTP(w, r)
	// default handler: api
	default:
		aapi.WithService(service).ServeHTTP(w, r)
	}
}

// Meta is a http.Handler that routes based on endpoint metadata
func Meta(namespace string) http.Handler {
	return &metaHandler{
		r: router.NewRouter(router.WithNamespace(namespace)),
	}
}
