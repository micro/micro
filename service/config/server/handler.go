package server

import (
	"context"
	"strings"
	"sync"

	"github.com/micro/go-micro/v3/config"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/config"
	"github.com/micro/micro/v3/service/errors"
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

type Config struct{}

func (c *Config) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Read", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Read", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Read", err.Error())
	}

	ch, err := store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return errors.NotFound("config.Config.Read", "Not found")
	} else if err != nil {
		return errors.BadRequest("config.Config.Read", "read error: %v: %v", err, req.Namespace)
	}

	rsp.Value = &pb.Value{
		Data: string(ch[0].Value),
	}

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	values, err := config.NewJSONValues(ch[0].Value)
	if err != nil {
		return err
	}

	// peel apart the path
	parts := strings.Split(req.Path, pathSplitter)

	// we just want to pass back bytes
	rsp.Value.Data = string(values.Get(parts...).Bytes())

	return nil
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
	if err == store.ErrNotFound {
		return errors.NotFound("config.Config.Read", "Not found")
	} else if err != nil {
		return errors.BadRequest("config.Config.Read", "read error: %v: %v", err, req.Namespace)
	}

	values, err := config.NewJSONValues(ch[0].Value)
	if err != nil {
		return err
	}

	// peel apart the path
	parts := strings.Split(req.Path, pathSplitter)
	values.Set(req.Value.Data, parts...)
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
		return errors.NotFound("config.Config.Read", "Not found")
	} else if err != nil {
		return errors.BadRequest("config.Config.Read", "read error: %v: %v", err, req.Namespace)
	}

	values, err := config.NewJSONValues(ch[0].Value)
	if err != nil {
		return err
	}

	// peel apart the path
	parts := strings.Split(req.Path, pathSplitter)
	values.Delete(parts...)
	return store.Write(&store.Record{
		Key:   req.Namespace,
		Value: values.Bytes(),
	})
}
