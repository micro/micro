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
// Original source: github.com/micro/go-micro/v3/store/file/blob.go

package file

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/micro/micro/v5/service/store"
	bolt "go.etcd.io/bbolt"
)

// NewBlobStore returns a blob file store
func NewBlobStore(opts ...store.StoreOption) (store.BlobStore, error) {
	// parse the options
	var options store.StoreOptions
	for _, o := range opts {
		o(&options)
	}

	var dir string
	if options.Context != nil {
		if d, ok := options.Context.Value(dirKey{}).(string); ok {
			dir = d
		}
	}
	if len(dir) == 0 {
		dir = DefaultDir
	}

	// ensure the parent directory exists
	os.MkdirAll(dir, 0700)

	return &blobStore{dir}, nil
}

type blobStore struct {
	dir string
}

func (b *blobStore) db() (*bolt.DB, error) {
	dbPath := filepath.Join(b.dir, "blob.db")
	return bolt.Open(dbPath, 0700, &bolt.Options{Timeout: 5 * time.Second})
}

func (b *blobStore) Read(key string, opts ...store.BlobOption) (io.Reader, error) {
	// validate the key
	if len(key) == 0 {
		return nil, store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// open a connection to the database
	db, err := b.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// execute the transaction
	var value []byte
	readValue := func(tx *bolt.Tx) error {
		// check for the namespaces bucket
		bucket := tx.Bucket([]byte(options.Namespace))
		if bucket == nil {
			return store.ErrNotFound
		}

		// look for the blob within the bucket
		res := bucket.Get([]byte(key))
		if res == nil {
			return store.ErrNotFound
		}

		// the res object is only valid for the duration of the blot transaction, see:
		// https://github.com/golang/go/issues/33047
		value = make([]byte, len(res))
		copy(value, res)

		return nil
	}
	if err := db.View(readValue); err != nil {
		return nil, err
	}

	// return the blob
	return bytes.NewBuffer(value), nil
}

func (b *blobStore) Write(key string, blob io.Reader, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// open a connection to the database
	db, err := b.db()
	if err != nil {
		return err
	}
	defer db.Close()

	// execute the transaction
	return db.Update(func(tx *bolt.Tx) error {
		// create the bucket
		bucket, err := tx.CreateBucketIfNotExists([]byte(options.Namespace))
		if err != nil {
			return err
		}

		// write to the bucket
		value, err := ioutil.ReadAll(blob)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(key), value)
	})
}

func (b *blobStore) Delete(key string, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// open a connection to the database
	db, err := b.db()
	if err != nil {
		return err
	}
	defer db.Close()

	// execute the transaction
	return db.Update(func(tx *bolt.Tx) error {
		// check for the namespaces bucket
		bucket := tx.Bucket([]byte(options.Namespace))
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(key))
	})
}

func (b *blobStore) List(opts ...store.BlobListOption) ([]string, error) {
	var options store.BlobListOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}
	// open a connection to the database
	db, err := b.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// execute the transaction
	keys := []string{}
	readValue := func(tx *bolt.Tx) error {
		// check for the namespaces bucket
		bucket := tx.Bucket([]byte(options.Namespace))
		if bucket == nil {
			return store.ErrNotFound
		}
		c := bucket.Cursor()
		for {
			k, _ := c.Next()
			if k == nil {
				break
			}
			kcopy := make([]byte, len(k))
			copy(kcopy, k)
			kstring := string(kcopy)
			if len(options.Prefix) == 0 || strings.HasPrefix(kstring, options.Prefix) {
				keys = append(keys, kstring)
			}

		}
		return nil
	}
	if err := db.View(readValue); err != nil {
		return nil, err
	}

	// return the keys
	return keys, nil
}
