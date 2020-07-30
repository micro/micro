package client

import (
	"context"
	"net/http"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config/source"
	"github.com/micro/micro/v3/service/client"
	proto "github.com/micro/micro/v3/service/config/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

var (
	defaultNamespace = "micro"
	defaultPath      = ""
	name             = "go.micro.config"
)

type srv struct {
	serviceName string
	namespace   string
	path        string
	opts        source.Options
	client      proto.ConfigService
}

func (m *srv) Read() (set *source.ChangeSet, err error) {
	client := proto.NewConfigService(m.serviceName)
	req, err := client.Read(context.Background(), &proto.ReadRequest{
		Namespace: m.namespace,
		Path:      m.path,
	}, goclient.WithServiceToken())
	if verr := errors.Parse(err); verr != nil && verr.Code == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toChangeSet(req.Change.ChangeSet), nil
}

func (m *srv) Watch() (w source.Watcher, err error) {
	client := proto.NewConfigService(m.serviceName)
	stream, err := client.Watch(context.Background(), &proto.WatchRequest{
		Namespace: m.namespace,
		Path:      m.path,
	}, goclient.WithServiceToken())
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error("watch err: ", err)
		}
		return
	}
	return newWatcher(stream)
}

// Write is unsupported
func (m *srv) Write(cs *source.ChangeSet) error {
	return nil
}

func (m *srv) String() string {
	return "service"
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
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
	}

	return s
}
