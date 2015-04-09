package sic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/asim/go-micro/client"
	"github.com/asim/go-micro/errors"
	log "github.com/golang/glog"
)

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var service, method string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	switch r.Header.Get("Content-Type") {
	case "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e := errors.BadRequest("go.micro.api", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var body map[string]interface{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			e := errors.BadRequest("go.micro.api", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		service = body["service"].(string)
		method = body["method"].(string)
		request = body["request"]
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		json.Unmarshal([]byte(r.Form.Get("request")), &request)
	}

	var response map[string]interface{}
	req := client.NewJsonRequest(service, method, request)
	err := client.Call(req, &response)
	if err != nil {
		log.Errorf("Error calling %s.%s: %v", service, method, err)
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
