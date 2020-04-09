package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/cli/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/util/jitter"
	muProfile "github.com/micro/micro/v2/runtime/profile"
)

type manager struct {
	// Our own id
	Id string
	// the internal runtime aka local, kubernetes
	Runtime runtime.Runtime
	// the storage for what should be running
	Store store.Store
	// a runtime Profile to set for the service
	Profile []string

	sync.RWMutex
	// internal cache of services
	services map[string]*runtimeService
	// internal cache of events
	events map[string]*series

	// process state
	running bool
	exit    chan bool

	// used to propagate events
	eventChan chan *event
}

// stored in store
type runtimeService struct {
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:"options"`
	Status  string                 `json:"status"`
	Error   error                  `json:"error"`
}

// runtime event is single event in the store
type runtimeEvent struct {
	Id string `json:"id"`
	// the event
	Event *event `json:event`
	// the last update
	Timestamp time.Time `json:"timestamp"`
}

// event is a unique event generated when create/update/delete is called
type event struct {
	Type    string                 `json:"type"`
	Service *runtime.Service       `json:"service"`
	Options *runtime.CreateOptions `json:options"`
}

// a timeseries of events
type series struct {
	// the timeseries id
	Id string `json:"id"`
	// the series of events
	Events []*runtimeEvent `json:"events"`
}

var (
	// status ticks for updating local service stauts
	statusTick = time.Second * 10
	// TODO: if events are racy lower updateTick
	// the time at which we check events
	eventTick = time.Minute
	// the time at which we read all records
	updateTick = time.Minute * 5
)

func copyService(s *runtimeService) *runtime.Service {
	cp := new(runtime.Service)
	cp.Name = s.Service.Name
	cp.Version = s.Service.Version
	cp.Source = s.Service.Source
	cp.Metadata = make(map[string]string)
	for k, v := range s.Service.Metadata {
		cp.Metadata[k] = v
	}
	cp.Metadata["status"] = s.Status
	if s.Error != nil {
		cp.Metadata["error"] = s.Error.Error()
	}
	return cp
}

func key(s *runtime.Service) string {
	return s.Name + ":" + s.Version
}

func eventKey(id string) string {
	// hour : id
	return fmt.Sprintf("runtime:events:%d:%s", time.Now().Truncate(time.Hour).Unix(), id)
}

// read all the events from the store
func (m *manager) readEvents() (map[string]*series, error) {
	// read back all events for past hour
	records, err := m.Store.Read(eventKey(""), store.ReadPrefix())
	if err != nil {
		return nil, err
	}

	// unmarshal our events
	events := make(map[string]*series)

	for _, record := range records {
		var event *series

		// dont care about the error
		if err := json.Unmarshal(record.Value, &event); err != nil {
			continue
		}

		events[event.Id] = event
	}

	return events, nil
}

// save events in the store
func (m *manager) saveEvents(oldEvents *series, newEvents []*event) error {
	// new event series
	events := &series{Id: m.Id}

	// apply old  events
	if oldEvents != nil {
		events.Events = oldEvents.Events
	}

	// append new events
	for _, e := range newEvents {
		events.Events = append(events.Events, &runtimeEvent{
			Id:        uuid.New().String(),
			Event:     e,
			Timestamp: time.Now(),
		})
	}

	// marshal the content
	b, err := json.Marshal(events)
	if err != nil {
		return err
	}

	// save the record and expire after an hour
	return m.Store.Write(&store.Record{
		Key:    eventKey(m.Id),
		Value:  b,
		Expiry: time.Minute * 10,
	})
}

// sendEvents to be saved and processed
func (m *manager) sendEvents(events ...*event) {
	for _, ev := range events {
		m.eventChan <- ev
	}
}

