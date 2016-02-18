package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/errors"

	"golang.org/x/net/context"
)

// RPC Handler passes on a JSON or form encoded RPC request to
// a service.
func RPC(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var service, method, address string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	switch r.Header.Get("Content-Type") {
	case "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e := errors.BadRequest("go.micro.rpc", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var body map[string]interface{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			e := errors.BadRequest("go.micro.rpc", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var ok bool

		service, ok = body["service"].(string)
		if !ok {
			e := errors.BadRequest("go.micro.rpc", "invalid service")
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		method, ok = body["method"].(string)
		if !ok {
			e := errors.BadRequest("go.micro.rpc", "invalid method")
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		address, _ = body["address"].(string)
		req, _ := body["request"].(string)
		json.Unmarshal([]byte(req), &request)
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		json.Unmarshal([]byte(r.Form.Get("request")), &request)
	}

	var response map[string]interface{}
	var err error
	req := (*cmd.DefaultOptions().Client).NewJsonRequest(service, method, request)

	// remote call
	if len(address) > 0 {
		err = (*cmd.DefaultOptions().Client).CallRemote(context.Background(), address, req, &response)
	} else {
		err = (*cmd.DefaultOptions().Client).Call(context.Background(), req, &response)
	}
	if err != nil {
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

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}
