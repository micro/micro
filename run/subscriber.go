package run

import (
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"

	proto "github.com/micro/micro/run/proto"
)

type subHandler struct {
	m *manager
}

// Run invokes the run command
func (h *subHandler) Run(ctx context.Context, ev *proto.RunRequest) error {
	if len(ev.Url) == 0 {
		return errors.BadRequest(Name+".run", "url is blank")
	}

	// TODO: what to do with error?
	go h.m.Run(ev.Url, ev.Restart, ev.Update)
	return nil
}

// Stop invokes the stop command
func (h *subHandler) Stop(ctx context.Context, ev *proto.StopRequest) error {
	if len(ev.Url) == 0 {
		return errors.BadRequest(Name+".stop", "url is blank")
	}

	// TODO: what to do with error?
	h.m.Stop(ev.Url)
	return nil
}
