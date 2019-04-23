package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/micro/web/common"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func rpc(w http.ResponseWriter, ctx context.Context, rpcReq *rpcRequest) {

	if len(rpcReq.Service) == 0 {
		writeError(w, "Service Is Not found")
		return
	}

	if len(rpcReq.Endpoint) == 0 {
		writeError(w, "Endpoint Is Not found")
		return
	}

	// decode rpc request param body
	if req, ok := rpcReq.Request.(string); ok {
		d := json.NewDecoder(strings.NewReader(req))
		d.UseNumber()

		if err := d.Decode(&rpcReq.Request); err != nil {
			writeError(w, "error decoding request string: "+err.Error())
			return
		}
	}

	// create request/response
	var response json.RawMessage
	var err error
	req := (*cmd.DefaultOptions().Client).NewRequest(rpcReq.Service, rpcReq.Endpoint, rpcReq.Request, client.WithContentType("application/json"))

	var opts []client.CallOption

	// set timeout
	if rpcReq.timeout > 0 {
		opts = append(opts, client.WithRequestTimeout(time.Duration(rpcReq.timeout)*time.Second))
	}

	// remote call
	if len(rpcReq.Address) > 0 {
		opts = append(opts, client.WithAddress(rpcReq.Address))
	}

	// remote call
	err = (*cmd.DefaultOptions().Client).Call(ctx, req, &response, opts...)
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

	if strings.Contains(rpcReq.URL, "/v1/rpc") {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		writeJsonData(w, response)
	}
}

func apiProxy() http.Handler {
	/*sel := selector.NewSelector(
		selector.Registry(*cmd.DefaultOptions().Registry),
	)*/

	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		parts := strings.Split(r.URL.Path, "/v1/api-stats")
		if len(parts) < 2 {
			kill()
			return
		}

		address := r.URL.Query().Get("address")
		/*port, err := strconv.Atoi(portStr)
		  if err != nil {
		  	kill()
		  	return
		  }

		  next, err := sel.Select(common.APINamespace, func(options *selector.SelectOptions) {

		  	filter := func(services []*registry.Service) []*registry.Service {
		  		ret := make([]*registry.Service, 0)
		  		for _, s := range services {
		  			for _, n := range s.Nodes {
		  				if n.Address == address && n.Metadata["port"] == portStr {
		  					ret = append(ret, s)
		  				}
		  			}
		  		}
		  		return ret
		  	}

		  	options.Filters = append(options.Filters, filter)
		  })

		  if err != nil {
		  	kill()
		  	return
		  }

		  s, err := next()
		  if err != nil {
		  	kill()
		  	return
		  }*/

		r.Header.Set(common.BasePathHeader, "/v1/api-stats")
		r.Header.Set("Content-Type", "application/json")
		r.URL.Host = fmt.Sprintf("%s", address)
		r.URL.Path = "/stats" + strings.Join(parts[2:], "/")
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &common.Proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}

func writeJsonData(w http.ResponseWriter, data interface{}) {

	rsp := &Rsp{
		Data:    data,
		Success: true,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeError(w http.ResponseWriter, msg string) {

	rsp := &Rsp{
		Error:   msg,
		Success: false,
	}

	b, err := json.Marshal(rsp)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
