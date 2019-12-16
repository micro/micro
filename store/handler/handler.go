package handler

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/store"
	pb "github.com/micro/go-micro/store/service/proto"
)

type Store struct {
	// The default store
	Default store.Store

	// Store initialiser
	New func(string, string) (store.Store, error)

	// Store map
	sync.RWMutex
	Stores map[string]store.Store
}

func (s *Store) get(ctx context.Context) (store.Store, error) {
	// lock (might be a race)
	s.Lock()
	defer s.Unlock()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return s.Default, nil
	}

	namespace := md["Micro-Namespace"]
	prefix := md["Micro-Prefix"]

	if len(namespace) == 0 && len(prefix) == 0 {
		return s.Default, nil
	}

	str, ok := s.Stores[namespace+":"+prefix]
	// got it
	if ok {
		return str, nil
	}

	// create a new store

	// either namespace is not blank or prefix is not blank
	st, err := s.New(namespace, prefix)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.store", err.Error())
	}

	// save store
	s.Stores[namespace+":"+prefix] = st

	return st, nil
}

func (s *Store) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}

	vals, err := st.Read(req.Keys...)
	if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	for _, val := range vals {
		rsp.Records = append(rsp.Records, &pb.Record{
			Key:    val.Key,
			Value:  val.Value,
			Expiry: int64(val.Expiry.Seconds()),
		})
	}
	return nil
}

func (s *Store) Write(ctx context.Context, req *pb.WriteRequest, rsp *pb.WriteResponse) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}

	records := make([]*store.Record, 0, len(req.Records))

	for _, record := range req.Records {
		records = append(records, &store.Record{
			Key:    record.Key,
			Value:  record.Value,
			Expiry: time.Duration(record.Expiry) * time.Second,
		})
	}

	if err := st.Write(records...); err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}
	if err := st.Delete(req.Keys...); err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	return nil
}

func (s *Store) List(ctx context.Context, req *pb.ListRequest, stream pb.Store_ListStream) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}

	var vals []*store.Record

	if len(req.Key) > 0 {
		vals, err = st.Read(req.Key)
	} else {
		vals, err = st.List()
	}
	if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	rsp := new(pb.ListResponse)

	// TODO: batch sync
	for _, val := range vals {
		rsp.Records = append(rsp.Records, &pb.Record{
			Key:    val.Key,
			Value:  val.Value,
			Expiry: int64(val.Expiry.Seconds()),
		})
	}

	err = stream.Send(rsp)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	return nil
}
