package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	cache "github.com/patrickmn/go-cache"
)

// New returns an initialised cache.
func New() *Cache {
	return &Cache{
		cache: cache.New(cache.NoExpiration, 30*time.Second),
	}
}

// Cache for responses
type Cache struct {
	cache *cache.Cache
}

type Options struct {
	// Expiry sets the cache expiry
	Expiry time.Duration
}

// used to store the options in context
type optionsKey struct{}

// Get a response from the cache
func (c *Cache) Get(ctx context.Context, req client.Request) (interface{}, bool) {
	return c.cache.Get(key(ctx, req))
}

// Set a response in the cache
func (c *Cache) Set(ctx context.Context, req client.Request, rsp interface{}, expiry time.Duration) {
	c.cache.Set(key(ctx, req), rsp, expiry)
}

// List the key value pairs in the cache
func (c *Cache) List() map[string]string {
	items := c.cache.Items()

	rsp := make(map[string]string, len(items))
	for k, v := range items {
		bytes, _ := json.Marshal(v.Object)
		rsp[k] = string(bytes)
	}

	return rsp
}

// key returns a hash for the context and request
func key(ctx context.Context, req client.Request) string {
	ns, _ := metadata.Get(ctx, "Micro-Namespace")

	bytes, _ := json.Marshal(map[string]interface{}{
		"namespace": ns,
		"request": map[string]interface{}{
			"service":  req.Service(),
			"endpoint": req.Endpoint(),
			"method":   req.Method(),
			"body":     req.Body(),
		},
	})

	h := fnv.New64()
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func SetOptions(ctx context.Context, opts *Options) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, optionsKey{}, opts)
}

func GetOptions(ctx context.Context) (*Options, bool) {
	if ctx == nil {
		return nil, false
	}
	opts, ok := ctx.Value(optionsKey{}).(*Options)
	return opts, ok
}

func CallOption(opts *Options) client.CallOption {
	return func(o *client.CallOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = SetOptions(o.Context, opts)
	}
}

func CallExpiry(t time.Duration) client.CallOption {
	return CallOption(&Options{
		Expiry: t,
	})
}
