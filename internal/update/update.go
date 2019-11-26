package update

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/micro/go-micro/runtime"
	"github.com/micro/go-micro/util/log"
)

var (
	// DefaultTick defines how often we poll for updates
	DefaultTick = 1 * time.Minute
	// DefaultURL defines url to poll for updates
	DefaultURL = "https://micro.mu/update"
)

// Build is service build
type Build struct {
	// Commit is git commit sha
	Commit string `json:"commit,omitempty"`
	// Image is Docker build timestamp
	Image string `json:"image"`
	// Release is micro release tag
	Release string `json:"release,omitempty"`
}

// notifier is http notifier
type notifier struct {
	sync.RWMutex
	// url to poll for updates
	url string
	// poll time to check for updates
	tick time.Duration
	// version is current version
	version time.Time
	// events is notifications channel
	events chan runtime.Event
	// indicates if we're running
	running bool
	// used to stop the runtime
	closed chan bool
}

// NewNotifier returns new runtime notifier
func NewNotifier(buildDate string) runtime.Notifier {
	// convert the build date to a time.Time value
	timestamp, err := strconv.ParseInt(buildDate, 10, 64)
	if err != nil {
		timestamp = time.Now().Unix()
	}

	// the current version
	version := time.Unix(timestamp, 0)

	// return a new notifier
	return newNotifier(DefaultURL, DefaultTick, version)
}

// NewHTTP creates HTTP poller and returns it
func newNotifier(url string, tick time.Duration, version time.Time) *notifier {
	return &notifier{
		url:     url,
		tick:    tick,
		version: version,
		closed:  make(chan bool),
	}
}

// Poll polls for updates and returns results
func (h *notifier) poll() (*Build, error) {
	// this should not return error, but lets make sure
	url, err := url.Parse(h.url)
	if err != nil {
		return nil, err
	}

	rsp, err := http.Get(url.String())
	if err != nil {
		log.Debugf("Notifier error polling updates: %v", err)
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != 200 {
		log.Debugf("Notifier error unexpected http response: %v", rsp.StatusCode)
		return nil, err
	}

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Debugf("Notifier error reading http response: %v", err)
		return nil, err
	}

	// encoding format is assumed to be json
	var build *Build
	if err := json.Unmarshal(b, &build); err != nil {
		log.Debugf("Notifier error unmarshalling response: %v", err)
		return nil, err
	}

	return build, nil
}

// run runs the notifier
func (h *notifier) run() {
	t := time.NewTicker(h.tick)
	defer t.Stop()

	for {
		select {
		case <-h.closed:
			return
		case <-t.C:
			log.Debugf("Notifier polling for new update: %s", h.url)
			resp, err := h.poll()
			if err != nil {
				log.Debugf("Notifier error polling for updates: %v", err)
				continue
			}
			// parse returned response to timestamp
			buildTime, err := time.Parse(time.RFC3339, resp.Image)
			if err != nil {
				log.Debugf("Notifier error parsing build time: %v", err)
				continue
			}

			// if the latest build is newer than the current emit Update event
			if !buildTime.After(h.version) {
				continue
			}

			// fire the event
			h.events <- runtime.Event{
				// new update
				Type: runtime.Update,
				// timestamp of the update
				Timestamp: buildTime,
			}

			// set the build time
			h.version = buildTime
		}
	}
}

// Notify polls for new build and returns a channel to consume the events
func (h *notifier) Notify() (<-chan runtime.Event, error) {
	h.Lock()
	defer h.Unlock()

	// already running
	if h.running {
		return h.events, nil
	}

	// set running
	h.running = true
	h.closed = make(chan bool)
	h.events = make(chan runtime.Event)

	// runt the notifier
	go h.run()

	return h.events, nil
}

// Close stops the notifier
func (h *notifier) Close() error {
	h.Lock()
	defer h.Unlock()

	if !h.running {
		return nil
	}

	select {
	case <-h.closed:
		return nil
	default:
		close(h.closed)
		// stop the event stream
		close(h.events)
		// set not running
		h.running = false
	}

	return nil
}

// String implements tringer interface
func (h *notifier) String() string {
	return "default"
}
