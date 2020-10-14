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
// Original source: github.com/micro/go-micro/v3/store/cache/cache_test.go

package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/file"
	"github.com/stretchr/testify/assert"
)

func cleanup(db string, s store.Store) {
	s.Close()
	dir := filepath.Join(file.DefaultDir, db+"/")
	os.RemoveAll(dir)
}

func TestRead(t *testing.T) {
	cf := NewStore(file.NewStore())
	cf.Init()
	cfInt := cf.(*cache)
	defer cleanup(file.DefaultDatabase, cf)

	_, err := cf.Read("key1")
	assert.Error(t, err, "Unexpected record")
	cfInt.b.Write(&store.Record{
		Key:   "key1",
		Value: []byte("foo"),
	})
	recs, err := cf.Read("key1")
	assert.NoError(t, err)
	assert.Len(t, recs, 1, "Expected a record to be pulled from file store")
	recs, err = cfInt.m.Read("key1")
	assert.NoError(t, err)
	assert.Len(t, recs, 1, "Expected a memory store to be populatedfrom file store")

}

func TestWrite(t *testing.T) {
	cf := NewStore(file.NewStore())
	cf.Init()
	cfInt := cf.(*cache)
	defer cleanup(file.DefaultDatabase, cf)

	cf.Write(&store.Record{
		Key:   "key1",
		Value: []byte("foo"),
	})
	recs, _ := cfInt.m.Read("key1")
	assert.Len(t, recs, 1, "Expected a record in the memory store")
	recs, _ = cfInt.b.Read("key1")
	assert.Len(t, recs, 1, "Expected a record in the file store")

}

func TestDelete(t *testing.T) {
	cf := NewStore(file.NewStore())
	cf.Init()
	cfInt := cf.(*cache)
	defer cleanup(file.DefaultDatabase, cf)

	cf.Write(&store.Record{
		Key:   "key1",
		Value: []byte("foo"),
	})
	recs, _ := cfInt.m.Read("key1")
	assert.Len(t, recs, 1, "Expected a record in the memory store")
	recs, _ = cfInt.b.Read("key1")
	assert.Len(t, recs, 1, "Expected a record in the file store")
	cf.Delete("key1")

	_, err := cfInt.m.Read("key1")
	assert.Error(t, err, "Expected no records in memory store")
	_, err = cfInt.b.Read("key1")
	assert.Error(t, err, "Expected no records in file store")

}

func TestList(t *testing.T) {
	cf := NewStore(file.NewStore())
	cf.Init()
	cfInt := cf.(*cache)
	defer cleanup(file.DefaultDatabase, cf)

	keys, err := cf.List()
	assert.NoError(t, err)
	assert.Len(t, keys, 0)
	cfInt.b.Write(&store.Record{
		Key:   "key1",
		Value: []byte("foo"),
	})

	cfInt.b.Write(&store.Record{
		Key:   "key2",
		Value: []byte("foo"),
	})
	keys, err = cf.List()
	assert.NoError(t, err)
	assert.Len(t, keys, 2)

}
