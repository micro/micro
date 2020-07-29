package proxy

import (
	"context"
	"strings"

	goauth "github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/metadata"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
)

// authHandler wraps a server handler to perform auth
func authHandler() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// Extract the token if the header is present. We will inspect the token regardless of if it's
			// present or not since noop auth will return a blank account upon Inspecting a blank token.
			var token string
			if header, ok := metadata.Get(ctx, "Authorization"); ok {
				// Ensure the correct scheme is being used
				if !strings.HasPrefix(header, goauth.BearerScheme) {
					return errors.Unauthorized(req.Service(), "invalid authorization header. expected Bearer schema")
				}

				// Strip the bearer scheme prefix
				token = strings.TrimPrefix(header, goauth.BearerScheme)
			}

			// Inspect the token and decode an account
			account, _ := auth.Inspect(token)

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
			res := &goauth.Resource{
				Type:     "service",
				Name:     req.Service(),
				Endpoint: req.Endpoint(),
			}

			// Verify the caller has access to the resource.
			err = auth.Verify(account, res, goauth.VerifyNamespace(ns))
			if err == goauth.ErrForbidden && account != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), account.ID)
			} else if err == goauth.ErrForbidden {
				return errors.Unauthorized(req.Service(), "Unauthorized call made to %v:%v", req.Service(), req.Endpoint())
			} else if err != nil {
				return errors.InternalServerError(req.Service(), "Error authorizing request: %v", err)
			}

			// The user is authorised, allow the call
			return h(ctx, req, rsp)
		}
	}
}
