package config

import (
	"errors"

	proto "github.com/micro/go-micro/config/source/mucp/proto"
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

// Watch created by a client RPC request
func Watch(id string) (*watcher, error) {
	mtx.Lock()
	w := &watcher{
		id:   id,
		exit: make(chan bool),
		next: make(chan *proto.WatchResponse),
	}
	watchers[id] = append(watchers[id], w)
	mtx.Unlock()
	return w, nil
}
