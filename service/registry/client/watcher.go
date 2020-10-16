package client

import (
	"time"

	pb "github.com/micro/micro/v3/proto/registry"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/registry/util"
)

type serviceWatcher struct {
	stream pb.Registry_WatchService
	closed chan bool
}

func (s *serviceWatcher) Next() (*registry.Result, error) {
	var i int

	for {
		// check if closed
		select {
		case <-s.closed:
			return nil, registry.ErrWatcherStopped
		default:
		}

		r, err := s.stream.Recv()
		if err != nil {
			return nil, err
		}

		// result is nil
		if r == nil {
			i++

			// only process for 3 attempts if nil
			if i > 3 {
				return nil, registry.ErrWatcherStopped
			}

			// wait a moment
			time.Sleep(time.Second)

			// otherwise continue
			continue
		}

		return &registry.Result{
			Action:  r.Action,
			Service: util.ToService(r.Service),
		}, nil
	}
}

func (s *serviceWatcher) Stop() {
	select {
	case <-s.closed:
		return
	default:
		close(s.closed)
		s.stream.Close()
	}
}

func newWatcher(stream pb.Registry_WatchService) registry.Watcher {
	return &serviceWatcher{
		stream: stream,
		closed: make(chan bool),
	}
}
