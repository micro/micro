package wrapper

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/micro/v2/internal/namespace"
	muauth "github.com/micro/micro/v2/service/auth"
)

type authWrapper struct {
	client.Client
}

func (a *authWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = a.wrapContext(ctx, opts...)
	return a.Client.Call(ctx, req, rsp, opts...)
}

func (a *authWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	ctx = a.wrapContext(ctx, opts...)
	return a.Client.Stream(ctx, req, opts...)
}

func (a *authWrapper) wrapContext(ctx context.Context, opts ...client.CallOption) context.Context {
	// parse the options
	var options client.CallOptions
	for _, o := range opts {
		o(&options)
	}

	// check to see if the authorization header has already been set.
	// We dont't override the header unless the ServiceToken option has
	// been specified or the header wasn't provided
	if _, ok := metadata.Get(ctx, "Authorization"); ok && !options.ServiceToken {
		return ctx
	}

	// if auth is nil we won't be able to get an access token, so we execute
	// the request without one.
	aa := muauth.DefaultAuth
	if aa == nil {
		return ctx
	}

	// set the namespace header if it has not been set (e.g. on a service to service request)
	if _, ok := metadata.Get(ctx, "Micro-Namespace"); !ok {
		ctx = metadata.Set(ctx, "Micro-Namespace", aa.Options().Issuer)
	}

	// check to see if we have a valid access token
	aaOpts := aa.Options()
	if aaOpts.Token != nil && !aaOpts.Token.Expired() {
		ctx = metadata.Set(ctx, "Authorization", auth.BearerScheme+aaOpts.Token.AccessToken)
		return ctx
	}

	// call without an auth token
	return ctx
}

// AuthClient wraps requests with the auth header
func AuthClient(c client.Client) client.Client {
	return &authWrapper{c}
}

// AuthHandler wraps a server handler to perform auth
func AuthHandler() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// get the auth.Auth interface
			a := muauth.DefaultAuth

			// Extract the token if the header is present. We will inspect the token regardless of if it's
			// present or not since noop auth will return a blank account upon Inspecting a blank token.
			var token string
			if header, ok := metadata.Get(ctx, "Authorization"); ok {
				// Ensure the correct scheme is being used
				if !strings.HasPrefix(header, auth.BearerScheme) {
					return errors.Unauthorized(req.Service(), "invalid authorization header. expected Bearer schema")
				}

				// Strip the bearer scheme prefix
				token = strings.TrimPrefix(header, auth.BearerScheme)
			}

			// Inspect the token and decode an account
			account, _ := a.Inspect(token)

			// ensure only accounts with the correct namespace can access this namespace,
			// since the auth package will verify access below, and some endpoints could
			// be public, we allow nil accounts access using the namespace.Public option.
			ns := a.Options().Issuer
			err := namespace.Authorize(ctx, ns, namespace.Public(ns))
			if err == namespace.ErrForbidden {
				return errors.Forbidden(req.Service(), err.Error())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), err.Error())
			}

			// construct the resource
			res := &auth.Resource{
				Type:     "service",
				Name:     req.Service(),
				Endpoint: req.Endpoint(),
			}

			// Verify the caller has access to the resource.
			err = a.Verify(account, res, auth.VerifyNamespace(ns))
			if err == auth.ErrForbidden && account != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), account.ID)
			} else if err == auth.ErrForbidden {
				return errors.Unauthorized(req.Service(), "Unauthorized call made to %v:%v", req.Service(), req.Endpoint())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), "Error authorizing request: %v", err)
			}

			// There is an account, set it in the context
			if account != nil {
				ctx = auth.ContextWithAccount(ctx, account)
			}

			// The user is authorised, allow the call
			return h(ctx, req, rsp)
		}
	}
}
