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
// Original source: github.com/micro/go-micro/v3/codec/protorpc/netstring.go

package protorpc

import (
	"encoding/binary"
	"io"
)

// WriteNetString writes data to a big-endian netstring on a Writer.
// Size is always a 32-bit unsigned int.
func WriteNetString(w io.Writer, data []byte) (written int, err error) {
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(data)))
	if written, err = w.Write(size); err != nil {
		return
	}
	return w.Write(data)
}

// ReadNetString reads data from a big-endian netstring.
func ReadNetString(r io.Reader) (data []byte, err error) {
	sizeBuf := make([]byte, 4)
	_, err = r.Read(sizeBuf)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(sizeBuf)
	if size == 0 {
		return nil, nil
	}
	data = make([]byte, size)
	_, err = r.Read(data)
	if err != nil {
		return nil, err
	}
	return
}
