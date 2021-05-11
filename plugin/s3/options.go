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
// Original source: github.com/micro/go-micro/v3/store/s3/options.go

package s3

import "crypto/tls"

// Options used to configure the s3 blob store
type Options struct {
	Bucket          string
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
	TLSConfig       *tls.Config
}

// Option configures one or more options
type Option func(o *Options)

// Endpoint sets the endpoint option
func Endpoint(e string) Option {
	return func(o *Options) {
		o.Endpoint = e
	}
}

// Region sets the region option
func Region(r string) Option {
	return func(o *Options) {
		o.Region = r
	}
}

// Credentials sets the AccessKeyID and SecretAccessKey options
func Credentials(id, secret string) Option {
	return func(o *Options) {
		o.AccessKeyID = id
		o.SecretAccessKey = secret
	}
}

// Bucket sets the bucket name option
func Bucket(name string) Option {
	return func(o *Options) {
		o.Bucket = name
	}
}

// Insecure sets the secure option to false. It is enabled by default.
func Insecure() Option {
	return func(o *Options) {
		o.Secure = false
	}
}

// TLSConfig sets the tls config for the client
func TLSConfig(c *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = c
	}
}
