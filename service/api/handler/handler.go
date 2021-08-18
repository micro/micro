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
// Original source: github.com/micro/go-micro/v3/api/handler/handler.go

// Package handler provides http handlers
package handler

import (
	"context"
	"net/http"

	"github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/api/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/util/auth/namespace"
)

// Handler represents a HTTP handler that manages a request
type Handler interface {
	// standard http handler
	http.Handler
	// name of handler
	String() string
}

type APIHandler struct {
}

func (a *APIHandler) ReadBlockList(ctx context.Context, request *api.ReadBlockListRequest, response *api.ReadBlockListResponse) error {
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "api.API.AddToBlockList"); err != nil {
		return err
	}
	blocked, err := auth.DefaultBlockList.List(ctx)
	if err != nil {
		return err
	}
	response.Ids = blocked
	return nil
}

func (a *APIHandler) AddToBlockList(ctx context.Context, request *api.AddToBlockListRequest, response *api.AddToBlockListResponse) error {
	// authorize the request
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "api.API.AddToBlockList"); err != nil {
		return err
	}

	if len(request.Id) == 0 {
		return errors.BadRequest("api.AddToBlockList", "Missing ID field")
	}
	if len(request.Namespace) == 0 {
		return errors.BadRequest("api.AddToBlockList", "Missing Namespace field")
	}

	return auth.DefaultBlockList.Add(ctx, request.Id, request.Namespace)
}

func (a *APIHandler) RemoveFromBlockList(ctx context.Context, request *api.RemoveFromBlockListRequest, response *api.RemoveFromBlockListResponse) error {
	if err := namespace.AuthorizeAdmin(ctx, namespace.DefaultNamespace, "api.API.AddToBlockList"); err != nil {
		return err
	}
	if len(request.Id) == 0 {
		return errors.BadRequest("api.RemoveFromBlockList", "Missing ID field")
	}
	if len(request.Namespace) == 0 {
		return errors.BadRequest("api.RemoveFromBlockList", "Missing Namespace field")
	}

	return auth.DefaultBlockList.Remove(ctx, request.Id, request.Namespace)
}
