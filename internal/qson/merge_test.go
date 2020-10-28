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
// Original source: github.com/micro/go-micro/v3/util/qson/merge_test.go

package qson

import "testing"

func TestMergeSlice(t *testing.T) {
	a := []interface{}{"a"}
	b := []interface{}{"b"}
	actual := mergeSlice(a, b)
	if len(actual) != 2 {
		t.Errorf("Expected size to be 2.")
	}
	if actual[0] != "a" {
		t.Errorf("Expected index 0 to have value a. Actual: %s", actual[0])
	}
	if actual[1] != "b" {
		t.Errorf("Expected index 1 to have value b. Actual: %s", actual[1])
	}
}

func TestMergeMap(t *testing.T) {
	a := map[string]interface{}{
		"a": "b",
	}
	b := map[string]interface{}{
		"b": "c",
	}
	actual := mergeMap(a, b)
	if len(actual) != 2 {
		t.Errorf("Expected size to be 2.")
	}
	if actual["a"] != "b" {
		t.Errorf("Expected key \"a\" to have value b. Actual: %s", actual["a"])
	}
	if actual["b"] != "c" {
		t.Errorf("Expected key \"b\" to have value c. Actual: %s", actual["b"])
	}
}
