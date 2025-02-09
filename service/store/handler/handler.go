package handler

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"

	pb "github.com/micro/micro/v5/proto/store"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/store"
	"github.com/micro/micro/v5/util/auth/namespace"
)

const (
	defaultDatabase = namespace.DefaultNamespace
	defaultTable    = namespace.DefaultNamespace
	internalTable   = "store"
)

type Store struct {
	// local Stores cache
	sync.RWMutex
	Stores map[string]bool
}

// List all the keys in a table
func (h *Store) List(ctx context.Context, req *pb.ListRequest, stream pb.Store_ListStream) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.ListOptions{}
	}
	if len(req.Options.Database) == 0 {
		req.Options.Database = defaultDatabase
	}
	if len(req.Options.Table) == 0 {
		req.Options.Table = defaultTable
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Database, "store.Store.List"); err != nil {
		return err
	}

	// setup the store
	if err := h.setupTable(req.Options.Database, req.Options.Table); err != nil {
		return errors.InternalServerError("store.Store.List", err.Error())
	}

	// setup the options
	opts := []store.ListOption{
		store.ListFrom(req.Options.Database, req.Options.Table),
	}
	if len(req.Options.Prefix) > 0 {
		opts = append(opts, store.ListPrefix(req.Options.Prefix))
	}
	if len(req.Options.Suffix) > 0 {
		opts = append(opts, store.ListSuffix(req.Options.Suffix))
	}
	if req.Options.Offset > 0 {
		opts = append(opts, store.ListOffset(uint(req.Options.Offset)))
	}
	if req.Options.Limit > 0 {
		opts = append(opts, store.ListLimit(uint(req.Options.Limit)))
	}
	if len(req.Options.Order) > 0 {
		order := store.OrderAsc
		if req.Options.Order == string(store.OrderDesc) {
			order = store.OrderDesc
		}
		opts = append(opts, store.ListOrder(order))
	}

	// list from the store
	vals, err := store.DefaultStore.List(opts...)
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("store.Store.List", err.Error())
	} else if err != nil {
		return errors.InternalServerError("store.Store.List", err.Error())
	}

	// serialize the response
	// TODO: batch sync
	rsp := new(pb.ListResponse)
	for _, val := range vals {
		rsp.Keys = append(rsp.Keys, val)
	}

	err = stream.Send(rsp)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return errors.InternalServerError("store.Store.List", err.Error())
	}
	return nil
}

// Read records from the store
func (h *Store) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.ReadOptions{}
	}
	if len(req.Options.Database) == 0 {
		req.Options.Database = defaultDatabase
	}
	if len(req.Options.Table) == 0 {
		req.Options.Table = defaultTable
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Database, "store.Store.Read"); err != nil {
		return err
	}

	// setup the store
	if err := h.setupTable(req.Options.Database, req.Options.Table); err != nil {
		return errors.InternalServerError("store.Store.Read", err.Error())
	}

	// setup the options
	opts := []store.ReadOption{
		store.ReadFrom(req.Options.Database, req.Options.Table),
	}
	if req.Options.Prefix {
		opts = append(opts, store.ReadPrefix())
	}
	if req.Options.Suffix {
		opts = append(opts, store.ReadSuffix())
	}
	if req.Options.Limit > 0 {
		opts = append(opts, store.ReadLimit(uint(req.Options.Limit)))
	}
	if req.Options.Offset > 0 {
		opts = append(opts, store.ReadOffset(uint(req.Options.Offset)))
	}
	if len(req.Options.Order) > 0 {
		order := store.OrderAsc
		if req.Options.Order == string(store.OrderDesc) {
			order = store.OrderDesc
		}
		opts = append(opts, store.ReadOrder(order))
	}

	// read from the database
	vals, err := store.DefaultStore.Read(req.Key, opts...)
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("store.Store.Read", err.Error())
	} else if err != nil {
		return errors.InternalServerError("store.Store.Read", err.Error())
	}

	// serialize the result
	for _, val := range vals {
		metadata := make(map[string]*pb.Field)
		for k, v := range val.Metadata {
			metadata[k] = &pb.Field{
				Type:  reflect.TypeOf(v).String(),
				Value: fmt.Sprintf("%v", v),
			}
		}
		rsp.Records = append(rsp.Records, &pb.Record{
			Key:      val.Key,
			Value:    val.Value,
			Expiry:   int64(val.Expiry.Seconds()),
			Metadata: metadata,
		})
	}
	return nil
}

