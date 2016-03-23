package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/micro/internal/helper"
)

type rpcRequest struct {
	Service string
	Method  string
	Address string
	Request interface{}
}

// RPC Handler passes on a JSON or form encoded RPC request to
// a service.
func RPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	badRequest := func(description string) {
		e := errors.BadRequest("go.micro.rpc", description)
		w.WriteHeader(400)
		w.Write([]byte(e.Error()))
	}

	var service, method, address string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	switch r.Header.Get("Content-Type") {
	case "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			badRequest(err.Error())
			return
		}

		var rpcReq rpcRequest

		if err = json.Unmarshal(b, &rpcReq); err != nil {
			badRequest(err.Error())
			return
		}

		service = rpcReq.Service
		method = rpcReq.Method
		address = rpcReq.Address
		request = rpcReq.Request

		// JSON as string
		if req, ok := rpcReq.Request.(string); ok {
			if err = json.Unmarshal([]byte(req), &request); err != nil {
				badRequest("error decoding request string: " + err.Error())
				return
			}
		}
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		if err := json.Unmarshal([]byte(r.Form.Get("request")), &request); err != nil {
			badRequest("error decoding request string: " + err.Error())
			return
		}
	}

	if len(service) == 0 {
		badRequest("invalid service")
		return
	}

	if len(method) == 0 {
		badRequest("invalid method")
		return
	}

	// create request/response
	var response map[string]interface{}
	var err error
	req := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)

	// create context
	ctx := helper.RequestToContext(r)

	// remote call
	if len(address) > 0 {
		err = (*cmd.DefaultOptions().Client).CallRemote(ctx, address, req, &response)
	} else {
		err = (*cmd.DefaultOptions().Client).Call(ctx, req, &response)
	}
	if err != nil {
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			// assuming it's totally screwed
			ce.Code = 500
			ce.Id = "go.micro.rpc"
			ce.Status = http.StatusText(500)
			ce.Detail = "error during request: " + ce.Detail
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}
