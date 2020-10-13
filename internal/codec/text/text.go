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
// Original source: github.com/micro/go-micro/v3/codec/text/text.go

// Package text reads any text/* content-type
package text

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/micro/micro/v3/internal/codec"
)

type Codec struct {
	Conn io.ReadWriteCloser
}

// Frame gives us the ability to define raw data to send over the pipes
type Frame struct {
	Data []byte
}

func (c *Codec) ReadHeader(m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *Codec) ReadBody(b interface{}) error {
	// read bytes
	buf, err := ioutil.ReadAll(c.Conn)
	if err != nil {
		return err
	}

	switch v := b.(type) {
	case *string:
		*v = string(buf)
	case *[]byte:
		*v = buf
	case *Frame:
		v.Data = buf
	default:
		return fmt.Errorf("failed to read body: %v is not type of *[]byte", b)
	}

	return nil
}

func (c *Codec) Write(m *codec.Message, b interface{}) error {
	var v []byte
	switch ve := b.(type) {
	case nil:
		return nil
	case *Frame:
		v = ve.Data
	case *[]byte:
		v = *ve
	case *string:
		v = []byte(*ve)
	case string:
		v = []byte(ve)
	case []byte:
		v = ve
	default:
		return fmt.Errorf("failed to write: %v is not type of *[]byte or []byte", b)
	}
	_, err := c.Conn.Write(v)
	return err
}

func (c *Codec) Close() error {
	return c.Conn.Close()
}

func (c *Codec) String() string {
	return "text"
}

func NewCodec(c io.ReadWriteCloser) codec.Codec {
	return &Codec{
		Conn: c,
	}
}
