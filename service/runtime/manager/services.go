package manager

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/internal/namespace"
)

// service is the object persisted in the store
type service struct {
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
}

const (
	// servicePrefix is prefixed to the key for service records
	servicePrefix = "service:"
)

// key to write the service to the store under, e.g:
// "service/foo/go.micro.service.bar:latest"
func (s *service) Key() string {
	return servicePrefix + s.Options.Namespace + ":" + s.Service.Name + ":" + s.Service.Version
}

// createService writes the service to the store
func (m *manager) createService(srv *runtime.Service, opts *runtime.CreateOptions) error {
	s := &service{srv, opts}

	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return m.options.Store.Write(&store.Record{Key: s.Key(), Value: bytes})
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

	recs, err := m.options.Store.Read(prefix, store.ReadPrefix())
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

// deleteSevice from the store
func (m *manager) deleteService(namespace string, srv *runtime.Service) error {
	obj := &service{srv, &runtime.CreateOptions{Namespace: namespace}}
	return m.options.Store.Delete(obj.Key())
}

// listNamespaces of the services in the store
func (m *manager) listNamespaces() ([]string, error) {
	recs, err := m.options.Store.Read(servicePrefix, store.ReadPrefix())
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
