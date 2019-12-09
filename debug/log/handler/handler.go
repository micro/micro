package handler

import (
	"context"

	pb "github.com/micro/micro/debug/log/proto"
)

type Log struct{}

func (l *Log) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	return nil
}
