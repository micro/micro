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
// Original source: github.com/micro/go-micro/v3/debug/log/kubernetes/stream.go

package kubernetes

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/micro/micro/v3/internal/debug/log"
)

func write(l log.Record) error {
	m, err := json.Marshal(l)
	if err == nil {
		_, err := fmt.Fprintf(os.Stderr, "%s", m)
		return err
	}
	return err
}

type kubeStream struct {
	// the k8s log stream
	stream chan log.Record
	sync.Mutex
	// the stop chan
	stop chan bool
}

func (k *kubeStream) Chan() <-chan log.Record {
	return k.stream
}

func (k *kubeStream) Stop() error {
	k.Lock()
	defer k.Unlock()
	select {
	case <-k.stop:
		return nil
	default:
		close(k.stop)
		close(k.stream)
	}
	return nil
}
