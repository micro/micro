package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/broker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type testObj struct {
	One string
	Two int64
}

func TestBroker(t *testing.T) {
	host := os.Getenv("REDIS_HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if len(port) == 0 {
		port = "6379"
	}
	b := NewBroker(broker.Addrs(fmt.Sprintf("redis://%s:%s", host, port)))

	assert.NoError(t, b.Connect())

	t.Run("Simple", func(t *testing.T) {
		count := 0
		doneCh := make(chan bool)
		h := func(message *broker.Message) error {
			count++
			if count == 10 {
				doneCh <- true
			}
			return nil
		}
		sub, err := b.Subscribe("foo", h)

		assert.NoError(t, err)
		assert.Equal(t, "foo", sub.Topic())

		// consumer 2
		count2 := 0
		doneCh2 := make(chan bool)
		h2 := func(message *broker.Message) error {
			assert.NotEmpty(t, message.Body)
			assert.NotEmpty(t, message.Header)
			var tobj testObj
			assert.NoError(t, json.Unmarshal(message.Body, &tobj))
			i, err := strconv.Atoi(tobj.One)
			assert.NoError(t, err)
			assert.Equal(t, int64(i), tobj.Two)

			count2++
			if count2 == 10 {
				doneCh2 <- true
			}
			return nil
		}
		sub2, err := b.Subscribe("foo", h2)

		assert.NoError(t, err)
		assert.Equal(t, "foo", sub2.Topic())

		for i := 0; i < 10; i++ {
			tobj := testObj{
				One: fmt.Sprintf("%d", i),
				Two: int64(i),
			}
			by, _ := json.Marshal(tobj)
			assert.NoError(t, b.Publish("foo", &broker.Message{
				Header: map[string]string{"fooheader": "bar"},
				Body:   by,
			}))
		}

		select {
		case <-doneCh:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for messages, recieved %d", count)
		}
		assert.NoError(t, sub.Unsubscribe())

		select {
		case <-doneCh2:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for messages, received %d", count2)
		}
		assert.NoError(t, sub2.Unsubscribe())

	})
	t.Run("ErrorHandler", func(t *testing.T) {
		doneCh := make(chan bool)
		h := func(message *broker.Message) error {
			return errors.New("BARF")
		}
		eh := func(message *broker.Message, err error) {
			doneCh <- true
		}
		sub, err := b.Subscribe("fooerror", h, broker.HandleError(eh))

		assert.NoError(t, err)
		assert.Equal(t, "fooerror", sub.Topic())

		tobj := testObj{
			One: "1",
			Two: 2,
		}
		by, _ := json.Marshal(tobj)
		assert.NoError(t, b.Publish("fooerror", &broker.Message{
			Header: map[string]string{"fooheader": "bar"},
			Body:   by,
		}))

		select {
		case <-doneCh:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for error processing")
		}

	})

	t.Run("WithQueue", func(t *testing.T) {
		count := 0
		count2 := 0
		doneCh := make(chan bool)
		allDoneCh := make(chan bool)
		h := func(message *broker.Message) error {
			if count == 0 {
				doneCh <- true
			}
			count++
			if count+count2 == 100 {
				allDoneCh <- true
			}
			return nil
		}
		sub, err := b.Subscribe("fooqueue", h, broker.Queue("queue1"))

		assert.NoError(t, err)
		assert.Equal(t, "fooqueue", sub.Topic())

		// consumer 2
		doneCh2 := make(chan bool)
		h2 := func(message *broker.Message) error {
			if count2 == 0 {
				doneCh2 <- true
			}
			count2++
			if count+count2 == 100 {
				allDoneCh <- true
			}
			return nil
		}
		sub2, err := b.Subscribe("fooqueue", h2, broker.Queue("queue1"))

		assert.NoError(t, err)
		assert.Equal(t, "fooqueue", sub2.Topic())

		for i := 0; i < 100; i++ {
			tobj := testObj{
				One: fmt.Sprintf("%d", i),
				Two: int64(i),
			}
			by, _ := json.Marshal(tobj)
			assert.NoError(t, b.Publish("fooqueue", &broker.Message{
				Header: map[string]string{"fooheader": "bar"},
				Body:   by,
			}))
		}

		select {
		case <-doneCh:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for messages for channel1")
		}

		select {
		case <-doneCh2:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for messages for channel2")
		}

		select {
		case <-allDoneCh:
		case <-time.After(5 * time.Second):
			t.Errorf("Timed out waiting for all the messages, received %d", count+count2)
		}

		assert.NoError(t, sub.Unsubscribe())
		assert.NoError(t, sub2.Unsubscribe())
	})

}
