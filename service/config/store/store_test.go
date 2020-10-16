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
// Original source: github.com/micro/go-micro/v3/config/store/store_test.go
package store

import (
	"reflect"
	"testing"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/config/secrets"
	"github.com/micro/micro/v3/service/store/memory"
)

type conf1 struct {
	A string  `json:"a"`
	B int64   `json:"b"`
	C float64 `json:"c"`
	D bool    `json:"d"`
}

func TestBasics(t *testing.T) {
	conf, err := NewConfig(memory.NewStore(), "micro")
	if err != nil {
		t.Fatal(err)
	}
	testBasics(conf, t)
	// We need to get a new config because existing config so
	conf, err = NewConfig(memory.NewStore(), "micro1")
	if err != nil {
		t.Fatal(err)
	}
	secrets, err := secrets.NewSecrets(conf, "somethingRandomButLongEnough32by")
	if err != nil {
		t.Fatal(err)
	}
	testBasics(secrets, t)
}

func testBasics(c config.Config, t *testing.T) {
	original := &conf1{
		"Hi", int64(42), float64(42.2), true,
	}
	err := c.Set("key", original)
	if err != nil {
		t.Fatal(err)
	}
	getted := &conf1{}
	val, err := c.Get("key")
	if err != nil {
		t.Fatal(err)
	}
	err = val.Scan(getted)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, getted) {
		t.Fatalf("Not equal: %v and %v", original, getted)
	}

	// Testing merges now
	err = c.Set("key", map[string]interface{}{
		"b": 55,
		"e": map[string]interface{}{
			"e1": true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{}
	val, err = c.Get("key")
	if err != nil {
		t.Fatal(err)
	}
	err = val.Scan(&m)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]interface{}{
		"a": "Hi",
		"b": float64(55),
		"c": float64(42.2),
		"d": true,
		"e": map[string]interface{}{
			"e1": true,
		},
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("Not equal: %v and %v", m, expected)
	}

	// Set just one value
	expected = map[string]interface{}{
		"a": "Hi",
		"b": float64(55),
		"c": float64(42.2),
		"d": true,
		"e": map[string]interface{}{
			"e1": float64(45),
		},
	}
	err = c.Set("key.e.e1", 45)
	if err != nil {
		t.Fatal(err)
	}

	m = map[string]interface{}{}
	val, err = c.Get("key")
	if err != nil {
		t.Fatal(err)
	}
	err = val.Scan(&m)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("Not equal: %v and %v", m, expected)
	}
}
