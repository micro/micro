package cache

import (
	"context"
	"testing"
	"time"

	metadata "github.com/micro/micro/v5/service/context"
)

func TestCache(t *testing.T) {
	ctx := context.TODO()
	req := &testRequest{service: "go.micro.service.foo", method: "Foo.Bar"}

	t.Run("CacheMiss", func(t *testing.T) {
		if _, ok := New().Get(ctx, req); ok {
			t.Errorf("Expected to get no result from Get")
		}
	})

	t.Run("CacheHit", func(t *testing.T) {
		c := New()

		rsp := "theresponse"
		c.Set(ctx, req, rsp, time.Minute)

		if res, ok := c.Get(ctx, req); !ok {
			t.Errorf("Expected a result, got nothing")
		} else if res != rsp {
			t.Errorf("Expected '%v' result, got '%v'", rsp, res)
		}
	})
}

func TestCacheKey(t *testing.T) {
	ctx := context.TODO()

	req1 := &testRequest{service: "go.micro.service.foo", method: "Foo.Bar"}
	req2 := &testRequest{service: "go.micro.service.foo", method: "Foo.Baz"}
	req3 := &testRequest{service: "go.micro.service.foo", method: "Foo.Bar", body: "customquery"}

	t.Run("IdenticalRequests", func(t *testing.T) {
		key1 := key(ctx, req1)
		key2 := key(ctx, req1)
		if key1 != key2 {
			t.Errorf("Expected the keys to match for identical requests and context")
		}
	})

	t.Run("DifferentRequestEndpoints", func(t *testing.T) {
		key1 := key(ctx, req1)
		key2 := key(ctx, req2)

		if key1 == key2 {
			t.Errorf("Expected the keys to differ for different request endpoints")
		}
	})

	t.Run("DifferentRequestBody", func(t *testing.T) {
		key1 := key(ctx, req2)
		key2 := key(ctx, req3)

		if key1 == key2 {
			t.Errorf("Expected the keys to differ for different request bodies")
		}
	})

	t.Run("DifferentMetadata", func(t *testing.T) {
		mdCtx := metadata.Set(context.TODO(), "Micro-Namespace", "bar")
		key1 := key(mdCtx, req1)
		key2 := key(ctx, req1)

		if key1 == key2 {
			t.Errorf("Expected the keys to differ for different metadata")
		}
	})
}
