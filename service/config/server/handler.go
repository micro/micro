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
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/internal/namespace"
	muclient "github.com/micro/micro/v3/service/client"
	pb "github.com/micro/micro/v3/service/config/proto"
)

const (
	defaultNamespace = "micro"
	pathSplitter     = "."
)

var (
	watchTopic = "go.micro.config.events"
	watchers   = make(map[string][]*watcher)

	// we now support json only
	reader = jr.NewReader()
	mtx    sync.RWMutex
)

type config struct {
	Store store.Store
}

func (c *config) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	ch, err := c.Store.Read(req.Namespace)
	if err == store.ErrNotFound {
		return errors.NotFound("go.micro.config.Read", "Not found")
	} else if err != nil {
		return errors.BadRequest("go.micro.config.Read", "read error: %v: %v", err, req.Namespace)
	}

	rsp.Change = new(pb.Change)

	// Unmarshal value
	if err = json.Unmarshal(ch[0].Value, rsp.Change); err != nil {
		return errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
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
		return errors.InternalServerError("go.micro.config.Read", err.Error())
	}

	// peel apart the path
	parts := strings.Split(req.Path, pathSplitter)

	// we just want to pass back bytes
	rsp.Change.ChangeSet.Data = string(values.Get(parts...).Bytes())

	return nil
}

func (c *config) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Change == nil || req.Change.ChangeSet == nil {
		return errors.BadRequest("go.micro.config.Create", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	if len(req.Change.Path) > 0 {
		vals, err := values(&source.ChangeSet{
			Format: "json",
		})
		if err != nil {
			return errors.InternalServerError("go.micro.config.Create", err.Error())
		}

		// peel apart the path
		parts := strings.Split(req.Change.Path, pathSplitter)
		// set the values
		vals.Set(req.Change.ChangeSet.Data, parts...)
		// change the changeset value
		req.Change.ChangeSet.Data = string(vals.Bytes())
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	record := &store.Record{
		Key: req.Change.Namespace,
	}

	var err error
	record.Value, err = json.Marshal(req.Change)
	if err != nil {
		return errors.BadRequest("go.micro.config.Create", "marshal error: %v", err)
	}

	if err := c.Store.Write(record); err != nil {
		return errors.BadRequest("go.micro.config.Create", "create new into db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *config) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if req.Change == nil || req.Change.ChangeSet == nil {
		return errors.BadRequest("go.micro.config.Update", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	// set the changeset timestamp
	req.Change.ChangeSet.Timestamp = time.Now().Unix()
	oldCh := &pb.Change{}

	// Get the current change set
	var record *store.Record
	records, err := c.Store.Read(req.Change.Namespace)
	if err != nil {
		if err.Error() != "not found" {
			return errors.BadRequest("go.micro.config.Update", "read old value error: %v", err)
		}
		// create new record
		record = new(store.Record)
		record.Key = req.Change.Namespace
	} else {
		// Unmarshal value
		if err := json.Unmarshal(records[0].Value, oldCh); err != nil {
			return errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
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
			return errors.InternalServerError("go.micro.config.Update", "error getting existing change: %v", err)
		}

		// Apply the data to the existing change
		values.Set(req.Change.ChangeSet.Data, strings.Split(req.Change.Path, pathSplitter)...)

		// Create a new change
		newChange, err = merge(&source.ChangeSet{Data: values.Bytes()})
		if err != nil {
			return errors.InternalServerError("go.micro.config.Update", "create a new change error: %v", err)
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
			return errors.BadRequest("go.micro.srv.config.Update", "merge all error: %v", err)
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
		return errors.BadRequest("go.micro.config.Update", "marshal error: %v", err)
	}

	if err := c.Store.Write(record); err != nil {
		return errors.BadRequest("go.micro.config.Update", "update into db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *config) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if req.Change == nil {
		return errors.BadRequest("go.micro.srv.Delete", "invalid change")
	}
	if len(req.Change.Namespace) == 0 {
		req.Change.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Change.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	if req.Change.ChangeSet == nil {
		req.Change.ChangeSet = &pb.ChangeSet{}
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// We're going to delete the record as we have no path and no data
	if len(req.Change.Path) == 0 {
		if err := c.Store.Delete(req.Change.Namespace); err != nil {
			return errors.BadRequest("go.micro.srv.Delete", "delete from db error: %v", err)
		}
		return nil
	}

	// We've got a path. Let's update the required path

	// Get the current change set
	records, err := c.Store.Read(req.Change.Namespace)
	if err != nil {
		if err.Error() != "not found" {
			return errors.BadRequest("go.micro.srv.Delete", "read old value error: %v", err)
		}
		return nil
	}

	ch := &pb.Change{}
	// Unmarshal value
	if err := json.Unmarshal(records[0].Value, ch); err != nil {
		return errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
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
		return errors.BadRequest("go.micro.srv.Delete", "Get the current config as values error: %v", err)
	}

	// Delete at the given path
	values.Del(strings.Split(req.Change.Path, pathSplitter)...)

	// Create a change record from the values
	change, err := merge(&source.ChangeSet{Data: values.Bytes()})
	if err != nil {
		return errors.BadRequest("go.micro.srv.Delete", "Create a change record from the values error: %v", err)
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
		return errors.BadRequest("go.micro.config.Update", "marshal error: %v", err)
	}

	if err := c.Store.Write(records[0]); err != nil {
		return errors.BadRequest("go.micro.srv.Delete", "update record set to db error: %v", err)
	}

	_ = publish(ctx, &pb.WatchResponse{Namespace: req.Change.Namespace, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *config) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) (err error) {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	list, err := c.Store.List(store.ListPrefix(req.Namespace))
	if err != nil {
		return errors.BadRequest("go.micro.config.List", "query value error: %v", err)
	}

	// TODO: optimise filtering for prefix listing
	for _, v := range list {
		rec, err := c.Store.Read(v)
		if err != nil {
			return err
		}

		ch := &pb.Change{}
		if err := json.Unmarshal(rec[0].Value, ch); err != nil {
			return errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
		}

		if ch.ChangeSet != nil {
			ch.ChangeSet.Data = string(ch.ChangeSet.Data)
		}

		rsp.Values = append(rsp.Values, ch)
	}

	return nil
}

func (c *config) Watch(ctx context.Context, req *pb.WatchRequest, stream pb.Config_WatchStream) error {
	if len(req.Namespace) == 0 {
		req.Namespace = defaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.config", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.config", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.config", err.Error())
	}

	watch, err := Watch(req.Namespace)
	if err != nil {
		return errors.BadRequest("go.micro.srv.Watch", "watch error: %v", err)
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
			return errors.BadRequest("go.micro.srv.Watch", "listen the Next error: %v", err)
		}
		if ch.ChangeSet != nil {
			ch.ChangeSet.Data = string(ch.ChangeSet.Data)
		}
		if err := stream.Send(ch); err != nil {
			return errors.BadRequest("go.micro.srv.Watch", "send the Change error: %v", err)
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
