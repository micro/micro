package server

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	gostore "github.com/micro/go-micro/v3/store"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
)

const (
	sourcePrefix        = "source://"
	blobNamespacePrefix = "micro/runtime"
)

// Source implements the proto source service interface
type Source struct{}

// Upload source to the server
func (s *Source) Upload(ctx context.Context, stream pb.Source_UploadStream) error {
	// authorize the request
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("runtime.Source.Upload", "An account is required to upload source")
	}
	namespace := acc.Issuer

	// recieve the source from the client
	buf := bytes.NewBuffer(nil)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.InternalServerError("runtime.Source.Upload", err.Error())
		}

		// write the bytes to the buffer
		if _, err := buf.Write(req.Data); err != nil {
			return err
		}
	}

	// ensure the blob was sent over the stream
	if buf == nil {
		return errors.BadRequest("runtime.Source.Upload", "No blob was sent")
	}

	// write the source to the store
	key := "source://" + uuid.New().String()
	opt := gostore.BlobNamespace(namespace)
	if err := store.DefaultBlobStore.Write(key, buf, opt); err != nil {
		return fmt.Errorf("Error writing source to blob store: %v", err)
	}

	// todo: implement cleanup logic. write to the normal store and have a cleanup func loop through
	// the records every X minutes. Note: we only want to do this if the builder is set.

	// close the stream
	return stream.SendAndClose(&pb.UploadResponse{Id: key})
}
