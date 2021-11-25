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
// Original source: github.com/micro/go-micro/v3/util/buf/buf.go

package buf

import (
	"bytes"
)

type buffer struct {
	*bytes.Buffer
}

func (b *buffer) Close() error {
	b.Buffer.Reset()
	return nil
}

func New(b *bytes.Buffer) *buffer {
	if b == nil {
		b = bytes.NewBuffer(nil)
	}
	return &buffer{b}
}
