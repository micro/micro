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
// Original source: github.com/micro/go-micro/v3/debug/trace/trace.go

// Package trace provides an interface for distributed tracing
package trace

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/context/metadata"
)

// Tracer is an interface for distributed tracing
type Tracer interface {
	// Start a trace
	Start(ctx context.Context, name string) (context.Context, *Span)
	// Finish the trace
	Finish(*Span) error
	// Read the traces
	Read(...ReadOption) ([]*Span, error)
}

// SpanType describe the nature of the trace span
type SpanType int

const (
	// SpanTypeRequestInbound is a span created when serving a request
	SpanTypeRequestInbound SpanType = iota
	// SpanTypeRequestOutbound is a span created when making a service call
	SpanTypeRequestOutbound
)

// Span is used to record an entry
type Span struct {
	// Id of the trace
	Trace string
	// name of the span
	Name string
	// id of the span
	Id string
	// parent span id
	Parent string
	// Start time
	Started time.Time
	// Duration in nano seconds
	Duration time.Duration
	// associated data
	Metadata map[string]string
	// Type
	Type SpanType
}

const (
	traceIDKey = "Micro-Trace-Id"
	spanIDKey  = "Micro-Span-Id"
)

// FromContext returns a span from context
func FromContext(ctx context.Context) (traceID string, parentSpanID string, isFound bool) {
	traceID, traceOk := metadata.Get(ctx, traceIDKey)
	microID, microOk := metadata.Get(ctx, "Micro-Id")
	if !traceOk && !microOk {
		isFound = false
		return
	}
	if !traceOk {
		traceID = microID
	}
	parentSpanID, ok := metadata.Get(ctx, spanIDKey)
	return traceID, parentSpanID, ok
}

// ToContext saves the trace and span ids in the context
func ToContext(ctx context.Context, traceID, parentSpanID string) context.Context {
	return metadata.MergeContext(ctx, map[string]string{
		traceIDKey: traceID,
		spanIDKey:  parentSpanID,
	}, true)
}

var (
	DefaultTracer Tracer = new(noop)
)

type noop struct{}

func (n *noop) Init(...Option) error {
	return nil
}

func (n *noop) Start(ctx context.Context, name string) (context.Context, *Span) {
	return nil, nil
}

func (n *noop) Finish(*Span) error {
	return nil
}

func (n *noop) Read(...ReadOption) ([]*Span, error) {
	return nil, nil
}
