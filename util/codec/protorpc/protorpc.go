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
// Original source: github.com/micro/go-micro/v3/codec/protorpc/protorpc.go

// Protorpc provides a net/rpc proto-rpc codec. See envelope.proto for the format.
package protorpc

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/micro/micro/v5/util/codec"
)

type flusher interface {
	Flush() error
}

type protoCodec struct {
	sync.Mutex
	rwc io.ReadWriteCloser
	mt  codec.MessageType
	buf *bytes.Buffer
}

func (c *protoCodec) Close() error {
	c.buf.Reset()
	return c.rwc.Close()
}

func (c *protoCodec) String() string {
	return "proto-rpc"
}

func id(id string) uint64 {
	p, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		p = 0
	}
	i := uint64(p)
	return i
}

func (c *protoCodec) Write(m *codec.Message, b interface{}) error {
	switch m.Type {
	case codec.Request:
		c.Lock()
		defer c.Unlock()
		// This is protobuf, of course we copy it.
		pbr := &Request{ServiceMethod: m.Method, Seq: id(m.Id)}
		data, err := proto.Marshal(pbr)
		if err != nil {
			return err
		}
		_, err = WriteNetString(c.rwc, data)
		if err != nil {
			return err
		}
		// dont trust or incoming message
		m, ok := b.(proto.Message)
		if !ok {
			return codec.ErrInvalidMessage
		}
		data, err = proto.Marshal(m)
		if err != nil {
			return err
		}
		_, err = WriteNetString(c.rwc, data)
		if err != nil {
			return err
		}
		if flusher, ok := c.rwc.(flusher); ok {
			if err = flusher.Flush(); err != nil {
				return err
			}
		}
	case codec.Response, codec.Error:
		c.Lock()
		defer c.Unlock()
		rtmp := &Response{ServiceMethod: m.Method, Seq: id(m.Id), Error: m.Error}
		data, err := proto.Marshal(rtmp)
		if err != nil {
			return err
		}
		_, err = WriteNetString(c.rwc, data)
		if err != nil {
			return err
		}
		if pb, ok := b.(proto.Message); ok {
			data, err = proto.Marshal(pb)
			if err != nil {
				return err
			}
		} else {
			data = nil
		}
		_, err = WriteNetString(c.rwc, data)
		if err != nil {
			return err
		}
		if flusher, ok := c.rwc.(flusher); ok {
			if err = flusher.Flush(); err != nil {
				return err
			}
		}
	case codec.Event:
		m, ok := b.(proto.Message)
		if !ok {
			return codec.ErrInvalidMessage
		}
		data, err := proto.Marshal(m)
		if err != nil {
			return err
		}
		c.rwc.Write(data)
	default:
		return fmt.Errorf("Unrecognised message type: %v", m.Type)
	}
	return nil
}

func (c *protoCodec) ReadHeader(m *codec.Message, mt codec.MessageType) error {
	c.buf.Reset()
	c.mt = mt

	switch mt {
	case codec.Request:
		data, err := ReadNetString(c.rwc)
		if err != nil {
			return err
		}
		rtmp := new(Request)
		err = proto.Unmarshal(data, rtmp)
		if err != nil {
			return err
		}
		m.Method = rtmp.GetServiceMethod()
		m.Id = fmt.Sprintf("%d", rtmp.GetSeq())
	case codec.Response:
		data, err := ReadNetString(c.rwc)
		if err != nil {
			return err
		}
		rtmp := new(Response)
		err = proto.Unmarshal(data, rtmp)
		if err != nil {
			return err
		}
		m.Method = rtmp.GetServiceMethod()
		m.Id = fmt.Sprintf("%d", rtmp.GetSeq())
		m.Error = rtmp.GetError()
	case codec.Event:
		_, err := io.Copy(c.buf, c.rwc)
		return err
	default:
		return fmt.Errorf("Unrecognised message type: %v", mt)
	}
	return nil
}

func (c *protoCodec) ReadBody(b interface{}) error {
	var data []byte
	switch c.mt {
	case codec.Request, codec.Response:
		var err error
		data, err = ReadNetString(c.rwc)
		if err != nil {
			return err
		}
	case codec.Event:
		data = c.buf.Bytes()
	default:
		return fmt.Errorf("Unrecognised message type: %v", c.mt)
	}
	if b != nil {
		return proto.Unmarshal(data, b.(proto.Message))
	}
	return nil
}

func NewCodec(rwc io.ReadWriteCloser) codec.Codec {
	return &protoCodec{
		buf: bytes.NewBuffer(nil),
		rwc: rwc,
	}
}
