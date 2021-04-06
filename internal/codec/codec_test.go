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
// Original source: github.com/micro/go-micro/v3/codec/codec_test.go

package codec_test

import (
	"io"
	"testing"

	"github.com/micro/micro/v3/internal/codec"
	"github.com/micro/micro/v3/internal/codec/bytes"
	"github.com/micro/micro/v3/internal/codec/grpc"
	"github.com/micro/micro/v3/internal/codec/json"
	"github.com/micro/micro/v3/internal/codec/jsonrpc"
	"github.com/micro/micro/v3/internal/codec/proto"
	"github.com/micro/micro/v3/internal/codec/protorpc"
	"github.com/micro/micro/v3/internal/codec/text"
)

type testRWC struct{}

func (rwc *testRWC) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (rwc *testRWC) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (rwc *testRWC) Close() error {
	return nil
}

func getCodecs(c io.ReadWriteCloser) map[string]codec.Codec {
	return map[string]codec.Codec{
		"bytes":    bytes.NewCodec(c),
		"grpc":     grpc.NewCodec(c),
		"json":     json.NewCodec(c),
		"jsonrpc":  jsonrpc.NewCodec(c),
		"proto":    proto.NewCodec(c),
		"protorpc": protorpc.NewCodec(c),
		"text":     text.NewCodec(c),
	}
}

func Test_WriteEmptyBody(t *testing.T) {
	for name, c := range getCodecs(&testRWC{}) {
		err := c.Write(&codec.Message{
			Type:   codec.Error,
			Header: map[string]string{},
		}, nil)
		if err != nil {
			t.Fatalf("codec %s - expected no error when writing empty/nil body: %s", name, err)
		}
	}
}
