package server

import (
	"context"
	"net/http"
	"strings"

	"micro.dev/v4/service/api"
	"micro.dev/v4/service/api/router/registry"
	"micro.dev/v4/service/auth"
	metadata "micro.dev/v4/service/context"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/server"
	inauth "micro.dev/v4/util/auth"
	"micro.dev/v4/util/ctx"
	"micro.dev/v4/util/namespace"
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
				if !strings.HasPrefix(header, inauth.BearerScheme) {
					return errors.Unauthorized(req.Service(), "invalid authorization header. expected Bearer schema")
				}

				// Strip the bearer scheme prefix
				token = strings.TrimPrefix(header, inauth.BearerScheme)
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

			// construct the resource
			res := &auth.Resource{
				Type:     "service",
				Name:     req.Service(),
				Endpoint: req.Endpoint(),
			}

			// Verify the caller has access to the resource.
			err := auth.Verify(account, res, auth.VerifyNamespace(ns))
			if err == auth.ErrForbidden && account != nil {
				return errors.Forbidden(req.Service(), "Forbidden call made to %v:%v by %v", req.Service(), req.Endpoint(), account.ID)
			} else if err == auth.ErrForbidden {
				return errors.Unauthorized(req.Service(), "Unauthorized call made to %v:%v", req.Service(), req.Endpoint())
			} else if err != nil {
				return errors.InternalServerError("proxy", "Error authorizing request: %v", err)
			}

			// The user is authorised, allow the call
			return h(ctx, req, rsp)
		}
	}
}

func authWrapper() api.Wrapper {
	resolver := registry.NewResolver()

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Determine the name of the service being requested
			endpoint := resolver.Resolve(req)
			ctrx := context.WithValue(req.Context(), registry.Endpoint{}, endpoint)
			*req = *req.Clone(ctrx)

			// Set the metadata so we can access it in micro api / web
			req = req.WithContext(ctx.FromRequest(req))

			// Extract the token from the request
			var token string
			if header := req.Header.Get("Authorization"); len(header) > 0 {
				// Extract the auth token from the request
				if strings.HasPrefix(header, inauth.BearerScheme) {
					token = header[len(inauth.BearerScheme):]
				}
			} else {
				// Get the token out the cookies if not provided in headers
				if c, err := req.Cookie("micro-token"); err == nil && c != nil {
					token = strings.TrimPrefix(c.Value, inauth.TokenCookieName+"=")
					req.Header.Set("Authorization", inauth.BearerScheme+token)
				}
			}

			// Get the account using the token, some are unauthenticated, so the lack of an
			// account doesn't necessarily mean a forbidden request
			acc, err := auth.Inspect(token)
			if err == nil {
				// inject into the context
				ctx := auth.ContextWithAccount(req.Context(), acc)
				*req = *req.Clone(ctx)
			}

			// Determine the namespace and set it in the header. If the user passed auth creds
			// on the request, use the namespace that issued the account, otherwise check for
			// the domain of the resolved endpoint.
			ns := req.Header.Get(namespace.NamespaceKey)
			if len(ns) == 0 && acc != nil {
				ns = acc.Issuer
				req.Header.Set(namespace.NamespaceKey, ns)
			} else if len(ns) == 0 {
				ns = endpoint.Domain
				req.Header.Set(namespace.NamespaceKey, ns)
			}

			// Ensure accounts only issued by the namespace are valid.
			if acc != nil && acc.Issuer != ns {
				acc = nil
			}

			// construct the resource name, e.g. home => foo.api.home
			resName := endpoint.Name
			resEndpoint := endpoint.Method

			// Options to use when verifying the request
			verifyOpts := []auth.VerifyOption{
				auth.VerifyContext(req.Context()),
				auth.VerifyNamespace(ns),
			}

			logger.Debugf("Resolving %v %v", resName, resEndpoint)

			// Perform the verification check to see if the account has access to
			// the resource they're requesting
			res := &auth.Resource{Type: "service", Name: resName, Endpoint: resEndpoint}
			if err := auth.Verify(acc, res, verifyOpts...); err == nil {
				// The account has the necessary permissions to access the resource
				h.ServeHTTP(w, req)
				return
			} else if err != auth.ErrForbidden {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// The account is set, but they don't have enough permissions, hence
			// we return a forbidden error.
			if acc != nil {
				http.Error(w, "Forbidden request", http.StatusForbidden)
				return
			}

			http.Error(w, "unauthorized request", http.StatusUnauthorized)
			return
		})
	}
}
