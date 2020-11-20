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
// Original source: github.com/micro/go-micro/v3/store/s3/s3.go

package s3

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/micro/micro/v3/service/store"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

var keyRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

// NewBlobStore returns an initialized s3 blob store
func NewBlobStore(opts ...Option) (store.BlobStore, error) {
	// parse the options
	options := Options{Secure: true}
	for _, o := range opts {
		o(&options)
	}
	minioOpts := &minio.Options{
		Secure: options.Secure,
	}
	if len(options.AccessKeyID) > 0 || len(options.SecretAccessKey) > 0 {
		minioOpts.Creds = credentials.NewStaticV2(options.AccessKeyID, options.SecretAccessKey, "")
	}

	// configure the transport to use custom tls config if provided
	if options.TLSConfig != nil {
		ts, err := minio.DefaultTransport(options.Secure)
		if err != nil {
			return nil, errors.Wrap(err, "Error setting up s3 blob store transport")
		}
		ts.TLSClientConfig = options.TLSConfig
		minioOpts.Transport = ts
	}

	// initialize minio client
	client, err := minio.New(options.Endpoint, minioOpts)
	if err != nil {
		return nil, errors.Wrap(err, "Error connecting to s3 blob store")
	}

	// return the blob store
	return &s3{client, &options}, nil
}

type s3 struct {
	client  *minio.Client
	options *Options
}

func (s *s3) Read(key string, opts ...store.BlobOption) (io.Reader, error) {
	// validate the key
	if len(key) == 0 {
		return nil, store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = keyRegex.ReplaceAllString(key, "-")

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// lookup the object
	var res *minio.Object
	var err error
	if len(s.options.Bucket) > 0 {
		res, err = s.client.GetObject(
			context.TODO(),                        // context
			s.options.Bucket,                      // bucket name
			filepath.Join(options.Namespace, key), // object name
			minio.GetObjectOptions{},              // options
		)
	} else {
		res, err = s.client.GetObject(
			context.TODO(),           // context
			options.Namespace,        // bucket name
			key,                      // object name
			minio.GetObjectOptions{}, // options
		)
	}

	// scaleway will return a 404 if the bucket doesn't exist
	if verr, ok := err.(minio.ErrorResponse); ok && verr.StatusCode == http.StatusNotFound {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// check the object info, if an error is returned the object could not be found
	_, err = res.Stat()
	if verr, ok := err.(minio.ErrorResponse); ok && verr.StatusCode == http.StatusNotFound {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// return the result
	return res, nil
}

func (s *s3) Write(key string, blob io.Reader, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = keyRegex.ReplaceAllString(key, "-")

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// if the bucket exists, write using the namespace as a filepath
	if len(s.options.Bucket) > 0 {
		_, err := s.client.PutObject(
			context.TODO(),                        // context
			s.options.Bucket,                      // bucket name
			filepath.Join(options.Namespace, key), // object name
			blob,                                  // reader
			int64(-1),                             // length of object
			minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			},
		)
		return err
	}

	// check the bucket exists, create it if not
	if exists, err := s.client.BucketExists(context.TODO(), options.Namespace); err != nil {
		return err
	} else if !exists {
		opts := minio.MakeBucketOptions{Region: s.options.Region}
		if err := s.client.MakeBucket(context.TODO(), options.Namespace, opts); err != nil {
			return err
		}
	}

	// create the object in the bucket
	_, err := s.client.PutObject(
		context.TODO(),    // context
		options.Namespace, // bucket name
		key,               // object name
		blob,              // reader
		int64(-1),         // length of object
		minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	return err
}

func (s *s3) Delete(key string, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = keyRegex.ReplaceAllString(key, "-")

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	if len(s.options.Bucket) > 0 {
		return s.client.RemoveObject(
			context.TODO(),                        // context
			s.options.Bucket,                      // bucket name
			filepath.Join(options.Namespace, key), // object name
			minio.RemoveObjectOptions{},           // options
		)
	}

	return s.client.RemoveObject(
		context.TODO(),              // context
		options.Namespace,           // bucket name
		key,                         // object name
		minio.RemoveObjectOptions{}, // options
	)
}
