package proxy

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/micro/v2/internal/namespace"
)

// authHandler wraps a server handler to perform auth
func authHandler(fn func() auth.Auth) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// get the auth.Auth interface
			a := fn()

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

			// Extract the namespace header
			ns, ok := metadata.Get(ctx, "Micro-Namespace")
			if !ok && account != nil {
				ns = account.Issuer
				ctx = metadata.Set(ctx, "Micro-Namespace", ns)
			} else if !ok {
				ns = namespace.DefaultNamespace
				ctx = metadata.Set(ctx, "Micro-Namespace", ns)
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
			err = a.Verify(account, res, auth.VerifyNamespace(ns))
			if err == auth.ErrForbidden && account != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), account.ID)
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
