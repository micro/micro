// Copyright 2020 Asim Aslam
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-plugins/v3/registry/etcd/watcher.go

package etcd

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.etcd.io/etcd/clientv3"
	"github.com/micro/micro/v3/service/registry"
)

type etcdWatcher struct {
	w       clientv3.WatchChan
	client  *clientv3.Client
	timeout time.Duration

	mtx    sync.Mutex
	stop   chan bool
	cancel func()
}

func newEtcdWatcher(c *clientv3.Client, timeout time.Duration, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}
	if len(wo.Domain) == 0 {
		wo.Domain = defaultDomain
	}

	watchPath := prefix
	if wo.Domain == registry.WildcardDomain {
		if len(wo.Service) > 0 {
			return nil, errors.New("Cannot watch a service across domains")
		}
		watchPath = prefix
	} else if len(wo.Service) > 0 {
		watchPath = servicePath(wo.Domain, wo.Service) + "/"
	}

	ctx, cancel := context.WithCancel(context.Background())
	w := c.Watch(ctx, watchPath, clientv3.WithPrefix(), clientv3.WithPrevKV())
	stop := make(chan bool, 1)

	return &etcdWatcher{
		cancel:  cancel,
		stop:    stop,
		w:       w,
		client:  c,
		timeout: timeout,
	}, nil
}

func (ew *etcdWatcher) Next() (*registry.Result, error) {
	for wresp := range ew.w {
		if wresp.Err() != nil {
			return nil, wresp.Err()
		}
		if wresp.Canceled {
			return nil, errors.New("could not get next")
		}
		for _, ev := range wresp.Events {
			service := decode(ev.Kv.Value)
			var action string

			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"

				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}

			if service == nil {
				continue
			}
			return &registry.Result{
				Action:  action,
				Service: service,
			}, nil
		}
	}
	return nil, errors.New("could not get next")
}

func (ew *etcdWatcher) Stop() {
	ew.mtx.Lock()
	defer ew.mtx.Unlock()

	select {
	case <-ew.stop:
		return
	default:
		close(ew.stop)
		ew.cancel()
		ew.client.Close()
	}
}
