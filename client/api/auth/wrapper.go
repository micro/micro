package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro/go-micro/v2/api/resolver"
	"github.com/micro/go-micro/v2/api/resolver/subdomain"
	"github.com/micro/go-micro/v2/api/server"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/util/ctx"
	inauth "github.com/micro/micro/v2/internal/auth"
	"github.com/micro/micro/v2/internal/namespace"
)

// Wrapper wraps a handler and authenticates requests
func Wrapper(r resolver.Resolver, prefix string) server.Wrapper {
	return func(h http.Handler) http.Handler {
		return authWrapper{
			handler:       h,
			resolver:      r,
			servicePrefix: prefix,
			auth:          auth.DefaultAuth,
		}
	}
}

type authWrapper struct {
	handler       http.Handler
	auth          auth.Auth
	resolver      resolver.Resolver
	servicePrefix string
}

func (a authWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Determine the name of the service being requested
	endpoint, err := a.resolver.Resolve(req)
	if err == resolver.ErrInvalidPath || err == resolver.ErrNotFound {
		// a file not served by the resolver has been requested (e.g. favicon.ico)
		endpoint = &resolver.Endpoint{Path: req.URL.Path}
	} else if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
		return
	} else {
		// set the endpoint in the context so it can be used to resolve
		// the request later
		ctx := context.WithValue(req.Context(), resolver.Endpoint{}, endpoint)
		*req = *req.Clone(ctx)
	}

	// If an error occured looking up the route, the domain isn't returned. TODO: Find a better way
	// of resolving network for non-standard requests, e.g. "/rpc".
	if r, ok := a.resolver.(*subdomain.Resolver); ok && len(endpoint.Domain) == 0 {
		endpoint.Domain = r.Domain(req)
	}

	// Determine the namespace and set it in the header
	ns := req.Header.Get(namespace.NamespaceKey)
	if len(ns) == 0 {
		ns = endpoint.Domain
		req.Header.Set(namespace.NamespaceKey, ns)
	}

	// Set the metadata so we can access it in micro api / web
	req = req.WithContext(ctx.FromRequest(req))

	// Extract the token from the request
	var token string
	if header := req.Header.Get("Authorization"); len(header) > 0 {
		// Extract the auth token from the request
		if strings.HasPrefix(header, auth.BearerScheme) {
			token = header[len(auth.BearerScheme):]
		}
	} else {
		// Get the token out the cookies if not provided in headers
		if c, err := req.Cookie("micro-token"); err == nil && c != nil {
			token = strings.TrimPrefix(c.Value, inauth.TokenCookieName+"=")
			req.Header.Set("Authorization", auth.BearerScheme+token)
		}
	}

	// Get the account using the token, some are unauthenticated, so the lack of an
	// account doesn't necesserially mean a forbidden request
	acc, _ := a.auth.Inspect(token)

	// Ensure accounts only issued by the namesace are valid
	if acc != nil && acc.Issuer != ns {
		acc = nil
	}

	// construct the resource name, e.g. home => go.micro.web.home
	resName := a.servicePrefix + "." + endpoint.Name

	// determine the resource path. there is an inconsistency in how resolvers
	// use method, some use it as Users.ReadUser (the rpc method), and others
	// use it as the HTTP method, e.g GET. TODO: Refactor this to make it consistent.
	resEndpoint := endpoint.Path
	if len(endpoint.Path) == 0 {
		resEndpoint = endpoint.Method
	}

	// Options to use when verifying the request
	verifyOpts := []auth.VerifyOption{
		auth.VerifyContext(req.Context()),
		auth.VerifyNamespace(ns),
	}

	// Perform the verification check to see if the account has access to
	// the resource they're requesting
	res := &auth.Resource{Type: "service", Name: resName, Endpoint: resEndpoint}
	if err := a.auth.Verify(acc, res, verifyOpts...); err == nil {
		// The account has the necessary permissions to access the resource
		a.handler.ServeHTTP(w, req)
		return
	}

	// The account is set, but they don't have enough permissions, hence
	// we return a forbidden error.
	if acc != nil {
		http.Error(w, "Forbidden request", 403)
		return
	}

	// If there is no auth login url set, 401
	loginURL := a.auth.Options().LoginURL
	if loginURL == "" {
		http.Error(w, "unauthorized request", 401)
		return
	}

	// Redirect to the login path
	params := url.Values{"redirect_to": {req.URL.String()}}
	loginWithRedirect := fmt.Sprintf("%v?%v", loginURL, params.Encode())
	http.Redirect(w, req, loginWithRedirect, http.StatusTemporaryRedirect)
}
