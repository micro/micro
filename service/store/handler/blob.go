package handler

import (
	"bytes"
	"context"
	"io"

	pb "github.com/micro/micro/v5/proto/store"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/store"
	authns "github.com/micro/micro/v5/util/auth/namespace"
	"github.com/micro/micro/v5/util/namespace"
)

const bufferSize = 1024

type BlobStore struct{}

func (b *BlobStore) Read(ctx context.Context, req *pb.BlobReadRequest, stream pb.BlobStore_ReadStream) error {
	// parse the options
	if ns := req.GetOptions().GetNamespace(); len(ns) == 0 {
		req.Options = &pb.BlobOptions{
			Namespace: namespace.FromContext(ctx),
		}
	}

	// authorize the request
	if err := authns.AuthorizeAdmin(ctx, req.Options.Namespace, "store.Blob.Read"); err != nil {
		return err
	}

	// execute the request
	blob, err := store.DefaultBlobStore.Read(req.Key, store.BlobNamespace(req.Options.Namespace))
	if err == store.ErrNotFound {
		return errors.NotFound("store.Blob.Read", "Blob not found")
	} else if err == store.ErrMissingKey {
		return errors.BadRequest("store.Blob.Read", "Missing key")
	} else if err != nil {
		return errors.InternalServerError("store.Blob.Read", err.Error())
	}

	// read from the blob and stream it to the client
	buffer := make([]byte, bufferSize)
	for {
		num, err := blob.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if err := stream.Send(&pb.BlobReadResponse{Blob: buffer[:num]}); err != nil {
			return err
		}
	}

	return stream.Close()
}

func (b *BlobStore) Write(ctx context.Context, stream pb.BlobStore_WriteStream) error {
	// the key and options are passed on each message but we only need to extract them once
	var buf *bytes.Buffer
	var key string
	var options *pb.BlobOptions

	// recieve the blob from the client
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.InternalServerError("store.Blob.Write", err.Error())
		}

		if buf == nil {
			// first message recieved from the stream
			buf = bytes.NewBuffer(req.Blob)
			key = req.Key
			options = req.Options

			// parse the options
			if options == nil || len(options.Namespace) == 0 {
				options.Namespace = namespace.FromContext(ctx)
			}

			// authorize the request. do this inside the loop so we fail fast
			if err := authns.AuthorizeAdmin(ctx, options.Namespace, "store.Blob.Write"); err != nil {
				return err
			}

		} else {
			// subsequent message recieved from the stream
			buf.Write(req.Blob)
		}
	}

	// ensure the blob was sent over the stream
	if buf == nil {
		return errors.BadRequest("store.Blob.Write", "No blob was sent")
	}

	// execute the request
	err := store.DefaultBlobStore.Write(key, buf, store.BlobNamespace(options.Namespace), store.BlobPublic(options.Public), store.BlobContentType(options.ContentType))
	if err == store.ErrMissingKey {
		return errors.BadRequest("store.Blob.Write", "Missing key")
	} else if err != nil {
		return errors.InternalServerError("store.Blob.Write", err.Error())
	}

	// close the stream
	return stream.SendAndClose(&pb.BlobWriteResponse{})
}

func (b *BlobStore) Delete(ctx context.Context, req *pb.BlobDeleteRequest, rsp *pb.BlobDeleteResponse) error {
	// parse the options
	if ns := req.GetOptions().GetNamespace(); len(ns) == 0 {
		req.Options = &pb.BlobOptions{
			Namespace: namespace.FromContext(ctx),
		}
	}

	// authorize the request
	if err := authns.AuthorizeAdmin(ctx, req.Options.Namespace, "store.Blob.Delete"); err != nil {
		return err
	}

	// execute the request
	err := store.DefaultBlobStore.Delete(req.Key, store.BlobNamespace(req.Options.Namespace))
	if err == store.ErrNotFound {
		return errors.NotFound("store.Blob.Delete", "Blob not found")
	} else if err == store.ErrMissingKey {
		return errors.BadRequest("store.Blob.Delete", "Missing key")
	} else if err != nil {
		return errors.InternalServerError("store.Blob.Delete", err.Error())
	}

	return nil
}

func (b *BlobStore) List(ctx context.Context, req *pb.BlobListRequest, rsp *pb.BlobListResponse) error {
	// parse the options
	if ns := req.GetOptions().GetNamespace(); len(ns) == 0 {
		req.Options = &pb.BlobListOptions{
			Namespace: namespace.FromContext(ctx),
			Prefix:    req.GetOptions().GetPrefix(),
		}
	}

	// authorize the request
	if err := authns.AuthorizeAdmin(ctx, req.Options.Namespace, "store.Blob.List"); err != nil {
		return err
	}

	// execute the request
	keys, err := store.DefaultBlobStore.List(
		store.BlobListNamespace(req.GetOptions().GetNamespace()),
		store.BlobListPrefix(req.GetOptions().GetPrefix()))
	if err != nil {
		return errors.InternalServerError("store.Blob.List", err.Error())
	}
	rsp.Keys = keys

	return nil

}
