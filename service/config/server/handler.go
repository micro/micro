package server

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	cr "github.com/micro/go-micro/v3/config/reader"
	jr "github.com/micro/go-micro/v3/config/reader/json"
	"github.com/micro/go-micro/v3/config/source"
	"github.com/micro/micro/v3/internal/namespace"
	pb "github.com/micro/micro/v3/proto/config"
	muclient "github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
)

const (
	defaultNamespace = "micro"
	pathSplitter     = "."
)

var (
	watchTopic = "config.events"
	watchers   = make(map[string][]*watcher)

	// we now support json only
	reader = jr.NewReader()
	mtx    sync.RWMutex
)

type Config struct{}

func (c *Config) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
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

	rsp.Change = new(pb.Change)

	// Unmarshal value
	if err = json.Unmarshal(ch[0].Value, rsp.Change); err != nil {
		return errors.BadRequest("config.Config.Read", "unmarshal value error: %v", err)
	}

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	rc := rsp.Change.ChangeSet

	// generate reader.Values from the changeset
	values, err := values(&source.ChangeSet{
		Timestamp: time.Unix(rc.Timestamp, 0),
		Data:      []byte(rc.Data),
		Checksum:  rc.Checksum,
		Format:    rc.Format,
		Source:    rc.Source,
	})
	if err != nil {
		return errors.InternalServerError("config.Config.Read", err.Error())
	}

	// peel apart the path
	parts := strings.Split(req.Path, pathSplitter)

	// we just want to pass back bytes
	rsp.Change.ChangeSet.Data = string(values.Get(parts...).Bytes())

	return nil
}

