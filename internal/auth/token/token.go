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
// Original source: github.com/micro/go-micro/v3/util/token/token.go

package token

import (
	"errors"
	"time"

	"github.com/micro/micro/v3/service/auth"
)

var (
	// ErrNotFound is returned when a token cannot be found
	ErrNotFound = errors.New("token not found")
	// ErrEncodingToken is returned when the service encounters an error during encoding
	ErrEncodingToken = errors.New("error encoding the token")
	// ErrInvalidToken is returned when the token provided is not valid
	ErrInvalidToken = errors.New("invalid token provided")
)

// Provider generates and inspects tokens
type Provider interface {
	Generate(account *auth.Account, opts ...GenerateOption) (*Token, error)
	Inspect(token string) (*auth.Account, error)
	String() string
}

type Token struct {
	// The actual token
	Token string `json:"token"`
	// Time of token creation
	Created time.Time `json:"created"`
	// Time of token expiry
	Expiry time.Time `json:"expiry"`
}
