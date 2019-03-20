package v1

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"
	"net/http"
	"strings"
	"time"
)

func rpc(w http.ResponseWriter, ctx context.Context, rpcReq *rpcRequest) {

	if len(rpcReq.Service) == 0 {
		writeError(w, "Service Is Not found")
	}

	if len(rpcReq.Endpoint) == 0 {
		writeError(w, "Endpoint Is Not found")
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

	writeJsonData(w, response)
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
