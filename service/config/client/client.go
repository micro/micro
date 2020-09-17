package client

import (
	"net/http"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config"
	proto "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/errors"
)

var (
	defaultNamespace = "micro"
	defaultPath      = ""
	name             = "config"
)

type srv struct {
	serviceName string
	opts        config.Options
	namespace   string
	path        string
	client      proto.ConfigService
}

func (m *srv) Get() (set *proto.Value, err error) {
	req, err := m.client.Get(context.DefaultContext, &proto.GetRequest{
		Namespace: m.namespace,
		Path:      m.path,
	}, goclient.WithAuthToken())
	if verr := errors.Parse(err); verr != nil && verr.Code == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return req.Value, nil
}

// Write is unsupported
func (m *srv) Write() error {
	return nil
}

func (m *srv) String() string {
	return "service"
}

func NewSource(opts ...config.Option) *srv {
	var options config.Options
	for _, o := range opts {
		o(&options)
	}

	addr := name
	namespace := defaultNamespace
	path := defaultPath

	if options.Context != nil {
		a, ok := options.Context.Value(serviceNameKey{}).(string)
		if ok {
			addr = a
		}

		k, ok := options.Context.Value(namespaceKey{}).(string)
		if ok {
			namespace = k
		}

		p, ok := options.Context.Value(pathKey{}).(string)
		if ok {
			path = p
		}
	}

	if options.Client == nil {
		options.Client = client.DefaultClient
	}

	s := &srv{
		serviceName: addr,
		opts:        options,
		namespace:   namespace,
		path:        path,
		client:      proto.NewConfigService(addr, options.Client),
	}

	return s
}
