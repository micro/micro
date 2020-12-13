package wrapper

import (
	"context"
	"strings"

	inauth "github.com/micro/micro/v3/internal/auth"
	"github.com/micro/micro/v3/internal/auth/namespace"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/server"
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

	// set the namespace header if it has not been set (e.g. on a service to service request)
	authOpts := auth.DefaultAuth.Options()
	if _, ok := metadata.Get(ctx, "Micro-Namespace"); !ok {
		ctx = metadata.Set(ctx, "Micro-Namespace", authOpts.Issuer)
	}

	// We dont't override the header unless the AuthToken option has been specified
	if !options.AuthToken {
		return ctx
	}

	// check to see if we have a valid access token
	if authOpts.Token != nil && !authOpts.Token.Expired() {
		ctx = metadata.Set(ctx, "Authorization", inauth.BearerScheme+authOpts.Token.AccessToken)
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
			// Extract the token if the header is present. We will inspect the token regardless of if it's
			// present or not since noop auth will return a blank account upon Inspecting a blank token.
			var token string
			if header, ok := metadata.Get(ctx, "Authorization"); ok {
				// Ensure the correct scheme is being used
				if !strings.HasPrefix(header, inauth.BearerScheme) {
					return errors.Unauthorized(req.Service(), "invalid authorization header. expected Bearer schema")
				}

				// Strip the bearer scheme prefix
				token = strings.TrimPrefix(header, inauth.BearerScheme)
			}

			// Determine the namespace
			ns := auth.DefaultAuth.Options().Issuer

			var acc *auth.Account
			if a, err := auth.Inspect(token); err == nil && a.Issuer == ns {
				// We only use accounts issued by the same namespace as the service when verifying against
				// the rule set.
				ctx = auth.ContextWithAccount(ctx, a)
				acc = a
			} else if err == nil && ns == namespace.DefaultNamespace {
				// for the default domain, we want to inject the account into the context so that the
				// server can access it (since it's designed for multi-tenancy), however we don't want to
				// use it when verifying against the auth rules, since this will allow any user access to the
				// services running in the micro namespace
				ctx = auth.ContextWithAccount(ctx, a)
			}

			// ensure only accounts with the correct namespace can access this namespace,
			// since the auth package will verify access below, and some endpoints could
			// be public, we allow nil accounts access using the namespace.Public option.
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
			err = auth.Verify(acc, res, auth.VerifyNamespace(ns))
			if err == auth.ErrForbidden && acc != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), acc.ID)
			} else if err == auth.ErrForbidden {
				return errors.Unauthorized(req.Service(), "Unauthorized call made to %v:%v", req.Service(), req.Endpoint())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), "Error authorizing request: %v", err)
			}

			// The user is authorised, allow the call
			return h(ctx, req, rsp)
		}
	}
}
