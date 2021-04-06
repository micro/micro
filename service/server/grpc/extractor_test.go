// Copyright 2020 Asim Aslam
//
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
// Original source: github.com/micro/go-micro/v3/server/grpc/extractor_test.go

package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestHandler struct{}

type TestRequest struct{}

type TestResponse struct{}

func (t *TestHandler) Test(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	return nil
}

type TestRequest2 struct {
	int64Arg        int64                         `json:"int64Arg,omitempty"`
	stringArg       string                        `json:"stringArg,omitempty"`
	stringSliceArg  []string                      `json:"stringSliceArg,omitempty"`
	pointerSliceArg []*TestRequest                `json:"pointerSliceArg,omitempty"`
	mapArg          map[string]bool               `json:"mapArg,omitempty"`
	mapPointerArg   map[*TestRequest]TestResponse `json:"mapPointerArg,omitempty"`
}
type TestResponse2 struct {
	int64Arg        int64                         `json:"int64Arg,omitempty"`
	stringArg       string                        `json:"stringArg,omitempty"`
	stringSliceArg  []string                      `json:"stringSliceArg,omitempty"`
	pointerSliceArg []*TestRequest                `json:"pointerSliceArg,omitempty"`
	mapArg          map[string]bool               `json:"mapArg,omitempty"`
	mapPointerArg   map[*TestRequest]TestResponse `json:"mapPointerArg,omitempty"`
}

func (t *TestHandler) Test2(ctx context.Context, req *TestRequest2, rsp *TestResponse2) error {
	return nil
}

func TestExtractEndpoint(t *testing.T) {
	handler := &TestHandler{}
	typ := reflect.TypeOf(handler)

	type param struct {
		name  string
		value string
	}
	tcs := []struct {
		name    string
		reqName string
		reqType string
		reqArgs []param
		rspName string
		rspType string
		rspArgs []param
	}{
		{
			name:    "Test",
			reqName: "TestRequest",
			reqType: "TestRequest",
			rspName: "TestResponse",
			rspType: "TestResponse",
		},
		{
			name:    "Test2",
			reqName: "TestRequest2",
			reqType: "TestRequest2",
			rspName: "TestResponse2",
			rspType: "TestResponse2",
			reqArgs: []param{
				{
					name:  "int64Arg",
					value: "int64",
				},
				{
					name:  "stringArg",
					value: "string",
				},
				{
					name:  "stringSliceArg",
					value: "[]string",
				},
				{
					name:  "pointerSliceArg",
					value: "[]TestRequest",
				},
				{
					name:  "mapArg",
					value: "map[string]bool",
				},
				{
					name:  "mapPointerArg",
					value: "map[TestRequest]TestResponse",
				},
			},
			rspArgs: []param{
				{
					name:  "int64Arg",
					value: "int64",
				},
				{
					name:  "stringArg",
					value: "string",
				},
				{
					name:  "stringSliceArg",
					value: "[]string",
				},
				{
					name:  "pointerSliceArg",
					value: "[]TestRequest",
				},
				{
					name:  "mapArg",
					value: "map[string]bool",
				},
				{
					name:  "mapPointerArg",
					value: "map[TestRequest]TestResponse",
				},
			},
		},
	}

	for _, tc := range tcs {
		m, ok := typ.MethodByName(tc.name)
		assert.True(t, ok)
		e := extractEndpoint(m)
		assert.Equal(t, tc.name, e.Name)
		assert.NotNil(t, e.Request)
		assert.NotNil(t, e.Response)
		assert.Equal(t, tc.reqName, e.Request.Name)
		assert.Equal(t, tc.reqType, e.Request.Type)
		assert.Equal(t, tc.rspName, e.Response.Name)
		assert.Equal(t, tc.rspType, e.Response.Type)
		for i, v := range tc.reqArgs {
			assert.Equal(t, v.name, e.Request.Values[i].Name)
			assert.Equal(t, v.value, e.Request.Values[i].Type)
		}
		for i, v := range tc.rspArgs {
			assert.Equal(t, v.name, e.Response.Values[i].Name)
			assert.Equal(t, v.value, e.Response.Values[i].Type)
		}
	}

}
