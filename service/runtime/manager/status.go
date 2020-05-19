package manager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store"
)

// statusPrefix is prefixed to every status key written to the memory store
const statusPrefix = "status:"

// serviceStatus contains the runtime specific information for a service
type serviceStatus struct {
	Status string
	Error  string
}

// statusPollFrequency is the max frequency the manager will check for new statuses in the runtime
var statusPollFrequency = time.Second * 30

// watchStatus calls the managed runtime, gets the serviceStatus for all services listed in the
// store and writes it to the memory store
func (m *manager) watchStatus() {
	ticker := time.NewTicker(statusPollFrequency)

loop:
	for {
		namespaces, err := m.listNamespaces()
		if err != nil {
			logger.Warnf("Error listing namespaces: %v", err)
			continue loop
		}

		for _, ns := range namespaces {
			srvs, err := m.Runtime.Read(runtime.ReadNamespace(ns))
			if err != nil {
				logger.Warnf("Error reading namespace %v: %v", ns, err)
				continue loop
			}

			for _, srv := range srvs {
				if err := m.cacheStatus(ns, srv); err != nil {
					logger.Warnf("Error caching status: %v", err)
					continue loop
				}
			}
		}

		<-ticker.C
	}
}

// cacheStatus writes a services status to the memory store which is then later returned in service
// metadata on Runtime.Read
func (m *manager) cacheStatus(ns string, srv *runtime.Service) error {
	// errors / status is returned from the underlying runtime using srv.Metadata. TODO: Consider
	// changing this so status / error are attributes on runtime.Service.
	if srv.Metadata == nil {
		return fmt.Errorf("Service %v:%v (%v) is missing metadata", srv.Name, srv.Version, ns)
	}

	key := fmt.Sprintf("%v%v:%v:%v", statusPrefix, ns, srv.Name, srv.Version)
	val := &serviceStatus{Status: srv.Metadata["status"], Error: srv.Metadata["error"]}

	bytes, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return m.cache.Write(&store.Record{Key: key, Value: bytes})
}

// listStautuses returns all the statuses for the services in a given namespace with 'name:version'
// as the format used for the keys in the map.
func (m *manager) listStatuses(ns string) (map[string]*serviceStatus, error) {
	recs, err := m.cache.Read(statusPrefix+ns+":", store.ReadPrefix())
	if err != nil {
		return nil, fmt.Errorf("Error listing statuses from the store for namespace %v: %v", ns, err)
	}

	statuses := make(map[string]*serviceStatus, len(recs))

	for _, rec := range recs {
		var status *serviceStatus
		if err := json.Unmarshal(rec.Value, &status); err != nil {
			return nil, err
		}

		// record keys are formatted: 'prefix:namespace:name:version'
		if comps := strings.Split(rec.Key, ":"); len(comps) == 4 {
			statuses[comps[2]+":"+comps[3]] = status
		} else {
			return nil, fmt.Errorf("Invalid key: %v", err)
		}
	}

	return statuses, nil
}
