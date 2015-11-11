package cli

import (
	"github.com/myodc/go-micro/registry"
)

type sortedServices struct {
	services []*registry.Service
}

func (s sortedServices) Len() int {
	return len(s.services)
}

func (s sortedServices) Less(i, j int) bool {
	return s.services[i].Name < s.services[j].Name
}

func (s sortedServices) Swap(i, j int) {
	s.services[i], s.services[j] = s.services[j], s.services[i]
}
