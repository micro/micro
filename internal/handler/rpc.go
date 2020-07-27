package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/micro/go-micro/v3/api/handler"
	"github.com/micro/go-micro/v3/api/resolver"
	"github.com/micro/go-micro/v3/api/resolver/subdomain"
	"github.com/micro/go-micro/v3/api/server/cors"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v2/internal/helper"
	muclient "github.com/micro/micro/v2/service/client"
)

type rpcRequest struct {
	Service  string
	Endpoint string
	Method   string
	Address  string
	Request  interface{}
}

type rpcHandler struct {
	resolver resolver.Resolver
}

func (h *rpcHandler) String() string {
	return "internal/rpc"
}

// ServeHTTP passes on a JSON or form encoded RPC request to a service.
func (h *rpcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		cors.SetHeaders(w, r)
		return
	}

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

	var service, endpoint, address string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	ct := r.Header.Get("Content-Type")

	// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}

	switch ct {
	case "application/json":
		var rpcReq rpcRequest

		d := json.NewDecoder(r.Body)
		d.UseNumber()

		if err := d.Decode(&rpcReq); err != nil {
			badRequest(err.Error())
			return
		}

		service = rpcReq.Service
		endpoint = rpcReq.Endpoint
		address = rpcReq.Address
		request = rpcReq.Request
		if len(endpoint) == 0 {
			endpoint = rpcReq.Method
		}

		// JSON as string
		if req, ok := rpcReq.Request.(string); ok {
			d := json.NewDecoder(strings.NewReader(req))
			d.UseNumber()

			if err := d.Decode(&request); err != nil {
				badRequest("error decoding request string: " + err.Error())
				return
			}
		}
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		endpoint = r.Form.Get("endpoint")
		address = r.Form.Get("address")
		if len(endpoint) == 0 {
			endpoint = r.Form.Get("method")
		}

		d := json.NewDecoder(strings.NewReader(r.Form.Get("request")))
		d.UseNumber()

		if err := d.Decode(&request); err != nil {
			badRequest("error decoding request string: " + err.Error())
			return
		}
	}

	if len(service) == 0 {
		badRequest("invalid service")
		return
	}

	if len(endpoint) == 0 {
		badRequest("invalid endpoint")
		return
	}

	// create request/response
	var response json.RawMessage
	var err error
	req := muclient.DefaultClient.NewRequest(service, endpoint, request, client.WithContentType("application/json"))

	// create context
	ctx := helper.RequestToContext(r)

	var opts []client.CallOption

	timeout, _ := strconv.Atoi(r.Header.Get("Timeout"))
	// set timeout
	if timeout > 0 {
		opts = append(opts, client.WithRequestTimeout(time.Duration(timeout)*time.Second))
	}

	// remote call
	if len(address) > 0 {
		opts = append(opts, client.WithAddress(address))
	}

	// since services can be running in many domains, we'll use the resolver to determine the domain
	// which should be used on the call
	if resolver, ok := h.resolver.(*subdomain.Resolver); ok {
		if dom := resolver.Domain(r); len(dom) > 0 {
			opts = append(opts, client.WithNetwork(dom))
		}
	}

	// remote call
	err = muclient.DefaultClient.Call(ctx, req, &response, opts...)
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

	b, _ := response.MarshalJSON()
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

// NewRPCHandler returns an initialized RPC handler
func NewRPCHandler(r resolver.Resolver) handler.Handler {
	return &rpcHandler{r}
}
