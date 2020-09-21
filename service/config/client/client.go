package client

import (
	"encoding/json"
	"net/http"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/config"
	proto "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
)

var (
	defaultNamespace = "micro"
	name             = "config"
)

type srv struct {
	opts      config.Options
	namespace string
	client    proto.ConfigService
}

func (m *srv) Get(path string, options ...config.Option) config.Value {
	req, err := m.client.Get(context.DefaultContext, &proto.GetRequest{
		Namespace: m.namespace,
		Path:      path,
	}, goclient.WithAuthToken())
	if verr := errors.Parse(err); verr != nil && verr.Code == http.StatusNotFound {
		return config.NewJSONValue([]byte("null"))
	} else if err != nil {
		logger.Error(err)
		return config.NewJSONValue([]byte("null"))
	}

	return config.NewJSONValue([]byte(req.Value.Data))
}

func (m *srv) Set(path string, value interface{}, options ...config.Option) {
	dat, _ := json.Marshal(value)
	_, err := m.client.Set(context.DefaultContext, &proto.SetRequest{
		Namespace: m.namespace,
		Path:      path,
		Value: &proto.Value{
			Data: string(dat),
		},
	}, goclient.WithAuthToken())
	if err != nil {
		logger.Error(err)
	}
}

func (m *srv) Delete(path string, options ...config.Option) {
	_, err := m.client.Delete(context.DefaultContext, &proto.DeleteRequest{
		Namespace: m.namespace,
		Path:      path,
	}, goclient.WithAuthToken())
	if err != nil {
		logger.Error(err)
	}
}

func (m *srv) String() string {
	return "service"
}

func NewConfig(namespace string) *srv {
	addr := name
	if len(namespace) == 0 {
		namespace = defaultNamespace
	}

	s := &srv{
		namespace: namespace,
		client:    proto.NewConfigService(addr, client.DefaultClient),
	}

	return s
}
