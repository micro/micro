package main

import (
	"math/rand"

	pb "github.com/micro/micro/v3/test/service/stream/proto"
)

func randomPoint(r *rand.Rand) *pb.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{Latitude: lat, Longitude: long}
}
