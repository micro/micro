package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/store"
	pb "github.com/micro/micro/v3/proto/store"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
)

const bufferSize = 100

// NewBlobStore returns a new store service implementation
func NewBlobStore() store.BlobStore {
	return &blob{
		client: pb.NewBlobStoreService("store", client.DefaultClient),
	}
}

type blob struct {
	client pb.BlobStoreService
}

func (b *blob) Read(key string, opts ...store.BlobOption) (io.Reader, error) {
	// validate the key
	if len(key) == 0 {
		return nil, store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}

	// execute the rpc
	stream, err := b.client.Read(context.TODO(), &pb.BlobReadRequest{
		Key: key,
		Options: &pb.BlobOptions{
			Namespace: options.Namespace,
		},
	}, goclient.WithAuthToken())

	// handle the error
	if verr := errors.Parse(err); verr != nil && verr.Code == http.StatusNotFound {
		return nil, store.ErrNotFound
	} else if verr != nil {
		return nil, verr
	} else if err != nil {
		return nil, err
	}

	// create a buffer to store the bytes in
	buf := bytes.NewBuffer(nil)

	// keep recieving bytes from the stream until it's closed by the server
	fmt.Println("Waiting to rec")
	for {
		res, err := stream.Recv()
		fmt.Println("Rec", string(res.Blob), err)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		buf.Write(res.Blob)
	}

	// return the bytes
	return buf, nil
}

func (b *blob) Write(key string, blob io.Reader, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// open the stream
	stream, err := b.client.Write(context.TODO(), goclient.WithAuthToken())
	if verr := errors.Parse(err); verr != nil {
		return verr
	} else if err != nil {
		return err
	}
	defer stream.Close()

	// read from the blob and stream it to the server
	buffer := make([]byte, bufferSize)
	for {
		if _, err := blob.Read(buffer); err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		err := stream.Send(&pb.BlobWriteRequest{
			Key: key,
			Options: &pb.BlobOptions{
				Namespace: options.Namespace,
			},
			Blob: buffer,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *blob) Delete(key string, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}

	// execute the rpc
	_, err := b.client.Delete(context.TODO(), &pb.BlobDeleteRequest{
		Key: key,
		Options: &pb.BlobOptions{
			Namespace: options.Namespace,
		},
	}, goclient.WithAuthToken())

	// handle the error
	if verr := errors.Parse(err); verr != nil && verr.Code == http.StatusNotFound {
		return store.ErrNotFound
	} else if verr != nil {
		return verr
	}

	return err
}
