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
// Original source: github.com/micro/go-micro/v3/auth/rules_test.go

package rules

import (
	"testing"

	"github.com/micro/micro/v3/service/auth"
)

func TestVerify(t *testing.T) {
	srvResource := &auth.Resource{
		Type:     "service",
		Name:     "go.micro.service.foo",
		Endpoint: "Foo.Bar",
	}

	webResource := &auth.Resource{
		Type:     "service",
		Name:     "go.micro.web.foo",
		Endpoint: "/foo/bar",
	}

	catchallResource := &auth.Resource{
		Type:     "*",
		Name:     "*",
		Endpoint: "*",
	}

	tt := []struct {
		Name     string
		Rules    []*auth.Rule
		Account  *auth.Account
		Resource *auth.Resource
		Error    error
		Options  []auth.VerifyOption
	}{
		{
			Name:     "NoRules",
			Rules:    []*auth.Rule{},
			Account:  nil,
			Resource: srvResource,
			Error:    auth.ErrForbidden,
		},
		{
			Name:     "CatchallPublicAccount",
			Account:  &auth.Account{},
			Resource: srvResource,
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "",
					Resource: catchallResource,
				},
			},
		},
		{
			Name:     "CatchallPublicNoAccount",
			Resource: srvResource,
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "",
					Resource: catchallResource,
				},
			},
		},
		{
			Name:     "CatchallPrivateAccount",
			Account:  &auth.Account{},
			Resource: srvResource,
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
				},
			},
		},
		{
			Name:     "CatchallPrivateNoAccount",
			Resource: srvResource,
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "CatchallServiceRuleMatch",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     srvResource.Type,
						Name:     srvResource.Name,
						Endpoint: "*",
					},
				},
			},
		},
		{
			Name:     "CatchallServiceRuleNoMatch",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     srvResource.Type,
						Name:     "wrongname",
						Endpoint: "*",
					},
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "ExactRuleValidScope",
			Resource: srvResource,
			Account: &auth.Account{
				Scopes: []string{"neededscope"},
			},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "neededscope",
					Resource: srvResource,
				},
			},
		},
		{
			Name:     "ExactRuleInvalidScope",
			Resource: srvResource,
			Account: &auth.Account{
				Scopes: []string{"neededscope"},
			},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "invalidscope",
					Resource: srvResource,
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "CatchallDenyWithAccount",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessDenied,
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "CatchallDenyWithNoAccount",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessDenied,
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "RulePriorityGrantFirst",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessGranted,
					Priority: 1,
				},
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessDenied,
					Priority: 0,
				},
			},
		},
		{
			Name:     "RulePriorityDenyFirst",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessGranted,
					Priority: 0,
				},
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
					Access:   auth.AccessDenied,
					Priority: 1,
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "WebExactEndpointValid",
			Resource: webResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: webResource,
				},
			},
		},
		{
			Name:     "WebExactEndpointInalid",
			Resource: webResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     webResource.Type,
						Name:     webResource.Name,
						Endpoint: "invalidendpoint",
					},
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "WebWildcardEndpoint",
			Resource: webResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     webResource.Type,
						Name:     webResource.Name,
						Endpoint: "*",
					},
				},
			},
		},
		{
			Name:     "WebWildcardPathEndpointValid",
			Resource: webResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     webResource.Type,
						Name:     webResource.Name,
						Endpoint: "/foo/*",
					},
				},
			},
		},
		{
			Name:     "WebWildcardPathEndpointInvalid",
			Resource: webResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope: "*",
					Resource: &auth.Resource{
						Type:     webResource.Type,
						Name:     webResource.Name,
						Endpoint: "/bar/*",
					},
				},
			},
			Error: auth.ErrForbidden,
		},
		{
			Name:     "CrossNamespaceForbidden",
			Resource: srvResource,
			Account:  &auth.Account{Issuer: "foo"},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
				},
			},
			Error:   auth.ErrForbidden,
			Options: []auth.VerifyOption{auth.VerifyNamespace("bar")},
		},
		{
			Name:     "CrossNamespaceNilAccountForbidden",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    "*",
					Resource: catchallResource,
				},
			},
			Error:   auth.ErrForbidden,
			Options: []auth.VerifyOption{auth.VerifyNamespace("bar")},
		},
		{
			Name:     "CrossNamespacePublic",
			Resource: srvResource,
			Account:  &auth.Account{Issuer: "foo"},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    auth.ScopePublic,
					Resource: catchallResource,
				},
			},
			Options: []auth.VerifyOption{auth.VerifyNamespace("bar")},
		},
		{
			Name:     "CrossNamespacePublicNilAccount",
			Resource: srvResource,
			Account:  &auth.Account{},
			Rules: []*auth.Rule{
				&auth.Rule{
					Scope:    auth.ScopePublic,
					Resource: catchallResource,
				},
			},
			Options: []auth.VerifyOption{auth.VerifyNamespace("bar")},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if err := VerifyAccess(tc.Rules, tc.Account, tc.Resource, tc.Options...); err != tc.Error {
				t.Errorf("Expected %v but got %v", tc.Error, err)
			}
		})
	}
}
