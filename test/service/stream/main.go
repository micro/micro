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

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() { // test 1, ordinary sub with autoack
		defer wg.Done()
		evChan, err := events.Subscribe("test1")
		if err != nil {
			logger.Fatalf("Error creating subscriber: %v", err)
		}

		err = events.Publish("test1", payload, goevents.WithMetadata(metadata))
		if err != nil {
			logger.Errorf("Error publishing event: %v", err)
		} else {
			logger.Infof("TEST1: Published event ok")
		}

		event, ok := <-evChan
		if !ok {
			logger.Fatal("Channel closed")
		}

		if event.Topic != "test1" {
			logger.Errorf("Incorrect topic: %v", event.Topic)
			return
		}
		if event.Metadata == nil || event.Metadata["foo"] != metadata["foo"] {
			logger.Errorf("Incorrect metadata: %v", event.Metadata)
			return
		}

		var result testPayload
		if err := event.Unmarshal(&result); err != nil {
			logger.Errorf("Error decoding result: %v. Payload was %v", err, string(event.Payload))
			return
		}
		if result.Message != payload.Message {
			logger.Errorf("Incorrect message: %v", result.Message)
			return
		}

		logger.Infof("TEST1: Received event ok")

	}()

	wg.Add(1)
	go func() { // test 2, sub with manual ack
		defer wg.Done()
		evChan, err := events.Subscribe("test2", goevents.WithAutoAck(false, 5*time.Second))
		if err != nil {
			logger.Fatalf("Error creating subscriber: %v", err)
		}

		err = events.Publish("test2", payload, goevents.WithMetadata(metadata))

		if err != nil {
			logger.Errorf("Error publishing event: %v", err)
		} else {
			logger.Infof("TEST2: Published event ok")
		}

		event, ok := <-evChan
		if !ok {
			logger.Fatal("Channel closed")
		}

		if event.Topic != "test2" {
			logger.Errorf("Incorrect topic: %v", event.Topic)
			return
		}
		if event.Metadata == nil || event.Metadata["foo"] != metadata["foo"] {
			logger.Errorf("Incorrect metadata: %v", event.Metadata)
			return
		}

		var result testPayload
		if err := event.Unmarshal(&result); err != nil {
			logger.Errorf("Error decoding result: %v. Payload was %v", err, string(event.Payload))
			return
		}
		if result.Message != payload.Message {
			logger.Errorf("Incorrect message: %v", result.Message)
			return
		}
		id := event.ID
		// nack the event to put it back on the queue
		event.Nack()

		select {
		case event := <-evChan:
			if event.ID != id {
				logger.Errorf("Unexpected event received, expected %s, received %s", id, event.ID)
			}
		case <-time.After(10 * time.Second):
			logger.Errorf("Timed out waiting for event")
		}

		logger.Infof("TEST2: Finished ok")

	}()

	wg.Wait()
}
