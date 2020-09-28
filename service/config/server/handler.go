package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/micro/go-micro/v3/config"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/config"
	merrors "github.com/micro/micro/v3/service/errors"
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
	var dec []byte
	var err error
	if len(key) == 0 {
		logger.Warn("No encryption key provided")
	} else {
		dec, err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			logger.Warnf("Error decoding key: %v", err)
		}
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
		return merrors.Forbidden("config.Config.Get", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return merrors.Unauthorized("config.Config.Get", err.Error())
	} else if err != nil {
		return merrors.InternalServerError("config.Config.Get", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return merrors.NotFound("config.Config.Get", "Not found")
	} else if err != nil {
		return merrors.BadRequest("config.Config.Get", "read error: %v: %v", err, req.Namespace)
	}

	// get secret from options
	var secret bool
	if req.GetOptions().GetSecret() {
		secret = true
	}

	rsp.Value = &pb.Value{}

	values := config.NewJSONValues(ch[0].Value)

	// we just want to pass back bytes
	bytes := values.Get(req.Path).Bytes()
	dat, err := leavesToValues(string(bytes), secret, string(c.secret))
	if err != nil {
		return merrors.InternalServerError("config.config.Get", "Error in config structure: %v", err)
	}

	response, _ := json.Marshal(dat)
	rsp.Value.Data = string(response)

	return nil
}

func leavesToValues(data string, decodeSecrets bool, encryptionKey string) (interface{}, error) {
	var m interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		return m, err
	}
	return traverse(m, decodeSecrets, encryptionKey)
}

func traverse(i interface{}, decodeSecrets bool, encryptionKey string) (interface{}, error) {
	switch v := i.(type) {
	case map[string]interface{}:
		if val, ok := v["leaf"].(bool); ok && val {
			isSecret, isSecretOk := v["secret"].(bool)
			if isSecretOk && isSecret && !decodeSecrets {
				return "[secret]", nil
			}
			marshalledValue, ok := v["value"].(string)
			if !ok {
				return nil, fmt.Errorf("Value field in leaf %v can't be found", v)
			}
			if isSecretOk && isSecret {
				if len(encryptionKey) == 0 {
					return nil, errors.New("Can't decode secret: secret key is not set")
				}
				dec, err := base64.StdEncoding.DecodeString(marshalledValue)
				if err != nil {
					fmt.Println("marshalled value", marshalledValue)
					return nil, errors.New("Badly encoded secret")
				}
				decrypted, err := decrypt(string(dec), []byte(encryptionKey))
				if err != nil {
					return nil, fmt.Errorf("Failed to decrypt: %v", err)
				}
				marshalledValue = decrypted
			}
			var value interface{}
			err := json.Unmarshal([]byte(marshalledValue), &value)
			return value, err
		}
		ret := map[string]interface{}{}
		for key, val := range v {
			value, err := traverse(val, decodeSecrets, encryptionKey)
			if err != nil {
				return ret, err
			}
			ret[key] = value
		}
		return ret, nil
	case []interface{}:
		for _, e := range v {
			ret := []interface{}{}
			value, err := traverse(e, decodeSecrets, encryptionKey)
			if err != nil {
				return ret, err
			}
			ret = append(ret, value)
			return ret, nil
		}
	default:
		return i, nil
	}
	return i, nil
}

func (c *Config) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	if req.Value == nil {
		return merrors.BadRequest("config.Config.Update", "invalid change")
	}
	ns := req.Namespace
	if len(ns) == 0 {
		ns = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, ns); err == namespace.ErrForbidden {
		return merrors.Forbidden("config.Config.Update", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return merrors.Unauthorized("config.Config.Update", err.Error())
	} else if err != nil {
		return merrors.InternalServerError("config.Config.Update", err.Error())
	}

	ch, err := store.Read(ns)
	dat := []byte{}
	if err == store.ErrNotFound {
		dat = []byte("{}")
	} else if err != nil {
		return merrors.BadRequest("config.Config.Set", "read error: %v: %v", err, ns)
	}

	if len(ch) > 0 {
		dat = ch[0].Value
	}
	values := config.NewJSONValues(dat)

	// get secret from options
	var secret bool
	if req.GetOptions().GetSecret() {
		secret = true
	}

	// req.Value.Data is a json encoded value
	data := req.Value.Data
	if secret {
		if len(c.secret) == 0 {
			return merrors.InternalServerError("config.Config.Set", "Can't encode secret: secret key is not set")
		}
		encrypted, err := encrypt(data, c.secret)
		if err != nil {
			return merrors.InternalServerError("config.Config.Set", "Failed to encrypt", err)
		}
		data = string(base64.StdEncoding.EncodeToString([]byte(encrypted)))
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
		return merrors.Forbidden("config.Config.Delete", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return merrors.Unauthorized("config.Config.Delete", err.Error())
	} else if err != nil {
		return merrors.InternalServerError("config.Config.Delete", err.Error())
	}

	ch, err := store.Read(ns)
	if err == store.ErrNotFound {
		return merrors.NotFound("config.Config.Delete", "Not found")
	} else if err != nil {
		return merrors.BadRequest("config.Config.Delete", "read error: %v: %v", err, ns)
	}

	values := config.NewJSONValues(ch[0].Value)

	values.Delete(req.Path)
	return store.Write(&store.Record{
		Key:   ns,
		Value: values.Bytes(),
	})
}
