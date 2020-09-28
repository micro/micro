package main

import (
	"sync"
	"time"

	goevents "github.com/micro/go-micro/v3/events"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

type testPayload struct {
	Message string
}

var (
	payload  = testPayload{Message: "HelloWorld"}
	metadata = map[string]string{"foo": "bar"}
)

func main() {
	srv := service.New()
	srv.Init()

	go func() { // test 1, ordinary sub with autoack
		evChan, err := events.Subscribe("test1")
		if err != nil {
			logger.Fatalf("Error creating subscriber: %v", err)
		}
		// Is there a race condition here with publish coming straight after subscribe? Seems to be according to this
		// test so sleep for a bit to wait for nats to register subscription properly
		time.Sleep(2 * time.Second)

		logger.Infof("TEST1: publishing event")

		err = events.Publish("test1", payload, goevents.WithMetadata(metadata))
		if err != nil {
			logger.Errorf("TEST1: Error publishing event: %v", err)
			return
		}

		logger.Infof("TEST1: waiting for event")
		event, ok := <-evChan
		if !ok {
			logger.Error("TEST1: Channel closed")
			return
		}

		if event.Topic != "test1" {
			logger.Errorf("TEST1: Incorrect topic: %v", event.Topic)
			return
		}
		if event.Metadata == nil || event.Metadata["foo"] != metadata["foo"] {
			logger.Errorf("TEST1: Incorrect metadata: %v", event.Metadata)
			return
		}

		var result testPayload
		if err := event.Unmarshal(&result); err != nil {
			logger.Errorf("TEST1: Error decoding result: %v. Payload was %v", err, string(event.Payload))
			return
		}
		if result.Message != payload.Message {
			logger.Errorf("TEST1: Incorrect message: %v", result.Message)
			return
		}

		logger.Infof("TEST1: Finished ok")

	}()

	go func() { // test 2, sub with manual ack
		evChan, err := events.Subscribe("test2", goevents.WithAutoAck(false, 5*time.Second))
		if err != nil {
			logger.Errorf("TEST2: Error creating subscriber: %v", err)
			return
		}

		// Is there a race condition here with publish coming straight after subscribe? Seems to be according to this
		// test so sleep for a bit to wait for nats to register subscription properly
		time.Sleep(2 * time.Second)

		logger.Infof("TEST2: publishing event")
		err = events.Publish("test2", payload, goevents.WithMetadata(metadata))

		if err != nil {
			logger.Errorf("TEST2: Error publishing event: %v", err)
			return
		}

		logger.Infof("TEST2: waiting for event")
		event, ok := <-evChan
		if !ok {
			logger.Errorf("TEST2: Channel closed")
			return
		}

		if event.Topic != "test2" {
			logger.Errorf("TEST2: Incorrect topic: %v", event.Topic)
			return
		}
		if event.Metadata == nil || event.Metadata["foo"] != metadata["foo"] {
			logger.Errorf("TEST2: Incorrect metadata: %v", event.Metadata)
			return
		}

		var result testPayload
		if err := event.Unmarshal(&result); err != nil {
			logger.Errorf("TEST2: Error decoding result: %v. Payload was %v", err, string(event.Payload))
			return
		}
		if result.Message != payload.Message {
			logger.Errorf("TEST2: Incorrect message: %v", result.Message)
			return
		}
		id := event.ID
		// nack the event to put it back on the queue
		logger.Infof("TEST2: Nacking the event")
		event.Nack()

		logger.Infof("TEST2: waiting for event")
		select {
		case event := <-evChan:
			if event.ID != id {
				logger.Errorf("Unexpected event received, expected %s, received %s", id, event.ID)
				return
			}
			logger.Infof("TEST2: Acking the event")
			event.Ack()

		case <-time.After(10 * time.Second):
			logger.Errorf("Timed out waiting for event")
			return
		}

		// we've acked so should receive nothing else
		logger.Infof("TEST2: Waiting for timeout")
		select {
		case event := <-evChan:
			logger.Errorf("Unexpected event received %s", event.ID)
			return
		case <-time.After(10 * time.Second):

		}

		logger.Infof("TEST2: Finished ok")

	}()
	wg := sync.WaitGroup{}
	wg.Add(1)
	// wait indefinitely so this only runs once
	wg.Wait()
}
