package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/micro/go-micro/v3/config"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

const (
	defaultNamespace = "micro"
	pathSplitter     = "."
)

var (
	// we now support json only
	mtx sync.RWMutex
)

type Config struct {
	secret []byte
}

func NewConfig(key string) *Config {
	if len(key) == 0 {
		logger.Fatalf("No encryption key provided")
	}
	dec, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		logger.Fatalf("Error decoding key: %v", err)
	}
	return &Config{
		secret: dec,
	}
}

func (c *Config) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Get", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Get", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Get", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return errors.NotFound("config.Config.Get", "Not found")
	} else if err != nil {
		return errors.BadRequest("config.Config.Get", "read error: %v: %v", err, req.Namespace)
	}

	rsp.Value = &pb.Value{
		Data: string(ch[0].Value),
	}

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	values := config.NewJSONValues(ch[0].Value)

	// we just want to pass back bytes
	rsp.Value.Data = string(values.Get(req.Path).Bytes())
	dat := leavesToValues(rsp.Value.Data, !req.Secret)

	if req.Secret {
		dec, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", dat))
		if err != nil {
			return errors.InternalServerError("config.Config.Get", "Badly encoded secret")
		}
		rsp.Value.Data = decrypt(string(dec), c.secret)
	}

	switch v := dat.(type) {
	case string:
		rsp.Value.Data = v
	case nil:
		rsp.Value.Data = "null"
	case int64:
		rsp.Value.Data = fmt.Sprintf("%v", v)
	case bool:
		rsp.Value.Data = fmt.Sprintf("%v", v)
	default:
		bs, _ := json.Marshal(dat)
		rsp.Value.Data = string(bs)
	}
	return nil
}

func leavesToValues(data string, maskSecrets bool) interface{} {
	if data == "null" {
		return nil
	}
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return data
	}
	return traverse(m, maskSecrets)
}

func traverse(i interface{}, maskSecrets bool) interface{} {
	switch v := i.(type) {
	case map[string]interface{}:
		if val, ok := v["leaf"].(bool); ok && val {
			isSecret, isSecretOk := v["secret"].(bool)
			if isSecretOk && isSecret && maskSecrets {
				return "[secret]"
			}
			value, valueOk := v["value"].(string)
			if valueOk {
				return value
			}
			return ""
		}
		ret := map[string]interface{}{}
		for key, val := range v {
			ret[key] = traverse(val, maskSecrets)
		}
		return ret
	case []interface{}:
		for _, e := range v {
			ret := []interface{}{}
			ret = append(ret, traverse(e, maskSecrets))
			return ret
		}
	default:
		return i
	}
	return i
}

func (c *Config) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	if req.Value == nil {
		return errors.BadRequest("config.Config.Update", "invalid change")
	}
	ns := req.Namespace
	if len(ns) == 0 {
		ns = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, ns); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Update", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Update", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Update", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	dat := []byte{}
	if err == store.ErrNotFound {
		dat = []byte("{}")
	} else if err != nil {
		return errors.BadRequest("config.Config.Set", "read error: %v: %v", err, req.Namespace)
	}

	if len(ch) > 0 {
		dat = ch[0].Value
	}
	values := config.NewJSONValues(dat)

	data := req.Value.Data
	if req.Secret {
		data = string(base64.StdEncoding.EncodeToString([]byte(encrypt(data, c.secret))))
		// Need to save metainformation with secret values too
		values.Set(req.Path, map[string]interface{}{
			"secret": true,
			"value":  data,
			"leaf":   true,
		})
	} else {
		values.Set(req.Path, map[string]interface{}{
			"value": data,
			"leaf":  true,
		})
	}

	return store.Write(&store.Record{
		Key:   req.Namespace,
		Value: values.Bytes(),
	})
}

func (c *Config) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	ns := req.Namespace
	if len(ns) == 0 {
		ns = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, ns); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Delete", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Delete", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Delete", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return errors.NotFound("config.Config.Delete", "Not found")
	} else if err != nil {
		return errors.BadRequest("config.Config.Delete", "read error: %v: %v", err, req.Namespace)
	}

	values := config.NewJSONValues(ch[0].Value)

	values.Delete(req.Path)
	return store.Write(&store.Record{
		Key:   req.Namespace,
		Value: values.Bytes(),
	})
}
