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
// Original source: github.com/micro/go-micro/v3/auth/jwt/jwt.go

// Package jwt is a jwt implementation of the auth interface
package jwt

import (
	"sync"
	"time"

	"github.com/micro/micro/v3/internal/auth/rules"
	"github.com/micro/micro/v3/internal/auth/token"
	"github.com/micro/micro/v3/internal/auth/token/jwt"
	"github.com/micro/micro/v3/service/auth"
)

// NewAuth returns a new instance of the Auth service
func NewAuth(opts ...auth.Option) auth.Auth {
	j := new(jwtAuth)
	j.Init(opts...)
	return j
}

type jwtAuth struct {
	options auth.Options
	token   token.Provider
	rules   []*auth.Rule

	sync.Mutex
}

func (j *jwtAuth) String() string {
	return "jwt"
}

func (j *jwtAuth) Init(opts ...auth.Option) {
	j.Lock()
	defer j.Unlock()

	for _, o := range opts {
		o(&j.options)
	}

	j.token = jwt.NewTokenProvider(
		token.WithPrivateKey(j.options.PrivateKey),
		token.WithPublicKey(j.options.PublicKey),
	)
}

func (j *jwtAuth) Options() auth.Options {
	j.Lock()
	defer j.Unlock()
	return j.options
}

func (j *jwtAuth) Generate(id string, opts ...auth.GenerateOption) (*auth.Account, error) {
	options := auth.NewGenerateOptions(opts...)
	if len(options.Issuer) == 0 {
		options.Issuer = j.Options().Issuer
	}
	name := options.Name
	if name == "" {
		name = id
	}
	account := &auth.Account{
		ID:       id,
		Type:     options.Type,
		Scopes:   options.Scopes,
		Metadata: options.Metadata,
		Issuer:   options.Issuer,
		Name:     name,
	}

	// generate a JWT secret which can be provided to the Token() method
	// and exchanged for an access token
	secret, err := j.token.Generate(account, token.WithExpiry(time.Hour*24*365))
	if err != nil {
		return nil, err
	}
	account.Secret = secret.Token

	// return the account
	return account, nil
}

func (j *jwtAuth) Grant(rule *auth.Rule) error {
	j.Lock()
	defer j.Unlock()
	j.rules = append(j.rules, rule)
	return nil
}

func (j *jwtAuth) Revoke(rule *auth.Rule) error {
	j.Lock()
	defer j.Unlock()

	rules := []*auth.Rule{}
	for _, r := range j.rules {
		if r.ID != rule.ID {
			rules = append(rules, r)
		}
	}

	j.rules = rules
	return nil
}

func (j *jwtAuth) Verify(acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	j.Lock()
	defer j.Unlock()

	return rules.VerifyAccess(j.rules, acc, res, opts...)
}

func (j *jwtAuth) Rules(opts ...auth.RulesOption) ([]*auth.Rule, error) {
	j.Lock()
	defer j.Unlock()
	return j.rules, nil
}

func (j *jwtAuth) Inspect(token string) (*auth.Account, error) {
	return j.token.Inspect(token)
}

func (j *jwtAuth) Token(opts ...auth.TokenOption) (*auth.AccountToken, error) {
	options := auth.NewTokenOptions(opts...)

	secret := options.RefreshToken
	if len(options.Secret) > 0 {
		secret = options.Secret
	}

	account, err := j.token.Inspect(secret)
	if err != nil {
		return nil, err
	}

	access, err := j.token.Generate(account, token.WithExpiry(options.Expiry))
	if err != nil {
		return nil, err
	}

	refresh, err := j.token.Generate(account, token.WithExpiry(options.Expiry+time.Hour))
	if err != nil {
		return nil, err
	}

	return &auth.AccountToken{
		Created:      access.Created,
		Expiry:       access.Expiry,
		AccessToken:  access.Token,
		RefreshToken: refresh.Token,
	}, nil
}