func (c *Config) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Change == nil || req.Change.ChangeSet == nil {
		return errors.BadRequest("config.Config.Create", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Create", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Create", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Create", err.Error())
	}

	if len(req.Change.Path) > 0 {
		vals, err := values(&source.ChangeSet{
			Format: "json",
		})
		if err != nil {
			return errors.InternalServerError("config.Config.Create", err.Error())
		}

		// peel apart the path
		parts := strings.Split(req.Change.Path, pathSplitter)
		// set the values
		vals.Set(req.Change.ChangeSet.Data, parts...)
		// change the changeset value
		req.Change.ChangeSet.Data = string(vals.Bytes())
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	record := &store.Record{Key: req.Change.Namespace}

	var err error
	record.Value, err = json.Marshal(req.Change)
	if err != nil {
		return errors.BadRequest("config.Config.Create", "marshal error: %v", err)
	}

	if err := store.Write(record); err != nil {
		return errors.BadRequest("config.Config.Create", "create new into db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if req.Change == nil || req.Change.ChangeSet == nil {
		return errors.BadRequest("config.Config.Update", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Update", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Update", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Update", err.Error())
	}

	// set the changeset timestamp
	req.Change.ChangeSet.Timestamp = time.Now().Unix()
	oldCh := &pb.Change{}

	// Get the current change set
	var record *store.Record
	records, err := store.Read(req.Change.Namespace)
	if err != nil {
		if err.Error() != "not found" {
			return errors.BadRequest("config.Config.Update", "read old value error: %v", err)
		}
		// create new record
		record = new(store.Record)
		record.Key = req.Change.Namespace
	} else {
		// Unmarshal value
		if err := json.Unmarshal(records[0].Value, oldCh); err != nil {
			return errors.BadRequest("config.Config.Update", "unmarshal value error: %v", err)
		}
		record = records[0]
	}

	// generate a new base changeset
	changeSet := &source.ChangeSet{
		Format: "json",
		Data:   []byte(`{}`),
	}

	if oldCh.ChangeSet != nil {
		changeSet = &source.ChangeSet{
			Timestamp: time.Unix(oldCh.ChangeSet.Timestamp, 0),
			Data:      []byte(oldCh.ChangeSet.Data),
			Checksum:  oldCh.ChangeSet.Checksum,
			Source:    oldCh.ChangeSet.Source,
			Format:    oldCh.ChangeSet.Format,
		}
	}

	var newChange *source.ChangeSet

	// Set the change at a particular path
	if len(req.Change.Path) > 0 {
		// Get values from existing change
		values, err := values(changeSet)
		if err != nil {
			return errors.InternalServerError("config.Config.Update", "error getting existing change: %v", err)
		}

		// Apply the data to the existing change
		values.Set(req.Change.ChangeSet.Data, strings.Split(req.Change.Path, pathSplitter)...)

		// Create a new change
		newChange, err = merge(&source.ChangeSet{Data: values.Bytes()})
		if err != nil {
			return errors.InternalServerError("config.Config.Update", "create a new change error: %v", err)
		}
	} else {
		// No path specified, business as usual
		newChange, err = merge(changeSet, &source.ChangeSet{
			Timestamp: time.Unix(req.Change.ChangeSet.Timestamp, 0),
			Data:      []byte(req.Change.ChangeSet.Data),
			Checksum:  req.Change.ChangeSet.Checksum,
			Source:    req.Change.ChangeSet.Source,
			Format:    req.Change.ChangeSet.Format,
		})
		if err != nil {
			return errors.BadRequest("config.Config.Update", "merge all error: %v", err)
		}
	}

	// update change set
	req.Change.ChangeSet = &pb.ChangeSet{
		Timestamp: newChange.Timestamp.Unix(),
		Data:      string(newChange.Data),
		Checksum:  newChange.Checksum,
		Source:    newChange.Source,
		Format:    newChange.Format,
	}

	record.Value, err = json.Marshal(req.Change)
	if err != nil {
		return errors.BadRequest("config.Config.Update", "marshal error: %v", err)
	}

	if err := store.Write(record); err != nil {
		return errors.BadRequest("config.Config.Update", "update into db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if req.Change == nil {
		return errors.BadRequest("config.Config.Delete", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Delete", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Delete", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Delete", err.Error())
	}

	if req.Change.ChangeSet == nil {
		req.Change.ChangeSet = &pb.ChangeSet{}
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// We're going to delete the record as we have no path and no data
	if len(req.Change.Path) == 0 {
		if err := store.Delete(req.Change.Namespace); err != nil {
			return errors.BadRequest("config.Config.Delete", "delete from db error: %v", err)
		}
		return nil
	}

	// We've got a path. Let's update the required path

	// Get the current change set
	records, err := store.Read(req.Change.Namespace)
	if err != nil {
		if err.Error() != "not found" {
			return errors.BadRequest("config.Config.Delete", "read old value error: %v", err)
		}
		return nil
	}

	ch := &pb.Change{}
	// Unmarshal value
	if err := json.Unmarshal(records[0].Value, ch); err != nil {
		return errors.BadRequest("config.Config.Delete", "unmarshal value error: %v", err)
	}

	// Get the current config as values
	values, err := values(&source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      []byte(ch.ChangeSet.Data),
		Checksum:  ch.ChangeSet.Checksum,
		Source:    ch.ChangeSet.Source,
		Format:    ch.ChangeSet.Format,
	})
	if err != nil {
		return errors.BadRequest("config.Config.Delete", "Get the current config as values error: %v", err)
	}

	// Delete at the given path
	values.Del(strings.Split(req.Change.Path, pathSplitter)...)

	// Create a change record from the values
	change, err := merge(&source.ChangeSet{Data: values.Bytes()})
	if err != nil {
		return errors.BadRequest("config.Config.Delete", "Create a change record from the values error: %v", err)
	}

	// Update change set
	req.Change.ChangeSet = &pb.ChangeSet{
		Timestamp: change.Timestamp.Unix(),
		Data:      string(change.Data),
		Checksum:  change.Checksum,
		Format:    change.Format,
		Source:    change.Source,
	}

	records[0].Value, err = json.Marshal(req.Change)
	if err != nil {
		return errors.BadRequest("config.Config.Delete", "marshal error: %v", err)
	}

	if err := store.Write(records[0]); err != nil {
		return errors.BadRequest("config.Config.Delete", "update record set to db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) (err error) {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.List", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.List", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.List", err.Error())
	}

	list, err := store.List(store.Prefix(req.Namespace))
	if err != nil {
		return errors.BadRequest("config.Config.List", "query value error: %v", err)
	}

	// TODO: optimise filtering for prefix listing
	for _, v := range list {
		rec, err := store.Read(v)
		if err != nil {
			return err
		}

		ch := &pb.Change{}
		if err := json.Unmarshal(rec[0].Value, ch); err != nil {
			return errors.BadRequest("config.Config.List", "unmarshal value error: %v", err)
		}

		if ch.ChangeSet != nil {
			ch.ChangeSet.Data = string(ch.ChangeSet.Data)
		}

		rsp.Values = append(rsp.Values, ch)
	}

	return nil
}

func (c *Config) Watch(ctx context.Context, req *pb.WatchRequest, stream pb.Config_WatchStream) error {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("config.Config.Watch", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("config.Config.Watch", err.Error())
	} else if err != nil {
		return errors.InternalServerError("config.Config.Watch", err.Error())
	}

	watch, err := Watch(req.Namespace)
	if err != nil {
		return errors.BadRequest("config.Config.Watch", "watch error: %v", err)
	}
	defer watch.Stop()

	go func() {
		select {
		case <-ctx.Done():
			watch.Stop()
			stream.Close()
		}
	}()

	for {
		ch, err := watch.Next()
		if err != nil {
			return errors.BadRequest("config.Config.Watch", "listen the Next error: %v", err)
		}
		if ch.ChangeSet != nil {
			ch.ChangeSet.Data = string(ch.ChangeSet.Data)
		}
		if err := stream.Send(ch); err != nil {
			return errors.BadRequest("config.Config.Watch", "send the Change error: %v", err)
		}
	}
}

// Used as a subscriber between config services for events
func Watcher(ctx context.Context, ch *pb.WatchResponse) error {
	mtx.RLock()
	for _, sub := range watchers[ch.Namespace] {
		select {
		case sub.next <- ch:
		case <-time.After(time.Millisecond * 100):
		}
	}
	mtx.RUnlock()
	return nil
}

func merge(ch ...*source.ChangeSet) (*source.ChangeSet, error) {
	return reader.Merge(ch...)
}

func values(ch *source.ChangeSet) (cr.Values, error) {
	return reader.Values(ch)
}

// publish a change
func publish(ctx context.Context, ch *pb.WatchResponse) error {
	req := muclient.NewMessage(watchTopic, ch)
	return muclient.Publish(ctx, req)
}
