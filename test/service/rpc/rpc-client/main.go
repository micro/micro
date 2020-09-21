package main

import (
	"context"
	"io"
	"math/rand"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/micro/v3/test/service/rpc/proto"
)

func main() {
	srv := service.New()
	cli := pb.NewRouteGuideService("rpc", srv.Client())

	// Looking for a valid feature
	logger.Infof("Testing Unary... Starting")
	printFeature(cli, &pb.Point{Latitude: 409146138, Longitude: -746188906})
	printFeature(cli, &pb.Point{Latitude: 0, Longitude: 0})
	logger.Infof("Testing Unary... Done")

	// Looking for features between 40, -75 and 42, -73.
	logger.Infof("Testing Server Streaming... Starting")
	printFeatures(cli, &pb.Rectangle{
		Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	})
	logger.Infof("Testing Server Streaming... Done")

	// RecordRoute
	logger.Infof("Testing Client Streaming... Starting")
	runRecordRoute(cli)
	logger.Infof("Testing Client Streaming... Done")

	// RouteChat
	logger.Infof("Testing Bidirectional Streaming... Starting")
	runRouteChat(cli)
	logger.Infof("Testing Bidirectional Streaming... Done")

	logger.Infof("Client completed ok")

	// prevent the client from restarting when running using micro run
	srv.Run()
}

// printFeature gets the feature for the given point.
func printFeature(client pb.RouteGuideService, point *pb.Point) {
	logger.Tracef("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		logger.Fatalf("%v.GetFeatures(_) = _, %v: ", client, err)
	}
	logger.Tracef("Got Features: %v", feature)
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client pb.RouteGuideService, rect *pb.Rectangle) {
	logger.Tracef("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		logger.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		logger.Tracef("Got Feature: %v", feature)
	}
}

// runRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func runRecordRoute(client pb.RouteGuideService) {
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	logger.Tracef("Traversing %d points.", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordRoute(ctx)
	if err != nil {
		logger.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			logger.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}

	summary, err := stream.CloseAndRecv()
	if err != nil {
		logger.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	logger.Tracef("Route summary: %v", summary)
}

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func runRouteChat(client pb.RouteGuideService) {
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		logger.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				logger.Fatalf("Failed to receive a note : %v", err)
			}
			logger.Tracef("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			logger.Fatalf("Failed to send a note: %v", err)
		}
	}
	if err := stream.Close(); err != nil {
		logger.Fatalf("%v.Close() got error: %v, want %v", stream, err, nil)
	}
	<-waitc
}

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}
