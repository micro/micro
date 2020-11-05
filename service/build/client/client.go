package client

import (
	"bytes"
	"io"

	pb "github.com/micro/micro/v3/proto/build"
	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
)

const bufferSize = 1024

// NewBuilder returns an initialized builder
func NewBuilder() build.Builder {
	return &builder{}
}

type builder struct {
	srv pb.BuildService
}

func (b *builder) Build(src io.Reader, opts ...build.Option) (io.Reader, error) {
	// parse the options
	var options build.Options
	for _, o := range opts {
		o(&options)
	}

	// start the stream
	stream, err := b.client().Build(context.WithNamespace("micro"), client.WithAuthToken())
	if err != nil {
		return nil, err
	}

	// setup the methods of communicating between the goroutines
	errChan := make(chan error, 1)
	doneChan := make(chan bool, 1)

	// send the source over the stream async
	go func() {
		var sentOptions bool
		for {
			buffer := make([]byte, bufferSize)
			for {
				num, err := src.Read(buffer)
				if err == io.EOF {
					stream.Close()
					return
				} else if err != nil {
					errChan <- err
					return
				}

				req := &pb.BuildRequest{
					Data: buffer[:num],
				}

				// pass the options on the first message only.
				if !sentOptions {
					req.Options = &pb.Options{Archive: options.Archive, Entrypoint: options.Entrypoint}
					sentOptions = true
				}

				// send the message over the stream
				if err := stream.Send(req); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// wait for the response from the server async
	result := bytes.NewBuffer(nil)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneChan <- true
				return
			} else if err != nil {
				errChan <- err
				return
			}

			if _, err := result.Write(res.Data); err != nil {
				errChan <- err
			}
		}
	}()

	select {
	case err := <-errChan:
		return nil, err
	case <-doneChan:
		return result, nil
	}
}

func (b *builder) client() pb.BuildService {
	if b.srv == nil {
		b.srv = pb.NewBuildService("build", client.DefaultClient)
	}
	return b.srv
}
