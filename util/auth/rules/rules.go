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
// Original source: github.com/micro/go-micro/v3/auth/rules.go

package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/micro/micro/v5/service/auth"
)

// VerifyAccess an account has access to a resource using the rules provided. If the account does not have
// access an error will be returned. If there are no rules provided which match the resource, an error
// will be returned
func VerifyAccess(rules []*auth.Rule, acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	// the rule is only to be applied if the type matches the resource or is catch-all (*)
	validTypes := []string{"*", res.Type}

	// the rule is only to be applied if the name matches the resource or is catch-all (*)
	validNames := []string{"*", res.Name}

	options := auth.VerifyOptions{}
	for _, o := range opts {
		o(&options)
	}

	// rules can have wildcard excludes on endpoints since this can also be a path for web services,
	// e.g. /foo/* would include /foo/bar. We also want to check for wildcards and the exact endpoint
	validEndpoints := []string{"*", res.Endpoint}
	if comps := strings.Split(res.Endpoint, "/"); len(comps) > 1 {
		for i := 1; i < len(comps)+1; i++ {
			wildcard := fmt.Sprintf("%v/*", strings.Join(comps[0:i], "/"))
			validEndpoints = append(validEndpoints, wildcard)
		}
	}

	// filter the rules to the ones which match the criteria above
	filteredRules := make([]*auth.Rule, 0)
	for _, rule := range rules {
		if !include(validTypes, rule.Resource.Type) {
			continue
		}
		if !include(validNames, rule.Resource.Name) {
			continue
		}
		if !include(validEndpoints, rule.Resource.Endpoint) {
			continue
		}
		filteredRules = append(filteredRules, rule)
	}

	// sort the filtered rules by priority, highest to lowest
	sort.SliceStable(filteredRules, func(i, j int) bool {
		return filteredRules[i].Priority > filteredRules[j].Priority
	})

	// loop through the rules and check for a rule which applies to this account
	for _, rule := range filteredRules {
		// a blank scope indicates the rule applies to everyone, even nil accounts
		if rule.Scope == auth.ScopePublic && rule.Access == auth.AccessDenied {
			return auth.ErrForbidden
		} else if rule.Scope == auth.ScopePublic && rule.Access == auth.AccessGranted {
			return nil
		}

		// all further checks require an account
		if acc == nil {
			continue
		}

		// TODO should this live here or further up?
		if rule.Scope != auth.ScopeAnyNamespaceAccount && acc.Issuer != options.Namespace {
			return auth.ErrForbidden
		}
		// TODO what does options.Context do?

		// this rule applies to any account
		if (rule.Scope == auth.ScopeAccount || rule.Scope == auth.ScopeAnyNamespaceAccount) && rule.Access == auth.AccessDenied {
			return auth.ErrForbidden
		} else if (rule.Scope == auth.ScopeAccount || rule.Scope == auth.ScopeAnyNamespaceAccount) && rule.Access == auth.AccessGranted {
			return nil
		}

		// if the account has the necessary scope
		if include(acc.Scopes, rule.Scope) && rule.Access == auth.AccessDenied {
			return auth.ErrForbidden
		} else if include(acc.Scopes, rule.Scope) && rule.Access == auth.AccessGranted {
			return nil
		}
	}

	// if no rules matched then return forbidden
	return auth.ErrForbidden
}

// include is a helper function which checks to see if the slice contains the value. includes is
// not case sensitive.
func include(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}
