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
// Original source: github.com/micro/go-micro/v3/runtime/local/local_test.go

package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntrypoint(t *testing.T) {
	wd, _ := os.Getwd()

	// test a service with a single .mu file
	result, err := Entrypoint(filepath.Join(wd, "test/bar"))
	assert.Nil(t, err, "Didn'expected entrypoint to return an error")
	assert.Equal(t, "bar.mu", result, "Expected entrypoint to return bar.mu")

	// test a service with multiple .mu files within the cmd folder
	result, err = Entrypoint(filepath.Join(wd, "test/foo"))
	assert.Error(t, err, "Expected entrypoint to return an error when multiple .mu files exist")
	assert.Equal(t, "", result, "Expected entrypoint to not return a result")

	// test a service with no .mu or main.go files
	result, err = Entrypoint(filepath.Join(wd, "test/empty"))
	assert.Error(t, err, "Expected entrypoint to return an error when no .mu files exist and no main.go files exist")
	assert.Equal(t, "", result, "Expected entrypoint to not return a result")

	// test a service with idiomatic folder structure and no .mu file
	result, err = Entrypoint(filepath.Join(wd, "test/qux"))
	assert.Nil(t, err, "Didn't expecte entrypoint to return an error")
	assert.Equal(t, "cmd/test/main.go", result, "Expected entrypoint to return cmd/test/main.go")
}
