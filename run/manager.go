package run

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/micro/go-run"
	"github.com/pborman/uuid"
)

type manager struct {
	runtime run.Runtime
	// uuid to process id
	sync.RWMutex
	services map[string]*service
}

type service struct {
	exit chan bool

	info    string
	updated time.Time
	process *run.Process
}

func (m *manager) run() {
	t := time.NewTicker(time.Minute)

	for _ = range t.C {
		m.Lock()
		for id, s := range m.services {
			// only process stopped
			if !(s.info == "stopped" || strings.HasPrefix(s.info, "error")) {
				continue
			}

			// get time since last update
			t := time.Since(s.updated)

			// delete stopped older than 10 minutes
			if t.Seconds() > 900 {
				delete(m.services, id)
			}
		}
		m.Unlock()
	}
}

func (m *manager) update(uid, info string) error {
	m.Lock()
	defer m.Unlock()

	srv, ok := m.services[uid]
	if !ok {
		return errors.New("does not exist")
	}

	return srv.update(info)
}

func (m *manager) setProc(uid string, proc *run.Process) {
	m.Lock()
	defer m.Unlock()

	if srv, ok := m.services[uid]; ok {
		srv.process = proc
	}
}

func (m *manager) Run(url string, re, u bool) string {
	uid := uuid.NewUUID().String()

	// save id
	m.Lock()
	m.services[uid] = &service{
		exit: make(chan bool),
		info: "pre-fetch",
	}
	m.Unlock()

	go func() {
		// get the source
		if err := m.update(uid, "fetching"); err != nil {
			return
		}

		src, err := m.runtime.Fetch(url, run.Update(u))
		if err != nil {
			m.update(uid, "error:"+err.Error())
			return
		}

		// build the binary
		if err := m.update(uid, "building"); err != nil {
			return
		}

		bin, err := m.runtime.Build(src)
		if err != nil {
			m.update(uid, "error:"+err.Error())
			return
		}

		for {
			// execute the binary
			if err := m.update(uid, "executing"); err != nil {
				return
			}

			proc, err := m.runtime.Exec(bin)
			if err != nil {
				m.update(uid, "error:"+err.Error())
				return
			}

			// set service process
			m.setProc(uid, proc)

			// wait till exit
			if err := m.update(uid, "running"); err != nil {
				return
			}

			// bail if not restarting
			if !re {
				if err := m.runtime.Wait(proc); err != nil {
					m.update(uid, "error:"+err.Error())
				}
				return
			}

			// log error since we manage the cycle
			if err := m.runtime.Wait(proc); err != nil {
				if err := m.update(uid, "error:"+err.Error()); err != nil {
					return
				}
			}

			// log restart
			if err := m.update(uid, "restarting"); err != nil {
				return
			}
		}

	}()

	return uid
}

func (m *manager) Status(id string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	srv, ok := m.services[id]
	if !ok {
		return "", errors.New(id + " does not exist")
	}

	return srv.info, nil
}

func (m *manager) Stop(id string) error {
	m.Lock()
	defer m.Unlock()

	srv, ok := m.services[id]
	if !ok {
		return errors.New(id + " does not exist")
	}

	// check if its already stopped
	if srv.info == "stopped" {
		return nil
	}
	// kill
	if srv.process != nil {
		m.runtime.Kill(srv.process)
	}

	// stop
	srv.stop()

	return nil
}

func (s *service) update(msg string) error {
	select {
	case <-s.exit:
		return errors.New("stopped")
	default:
		s.info = msg
		s.updated = time.Now()
	}
	return nil
}

func (s *service) stop() {
	select {
	case <-s.exit:
		return
	default:
		close(s.exit)
		s.info = "stopped"
		s.updated = time.Now()
	}
}

func newManager(r run.Runtime) *manager {
	m := &manager{
		runtime:  r,
		services: make(map[string]*service),
	}
	go m.run()
	return m
}
