// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/micro/v3/metrics/wrapper/metrics_wrapper.go

package wrapper

import (
	"time"

	"context"

	"github.com/micro/micro/v3/service/metrics"
	"github.com/micro/micro/v3/service/server"
)

// Wrapper provides a HandlerFunc for metrics.Reporter implementations:
type Wrapper struct {
	reporter metrics.Reporter
}

// New returns a *Wrapper configured with the given metrics.Reporter:
func New(reporter metrics.Reporter) *Wrapper {
	return &Wrapper{
		reporter: reporter,
	}
}

// HandlerFunc instruments handlers registered to a service:
func (w *Wrapper) HandlerFunc(handlerFunction server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {

		// Build some tags to describe the call:
		tags := metrics.Tags{
			"method": req.Method(),
		}

		// Start the clock:
		callTime := time.Now()

		// Run the handlerFunction:
		err := handlerFunction(ctx, req, rsp)

		// Add a result tag:
		if err != nil {
			tags["result"] = "failure"
		} else {
			tags["result"] = "success"
		}

		// Instrument the result (if the DefaultClient has been configured):
		w.reporter.Timing("service.handler", time.Since(callTime), tags)

		return err
	}
}
