package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/config"
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

	var bs []byte
	if len(req.Path) > 0 {
		bs = values.Get(req.Path).Bytes()
	} else {
		bs = values.Bytes()
	}
	dat, err := leavesToValues(string(bs), secret, string(c.secret))
	if err != nil {
		return merrors.InternalServerError("config.config.Get", "Error in config structure: %v", err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(dat)
	if err != nil {
		return merrors.BadRequest("config.Config.Get", "JSOn encode error: %v", err)
	}
	rsp.Value.Data = strings.TrimSpace(buf.String())

	return nil
}

// Read method is only here for backwards compatibility
func (c *Config) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	logger.Info("doing config read", req.Path, req.Namespace)
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return merrors.Forbidden("config.Config.Read", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return merrors.Unauthorized("config.Config.Read", err.Error())
	} else if err != nil {
		return merrors.InternalServerError("config.Config.Read", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return merrors.NotFound("config.Config.Read", "Not found")
	} else if err != nil {
		return merrors.BadRequest("config.Config.Read", "read error: %v: %v", err, req.Namespace)
	}

	rsp.Change = &pb.Change{
		Namespace: req.Namespace,
		Path:      req.Path,
		ChangeSet: &pb.ChangeSet{},
	}

	values := config.NewJSONValues(ch[0].Value)

	var bs []byte
	if len(req.Path) > 0 {
		bs = values.Get(req.Path).Bytes()
	} else {
		bs = values.Bytes()
	}

	dat, err := leavesToValues(string(bs), false, string(c.secret))
	if err != nil {
		return merrors.InternalServerError("config.config.Read", "Error in config structure: %v", err)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(dat)
	if err != nil {
		return merrors.BadRequest("config.Config.Read", "JSOn encode error: %v", err)
	}
	rsp.Change.ChangeSet.Data = strings.TrimSpace(buf.String())
	rsp.Change.ChangeSet.Format = "json"

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

func traverseMaps(m map[string]interface{}, paths []string, callback func(path string, value interface{}) error) error {
	for k, v := range m {
		val, ok := v.(map[string]interface{})
		if !ok {
			err := callback(strings.Join(append(paths, k), "."), v)
			if err != nil {
				return err
			}
			continue
		}
		err := traverseMaps(val, append(paths, k), callback)
		if err != nil {
			return err
		}
	}
	return nil
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

	var secret bool
	if req.GetOptions().GetSecret() {
		secret = true
	}

	// req.Value.Data is a json encoded value
	data := req.Value.Data
	var i interface{}
	err = json.Unmarshal([]byte(data), &i)
	if err != nil {
		return merrors.BadRequest("config.Config.Set", "Request is invalid JSON: %v", err)
	}
	m, ok := i.(map[string]interface{})
	// If it's a map, we do a merge
	if ok {
		// Need to nuke top level metadata as traverseMaps won't handle this
		cleanNode(values, req.Path)
		err = traverseMaps(m, strings.Split(req.Path, "."), func(p string, value interface{}) error {
			val, err := json.Marshal(value)
			if err != nil {
				return err
			}
			return c.setValue(values, secret, p, string(val))
		})
	} else {
		err = c.setValue(values, secret, req.Path, data)
	}

	return store.Write(&store.Record{
		Key:   req.Namespace,
		Value: values.Bytes(),
	})
}

func cleanNode(values *config.JSONValues, path string) {
	// Whatever the new value is being set, we need to delete
	// old metadat to prevent making a mess and introducing weird bugs,
	// ie. see `TestConfig/Test_plain_old_type_being_overwritten_by_map`
	values.Delete(path + ".leaf")
	values.Delete(path + ".value")
	values.Delete(path + ".secret")
}

func (c *Config) setValue(values *config.JSONValues, secret bool, path, data string) error {
	cleanNode(values, path)
	if secret {
		if len(c.secret) == 0 {
			return merrors.InternalServerError("config.Config.Set", "Can't encode secret: secret key is not set")
		}
		encrypted, err := encrypt(data, c.secret)
		if err != nil {
			return merrors.InternalServerError("config.Config.Set", "Failed to encrypt: %v", err)
		}
		data = string(base64.StdEncoding.EncodeToString([]byte(encrypted)))
		// Need to save metainformation with secret values too
		values.Set(path, map[string]interface{}{
			"secret": true,
			"value":  data,
			"leaf":   true,
		})
	} else {
		values.Set(path, map[string]interface{}{
			"value": data,
			"leaf":  true,
		})
	}
	return nil
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
