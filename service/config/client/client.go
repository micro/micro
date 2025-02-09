package client

import (
	"encoding/json"
	"net/http"

	proto "github.com/micro/micro/v5/proto/config"
	"github.com/micro/micro/v5/service/client"
	"github.com/micro/micro/v5/service/config"
	"github.com/micro/micro/v5/service/context"
	"github.com/micro/micro/v5/service/errors"
)

var (
	defaultNamespace = "micro"
)

type srv struct {
	opts      config.Options
	namespace string
	client    proto.ConfigService
}

func (m *srv) Get(path string, options ...config.Option) (config.Value, error) {
	o := config.Options{}
	for _, option := range options {
		option(&o)
	}
	nullValue := config.NewJSONValue([]byte("null"))
	req, err := m.client.Get(context.DefaultContext, &proto.GetRequest{
		Namespace: m.namespace,
		Path:      path,
		Options: &proto.Options{
			Secret: o.Secret,
		},
	}, client.WithAuthToken())
	if verr := errors.FromError(err); verr != nil && verr.Code == http.StatusNotFound {
		return nullValue, nil
	} else if err != nil {
		return nullValue, err
	}

	return config.NewJSONValue([]byte(req.Value.Data)), nil
}

func (m *srv) Set(path string, value interface{}, options ...config.Option) error {
	o := config.Options{}
	for _, option := range options {
		option(&o)
	}
	dat, _ := json.Marshal(value)
	_, err := m.client.Set(context.DefaultContext, &proto.SetRequest{
		Namespace: m.namespace,
		Path:      path,
		Value: &proto.Value{
			Data: string(dat),
		},
		Options: &proto.Options{
			Secret: o.Secret,
		},
	}, client.WithAuthToken())
	return err
}

func (m *srv) Delete(path string, options ...config.Option) error {
	_, err := m.client.Delete(context.DefaultContext, &proto.DeleteRequest{
		Namespace: m.namespace,
		Path:      path,
	}, client.WithAuthToken())
	return err
}

func (m *srv) String() string {
	return "service"
}

func NewConfig() *srv {
	s := &srv{
		namespace: defaultNamespace,
		client:    proto.NewConfigService("config", client.DefaultClient),
	}

	return s
}