// Write to the store
func (h *Store) Write(ctx context.Context, req *pb.WriteRequest, rsp *pb.WriteResponse) error {
	// validate the request
	if req.Record == nil {
		return errors.BadRequest("store.Store.Write", "no record specified")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.WriteOptions{}
	}
	if len(req.Options.Database) == 0 {
		req.Options.Database = defaultDatabase
	}
	if len(req.Options.Table) == 0 {
		req.Options.Table = defaultTable
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Database, "store.Store.Write"); err != nil {
		return err
	}

	// setup the store
	if err := h.setupTable(req.Options.Database, req.Options.Table); err != nil {
		return errors.InternalServerError("store.Store.Write", err.Error())
	}

	// setup the options
	opts := []store.WriteOption{
		store.WriteTo(req.Options.Database, req.Options.Table),
	}

	// construct the record
	metadata := make(map[string]interface{})
	for k, v := range req.Record.Metadata {
		metadata[k] = v.Value
	}
	record := &store.Record{
		Key:      req.Record.Key,
		Value:    req.Record.Value,
		Expiry:   time.Duration(req.Record.Expiry) * time.Second,
		Metadata: metadata,
	}

	// write to the store
	err := store.DefaultStore.Write(record, opts...)
	if err != nil && err == store.ErrNotFound {
		return errors.NotFound("store.Store.Write", err.Error())
	} else if err != nil {
		return errors.InternalServerError("store.Store.Write", err.Error())
	}

	return nil
}

func (h *Store) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.DeleteOptions{}
	}
	if len(req.Options.Database) == 0 {
		req.Options.Database = defaultDatabase
	}
	if len(req.Options.Table) == 0 {
		req.Options.Table = defaultTable
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Options.Database, "store.Store.Delete"); err != nil {
		return err
	}

	// setup the store
	if err := h.setupTable(req.Options.Database, req.Options.Table); err != nil {
		return errors.InternalServerError("store.Store.Delete", err.Error())
	}

	// setup the options
	opts := []store.DeleteOption{
		store.DeleteFrom(req.Options.Database, req.Options.Table),
	}

	// delete from the store
	if err := store.DefaultStore.Delete(req.Key, opts...); err == store.ErrNotFound {
		return errors.NotFound("store.Store.Delete", err.Error())
	} else if err != nil {
		return errors.InternalServerError("store.Store.Delete", err.Error())
	}

	return nil
}

// Databases lists all the databases
func (h *Store) Databases(ctx context.Context, req *pb.DatabasesRequest, rsp *pb.DatabasesResponse) error {
	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, defaultDatabase, "store.Store.Database"); err != nil {
		return err
	}

	// read the databases from the store
	opts := []store.ReadOption{
		store.ReadPrefix(),
		store.ReadFrom(defaultDatabase, internalTable),
	}
	recs, err := store.DefaultStore.Read("databases/", opts...)
	if err != nil {
		return errors.InternalServerError("store.Store.Databases", err.Error())
	}

	// serialize the response
	rsp.Databases = make([]string, len(recs))
	for i, r := range recs {
		rsp.Databases[i] = strings.TrimPrefix(r.Key, "databases/")
	}
	return nil
}

// Tables returns all the tables in a database
func (h *Store) Tables(ctx context.Context, req *pb.TablesRequest, rsp *pb.TablesResponse) error {
	// set defaults
	if len(req.Database) == 0 {
		req.Database = defaultDatabase
	}

	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, req.Database, "store.Store.Tables"); err != nil {
		return err
	}

	// construct the options
	opts := []store.ReadOption{
		store.ReadPrefix(),
		store.ReadFrom(defaultDatabase, internalTable),
	}

	// perform the query
	query := fmt.Sprintf("tables/%v/", req.Database)
	recs, err := store.DefaultStore.Read(query, opts...)
	if err != nil {
		return errors.InternalServerError("store.Store.Tables", err.Error())
	}

	// serialize the response
	rsp.Tables = make([]string, len(recs))
	for i, r := range recs {
		rsp.Tables[i] = strings.TrimPrefix(r.Key, "tables/"+req.Database+"/")
	}
	return nil
}

func (h *Store) setupTable(database, table string) error {
	// lock (might be a race)
	h.Lock()
	defer h.Unlock()

	// attempt to get the database
	if _, ok := h.Stores[database+":"+table]; ok {
		return nil
	}

	// record the new database in the internal store
	opt := store.WriteTo(defaultDatabase, internalTable)
	dbRecord := &store.Record{Key: "databases/" + database, Value: []byte{}}
	if err := store.DefaultStore.Write(dbRecord, opt); err != nil {
		return fmt.Errorf("Error writing new database to internal table: %v", err)
	}

	// record the new table in the internal store
	tableRecord := &store.Record{Key: "tables/" + database + "/" + table, Value: []byte{}}
	if err := store.DefaultStore.Write(tableRecord, opt); err != nil {
		return fmt.Errorf("Error writing new table to internal table: %v", err)
	}

	// write to the cache
	h.Stores[database+":"+table] = true
	return nil
}
