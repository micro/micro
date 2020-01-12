package config

import (
	"errors"

	proto "github.com/micro/micro/config/proto/config"
)

type watcher struct {
	id   string
	exit chan bool
	next chan *proto.WatchResponse
}

func (w *watcher) Next() (*proto.WatchResponse, error) {
	select {
	case c := <-w.next:
		return c, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
		return errors.New("already stopped")
	default:
		close(w.exit)
	}

	mtx.Lock()
	var wslice []*watcher

	for _, watch := range watchers[w.id] {
		if watch != w {
			wslice = append(wslice, watch)
		}
	}

	watchers[w.id] = wslice
	mtx.Unlock()

	return nil
}