// processEvent will execute an event immediately
func (m *manager) processEvent(ev *event) {
	// events to be process immediately
	var err error

	delete(ev.Service.Metadata, "status")
	delete(ev.Service.Metadata, "error")
	switch ev.Type {
	case "delete":
		log.Infof("Procesing deletion event %s", key(ev.Service))
		err = m.Runtime.Delete(ev.Service)
	case "update":
		log.Infof("Processing update event %s", key(ev.Service))
		err = m.Runtime.Update(ev.Service)
	case "create":
		// generate the runtime environment
		env := m.runtimeEnv(ev.Options)

		opts := []runtime.CreateOption{
			runtime.WithCommand(ev.Options.Command...),
			runtime.WithArgs(ev.Options.Args...),
			runtime.WithEnv(env),
			runtime.CreateType(ev.Options.Type),
			runtime.CreateImage(ev.Options.Image),
		}

		log.Infof("Processing create event %s", key(ev.Service))
		err = m.Runtime.Create(ev.Service, opts...)
	}

	if err != nil {
		log.Errorf("Erroring executing event %s for %s: %v", ev.Type, ev.Service.Name, err)

		// save the error
		// hacking, its a pointer
		m.Lock()
		v, ok := m.services[key(ev.Service)]
		if ok {
			v.Status = "error"
			v.Error = err
		}
		m.Unlock()
	}
}

// processEvents will read events and apply them
// - pull events from the store
// - grab past local events cached
// - compare and apply
func (m *manager) processEvents(newEvents []*event) error {
	// TODO: retry event processing

	// get the existing events
	events, err := m.readEvents()
	if err != nil {
		log.Errorf("Failed to list events from store: %v", err)
		return err
	}

	// save our events now for other to process
	if err := m.saveEvents(events[m.Id], newEvents); err != nil {
		log.Errorf("Failed to save events in store: %v", err)
		return err
	}

	// lock the whole op
	m.Lock()
	defer m.Unlock()

	// a set of events we're going to play
	playEvents := make(map[string]*series)

	// get existing events
	pastEvents := m.events

	// compare to new events
	for _, s := range events {
		// skip our own events
		if s.Id == m.Id {
			continue
		}

		// do we have an existing event series
		past, ok := pastEvents[s.Id]
		if !ok || past == nil || len(past.Events) == 0 {
			// if not ok we need to apply all
			playEvents[s.Id] = s
			continue
		}

		// new set of events to play
		var play []*runtimeEvent

		// last event
		last := past.Events[len(past.Events)-1]

		// range the events from this guy
		for _, ev := range s.Events {
			// skip anything older than the previous past
			if ev.Timestamp.Sub(last.Timestamp) <= time.Duration(0) {
				continue
			}

			// skip that old event
			if ev.Id == last.Id {
				continue
			}

			// apply all other evenst
			play = append(play, ev)
		}

		// save events to play from series
		playEvents[s.Id] = &series{Id: s.Id, Events: play}
	}

	// play the events in the series
	for _, series := range playEvents {
		// play all the events
		for _, ev := range series.Events {
			go m.processEvent(ev.Event)
		}
	}

	// save the played events
	m.events = playEvents

	return nil
}

func (m *manager) updateStatus() error {
	services, err := m.Runtime.List()
	if err != nil {
		log.Errorf("Failed to list runtime services: %v", err)
		return err
	}

	m.Lock()
	defer m.Unlock()

	running := make(map[string]*runtime.Service)

	// update running status
	for _, service := range services {
		k := key(service)
		// create running map
		running[k] = service
	}

	// delete from local cache
	for k, v := range m.services {
		srv, ok := running[k]
		if !ok {
			delete(m.services, k)
			continue
		}

		// update the service
		v.Service = srv
		m.services[k] = v
	}

	return nil
}

