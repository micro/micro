// Package notifier is a file server notifer
package notifier

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/micro/go-micro/runtime"
	"github.com/micro/go-micro/util/log"
)

type notifier struct {
	service string
	version string
	path    string

	once sync.Once
	sync.Mutex

	fsnotify *fsnotify.Watcher
	notify   chan runtime.Event
	update   chan fsnotify.Event
	exit     chan bool
}

func (n *notifier) run() {
	for {
		select {
		case <-n.exit:
			return
		case <-n.update:
			select {
			case n.notify <- runtime.Event{
				Type:      runtime.Update,
				Timestamp: time.Now(),
				Service:   n.service,
				Version:   fmt.Sprintf("%d", time.Now().Unix()),
			}:
			default:
				// bail out
			}
		case ev := <-n.fsnotify.Events:
			select {
			case n.update <- ev:
			default:
				// bail out
			}
		case <-n.fsnotify.Errors:
			// replace the watcher
			n.fsnotify, _ = fsnotify.NewWatcher()
			n.fsnotify.Add(n.path)
		}
	}
}

func (n *notifier) Notify() (<-chan runtime.Event, error) {
	select {
	case <-n.exit:
		return nil, errors.New("closed")
	default:
	}

	n.once.Do(func() {
		go n.run()
	})

	return n.notify, nil
}

func (n *notifier) Close() error {
	n.Lock()
	defer n.Unlock()
	select {
	case <-n.exit:
		return nil
	default:
		close(n.exit)
		n.fsnotify.Close()
		return nil
	}
}

// New returns a new notifier which watches the source
func New(service, version, source string) runtime.Notifier {
	n := &notifier{
		path:    filepath.Dir(source),
		exit:    make(chan bool),
		notify:  make(chan runtime.Event, 32),
		update:  make(chan fsnotify.Event, 32),
		service: service,
		version: version,
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	w.Add(n.path)
	// set the watcher
	n.fsnotify = w

	return n
}
