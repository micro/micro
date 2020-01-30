package handler

import (
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/micro/go-micro/v2/client"
	cr "github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/reader/json"
	"github.com/micro/go-micro/v2/config/source"
	mp "github.com/micro/go-micro/v2/config/source/service/proto"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/config/db"
	"golang.org/x/net/context"
)

var (
	PathSplitter = "/"
	WatchTopic   = "go.micro.config.events"
	watchers     = make(map[string][]*watcher)

	// we now support json only
	reader = json.NewReader()
	mtx    sync.RWMutex
)

type Handler struct{}

func (c *Handler) Read(ctx context.Context, req *mp.ReadRequest, rsp *mp.ReadResponse) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	if len(req.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Read", "invalid id")
		return err
	}

	ch, err := db.Read(req.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "read error: %v", err)
		return err
	}
	rsp.Change = &mp.Change{}
	// Unmarshal value
	err = proto.Unmarshal(ch.Value, rsp.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
		return err
	}

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	rcc := rsp.Change.ChangeSet
	values, err := values(&source.ChangeSet{
		Timestamp: time.Unix(rcc.Timestamp, 0),
		Data:      rcc.Data,
		Checksum:  rcc.Checksum,
		Format:    rcc.Format,
		Source:    rcc.Source,
	})
	if err != nil {
		err = errors.InternalServerError("go.micro.config.Read", err.Error())
		return err
	}

	parts := strings.Split(req.Path, PathSplitter)

	// we just want to pass back bytes
	rsp.Change.ChangeSet.Data = values.Get(parts...).Bytes()

	return nil
}

