package namespace

import (
	"context"
	"testing"

	"github.com/micro/micro/v5/service/auth"

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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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
				Type:     "user",
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

func TestHasTypeAndScope(t *testing.T) {
	tcs := []struct {
		name   string
		atype  string
		scope  string
		scopes []string
		result bool
	}{
		{
			name:  "hasScope",
			scope: "admin",

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
			acc := auth.Account{
				ID:     tc.name,
				Type:   tc.atype,
				Issuer: "foobar",
				Scopes: tc.scopes,
				Name:   tc.name,
			}
			assert.Equal(t, tc.result, hasTypeAndScope(tc.atype, tc.scope, &acc))
		})

	}
}
