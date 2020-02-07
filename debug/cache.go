package debug

import (
	"sync"
	"time"

	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/cache"
	"github.com/micro/go-micro/v2/util/log"
)

// Stats is the Debug.Stats handler
type cached struct {
	registry registry.Registry

	sync.RWMutex
	serviceCache []*registry.Service
}

func newCache(done <-chan bool) *cached {
	c := &cached{
		registry: cache.New(*cmd.DefaultOptions().Registry),
	}

	// first scan
	if err := c.scan(); err != nil {
		return nil
	}

	go c.run(done)

	return c
}

func (c *cached) services() []*registry.Service {
	c.RLock()
	defer c.RUnlock()
	return c.serviceCache
}

func (c *cached) run(done <-chan bool) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-done:
			return
		case <-t.C:
			if err := c.scan(); err != nil {
				log.Debug(err)
			}
		}
	}
}

func (c *cached) scan() error {
	services, err := c.registry.ListServices()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	serviceMap := make(map[string]*registry.Service)

	// check each service has nodes
	for _, service := range services {
		if len(service.Nodes) > 0 {
			serviceMap[service.Name+service.Version] = service
			continue
		}

		// get nodes that does not exist
		newServices, err := c.registry.GetService(service.Name)
		if err != nil {
			continue
		}

		// store service by version
		for _, service := range newServices {
			serviceMap[service.Name+service.Version] = service
		}
	}

	// flatten the map
	var serviceList []*registry.Service

	for _, service := range serviceMap {
		serviceList = append(serviceList, service)
	}

	// save the list
	c.Lock()
	c.serviceCache = serviceList
	c.Unlock()
	return nil
}
