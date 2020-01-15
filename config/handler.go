package config

import (
	"github.com/gogo/protobuf/proto"
	"github.com/micro/go-micro/store"
	"strings"
	"time"

	"github.com/micro/go-micro/config/source"
	mp "github.com/micro/go-micro/config/source/mucp/proto"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/config/db"
	"golang.org/x/net/context"
)

type Config struct{}

func (c *Config) Read(ctx context.Context, req *mp.ReadRequest, rsp *mp.ReadResponse) (err error) {
	if len(req.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Read", "invalid id")
		log.Error(err)
		return err
	}

	ch, err := db.Read(req.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "read error: %s", err)
		log.Error(err)
		return err
	}

	// Unmarshal value
	err = proto.Unmarshal(ch.Value, rsp.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %s", err)
		log.Error(err)
		return err
	}

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	rcc := rsp.Change.ChangeSet
	values, err := Values(&source.ChangeSet{
		Timestamp: time.Unix(rcc.Timestamp, 0),
		Data:      rcc.Data,
		Checksum:  rcc.Checksum,
		Format:    rcc.Format,
		Source:    rcc.Source,
	})
	if err != nil {
		err = errors.InternalServerError("go.micro.config.Read", err.Error())
		log.Error(err)
		return err
	}

	parts := strings.Split(req.Path, PathSplitter)

	// we just want to pass back bytes
	rsp.Change.ChangeSet.Data = values.Get(parts...).Bytes()

	return nil
}

func (c *Config) Create(ctx context.Context, req *mp.CreateRequest, rsp *mp.CreateResponse) (err error) {
	if req.Change == nil || req.Change.ChangeSet == nil {
		err = errors.BadRequest("go.micro.config.Create", "invalid change")
		log.Error(err)
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Create", "invalid id")
		log.Error(err)
		return err
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	record := &store.Record{}
	record.Value, err = proto.Marshal(req.Change)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Create", "marshal error: ", err)
		log.Error(err)
		return err
	}

	record.Key = req.Change.Key

	if err := db.Create(record); err != nil {
		err = errors.BadRequest("go.micro.config.Create", "create new into db error: ", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Update(ctx context.Context, req *mp.UpdateRequest, rsp *mp.UpdateResponse) (err error) {
	if req.Change == nil || req.Change.ChangeSet == nil {
		err = errors.BadRequest("go.micro.config.Update", "invalid change")
		log.Error(err)
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.config.Update", "invalid id")
		log.Error(err)
		return err
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// Get the current change set
	record, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "read old value error: ", err)
		log.Error(err)
		return err
	}

	ch := &mp.Change{}
	// Unmarshal value
	err = proto.Unmarshal(record.Value, ch)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %s", err)
		log.Error(err)
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
		vals, err := Values(&source.ChangeSet{Data: req.Change.ChangeSet.Data, Format: ch.ChangeSet.Format})
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "values error: %s", err)
			log.Error(err)
			return err
		}
		if err := vals.Get().Scan(&data); err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "scan data error: %s", err)
			log.Error(err)
			return err
		}

		// Get values from existing change
		values, err := Values(change)
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "get values from existing change error: %s", err)
			log.Error(err)
			return err
		}

		// Apply the data to the existing change
		values.Set(data, strings.Split(req.Change.Path, PathSplitter)...)

		// Create a new change
		newChange, err = Merge(&source.ChangeSet{Data: values.Bytes()})
		if err != nil {
			err = errors.InternalServerError("go.micro.config.Update", "create a new change error: %s", err)
			log.Error(err)
			return err
		}
	} else {
		// No path specified, business as usual
		newChange, err = Merge(change, &source.ChangeSet{
			Timestamp: time.Unix(req.Change.ChangeSet.Timestamp, 0),
			Data:      req.Change.ChangeSet.Data,
			Checksum:  req.Change.ChangeSet.Checksum,
			Source:    req.Change.ChangeSet.Source,
			Format:    req.Change.ChangeSet.Format,
		})
		if err != nil {
			err = errors.BadRequest("go.micro.srv.config.Update", "merge all error: ", err)
			log.Error(err)
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
		err = errors.BadRequest("go.micro.config.Update", "marshal error: ", err)
		log.Error(err)
		return err
	}

	if err := db.Update(record); err != nil {
		err = errors.BadRequest("go.micro.config.Update", "update into db error: ", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Delete(ctx context.Context, req *mp.DeleteRequest, rsp *mp.DeleteResponse) (err error) {
	if req.Change == nil {
		err = errors.BadRequest("go.micro.srv.Delete", "invalid change")
		log.Error(err)
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.Delete", "invalid id")
		log.Error(err)
		return err
	}

	if req.Change.ChangeSet == nil {
		req.Change.ChangeSet = &mp.ChangeSet{}
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// We're going to delete the record as we have no path and no data
	if len(req.Change.Path) == 0 {
		if err := db.Delete(req.Change.Key); err != nil {
			err = errors.BadRequest("go.micro.srv.Delete", "delete from db error: %s", err)
			log.Error(err)
			return err
		}
		return nil
	}

	// We've got a path. Let's update the required path

	// Get the current change set
	record, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "read old value error: ", err)
		log.Error(err)
		return err
	}

	ch := &mp.Change{}
	// Unmarshal value
	err = proto.Unmarshal(record.Value, ch)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %s", err)
		log.Error(err)
		return err
	}

	// Get the current config as values
	values, err := Values(&source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      ch.ChangeSet.Data,
		Checksum:  ch.ChangeSet.Checksum,
		Source:    ch.ChangeSet.Source,
		Format:    ch.ChangeSet.Format,
	})
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "Get the current config as values error: %s", err)
		log.Error(err)
		return err
	}

	// Delete at the given path
	values.Del(strings.Split(req.Change.Path, PathSplitter)...)

	// Create a change record from the values
	change, err := Merge(&source.ChangeSet{Data: values.Bytes()})
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "Create a change record from the values error: %s", err)
		log.Error(err)
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
		err = errors.BadRequest("go.micro.config.Update", "marshal error: ", err)
		log.Error(err)
		return err
	}

	if err := db.Update(record); err != nil {
		err = errors.BadRequest("go.micro.srv.Delete", "update record set to db error: %s", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &mp.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) List(ctx context.Context, req *mp.ListRequest, rsp *mp.ListResponse) error {
	list, err := db.List()
	if err != nil {
		err = errors.BadRequest("go.micro.config.List", "query value error: %s", err)
		log.Error(err)
		return err
	}

	for _, v := range list {

	}
	// Unmarshal value
	err := proto.Unmarshal(list, &rsp.Configs)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Read", "unmarshal value error: %s", err)
		log.Error(err)
		return err
	}
	return nil
}

func (c *Config) Watch(ctx context.Context, req *mp.WatchRequest, stream mp.Source_WatchStream) (err error) {
	if len(req.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.Watch", "invalid id")
		log.Error(err)
		return err
	}

	watch, err := Watch(req.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.srv.Watch", "watch error: %s", err)
		log.Error(err)
		return err
	}
	defer watch.Stop()

	for {
		ch, err := watch.Next()
		if err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.Watch", "listen the Next error: %s", err)
			log.Error(err)
			return err
		}

		if err := stream.Send(ch); err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.Watch", "send the Change error: %s", err)
			log.Error(err)
			return err
		}
	}
}
