package blocklist

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const (
	prefix = "api-blocklist"
)

type BlockList struct {
	redisClient *redis.Client
}

func New(addr, user, pass string, tlsConf *tls.Config) *BlockList {
	return &BlockList{
		redisClient: redis.NewClient(&redis.Options{
			Addr:      addr,
			Username:  user,
			Password:  pass,
			TLSConfig: tlsConf,
		}),
	}
}

func (b *BlockList) IsBlocked(ctx context.Context, id, namespace string) (bool, error) {
	err := b.redisClient.Get(ctx, fmt.Sprintf("%s:%s:%s", prefix, namespace, id)).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *BlockList) Add(ctx context.Context, id, namespace string) error {
	return b.redisClient.Set(ctx, fmt.Sprintf("%s:%s:%s", prefix, namespace, id), "true", 0).Err()
}

func (b *BlockList) Remove(ctx context.Context, id, namespace string) error {
	return b.redisClient.Del(ctx, fmt.Sprintf("%s:%s:%s", prefix, namespace, id)).Err()
}

func (b *BlockList) List(ctx context.Context) ([]string, error) {
	res, err := b.redisClient.Keys(ctx, fmt.Sprintf("%s:*", prefix)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return res, nil
}
