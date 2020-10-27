package cache

import (
	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/service/client"
)

type testRequest struct {
	service     string
	method      string
	endpoint    string
	contentType string
	codec       codec.Codec
	body        interface{}
	opts        client.RequestOptions
}

func newRequest(service, endpoint string, request interface{}, contentType string, reqOpts ...client.RequestOption) client.Request {
	var opts client.RequestOptions

	for _, o := range reqOpts {
		o(&opts)
	}

	// set the content-type specified
	if len(opts.ContentType) > 0 {
		contentType = opts.ContentType
	}

	return &testRequest{
		service:     service,
		method:      endpoint,
		endpoint:    endpoint,
		body:        request,
		contentType: contentType,
		opts:        opts,
	}
}

func (r *testRequest) ContentType() string {
	return r.contentType
}

func (r *testRequest) Service() string {
	return r.service
}

func (r *testRequest) Method() string {
	return r.method
}

func (r *testRequest) Endpoint() string {
	return r.endpoint
}

func (r *testRequest) Body() interface{} {
	return r.body
}

func (r *testRequest) Codec() codec.Writer {
	return r.codec
}

func (r *testRequest) Stream() bool {
	return r.opts.Stream
}
