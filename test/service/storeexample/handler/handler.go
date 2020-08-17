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

// Call is a single request handler called via client.Call or the generated client code
func (e *Example) TestExpiry(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	if err := mstore.Write(&store.Record{Key: "foo", Value: []byte("bar")},
		store.WriteExpiry(time.Now().Add(5*time.Second))); err != nil {
		log.Errorf("Error writing %s", err)
		return fmt.Errorf("Error writing record with expiry %s", err)
	}

	recs, err := mstore.Read("foo")
	if err != nil {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record with expiry %s", err)
	}
	if len(recs) != 1 {
		return fmt.Errorf("Error reading record, expected 1 record. Received %d", len(recs))
	}

	time.Sleep(6 * time.Second)
	recs, err = mstore.Read("foo")
	if err != store.ErrNotFound {
		log.Errorf("Error reading %s", err)
		return fmt.Errorf("Error reading record. Expected not found. Received %s and %d records", err, len(recs))
	}

	rsp.Msg = "Success"
	return nil
}
