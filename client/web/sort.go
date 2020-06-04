package web

import (
	"github.com/micro/go-micro/v2/registry"
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
