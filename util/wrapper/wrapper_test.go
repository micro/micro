package wrapper

import (
	"context"
	"testing"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/util/codec"

	. "github.com/onsi/gomega"
)

type dummyAuth struct {
	opts auth.Options
}

func (d dummyAuth) Init(opts ...auth.Option) {
}

func (d dummyAuth) Options() auth.Options {
	return d.opts
}

func (d dummyAuth) Generate(id string, opts ...auth.GenerateOption) (*auth.Account, error) {
	options := auth.NewGenerateOptions(opts...)
	name := options.Name
	if name == "" {
		name = id
	}
	return &auth.Account{
		ID:       id,
		Secret:   options.Secret,
		Metadata: options.Metadata,
		Scopes:   options.Scopes,
		Issuer:   d.Options().Issuer,
		Name:     name,
	}, nil
}

func (d dummyAuth) Verify(acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	return nil
}

func (d dummyAuth) Inspect(token string) (*auth.Account, error) {
	return &auth.Account{ID: token, Issuer: d.Options().Issuer}, nil
}

func (d dummyAuth) Token(opts ...auth.TokenOption) (*auth.AccountToken, error) {
	return &auth.AccountToken{}, nil
}

func (d dummyAuth) Grant(rule *auth.Rule) error {
	return nil
}

func (d dummyAuth) Revoke(rule *auth.Rule) error {
	return nil
}

func (d dummyAuth) Rules(option ...auth.RulesOption) ([]*auth.Rule, error) {
	return nil, nil
}

func (d dummyAuth) String() string {
	return "dummyAuth"
}

type dummyReq struct {
}

func (d dummyReq) Service() string {
	return "dummy"
}

func (d dummyReq) Method() string {
	return "dummy"
}

func (d dummyReq) Endpoint() string {
	return "dummy"
}

func (d dummyReq) ContentType() string {
	return "application/json"
}

func (d dummyReq) Header() map[string]string {
	return map[string]string{}
}

func (d dummyReq) Body() interface{} {
	return nil
}

func (d dummyReq) Read() ([]byte, error) {
	panic("implement me")
}

func (d dummyReq) Codec() codec.Reader {
	panic("implement me")
}

func (d dummyReq) Stream() bool {
	panic("implement me")
}

func TestAuthWrapper(t *testing.T) {
	g := NewWithT(t)

	tcs := []struct {
		name    string
		authHdr string
		err     error
		tok     string
	}{
		{
			name:    "Bearer auth",
			authHdr: "Bearer 123355",
			tok:     "123355",
			err:     nil,
		},
		{
			name:    "Basic auth",
			authHdr: "Basic Zm9vOmJhcg==",
			tok:     "bar",
			err:     nil,
		},
		{
			name:    "Basic auth but wrong format of user:pass",
			authHdr: "Basic Zm9vYmFyCg==",
			err:     errors.Unauthorized("dummy", "invalid authorization header. Incorrect format"),
		},
		{
			name:    "Unknown auth",
			authHdr: "Foobar 11111",
			err:     errors.Unauthorized("dummy", "invalid authorization header. Expected Bearer or Basic schema"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			auth.DefaultAuth = dummyAuth{}

			w := AuthHandler()
			dummyInvoked := false
			tok := ""
			dummy := func(ctx context.Context, req server.Request, rsp interface{}) error {
				dummyInvoked = true
				acc, _ := auth.AccountFromContext(ctx)
				// dummyAuth sets acc ID to the token string so we can easily test it
				tok = acc.ID
				return nil
			}
			ctx := context.Background()
			ctx = metadata.Set(ctx, "Authorization", tc.authHdr)
			err := w(dummy)(ctx, &dummyReq{}, nil)

			if tc.err == nil {
				g.Expect(dummyInvoked).To(BeTrue())
				g.Expect(tok).To(Equal(tc.tok))
			} else {
				g.Expect(dummyInvoked).To(BeFalse())
				g.Expect(err).To(Equal(tc.err))
			}

		})
	}
}
