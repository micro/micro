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
// Original source: github.com/micro/go-micro/v3/server/grpc/error.go

package grpc

import (
	"net/http"

	"github.com/micro/micro/v3/service/errors"
	"google.golang.org/grpc/codes"
)

var errMapping = map[int32]codes.Code{
	http.StatusOK:                  codes.OK,
	http.StatusBadRequest:          codes.InvalidArgument,
	http.StatusRequestTimeout:      codes.DeadlineExceeded,
	http.StatusNotFound:            codes.NotFound,
	http.StatusConflict:            codes.AlreadyExists,
	http.StatusForbidden:           codes.PermissionDenied,
	http.StatusUnauthorized:        codes.Unauthenticated,
	http.StatusPreconditionFailed:  codes.FailedPrecondition,
	http.StatusNotImplemented:      codes.Unimplemented,
	http.StatusInternalServerError: codes.Internal,
	http.StatusServiceUnavailable:  codes.Unavailable,
}

func microError(err *errors.Error) codes.Code {
	if err == nil {
		return codes.OK
	}

	if code, ok := errMapping[err.Code]; ok {
		return code
	}
	return codes.Unknown
}
