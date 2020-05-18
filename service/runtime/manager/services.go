package manager

import (
	"encoding/json"
	"strings"

	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store"
)

// service is the object persisted in the store
type service struct {
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
}

const (
	// servicePrefix is prefixed to the key for service records
	servicePrefix = "service/"
)

// key to write the service to the store under, e.g:
// "service/foo/go.micro.service.bar:latest"
func (s *service) Key() string {
	return servicePrefix + s.Options.Namespace + "/" + s.Service.Name + ":" + s.Service.Version
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
func (m *manager) readServices(namespace string, srv *runtime.Service) ([]*runtime.Service, error) {
	prefix := servicePrefix + namespace + "/"
	if len(srv.Name) > 0 {
		prefix += srv.Name + ":"
	}
	if len(srv.Version) > 0 {
		prefix += srv.Version
	}

	recs, err := m.options.Store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return nil, err
	} else if len(recs) == 0 {
		return nil, runtime.ErrNotFound
	}

	srvs := make([]*runtime.Service, 0, len(recs))
	for _, r := range recs {
		var s *service
		if err := json.Unmarshal(r.Value, &s); err != nil {
			return nil, err
		}
		srvs = append(srvs, s.Service)
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
	keys, err := m.options.Store.List(store.ListPrefix(servicePrefix))
	if err != nil {
		return nil, err
	}

	namespaces := make([]string, 0, len(keys))
	for _, key := range keys {
		// key is formatted 'prefix/namespace/name:version'
		if comps := strings.Split(key, "/"); len(comps) == 3 {
			namespaces = append(namespaces, comps[1])
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