// full refresh of the service list
func (m *manager) processServices() error {
	// list the keys from store
	records, err := m.Store.Read("", store.ReadPrefix())
	if err != nil {
		log.Errorf("Failed to list records from store: %v", err)
		return err
	}

	// list whats already runnning
	// TODO: change to read service: prefix
	services, err := m.Runtime.List()
	if err != nil {
		log.Errorf("Failed to list runtime services: %v", err)
		return err
	}

	// generate service map of running things
	running := make(map[string]*runtime.Service)

	for _, service := range services {
		k := key(service)
		running[k] = service
	}

	// create a map of services that should actually run
	shouldRun := make(map[string]*runtimeService)

	// iterate through and see what we need to run
	for _, record := range records {
		// decode the record
		var rs *runtimeService
		if err := json.Unmarshal(record.Value, &rs); err != nil {
			continue
		}

		// skip event records
		if strings.HasPrefix(record.Key, "runtime:events:") {
			continue
		}

		// things to run
		shouldRun[record.Key] = rs

		// check if its already running
		if v, ok := running[record.Key]; ok {
			// TODO: have actual runtime status
			rs.Status = v.Metadata["status"]
			if e := v.Metadata["error"]; len(e) > 0 {
				rs.Error = errors.New(e)
			}
			// replace service entry
			rs.Service = v
			continue
		}

		// generate the runtime environment
		env := m.runtimeEnv(rs.Options)

		// create a new set of options to use
		opts := []runtime.CreateOption{
			runtime.WithCommand(rs.Options.Command...),
			runtime.WithArgs(rs.Options.Args...),
			runtime.WithEnv(env),
			runtime.CreateType(rs.Options.Type),
			runtime.CreateImage(rs.Options.Image),
		}

		// set the status to starting
		rs.Status = "started"
		// service does not exist so start it
		if err := m.Runtime.Create(rs.Service, opts...); err != nil {
			if err != runtime.ErrAlreadyExists {
				log.Errorf("Error running %s: %v", key(rs.Service), err)

				// save the error
				rs.Status = "error"
				rs.Error = err
			}
		}
	}

	// check what we need to stop from the running list
	for _, service := range services {
		k := key(service)

		// check if it should be running
		if _, ok := shouldRun[k]; ok {
			continue
		}

		log.Infof("Stopping %s", k)

		// should not be running
		m.Runtime.Delete(service)
	}

	// save the current list of running things
	m.Lock()
	m.services = shouldRun
	m.Unlock()

	return nil
}

func (m *manager) runtimeEnv(options *runtime.CreateOptions) []string {
	setEnv := func(p []string, env map[string]string) {
		for _, v := range p {
			parts := strings.Split(v, "=")
			if len(parts) <= 1 {
				continue
			}
			env[parts[0]] = strings.Join(parts[1:], "=")
		}
	}

	// overwrite any values
	env := map[string]string{}

	// set the env vars provided
	setEnv(options.Env, env)

	// override with vars from the Profile
	setEnv(m.Profile, env)

	// create a new env
	var vars []string
	for k, v := range env {
		vars = append(vars, k+"="+v)
	}

	// setup the runtime env
	return vars
}

// TODO: watch events rather than poll
func (m *manager) run() {
	// when we publish, process and apply events
	t1 := time.NewTicker(eventTick)
	defer t1.Stop()

	// when we do a full refresh of records
	t2 := time.NewTicker(updateTick)
	defer t2.Stop()

	// status tick for updating status
	t3 := time.NewTicker(statusTick)
	defer t3.Stop()

	// save the existing set of events since on startup
	// we dont want to apply deltas
	m.Lock()
	m.events, _ = m.readEvents()
	m.Unlock()

	// save existing services to run
	m.processServices()

	// batch events to save for other regions to read
	var events []*event

	for {
		select {
		case <-t1.C:
			// jitter between 0 and 30 seconds
			time.Sleep(jitter.Do(time.Second * 30))
			// save and apply events
			if err := m.processEvents(events); err == nil {
				// clear the batch
				events = nil
			}
		case <-t2.C:
			// jitter between 0 and 1 minute
			time.Sleep(jitter.Do(time.Minute))
			// checks services to run in the store
			m.processServices()
		case <-t3.C:
			m.updateStatus()
		case ev := <-m.eventChan:
			// save an event
			events = append(events, ev)
		case <-m.exit:
			return
		}
	}
}

func (m *manager) String() string {
	return "manager"
}

func (m *manager) Init(opts ...runtime.Option) error {
	return nil
}

func (m *manager) Create(s *runtime.Service, opts ...runtime.CreateOption) error {
	m.Lock()
	defer m.Unlock()

	// we need to parse options to get the env vars
	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}

	// create service key
	k := key(s)

	rs := &runtimeService{
		Service: s,
		Options: &options,
		Status:  "starting",
	}

	// save locally
	m.services[k] = rs

	// create a new event
	ev := &event{
		Type:    "create",
		Service: s,
		Options: &options,
	}

	// send event
	go m.sendEvents(ev)

	// process the event immediately
	go m.processEvent(ev)

	// marshall the content
	b, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	// save the record
	return m.Store.Write(&store.Record{
		Key:   k,
		Value: b,
	})
}

