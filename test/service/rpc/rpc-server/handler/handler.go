package handler

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	pb "github.com/micro/micro/v3/test/service/rpc/proto"
)

// RouteGuide implements the route guide handler interface
type RouteGuide struct {
	Features  []*pb.Feature
	Notes     map[string][]*pb.RouteNote
	NotesLock sync.Mutex
}

// GetFeature obtains the feature at a given position.
func (r *RouteGuide) GetFeature(ctx context.Context, point *pb.Point, feature *pb.Feature) error {
	for _, f := range r.Features {
		if proto.Equal(f.Location, point) {
			*feature = *f
			return nil
		}
	}

	// No feature was found, return an unnamed feature
	feature.Location = point
	return nil
}

// ListFeatures obtains the Features available within the given Rectangle.  Results are
// streamed rather than returned at once (e.g. in a response message with a
// repeated field), as the rectangle may cover a large area and contain a
// huge number of features.
func (r *RouteGuide) ListFeatures(ctx context.Context, rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesStream) error {
	for _, f := range r.Features {
		if inRange(f.Location, rect) {
			if err := stream.Send(f); err != nil {
				return err
			}
		}
	}

	return nil
}

// RecordRoute accepts a stream of Points on a route being traversed, returning a
// RouteSummary when traversal is completed.
func (r *RouteGuide) RecordRoute(ctx context.Context, stream pb.RouteGuide_RecordRouteStream) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()

	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()

			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, f := range r.Features {
			if proto.Equal(f.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

// RouteChat accepts a stream of RouteNotes sent while a route is being traversed,
// while receiving other RouteNotes (e.g. from other users).
func (r *RouteGuide) RouteChat(ctx context.Context, stream pb.RouteGuide_RouteChatStream) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)

		r.NotesLock.Lock()
		r.Notes[key] = append(r.Notes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*pb.RouteNote, len(r.Notes[key]))
		copy(rn, r.Notes[key])
		r.NotesLock.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}
