package main

import (
	"time"

	goevents "github.com/micro/go-micro/v3/events"
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
	go func() {
		ticker := time.NewTicker(time.Second)

		for {
			err := events.Publish("test",
				goevents.WithPayload(payload),
				goevents.WithMetadata(metadata),
			)

			if err != nil {
				logger.Errorf("Error publishing event: %v", err)
			} else {
				logger.Infof("Published event ok")
			}

			<-ticker.C
		}
	}()

	evChan, err := events.Subscribe(goevents.WithTopic("test"))
	if err != nil {
		logger.Fatalf("Error creating subscriber: %v", err)
	}

	for {
		event, ok := <-evChan
		if !ok {
			logger.Fatal("Channel closed")
		}

		if event.Topic != "test" {
			logger.Errorf("Incorrect topic: %v", event.Topic)
			continue
		}
		if event.Metadata == nil || event.Metadata["foo"] != metadata["foo"] {
			logger.Errorf("Incorrect metadata: %v", event.Metadata)
			continue
		}

		var result testPayload
		if err := event.Unmarshal(&result); err != nil {
			logger.Errorf("Error decoding result: %v. Payload was %v", err, string(event.Payload))
			continue
		}
		if result.Message != payload.Message {
			logger.Errorf("Incorrect message: %v", result.Message)
			continue
		}

		logger.Infof("Recieved event ok")
	}
}
