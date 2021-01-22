package manager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/store"
)

// service is the object persisted in the store
type service struct {
	Service   *runtime.Service       `json:"service"`
	Options   *runtime.CreateOptions `json:"options"`
	Status    runtime.ServiceStatus  `json:"status"`
	UpdatedAt time.Time              `json:"last_updated"`
	Error     string                 `json:"error"`
}

const (
	// servicePrefix is prefixed to the key for service records
	servicePrefix = "service:"
)

// key to write the service to the store under, e.g:
// "service/foo/bar:latest"
func (s *service) Key() string {
	return servicePrefix + s.Options.Namespace + ":" + s.Service.Name + ":" + s.Service.Version
}

// writeService to the store
func (m *manager) writeService(srv *service) error {
	bytes, err := json.Marshal(srv)
	if err != nil {
		return err
	}

	return store.Write(&store.Record{Key: srv.Key(), Value: bytes})
}

// deleteService from the store
func (m *manager) deleteService(srv *service) error {
	return store.Delete(srv.Key())
}

// readServices returns all the services in a given namespace. If a service name and
// version are provided it will filter using these as well
func (m *manager) readServices(namespace string, srv *runtime.Service) ([]*service, error) {
	prefix := servicePrefix + namespace + ":"
	if len(srv.Name) > 0 {
		prefix += srv.Name + ":"
	}
	if len(srv.Name) > 0 && len(srv.Version) > 0 {
		prefix += srv.Version
	}

	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	} else if len(recs) == 0 {
		return make([]*service, 0), nil
	}

	srvs := make([]*service, 0, len(recs))
	for _, r := range recs {
		var s *service
		if err := json.Unmarshal(r.Value, &s); err != nil {
			return nil, err
		}
		srvs = append(srvs, s)
	}

	return srvs, nil
}

// listNamespaces of the services in the store. todo: remove this and have the watchServices func
// query the store directly
func (m *manager) listNamespaces() ([]string, error) {
	recs, err := store.Read(servicePrefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return []string{namespace.DefaultNamespace}, nil
	}

	namespaces := make([]string, 0, len(recs))
	for _, rec := range recs {
		// key is formatted 'prefix:namespace:name:version'
		if comps := strings.Split(rec.Key, ":"); len(comps) == 4 {
			namespaces = append(namespaces, comps[1])
		} else {
			return nil, fmt.Errorf("Invalid key: %v", rec.Key)
		}
	}

	return unique(namespaces), nil
}

// unique is a helper method to filter a slice of strings
// down to unique entries
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
