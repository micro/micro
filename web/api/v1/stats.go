package v1

import (
	"github.com/micro/micro/internal/helper"
	"net/http"
)

func (api *API) stats(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	rpcReq := &rpcRequest{
		Service:  r.URL.Query().Get("service"),
		Endpoint: "Debug.Stats",
		Request:  "{}",
		Address:  r.URL.Query().Get("address"),
	}

	rpc(w, helper.RequestToContext(r), rpcReq)
	return
}

func (api *API) apiStats(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	rpcReq := &rpcRequest{
		Service:  "go.micro.api",
		Endpoint: "stats",
		Request:  "{}",
		Address:  r.URL.Query().Get("address"),
	}

	rpc(w, helper.RequestToContext(r), rpcReq)
	return
}
