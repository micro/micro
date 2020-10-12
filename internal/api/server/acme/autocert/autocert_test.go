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
// Original source: github.com/micro/go-micro/v3/api/server/acme/autocert/autocert_test.go

package autocert

import (
	"testing"
)

func TestAutocert(t *testing.T) {
	l := NewProvider()
	if _, ok := l.(*autocertProvider); !ok {
		t.Error("NewProvider() didn't return an autocertProvider")
	}
	// TODO: Travis CI doesn't let us bind :443
	// if _, err := l.NewListener(); err != nil {
	// 	t.Error(err.Error())
	// }
}
