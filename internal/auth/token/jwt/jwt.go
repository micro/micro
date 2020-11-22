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
// Original source: github.com/micro/go-micro/v3/util/token/jwt/jwt.go

package jwt

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/micro/micro/v3/internal/auth/token"
	"github.com/micro/micro/v3/service/auth"
)

// authClaims to be encoded in the JWT
type authClaims struct {
	Type     string            `json:"type"`
	Scopes   []string          `json:"scopes"`
	Metadata map[string]string `json:"metadata"`
	Name     string            `json:"name"`

	jwt.StandardClaims
}

// JWT implementation of token provider
type JWT struct {
	opts token.Options
}

// NewTokenProvider returns an initialized basic provider
func NewTokenProvider(opts ...token.Option) token.Provider {
	return &JWT{
		opts: token.NewOptions(opts...),
	}
}

// Generate a new JWT
func (j *JWT) Generate(acc *auth.Account, opts ...token.GenerateOption) (*token.Token, error) {
	var priv []byte
	if strings.HasPrefix(j.opts.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----") {
		priv = []byte(j.opts.PrivateKey)
	} else {
		var err error
		priv, err = base64.StdEncoding.DecodeString(j.opts.PrivateKey)
		if err != nil {
			return nil, err
		}
	}

	// parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, token.ErrEncodingToken
	}

	// parse the options
	options := token.NewGenerateOptions(opts...)

	// backwards compatibility
	name := acc.Name
	if name == "" {
		name = acc.ID
	}

	// generate the JWT
	expiry := time.Now().Add(options.Expiry)
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, authClaims{
		Type: acc.Type, Scopes: acc.Scopes, Metadata: acc.Metadata, Name: name,
		StandardClaims: jwt.StandardClaims{
			Subject:   acc.ID,
			Issuer:    acc.Issuer,
			ExpiresAt: expiry.Unix(),
		},
	})
	tok, err := t.SignedString(key)
	if err != nil {
		return nil, err
	}

	// return the token
	return &token.Token{
		Token:   tok,
		Expiry:  expiry,
		Created: time.Now(),
	}, nil
}

// Inspect a JWT
func (j *JWT) Inspect(t string) (*auth.Account, error) {
	var pub []byte
	if strings.HasPrefix(j.opts.PublicKey, "-----BEGIN CERTIFICATE-----") {
		pub = []byte(j.opts.PublicKey)
	} else {
		var err error
		pub, err = base64.StdEncoding.DecodeString(j.opts.PublicKey)
		if err != nil {
			return nil, err
		}
	}

	// parse the public key
	res, err := jwt.ParseWithClaims(t, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pub)
	})
	if err != nil {
		return nil, token.ErrInvalidToken
	}

	// validate the token
	if !res.Valid {
		return nil, token.ErrInvalidToken
	}
	claims, ok := res.Claims.(*authClaims)
	if !ok {
		return nil, token.ErrInvalidToken
	}

	// backwards compatibility
	name := claims.Name
	if name == "" {
		name = claims.Subject
	}

	// return the token
	return &auth.Account{
		ID:       claims.Subject,
		Issuer:   claims.Issuer,
		Type:     claims.Type,
		Scopes:   claims.Scopes,
		Metadata: claims.Metadata,
		Name:     name,
	}, nil
}

// String returns JWT
func (j *JWT) String() string {
	return "jwt"
}
