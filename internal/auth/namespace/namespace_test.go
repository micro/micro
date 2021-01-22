package namespace

import (
	"context"
	"testing"

	"github.com/micro/micro/v3/service/auth"

	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	tcs := []struct {
		name   string
		acc    *auth.Account
		ns     string
		method string
		err    string
	}{
		{
			name: "MicroAdminAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroAdminAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroUserAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroUserAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooAdminAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooAdminAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "foo",
			method: "foo.Bar",
		},
		{
			name: "FooUserAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooUserAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "",
		},

		{
			name: "MicroServiceAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroServiceAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "FooServiceAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooServiceAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "foo",
			method: "foo.Bar",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.TODO()
			err := Authorize(auth.ContextWithAccount(ctx, tc.acc), tc.ns, tc.method)
			if tc.err != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthorizeAdmin(t *testing.T) {
	tcs := []struct {
		name   string
		acc    *auth.Account
		ns     string
		method string
		err    string
	}{
		{
			name: "MicroAdminAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroAdminAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroUserAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Unauthorized",
		},
		{
			name: "MicroUserAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooAdminAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooAdminAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"admin"},
			},
			ns:     "foo",
			method: "foo.Bar",
		},
		{
			name: "FooUserAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "user",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooUserAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "Unauthorized",
		},

		{
			name: "MicroServiceAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "MicroServiceAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "micro",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "foo",
			method: "foo.Bar",
			err:    "",
		},
		{
			name: "FooServiceAccessingMicro",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "micro",
			method: "foo.Bar",
			err:    "Forbidden",
		},
		{
			name: "FooServiceAccessingFoo",
			acc: &auth.Account{
				ID:       "1",
				Type:     "service",
				Issuer:   "foo",
				Metadata: nil,
				Scopes:   []string{"service"},
			},
			ns:     "foo",
			method: "foo.Bar",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.TODO()
			err := AuthorizeAdmin(auth.ContextWithAccount(ctx, tc.acc), tc.ns, tc.method)
			if tc.err != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestHasScope(t *testing.T) {
	tcs := []struct {
		name   string
		scope  string
		scopes []string
		result bool
	}{
		{
			name:   "hasScope",
			scope:  "admin",
			scopes: []string{"developer", "admin", "analyst"},
			result: true,
		},
		{
			name:   "hasScopeSingle",
			scope:  "admin",
			scopes: []string{"admin"},
			result: true,
		},
		{
			name:   "noScope",
			scope:  "admin",
			scopes: []string{},
			result: false,
		},
		{
			name:   "noMatch",
			scope:  "developer",
			scopes: []string{"admin", "analyst"},
			result: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.result, hasScope(tc.scope, tc.scopes))
		})

	}
}
