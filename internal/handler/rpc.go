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

		var body struct {
			Service string
			Method  string
			Address string
			Request interface{}
		}

		if err = json.Unmarshal(b, &body); err != nil {
			badRequest(err.Error())
			return
		}

		service = body.Service
		method = body.Method
		address = body.Address

		if reqString, ok := body.Request.(string); ok {
			// for backwards compatibility also accept JSON request objects wrapped as strings...
			if err = json.Unmarshal([]byte(reqString), &request); err != nil {
				badRequest("while decoding request string: " + err.Error())
				return
			}
		} else {
			request = body.Request
		}
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		if err := json.Unmarshal([]byte(r.Form.Get("request")), &request); err != nil {
			badRequest("while decoding request string: " + err.Error())
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
