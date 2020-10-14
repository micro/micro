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
// Original source: github.com/micro/go-micro/v3/registry/memory/watcher.go

package memory

import (
	"errors"

	"github.com/micro/micro/v3/service/registry"
)

type Watcher struct {
	id   string
	wo   registry.WatchOptions
	res  chan *registry.Result
	exit chan bool
}

func (m *Watcher) Next() (*registry.Result, error) {
	for {
		select {
		case r := <-m.res:
			if r.Service == nil {
				continue
			}

			if len(m.wo.Service) > 0 && m.wo.Service != r.Service.Name {
				continue
			}

			// extract domain from service metadata
			var domain string
			if r.Service.Metadata != nil && len(r.Service.Metadata["domain"]) > 0 {
				domain = r.Service.Metadata["domain"]
			} else {
				domain = registry.DefaultDomain
			}

			// only send the event if watching the wildcard or this specific domain
			if m.wo.Domain == registry.WildcardDomain || m.wo.Domain == domain {
				return r, nil
			}
		case <-m.exit:
			return nil, errors.New("watcher stopped")
		}
	}
}

func (m *Watcher) Stop() {
	select {
	case <-m.exit:
		return
	default:
		close(m.exit)
	}
}
