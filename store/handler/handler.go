package handler

import (
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/go-micro/v2/store/service/proto"
)

type Store struct {
	// The default store
	Default store.Store

	// The internal store where databases and table information is kept
	Internal store.Store

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

	database, _ := md.Get("Micro-Database")
	table, _ := md.Get("Micro-Table")

	if len(database) == 0 {
		database = s.Default.Options().Database
	}

	if len(table) == 0 {
		table = s.Default.Options().Table
	}

	// just use the default if nothing is specified
	if len(database) == 0 && len(table) == 0 {
		return s.Default, nil
	}

	// attempt to get the database
	str, ok := s.Stores[database+":"+table]
	// got it
	if ok {
		return str, nil
	}

	// create a new store
	// either database is not blank or table is not blank
	st, err := s.New(database, table)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.store", "failed to setup store: %s", err.Error())
	}

	// save store
	s.Stores[database+":"+table] = st

	return st, nil
}

func (s *Store) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}

	var opts []store.ReadOption
	if req.Options != nil && req.Options.Prefix {
		opts = append(opts, store.ReadPrefix())
	}

	vals, err := st.Read(req.Key, opts...)
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("go.micro.store", err.Error())
	} else if err != nil {
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

	if req.Record == nil {
		return errors.BadRequest("go.micro.store", "no record specified")
	}

	record := &store.Record{
		Key:    req.Record.Key,
		Value:  req.Record.Value,
		Expiry: time.Duration(req.Record.Expiry) * time.Second,
	}

	err = st.Write(record)
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("go.micro.store", err.Error())
	} else if err != nil {
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
	if err := st.Delete(req.Key); err == store.ErrNotFound {
		return errors.NotFound("go.micro.store", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	return nil
}

func (s *Store) Databases(ctx context.Context, req *pb.DatabasesRequest, rsp *pb.DatabasesResponse) error {
	recs, err := s.Internal.Read("databases/", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	rsp.Databases = make([]string, len(recs))
	for i, r := range recs {
		rsp.Databases[i] = strings.TrimPrefix(r.Key, "databases/")
	}
	return nil
}

func (s *Store) Tables(ctx context.Context, req *pb.TablesRequest, rsp *pb.TablesResponse) error {
	recs, err := s.Internal.Read("tables/"+req.Database+"/", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}
	rsp.Tables = make([]string, len(recs))
	for i, r := range recs {
		rsp.Tables[i] = strings.TrimPrefix(r.Key, "tables/"+"req.Database"+"/")
	}
	return nil
}

func (s *Store) List(ctx context.Context, req *pb.ListRequest, stream pb.Store_ListStream) error {
	// get new store
	st, err := s.get(ctx)
	if err != nil {
		return err
	}

	vals, err := st.List()
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("go.micro.store", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.store", err.Error())
	}

	rsp := new(pb.ListResponse)

	// TODO: batch sync
	for _, val := range vals {
		rsp.Keys = append(rsp.Keys, val)
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
