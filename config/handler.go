package config

import (
	"strings"
	"time"

	"github.com/micro/go-micro/config/source"
	proto "github.com/micro/go-micro/config/source/mucp/proto"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/micro/config/db"
	"golang.org/x/net/context"
)

type Config struct{}

func (c *Config) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) (err error) {
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
	// Set response
	rsp.Change = ch

	// if dont need path, we return all of the data
	if len(req.Path) == 0 {
		return nil
	}

	values, err := Values(&source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      ch.ChangeSet.Data,
		Checksum:  ch.ChangeSet.Checksum,
		Format:    ch.ChangeSet.Format,
		Source:    ch.ChangeSet.Source,
	})
	if err != nil {
		err = errors.InternalServerError("go.micro.srv.config.Read", err.Error())
		log.Error(err)
		return err
	}

	parts := strings.Split(req.Path, PathSplitter)

	// we just want to pass back bytes
	rsp.Change.ChangeSet.Data = values.Get(parts...).Bytes()

	return nil
}

func (c *Config) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) (err error) {
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

	if err := db.Create(req.Change); err != nil {
		err = errors.BadRequest("go.micro.config.Create", "create new into db error: ", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &proto.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) (err error) {
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
	ch, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.config.Update", "read old value error: ", err)
		log.Error(err)
		return err
	}

	csFormat := ch.ChangeSet.Format

	change := &source.ChangeSet{
		Timestamp: time.Unix(ch.ChangeSet.Timestamp, 0),
		Data:      ch.ChangeSet.Data,
		Checksum:  ch.ChangeSet.Checksum,
		Source:    ch.ChangeSet.Source,
		Format:    csFormat,
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
			Format:    csFormat,
		})
		if err != nil {
			err = errors.BadRequest("go.micro.srv.config.Update", "merge all error: ", err)
			log.Error(err)
			return err
		}
	}

	// update change set
	req.Change.ChangeSet = &proto.ChangeSet{
		Timestamp: newChange.Timestamp.Unix(),
		Data:      newChange.Data,
		Checksum:  newChange.Checksum,
		Source:    newChange.Source,
		Format:    csFormat,
	}

	if err := db.Update(req.Change); err != nil {
		err = errors.BadRequest("go.micro.srv.config.Update", "update into db error: ", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &proto.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) (err error) {
	if req.Change == nil {
		err = errors.BadRequest("go.micro.srv.config.Delete", "invalid change")
		log.Error(err)
		return err
	}

	if len(req.Change.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.config.Delete", "invalid id")
		log.Error(err)
		return err
	}

	if req.Change.ChangeSet == nil {
		req.Change.ChangeSet = &proto.ChangeSet{}
	}

	req.Change.ChangeSet.Timestamp = time.Now().Unix()

	// We're going to delete the record as we have no path and no data
	if len(req.Change.Path) == 0 {
		if err := db.Delete(req.Change); err != nil {
			err = errors.BadRequest("go.micro.srv.config.Delete", "delete from db error: %s", err)
			log.Error(err)
			return err
		}
		return nil
	}

	// We've got a path. Let's update the required path

	// Get the current change set
	ch, err := db.Read(req.Change.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.srv.config.Delete", "read the old from db error: %s", err)
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
		err = errors.BadRequest("go.micro.srv.config.Delete", "Get the current config as values error: %s", err)
		log.Error(err)
		return err
	}

	// Delete at the given path
	values.Del(strings.Split(req.Change.Path, PathSplitter)...)

	// Create a change record from the values
	change, err := Merge(&source.ChangeSet{Data: values.Bytes()})
	if err != nil {
		err = errors.BadRequest("go.micro.srv.config.Delete", "Create a change record from the values error: %s", err)
		log.Error(err)
		return err
	}

	// Update change set
	req.Change.ChangeSet = &proto.ChangeSet{
		Timestamp: change.Timestamp.Unix(),
		Data:      change.Data,
		Checksum:  change.Checksum,
		Format:    change.Format,
		Source:    change.Source,
	}

	if err := db.Update(req.Change); err != nil {
		err = errors.BadRequest("go.micro.srv.config.Delete", "update record set to db error: %s", err)
		log.Error(err)
		return err
	}

	_ = Publish(ctx, &proto.WatchResponse{Key: req.Change.Key, ChangeSet: req.Change.ChangeSet})

	return nil
}

func (c *Config) Watch(ctx context.Context, req *proto.WatchRequest, stream proto.Source_WatchStream) (err error) {
	if len(req.Key) == 0 {
		err = errors.BadRequest("go.micro.srv.config.Watch", "invalid id")
		log.Error(err)
		return err
	}

	watch, err := Watch(req.Key)
	if err != nil {
		err = errors.BadRequest("go.micro.srv.config.Watch", "watch error: %s", err)
		log.Error(err)
		return err
	}
	defer watch.Stop()

	for {
		ch, err := watch.Next()
		if err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.config.Watch", "listen the Next error: %s", err)
			log.Error(err)
			return err
		}

		if err := stream.Send(ch); err != nil {
			_ = stream.Close()
			err = errors.BadRequest("go.micro.srv.config.Watch", "send the Change error: %s", err)
			log.Error(err)
			return err
		}
	}
}
