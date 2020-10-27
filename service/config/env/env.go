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
// Original source: github.com/micro/go-micro/v3/config/env/env.go

// Package env provides config from environment variables
package env

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/micro/micro/v3/service/config"
)

type envConfig struct{}

// NewConfig returns new config
func NewConfig() (*envConfig, error) {
	return new(envConfig), nil
}

func formatKey(v string) string {
	if len(v) == 0 {
		return ""
	}

	v = strings.ToUpper(v)
	return strings.Replace(v, ".", "_", -1)
}

func (c *envConfig) Get(path string, options ...config.Option) (config.Value, error) {
	v := os.Getenv(formatKey(path))
	if len(v) == 0 {
		v = "{}"
	}
	return config.NewJSONValue([]byte(v)), nil
}

func (c *envConfig) Set(path string, val interface{}, options ...config.Option) error {
	key := formatKey(path)
	v, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return os.Setenv(key, string(v))
}

func (c *envConfig) Delete(path string, options ...config.Option) error {
	v := formatKey(path)
	return os.Unsetenv(v)
}
