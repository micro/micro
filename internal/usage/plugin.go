package usage

import (
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/registry"
	"github.com/micro/micro/plugin"
)

func init() {
	plugin.Register(Plugin())
}

func Plugin() plugin.Plugin {
	var requests uint64

	// create rand
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return plugin.NewPlugin(
		plugin.WithName("usage"),
		plugin.WithInit(func(c *cli.Context) error {
			// only do if enabled
			if !c.Bool("report_usage") {
				os.Setenv("MICRO_REPORT_USAGE", "false")
				return nil
			}

			var service string

			// set service name
			if c.Args().Len() > 0 && len(c.Args().Get(0)) > 0 {
				service = c.Args().Get(0)
			}

			// kick off the tracker
			go func() {
				// new report
				u := New(service)

				// initial publish in 30-60 seconds
				d := 30 + r.Intn(30)
				time.Sleep(time.Second * time.Duration(d))

				for {
					// get service list
					s, _ := registry.ListServices()
					// get requests
					reqs := atomic.LoadUint64(&requests)
					srvs := uint64(len(s))

					// reset requests
					atomic.StoreUint64(&requests, 0)

					// set metrics
					u.Metrics.Count["instances"] = uint64(1)
					u.Metrics.Count["requests"] = reqs
					u.Metrics.Count["services"] = srvs

					// send report
					Report(u)

					// now sleep 24 hours
					time.Sleep(time.Hour * 24)
				}
			}()

			return nil
		}),
		plugin.WithHandler(func(h http.Handler) http.Handler {
			// only enable if set
			if v := os.Getenv("MICRO_REPORT_USAGE"); v == "false" {
				return h
			}

			// return usage recorder
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// count requests
				atomic.AddUint64(&requests, 1)
				// serve the request
				h.ServeHTTP(w, r)
			})
		}),
	)
}
