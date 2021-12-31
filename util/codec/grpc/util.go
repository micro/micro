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
// Original source: github.com/micro/go-micro/v3/codec/grpc/util.go

package grpc

import (
	"encoding/binary"
	"fmt"
	"io"
)

var (
	MaxMessageSize = 1024 * 1024 * 4 // 4Mb
	maxInt         = int(^uint(0) >> 1)
)

func decode(r io.Reader) (uint8, []byte, error) {
	header := make([]byte, 5)

	// read the header
	if _, err := r.Read(header[:]); err != nil {
		return uint8(0), nil, err
	}

	// get encoding format e.g compressed
	cf := uint8(header[0])

	// get message length
	length := binary.BigEndian.Uint32(header[1:])

	// no encoding format
	if length == 0 {
		return cf, nil, nil
	}

	//
	if int64(length) > int64(maxInt) {
		return cf, nil, fmt.Errorf("grpc: received message larger than max length allowed on current machine (%d vs. %d)", length, maxInt)
	}
	if int(length) > MaxMessageSize {
		return cf, nil, fmt.Errorf("grpc: received message larger than max (%d vs. %d)", length, MaxMessageSize)
	}

	msg := make([]byte, int(length))

	if _, err := r.Read(msg); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return cf, nil, err
	}

	return cf, msg, nil
}

func encode(cf uint8, buf []byte, w io.Writer) error {
	header := make([]byte, 5)

	// set compression
	header[0] = byte(cf)

	// write length as header
	binary.BigEndian.PutUint32(header[1:], uint32(len(buf)))

	// read the header
	if _, err := w.Write(header[:]); err != nil {
		return err
	}

	// write the buffer
	_, err := w.Write(buf)
	return err
}
