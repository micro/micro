package config

import (
	"context"
	"time"

	"github.com/micro/go-micro/client"
	cr "github.com/micro/go-micro/config/reader"
	"github.com/micro/go-micro/config/source"
	proto "github.com/micro/micro/config/proto/config"
)

// Used as a subscriber between config services for events
func Watcher(ctx context.Context, ch *proto.WatchResponse) error {
	mtx.RLock()
	for _, sub := range watchers[ch.Id] {
		select {
		case sub.next <- ch:
		case <-time.After(time.Millisecond * 100):
		}
	}
	mtx.RUnlock()
	return nil
}

func Merge(ch ...*source.ChangeSet) (*source.ChangeSet, error) {
	return reader.Merge(ch...)
}

func Values(ch *source.ChangeSet) (cr.Values, error) {
	return reader.Values(ch)
}

// Publish a change
func Publish(ctx context.Context, ch *proto.WatchResponse) error {
	req := client.NewMessage(WatchTopic, ch)
	return client.Publish(ctx, req)
}
