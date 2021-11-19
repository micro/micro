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
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	sthree "github.com/aws/aws-sdk-go/service/s3"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

var doubleSlash = regexp.MustCompile("/+")
var removeCol = regexp.MustCompile(":")

func cleanKey(s string) string {
	return doubleSlash.ReplaceAllLiteralString(removeCol.ReplaceAllLiteralString(s, "/"), "/")
}

// NewBlobStore returns an initialized s3 blob store
func NewBlobStore(opts ...Option) (store.BlobStore, error) {
	// parse the options
	options := Options{Secure: true}
	for _, o := range opts {
		o(&options)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:    &options.Endpoint,
		Region:      &options.Region,
		Credentials: credentials.NewStaticCredentials(options.AccessKeyID, options.SecretAccessKey, ""),
	}))
	client := sthree.New(sess)

	// return the blob store
	return &s3{client, &options}, nil
}

type s3 struct {
	client  *sthree.S3
	options *Options
}

func (s *s3) Read(key string, opts ...store.BlobOption) (io.Reader, error) {
	// validate the key
	if len(key) == 0 {
		return nil, store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = cleanKey(key)

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	var err error
	var res *sthree.GetObjectOutput
	if len(s.options.Bucket) > 0 {
		k := filepath.Join(options.Namespace, key)
		res, err = s.client.GetObject(&sthree.GetObjectInput{
			Bucket: &s.options.Bucket, // bucket name
			Key:    &k,                // object name
		})
	} else {
		res, err = s.client.GetObject(&sthree.GetObjectInput{
			Bucket: &options.Namespace, // bucket name
			Key:    &key,               // object name
		})
	}

	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer([]byte{})
	_, err = io.Copy(out, res.Body)

	// return the result
	return out, nil
}

func (s *s3) Write(key string, blob io.Reader, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = cleanKey(key)

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	// if the bucket exists, write using the namespace as a filepath
	buf := new(strings.Builder)
	_, err := io.Copy(buf, blob)
	if err != nil {
		return err
	}
	acl := "private"
	if options.Public {
		acl = "public-read"
	}
	logger.Infof("Saving file %v with ACL %v into namespace %v", key, acl, options.Namespace)
	if len(s.options.Bucket) > 0 {
		k := filepath.Join(options.Namespace, key)
		object := sthree.PutObjectInput{
			Bucket:      &s.options.Bucket,
			Key:         &k,
			Body:        strings.NewReader(buf.String()),
			ACL:         aws.String(acl),
			ContentType: &options.ContentType,
		}
		_, err := s.client.PutObject(&object)
		return err
	}

	s.client.CreateBucket(&sthree.CreateBucketInput{
		Bucket: &options.Namespace,
	})

	k := filepath.Join(options.Namespace, key)
	object := sthree.PutObjectInput{
		Bucket:      &s.options.Bucket,
		Key:         &k,
		Body:        strings.NewReader(buf.String()),
		ACL:         aws.String(acl),
		ContentType: &options.ContentType,
	}
	_, err = s.client.PutObject(&object)
	return err
}

func (s *s3) Delete(key string, opts ...store.BlobOption) error {
	// validate the key
	if len(key) == 0 {
		return store.ErrMissingKey
	}

	// make the key safe for use with s3
	key = cleanKey(key)

	// parse the options
	var options store.BlobOptions
	for _, o := range opts {
		o(&options)
	}
	if len(options.Namespace) == 0 {
		options.Namespace = "micro"
	}

	if len(s.options.Bucket) > 0 {
		k := filepath.Join(options.Namespace, key) // object name
		_, err := s.client.DeleteObject(&sthree.DeleteObjectInput{
			Bucket: &s.options.Bucket, // bucket name
			Key:    &k,
		})
		return err
	}

	k := filepath.Join(options.Namespace, key) // object name
	_, err := s.client.DeleteObject(&sthree.DeleteObjectInput{
		Bucket: &options.Namespace, // bucket name
		Key:    &k,
	})
	return err
}
