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
// Original source: github.com/micro/go-micro/v3/util/ring/buffer.go

// Package ring provides a simple ring buffer for storing local data
package ring

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Buffer is ring buffer
type Buffer struct {
	size int

	sync.RWMutex
	vals    []*Entry
	streams map[string]*Stream
}

// Entry is ring buffer data entry
type Entry struct {
	Value     interface{}
	Timestamp time.Time
}

// Stream is used to stream the buffer
type Stream struct {
	// Id of the stream
	Id string
	// Buffered entries
	Entries chan *Entry
	// Stop channel
	Stop chan bool
}

// Put adds a new value to ring buffer
func (b *Buffer) Put(v interface{}) {
	b.Lock()
	defer b.Unlock()

	// append to values
	entry := &Entry{
		Value:     v,
		Timestamp: time.Now(),
	}
	b.vals = append(b.vals, entry)

	// trim if bigger than size required
	if len(b.vals) > b.size {
		b.vals = b.vals[1:]
	}

	// send to every stream
	for _, stream := range b.streams {
		select {
		case <-stream.Stop:
			delete(b.streams, stream.Id)
			close(stream.Entries)
		case stream.Entries <- entry:
		}
	}
}

// Get returns the last n entries
func (b *Buffer) Get(n int) []*Entry {
	b.RLock()
	defer b.RUnlock()

	// reset any invalid values
	if n > len(b.vals) || n < 0 {
		n = len(b.vals)
	}

	// create a delta
	delta := len(b.vals) - n

	// return the delta set
	return b.vals[delta:]
}

// Return the entries since a specific time
func (b *Buffer) Since(t time.Time) []*Entry {
	b.RLock()
	defer b.RUnlock()

	// return all the values
	if t.IsZero() {
		return b.vals
	}

	// if its in the future return nothing
	if time.Since(t).Seconds() < 0.0 {
		return nil
	}

	for i, v := range b.vals {
		// find the starting point
		d := v.Timestamp.Sub(t)

		// return the values
		if d.Seconds() > 0.0 {
			return b.vals[i:]
		}
	}

	return nil
}

// Stream logs from the buffer
// Close the channel when you want to stop
func (b *Buffer) Stream() (<-chan *Entry, chan bool) {
	b.Lock()
	defer b.Unlock()

	entries := make(chan *Entry, 128)
	id := uuid.New().String()
	stop := make(chan bool)

	b.streams[id] = &Stream{
		Id:      id,
		Entries: entries,
		Stop:    stop,
	}

	return entries, stop
}

// Size returns the size of the ring buffer
func (b *Buffer) Size() int {
	return b.size
}

// New returns a new buffer of the given size
func New(i int) *Buffer {
	return &Buffer{
		size:    i,
		streams: make(map[string]*Stream),
	}
}
