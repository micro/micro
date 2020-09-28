package handler

import (
	"context"
	"fmt"
	"time"

	log "github.com/micro/go-micro/v3/logger"
	store "github.com/micro/micro/v3/service/store"

	pb "example/proto"
)

type Example struct{}

// TestExpiry tests writing a record with expiry
func (e *Example) TestExpiry(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if err := writeWithExpiry("WriteExpiry", "bar", 5*time.Second); err != nil {
		return err
	}

	recs, err := store.Read("WriteExpiry")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteExpiry with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record WriteExpiry, expected 1 record. Received %d", len(recs))
	}

	if err := writeWithExpiry("Record.Expiry", "bar", 5*time.Second); err != nil {
		return err
	}

	recs, err = store.Read("Record.Expiry")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record Record.Expiry with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record Record.Expiry, expected 1 record. Received %d", len(recs))
	}

	if err := writeWithExpiry("WriteTTL", "bar", 5*time.Second); err != nil {
		return err
	}

	recs, err = store.Read("WriteTTL")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record WriteTTL with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record WriteTTL, expected 1 record. Received %d", len(recs))
	}

	time.Sleep(6 * time.Second)
	recs, err = store.Read("Record.Expiry")
	if err != store.ErrNotFound {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record Record.Expiry. Expected not found. Received %s and %d records", err, len(recs))
	}

	rsp.Msg = "Success"
	return nil
}

func writeWithExpiry(key, val string, duration time.Duration) error {
	if err := store.Write(&store.Record{Key: key, Value: []byte(val), Expiry: duration}); err != nil {
		log.Errorf("Error writing %s", err)
		return fmt.Errorf("Error writing record %s with expiry %s", key, err)
	}
	return nil
}

func (e *Example) TestList(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// Test Limit()
	for i := 0; i < 3; i++ {
		if err := writeWithExpiry(fmt.Sprintf("TestList%d", i), "bar", 5*time.Second); err != nil {
			return err
		}
	}

	recs, err := store.List(store.Prefix("TestList"))
	if err != nil {
		return fmt.Errorf("Error listing from store %s", err)
	}
	log.Infof("Recs %+v", recs)
	if len(recs) != 3 {
		return fmt.Errorf("Error listing records, expected 3, received %d", len(recs))
	}
	rsp.Msg = "Success"
	return nil
}

func (e *Example) TestListLimit(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	for i := 0; i < 10; i++ {
		if err := writeWithExpiry(fmt.Sprintf("TestLimit%d", i), "bar", 5*time.Second); err != nil {
			return err
		}
	}

	recs, err := store.List(store.Prefix("TestLimit"), store.Limit(2))
	if err != nil {
		return fmt.Errorf("Error listing from store %s", err)
	}
	log.Infof("Recs limit %+v", recs)
	if len(recs) != 2 {
		return fmt.Errorf("Error listing records with limit, expected 2, received %d", len(recs))
	}
	rsp.Msg = "Success"
	return nil
}

func (e *Example) TestListOffset(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	for i := 0; i < 20; i++ {
		if err := writeWithExpiry(fmt.Sprintf("TestOffset%d", i), "bar", 5*time.Second); err != nil {
			return err
		}
	}

	recs, err := store.List(store.Prefix("TestOffset"), store.Offset(5))
	if err != nil {
		return fmt.Errorf("Error listing from store %s", err)
	}
	log.Infof("Recs offset %+v", recs)
	if len(recs) != 15 {
		return fmt.Errorf("Error listing records with offset, expected 15, received %d", len(recs))
	}

	rsp.Msg = "Success"
	return nil

}
