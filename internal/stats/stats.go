package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type stats struct {
	mux *http.ServeMux

	sync.RWMutex

	Started int64  `json:"started"`
	Memory  string `json:"memory"`
	Threads int    `json:"threads"`
	GC      string `json:"gc_pause"`

	Counters []*counter `json:"counters"`

	running bool
	exit    chan bool
}

type counter struct {
	// time created
	Timestamp int64 `json:"timestamp"`
	// counters
	Status map[string]int `json:"status_codes"`
	Total  int            `json:"total_reqs"`
}

var (
	// 5 second window
	window = time.Second * 5
	// 120 seconds total
	total = 24
)

func (s *stats) handler(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	b, err := json.Marshal(s)
	s.RUnlock()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (s *stats) run() {
	t := time.NewTicker(window)
	w := 0

	for {
		select {
		case <-s.exit:
			t.Stop()
			return
		case <-t.C:
			// roll
			s.Lock()
			s.Counters = append(s.Counters, &counter{
				Timestamp: time.Now().Unix(),
				Status:    make(map[string]int),
			})
			if len(s.Counters) >= total {
				s.Counters = s.Counters[1:]
			}

			w++
			if w >= 2 {
				var mstat runtime.MemStats
				runtime.ReadMemStats(&mstat)
				s.Threads = runtime.NumGoroutine()
				s.Memory = fmt.Sprintf("%.2fmb", float64(mstat.Alloc)/float64(1024*1024))
				s.GC = fmt.Sprintf("%.3fms", float64(mstat.PauseTotalNs)/(1000*1000))
				w = 0
			}
			s.Unlock()
		}
	}
}

func (s *stats) Record(c string, t int) {
	s.Lock()
	counter := s.Counters[len(s.Counters)-1]
	counter.Status[c] += t
	counter.Total += t
	s.Counters[len(s.Counters)-1] = counter
	s.Unlock()
}

func (s *stats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var code string
	rw := &writer{w, 200}
	s.mux.ServeHTTP(rw, r)

	switch {
	case rw.status >= 500:
		code = "50x"
	case rw.status >= 400:
		code = "40x"
	case rw.status >= 300:
		code = "30x"
	case rw.status >= 200:
		code = "20x"
	}

	s.Record(code, 1)
}

func (s *stats) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	s.Started = time.Now().Unix()
	s.exit = make(chan bool)
	go s.run()
	return nil
}

func (s *stats) Stop() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	close(s.exit)
	s.Started = 0
	return nil
}

func New(p string, m *http.ServeMux) *stats {
	var mstat runtime.MemStats
	runtime.ReadMemStats(&mstat)

	s := &stats{
		mux:     m,
		Threads: runtime.NumGoroutine(),
		Memory:  fmt.Sprintf("%.2fmb", float64(mstat.Alloc)/float64(1024*1024)),
		GC:      fmt.Sprintf("%.3fms", float64(mstat.PauseTotalNs)/(1000*1000)),
		Counters: []*counter{
			&counter{
				Timestamp: time.Now().Unix(),
				Status:    make(map[string]int),
			},
		},
	}

	m.HandleFunc(p, s.handler)
	return s
}
