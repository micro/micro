package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/micro/v2/internal/config"
)

// NewRuntime returns an initialized Platform struct
func NewRuntime() runtime.Runtime {
	return new(Platform)
}

// Platform implements the runtime interface. It is being
// used temporarily to connect the CLI to the platform API,
// once the micro proxy has been implemented, this package
// will be ddepricated.
type Platform struct{}

// String describes runtime
func (p *Platform) String() string {
	return "platform"
}

// Create registers a service
func (p *Platform) Create(srv *runtime.Service, opts ...runtime.CreateOption) error {
	_, err := call("POST", "platform/CreateService", map[string]interface{}{
		"service": map[string]interface{}{
			"name":     srv.Name,
			"source":   srv.Source,
			"version":  srv.Version,
			"metadata": srv.Metadata,
		},
	})

	if err == nil {
		fmt.Printf("[Platform] Service %v:%v created\n", srv.Name, srv.Version)
	}

	return err
}

// Read returns the service
func (p *Platform) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	options := runtime.ReadOptions{}
	for _, o := range opts {
		o(&options)
	}

	resp, err := call("GET", "platform/ReadService", map[string]interface{}{
		"service": map[string]interface{}{
			"type":    options.Type,
			"name":    options.Service,
			"version": options.Version,
		},
	})
	if err != nil {
		return []*runtime.Service{}, err
	}

	// unmarashal the response
	var data struct {
		Services []*runtime.Service
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	// return the services
	return data.Services, nil
}

// Delete a service
func (p *Platform) Delete(srv *runtime.Service) error {
	_, err := call("POST", "platform/DeleteService", map[string]interface{}{
		"service": map[string]interface{}{
			"name":    srv.Name,
			"version": srv.Version,
		},
	})

	if err == nil {
		fmt.Printf("[Platform] Service %v:%v deleted\n", srv.Name, srv.Version)
	}

	return err
}

// List the managed services
func (p *Platform) List() ([]*runtime.Service, error) {
	resp, err := call("GET", "platform/ListServices", map[string]interface{}{})
	if err != nil {
		return []*runtime.Service{}, err
	}

	// unmarashal the response
	var data struct {
		Services []*runtime.Service
	}
	if err := json.Unmarshal(resp, &data); err != nil {
		return nil, err
	}

	// return the services
	return data.Services, nil
}

// call performs a HTTP request
func call(method, path string, data map[string]interface{}) ([]byte, error) {
	// Construct the JSON for the HTTP request:
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(jsonStr)

	// Get the auth token
	token, _ := config.Get("token")

	// Construct the request
	url := "https://api.micro.mu/" + path
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))

	// Execute the request
	resp, err := http.DefaultClient.Do(req)

	// Check for errors
	if resp.StatusCode == 401 {
		fmt.Println("[Platform] You must be logged in to perform this request")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[Platform] API Error: %v", resp.StatusCode)
	}

	// Return the result and error
	return ioutil.ReadAll(resp.Body)
}

// Update the service in place
func (p *Platform) Update(srv *runtime.Service) error {
	if err := p.Delete(srv); err != nil {
		return err
	}

	if err := p.Create(srv); err != nil {
		return err
	}

	return nil
}

// The following functions aren't required / supported
// by this implementation but have been added to satisfy
// the runtime interface

// Init initializes runtime
func (p *Platform) Init(...runtime.Option) error {
	return nil
}

// Start starts the runtime
func (p *Platform) Start() error {
	return nil
}

// Stop shuts down the runtime
func (p *Platform) Stop() error {
	return nil
}
