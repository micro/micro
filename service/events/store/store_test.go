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
// Original source: github.com/micro/go-micro/v3/events/store/store_test.go

package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/events"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	store := NewStore()

	testData := []events.Event{
		{ID: uuid.New().String(), Topic: "foo"},
		{ID: uuid.New().String(), Topic: "foo"},
		{ID: uuid.New().String(), Topic: "bar"},
	}

	// write the records to the store
	t.Run("Write", func(t *testing.T) {
		for _, event := range testData {
			err := store.Write(&event)
			assert.Nilf(t, err, "Writing an event should not return an error")
		}
	})

	// should not be able to read events from a blank topic
	t.Run("ReadMissingTopic", func(t *testing.T) {
		evs, err := store.Read("")
		assert.Equal(t, err, events.ErrMissingTopic, "Reading a blank topic should return an error")
		assert.Nil(t, evs, "No events should be returned")
	})

	// should only get the events from the topic requested
	t.Run("ReadTopic", func(t *testing.T) {
		evs, err := store.Read("foo")
		assert.Nilf(t, err, "No error should be returned")
		assert.Len(t, evs, 2, "Only the events for this topic should be returned")
	})

	// limits should be honoured
	t.Run("ReadTopicLimit", func(t *testing.T) {
		evs, err := store.Read("foo", events.ReadLimit(1))
		assert.Nilf(t, err, "No error should be returned")
		assert.Len(t, evs, 1, "The result should include no more than the read limit")
	})
}
