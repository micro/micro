package wrapper

import (
	"context"
	"reflect"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/client/cache"
)

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
