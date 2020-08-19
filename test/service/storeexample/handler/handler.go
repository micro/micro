package handler

import (
	"context"
	"fmt"
	"time"

	log "github.com/micro/go-micro/v3/logger"

	"github.com/micro/go-micro/v3/store"

	mstore "github.com/micro/micro/v3/service/store"

	pb "example/proto"
)

type Example struct{}

// TestExpiry tests writing a record with expiry / TTL in all the ways that we support
// - Record.Expiry
// - WriteExpiry
// - WriteTTL
func (e *Example) TestExpiry(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if err := mstore.Write(&store.Record{Key: "WriteExpiry", Value: []byte("bar")},
		store.WriteExpiry(time.Now().Add(5*time.Second))); err != nil {
		log.Errorf("Error writing %s", err)
		return fmt.Errorf("Error writing record WriteExpiry with expiry %s", err)
	}

	recs, err := mstore.Read("WriteExpiry")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteExpiry with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record WriteExpiry, expected 1 record. Received %d", len(recs))
	}

	if err := mstore.Write(&store.Record{Key: "Record.Expiry", Value: []byte("bar"), Expiry: 5 * time.Second}); err != nil {
		log.Errorf("Error writing %s", err)
		return fmt.Errorf("Error writing record Record.Expiry with expiry %s", err)
	}

	recs, err = mstore.Read("Record.Expiry")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record Record.Expiry with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record Record.Expiry, expected 1 record. Received %d", len(recs))
	}

	if err := mstore.Write(&store.Record{Key: "WriteTTL", Value: []byte("bar")}, store.WriteTTL(5*time.Second)); err != nil {
		log.Errorf("Error writing %s", err)
		return fmt.Errorf("Error writing record WriteTTL with expiry %s", err)
	}

	recs, err = mstore.Read("WriteTTL")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteTTL with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record WriteTTL, expected 1 record. Received %d", len(recs))
	}

	time.Sleep(6 * time.Second)
	recs, err = mstore.Read("WriteExpiry")
	if err != store.ErrNotFound {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteExpiry. Expected not found. Received %s and %d records", err, len(recs))
	}
	recs, err = mstore.Read("Record.Expiry")
	if err != store.ErrNotFound {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record Record.Expiry. Expected not found. Received %s and %d records", err, len(recs))
	}
	recs, err = mstore.Read("WriteTTL")
	if err != store.ErrNotFound {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteTTL. Expected not found. Received %s and %d records", err, len(recs))
	}

	rsp.Msg = "Success"
	return nil
}

func (e *Example) TestList(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	mstore.List()
}
