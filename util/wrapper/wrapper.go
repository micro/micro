package wrapper

import (
	"context"
	"encoding/base64"
	"reflect"
	"strings"

	"micro.dev/v4/service/auth"
	"micro.dev/v4/service/client"
	metadata "micro.dev/v4/service/context"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/server"
	inauth "micro.dev/v4/util/auth"
	"micro.dev/v4/util/cache"
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
				switch {
				case strings.HasPrefix(header, inauth.BearerScheme):
					// Strip the bearer scheme prefix
					token = strings.TrimPrefix(header, inauth.BearerScheme)
				case strings.HasPrefix(header, "Basic "):
					b, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
					if err != nil {
						return errors.Unauthorized(req.Service(), "invalid authorization header. Incorrect format")
					}
					parts := strings.SplitN(string(b), ":", 2)
					if len(parts) != 2 {
						return errors.Unauthorized(req.Service(), "invalid authorization header. Incorrect format")
					}

					token = parts[1]
				default:
					return errors.Unauthorized(req.Service(), "invalid authorization header. Expected Bearer or Basic schema")
				}
			}

			// Determine the namespace
			ns := auth.DefaultAuth.Options().Issuer

			var acc *auth.Account
			if a, err := auth.Inspect(token); err == nil {
				ctx = auth.ContextWithAccount(ctx, a)
				acc = a
			}

			// construct the resource
			res := &auth.Resource{
				Type:     "service",
				Name:     req.Service(),
				Endpoint: req.Endpoint(),
			}

			// Verify the caller has access to the resource.
			err := auth.Verify(acc, res, auth.VerifyNamespace(ns))
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

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	logger.Debugf("Calling service %s endpoint %s", req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp, opts...)
}

func (l *logWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	logger.Debugf("Streaming service %s endpoint %s", req.Service(), req.Endpoint())
	return l.Client.Stream(ctx, req, opts...)
}

func LogClient(c client.Client) client.Client {
	return &logWrapper{c}
}

func LogHandler() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			logger.Debugf("Serving request for service %s endpoint %s", req.Service(), req.Endpoint())
			return h(ctx, req, rsp)
		}
	}
}

type cacheWrapper struct {
	Cache *cache.Cache
	client.Client
}

// Call executes the request. If the CacheExpiry option was set, the response will be cached using
// a hash of the metadata and request as the key.
func (c *cacheWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// parse the options
	var options client.CallOptions
	for _, o := range opts {
		o(&options)
	}

	// if the client doesn't have a cacbe setup don't continue
	if c.Cache == nil {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	cacheOpts, ok := cache.GetOptions(options.Context)
	if !ok {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// if the cache expiry is not set, execute the call without the cache
	if cacheOpts.Expiry == 0 || rsp == nil {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	// check to see if there is a response cached, if there is assign it
	if r, ok := c.Cache.Get(ctx, req); ok {
		val := reflect.ValueOf(rsp).Elem()
		val.Set(reflect.ValueOf(r).Elem())
		return nil
	}

	// don't cache the result if there was an error
	if err := c.Client.Call(ctx, req, rsp, opts...); err != nil {
		return err
	}

	// set the result in the cache
	c.Cache.Set(ctx, req, rsp, cacheOpts.Expiry)
	return nil
}

// CacheClient wraps requests with the cache wrapper
func CacheClient(c client.Client) client.Client {
	return &cacheWrapper{
		Cache:  cache.New(),
		Client: c,
	}
}
