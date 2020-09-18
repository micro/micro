package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"

	pb "github.com/micro/micro/v3/test/service/stream/proto"
	"github.com/micro/micro/v3/test/service/stream/server/handler"
)

func main() {
	// Create the service
	srv := service.New(
		service.Name("stream"),
		service.Version("latest"),
	)

	// Load the test data
	features, err := loadFeatures()
	if err != nil {
		logger.Fatalf("Error loading features: %v", err)
	}

	// Register the handler
	pb.RegisterRouteGuideHandler(srv.Server(), &handler.RouteGuide{
		Features: features,
		Notes:    make(map[string][]*pb.RouteNote),
	})

	// Run the service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}

// loadFeatures loads features from the JSON file.
func loadFeatures() ([]*pb.Feature, error) {
	data, err := ioutil.ReadFile("./data/features.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to load default features: %v", err)
	}

	var result []*pb.Feature
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("Failed to load default features: %v", err)
	}

	return result, nil
}
