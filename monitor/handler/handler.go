package handler

import (
	"context"

	"github.com/micro/go-micro/v2/monitor"
	pb "github.com/micro/micro/v2/monitor/proto"
)

type Monitor struct {
	Monitor monitor.Monitor
}

func (m *Monitor) Check(ctx context.Context, req *pb.CheckRequest, rsp *pb.CheckResponse) error {
	err := m.Monitor.Check(req.Service)
	if err != nil {
		rsp.Status = err.Error()
		return nil
	}
	rsp.Status = "ok"
	return nil
}
