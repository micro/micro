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
// Original source: github.com/micro/go-micro/v3/debug/log/memory/memory_test.go

package memory

import (
	"reflect"
	"testing"

	"github.com/micro/micro/v3/internal/debug/log"
)

func TestLogger(t *testing.T) {
	// set size to some value
	size := 100
	// override the global logger
	lg := NewLog(log.Size(size))
	// make sure we have the right size of the logger ring buffer
	if lg.(*memoryLog).Size() != size {
		t.Errorf("expected buffer size: %d, got: %d", size, lg.(*memoryLog).Size())
	}

	// Log some cruft
	lg.Write(log.Record{Message: "foobar"})
	lg.Write(log.Record{Message: "foo bar"})

	// Check if the logs are stored in the logger ring buffer
	expected := []string{"foobar", "foo bar"}
	entries, _ := lg.Read(log.Count(len(expected)))
	for i, entry := range entries {
		if !reflect.DeepEqual(entry.Message, expected[i]) {
			t.Errorf("expected %s, got %s", expected[i], entry.Message)
		}
	}
}
