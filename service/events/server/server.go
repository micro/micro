package server

import (
	"time"

	goevents "github.com/micro/micro/v3/service/events"
	pb "github.com/micro/micro/v3/proto/events"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/urfave/cli/v2"
)

var systemTopics = []string{"runtime"}

// Run the micro broker
func Run(ctx *cli.Context) error {
	// new service
	srv := service.New(
		service.Name("events"),
	)

	// register the handlers
	pb.RegisterStreamHandler(srv.Server(), new(Stream))
	pb.RegisterStoreHandler(srv.Server(), new(Store))

	// subscribe to the system topics
	for _, topic := range systemTopics {
		go watch(topic)
	}

	// run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	return nil
}

// watch a topic and store the events published in the store
func watch(topic string) {
	stream, err := events.Subscribe(topic, goevents.WithQueue("events"))
	if err != nil {
		logger.Errorf("Error subscribing to topic %v: %v", topic, err)
		return
	}
	logger.Infof("Watching system topic: %v", topic)

	for {
		event, ok := <-stream
		if !ok {
			logger.Debugf("Stream closed for topic %v", topic)
			return
		}

		if err := events.DefaultStore.Write(&event, goevents.WithTTL(time.Hour*24)); err != nil {
			logger.Errorf("Error writing event %v to store: %v", event.ID, err)
		}
	}
}