func (m *manager) Read(opts ...runtime.ReadOption) ([]*runtime.Service, error) {
	var options runtime.ReadOptions
	for _, o := range opts {
		o(&options)
	}

	var services []*runtime.Service

	m.RLock()
	defer m.RUnlock()

	for _, rs := range m.services {
		srv := options.Service
		ver := options.Version
		typ := options.Type

		if len(srv) > 0 && rs.Service.Name != srv {
			continue
		}

		if len(ver) > 0 && rs.Service.Version != ver {
			continue
		}

		if len(typ) > 0 && rs.Service.Metadata["type"] != typ {
			continue
		}

		services = append(services, copyService(rs))
	}

	return services, nil
}

func (m *manager) Update(s *runtime.Service) error {
	m.Lock()
	defer m.Unlock()

	// create the key
	k := key(s)

	// read the existing record
	r, err := m.Store.Read(k)
	if err != nil {
		return err
	}

	// no service
	if len(r) == 0 {
		return errors.New("service not found")
	}

	var rs runtimeService
	if err := json.Unmarshal(r[0].Value, &rs); err != nil {
		return err
	}
	// set the service
	rs.Service = s
	// TODO: allow setting opts

	// if not running then run it
	evType := "update"

	// check if it exists
	if _, ok := m.services[k]; !ok {
		// set starting status
		rs.Status = "started"
		evType = "create"
		m.services[k] = &rs
	}

	// create event
	ev := &event{
		Type:    evType,
		Service: rs.Service,
		Options: rs.Options,
	}

	// fire an update
	go m.sendEvents(ev)

	// process the event immediately
	go m.processEvent(ev)

	// marshall the content
	b, err := json.Marshal(rs)
	if err != nil {
		return err
	}

	// save the record
	return m.Store.Write(&store.Record{
		Key:   k,
		Value: b,
	})
}

func (m *manager) Delete(s *runtime.Service) error {
	m.Lock()
	defer m.Unlock()

	k := key(s)

	// save local status
	v, ok := m.services[k]
	if !ok {
		return nil
	}

	// set status
	v.Status = "stopped"
	v.Service.Metadata["status"] = "stopped"

	// create new event
	ev := &event{
		Type:    "delete",
		Service: v.Service,
	}

	// fire an update
	go m.sendEvents(ev)

	// process the event immediately
	go m.processEvent(ev)

	// delete from store
	return m.Store.Delete(k)
}

func (m *manager) List() ([]*runtime.Service, error) {
	m.RLock()
	defer m.RUnlock()

	services := make([]*runtime.Service, 0, len(m.services))

	for _, service := range m.services {
		services = append(services, copyService(service))
	}

	return services, nil
}

func (m *manager) Start() error {
	m.Lock()
	defer m.Unlock()

	if m.running {
		return nil
	}

	// reset the exit channel
	m.exit = make(chan bool)

	// start the runtime
	if err := m.Runtime.Start(); err != nil {
		return err
	}

	// start the internal manager
	go m.run()

	// set to running
	m.running = true

	return nil
}

func (m *manager) Stop() error {
	m.Lock()
	defer m.Unlock()

	if !m.running {
		return nil
	}

	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
		m.Runtime.Stop()
		m.running = false
	}

	return nil
}

func (m *manager) Logs(s *runtime.Service, options ...runtime.LogsOption) (runtime.LogStream, error) {
	return m.Runtime.Logs(s, options...)
}

func newManager(ctx *cli.Context, r runtime.Runtime, s store.Store) *manager {
	var profile []string
	// peel out the env
	switch ctx.String("profile") {
	case "platform":
		profile = muProfile.Platform()
	case "server":
		profile = muProfile.Server()
	}

	return &manager{
		Id:        uuid.New().String(),
		Runtime:   r,
		Store:     s,
		Profile:   profile,
		services:  make(map[string]*runtimeService),
		events:    make(map[string]*series),
		eventChan: make(chan *event, 64),
		exit:      make(chan bool),
	}
}
