// Package handler provides an RPC handler which outbounds requests
package handler

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/micro/go-micro/errors"
	proto "github.com/micro/micro/api/proto"
)

type API struct{}

func (a *API) Call(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	url := req.Url
	if len(req.Path) > 0 {
		url = req.Url + req.Path
	}

	// TODO: set get/post params

	// create new request
	r, err := http.NewRequest(req.Method, url, bytes.NewReader(req.Body))
	if err != nil {
		return errors.BadRequest("go.micro.api", err.Error())
	}

	// set the request headers
	for k, v := range req.Header {
		r.Header.Set(k, strings.Join(v.Values, ","))
	}

	// make the request
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return errors.InternalServerError("go.micro.api", err.Error())
	}

	// set the status code
	rsp.StatusCode = int32(resp.StatusCode)

	// set the headers
	for k, v := range resp.Header {
		rsp.Header[k] = &proto.Pair{
			Key:    k,
			Values: v,
		}
	}

	// read response body
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.InternalServerError("go.micro.api", err.Error())
	}
	defer resp.Body.Close()

	// set the response body
	rsp.Body = b
	return nil
}