func (c *Handler) Create(ctx context.Context, req *mp.CreateRequest, rsp *mp.CreateResponse) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	if req.Change == nil || req.Change.ChangeSet == nil {
		err = errors.BadRequest("go.micro.config.Create", "invalid change")
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Create", "invalid id")
		return err
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	record := &store.Record{}
	record.Value, err = proto.Marshal(req.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Create", "marshal error: %v", err)
		return err
	}

	record.Key = req.Change.Key

	if err := db.Create(record); err != nil {
		err = errors.BadRequest("go.micro.config.Create", "create new into db error: %v", err)
		return err
	}

	_ = publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Handler) Update(ctx context.Context, req *mp.UpdateRequest, rsp *mp.UpdateResponse) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	if req.Change == nil || req.Change.ChangeSet == nil {
		err = errors.BadRequest("go.micro.config.Update", "invalid change")
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Update", "invalid id")
		return err
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// Get the current change set
	record, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "read old value error: %v", err)
		return err
	}

	ch := &mp.Change{}
	// Unmarshal value
	err = proto.Unmarshal(record.Value, ch)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
		return err
	}

	chc := ch.ChangeSet
	change := &source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      chc.Data,
		Checksum:  chc.Checksum,
		Source:    chc.Source,
		Format:    chc.Format,
	}

	var newChange *source.ChangeSet

	// Set the change at a particular path
	if len(req.Change.Path) > 0 {
		// Unpack the data as a go type
		var data interface{}
		vals, err := values(&source.ChangeSet{Data: req.Change.ChangeSet.Data, Format: ch.ChangeSet.Format})
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "values error: %v", err)
			return err
		}
		if err := vals.Get().Scan(&data); err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "scan data error: %v", err)
			return err
		}

		// Get values from existing change
		values, err := values(change)
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "get values from existing change error: %v", err)
			return err
		}

		// Apply the data to the existing change
		values.Set(data, strings.Split(req.Change.Path, PathSplitter)...)

		// Create a new change
		newChange, err = merge(&source.ChangeSet{Data: values.Bytes()})
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "create a new change error: %v", err)
			return err
		}
	} else {
		// No path specified, business as usual
		newChange, err = merge(change, &source.ChangeSet{
			Timestamp: time.Unix(req.Change.ChangeSet.Timestamp, 0),
			Data:      req.Change.ChangeSet.Data,
			Checksum:  req.Change.ChangeSet.Checksum,
			Source:    req.Change.ChangeSet.Source,
			Format:    req.Change.ChangeSet.Format,
		})
		if err != nil {
			err = errors.BadRequest("go.micro.srv.config.Update", "merge all error: %v", err)
			return err
		}
	}

	// update change set
	req.Change.ChangeSet = &mp.ChangeSet{
		Timestamp: newChange.Timestamp.Unix(),
		Data:      newChange.Data,
		Checksum:  newChange.Checksum,
		Source:    newChange.Source,
		Format:    newChange.Format,
	}

	record.Value, err = proto.Marshal(req.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "marshal error: %v", err)
		return err
	}

	if err := db.Update(record); err != nil {
		err = errors.BadRequest("go.micro.config.Update", "update into db error: %v", err)
		return err
	}

	_ = publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Handler) Delete(ctx context.Context, req *mp.DeleteRequest, rsp *mp.DeleteResponse) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	if req.Change == nil {
		err = errors.BadRequest("go.micro.srv.Delete", "invalid change")
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.Delete", "invalid id")
		return err
	}

	if req.Change.ChangeSet == nil {
		req.Change.ChangeSet = &mp.ChangeSet{}
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// We're going to delete the record as we have no path and no data
	if len(req.Change.Path) == 0 {
		if err := db.Delete(req.Change.Key); err != nil {
			err = errors.BadRequest("go.micro.srv.Delete", "delete from db error: %v", err)
			log.Error(err)
			return err
		}
		return nil
	}

	// We've got a path. Let's update the required path

	// Get the current change set
	record, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "read old value error: %v", err)
		return err
	}

	ch := &mp.Change{}
	// Unmarshal value
	err = proto.Unmarshal(record.Value, ch)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
		log.Error(err)
		return err
	}

	// Get the current config as values
	values, err := values(&source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      ch.ChangeSet.Data,
		Checksum:  ch.ChangeSet.Checksum,
		Source:    ch.ChangeSet.Source,
		Format:    ch.ChangeSet.Format,
	})
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "Get the current config as values error: %v", err)
		return err
	}

	// Delete at the given path
	values.Del(strings.Split(req.Change.Path, PathSplitter)...)

	// Create a change record from the values
	change, err := merge(&source.ChangeSet{Data: values.Bytes()})
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "Create a change record from the values error: %v", err)
		return err
	}

	// Update change set
	req.Change.ChangeSet = &mp.ChangeSet{
		Timestamp: change.Timestamp.Unix(),
		Data:      change.Data,
		Checksum:  change.Checksum,
		Format:    change.Format,
		Source:    change.Source,
	}

	record.Value, err = proto.Marshal(req.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "marshal error: %v", err)
		return err
	}

	if err := db.Update(record); err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "update record set to db error: %v", err)
		return err
	}

	_ = publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Handler) List(ctx context.Context, req *mp.ListRequest, rsp *mp.ListResponse) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	list, err := db.List()
	if err != nil {
		err = errors.BadRequest("go.micro.config.List", "query value error: %v", err)
		return err
	}

	for _, v := range list {
		ch := &mp.Change{}
		err := proto.Unmarshal(v.Value, ch)
		if err != nil {
			err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %v", err)
			return err
		}
		rsp.Values = append(rsp.Values, ch)
	}

	return nil
}

func (c *Handler) Watch(ctx context.Context, req *mp.WatchRequest, stream mp.Config_WatchStream) (err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	if len(req.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.Watch", "invalid id")
		return err
	}

	watch, err := Watch(req.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Watch", "watch error: %v", err)
		return err
	}
	defer watch.Stop()

	for {
		ch, err := watch.Next()
		if err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.Watch", "listen the Next error: %v", err)
			return err
		}

		if err := stream.Send(ch); err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.Watch", "send the Change error: %v", err)
			return err
		}
	}
}

// Used as a subscriber between config services for events
func Watcher(ctx context.Context, ch *mp.WatchResponse) error {
	mtx.RLock()
	for _, sub := range watchers[ch.Key] {
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
func publish(ctx context.Context, ch *mp.WatchResponse) error {
	req := client.NewMessage(WatchTopic, ch)
	return client.Publish(ctx, req)
}
