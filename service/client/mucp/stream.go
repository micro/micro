// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/client/mucp/mucp_stream.go

package mucp

import (
	"context"
	"io"
	"sync"

	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/service/client"
)

// Implements the streamer interface
type rpcStream struct {
	sync.RWMutex
	id       string
	closed   chan bool
	err      error
	request  client.Request
	response client.Response
	codec    codec.Codec
	context  context.Context

	// signal whether we should send EOS
	sendEOS bool

	// release releases the connection back to the pool
	release func(err error)
}

func (r *rpcStream) isClosed() bool {
	select {
	case <-r.closed:
		return true
	default:
		return false
	}
}

func (r *rpcStream) Context() context.Context {
	return r.context
}

func (r *rpcStream) Request() client.Request {
	return r.request
}

func (r *rpcStream) Response() client.Response {
	return r.response
}

func (r *rpcStream) Send(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	if r.isClosed() {
		r.err = errShutdown
		return errShutdown
	}

	req := codec.Message{
		Id:       r.id,
		Target:   r.request.Service(),
		Method:   r.request.Method(),
		Endpoint: r.request.Endpoint(),
		Type:     codec.Request,
	}

	if err := r.codec.Write(&req, msg); err != nil {
		r.err = err
		return err
	}

	return nil
}

func (r *rpcStream) Recv(msg interface{}) error {
	r.Lock()
	defer r.Unlock()

	if r.isClosed() {
		r.err = errShutdown
		return errShutdown
	}

	var resp codec.Message

	r.Unlock()
	err := r.codec.ReadHeader(&resp, codec.Response)
	r.Lock()
	if err != nil {
		if err == io.EOF && !r.isClosed() {
			r.err = io.ErrUnexpectedEOF
			return io.ErrUnexpectedEOF
		}
		r.err = err
		return err
	}

	switch {
	case len(resp.Error) > 0:
		// We've got an error response. Give this to the request;
		// any subsequent requests will get the ReadResponseBody
		// error if there is one.
		if resp.Error != lastStreamResponseError {
			r.err = serverError(resp.Error)
		} else {
			r.err = io.EOF
		}
		r.Unlock()
		err = r.codec.ReadBody(nil)
		r.Lock()
		if err != nil {
			r.err = err
		}
	default:
		r.Unlock()
		err = r.codec.ReadBody(msg)
		r.Lock()
		if err != nil {
			r.err = err
		}
	}

	return r.err
}

func (r *rpcStream) Error() error {
	r.RLock()
	defer r.RUnlock()
	return r.err
}

func (r *rpcStream) Close() error {
	r.Lock()

	select {
	case <-r.closed:
		r.Unlock()
		return nil
	default:
		close(r.closed)
		r.Unlock()

		// send the end of stream message
		if r.sendEOS {
			// no need to check for error
			r.codec.Write(&codec.Message{
				Id:       r.id,
				Target:   r.request.Service(),
				Method:   r.request.Method(),
				Endpoint: r.request.Endpoint(),
				Type:     codec.Error,
				Error:    lastStreamResponseError,
			}, nil)
		}

		err := r.codec.Close()

		// release the connection
		r.release(r.Error())

		// return the codec error
		return err
	}
}
