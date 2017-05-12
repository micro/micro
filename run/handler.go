package run

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	gorun "github.com/micro/go-run"
	"golang.org/x/net/context"

	proto "github.com/micro/micro/run/proto"
)

type runtimeHandler struct {
	r gorun.Runtime
}

type serviceHandler struct {
	m     *manager
	prun  micro.Publisher
	pstop micro.Publisher
}

func (h *runtimeHandler) Fetch(ctx context.Context, req *proto.FetchRequest, rsp *proto.FetchResponse) error {
	if len(req.Url) == 0 {
		return errors.BadRequest(Name+".fetch", "url is blank")
	}
	src, err := h.r.Fetch(req.Url, gorun.Update(req.Update))
	if err != nil {
		return errors.InternalServerError(Name+".fetch", err.Error())
	}
	rsp.Source = &proto.Source{
		Url: src.URL,
		Dir: src.Dir,
	}
	return nil
}

func (h *runtimeHandler) Build(ctx context.Context, req *proto.BuildRequest, rsp *proto.BuildResponse) error {
	if req.Source == nil {
		return errors.BadRequest(Name+".build", "source is nil")
	}
	bin, err := h.r.Build(&gorun.Source{
		URL: req.Source.Url,
		Dir: req.Source.Dir,
	})
	if err != nil {
		return errors.InternalServerError(Name+".build", err.Error())
	}
	rsp.Binary = &proto.Binary{
		Path:   bin.Path,
		Source: req.Source,
	}
	return nil
}

func (h *runtimeHandler) Exec(ctx context.Context, req *proto.ExecRequest, rsp *proto.ExecResponse) error {
	if req.Binary == nil {
		return errors.BadRequest(Name+".exec", "binary is nil")
	}
	if req.Binary.Source == nil {
		return errors.BadRequest(Name+".exec", "binary.Source is nil")
	}
	proc, err := h.r.Exec(&gorun.Binary{
		Path: req.Binary.Path,
		Source: &gorun.Source{
			URL: req.Binary.Source.Url,
			Dir: req.Binary.Source.Dir,
		},
	})
	if err != nil {
		return errors.InternalServerError(Name+".exec", err.Error())
	}
	rsp.Process = &proto.Process{
		Id:     proc.ID,
		Binary: req.Binary,
	}
	return nil
}

func (h *runtimeHandler) Kill(ctx context.Context, req *proto.KillRequest, rsp *proto.KillResponse) error {
	if req.Process == nil {
		return errors.BadRequest(Name+".kill", "process is nil")
	}

	if req.Process.Binary == nil {
		return errors.BadRequest(Name+".kill", "process.Binary is nil")
	}

	if req.Process.Binary.Source == nil {
		return errors.BadRequest(Name+".kill", "process.Binary.Source is nil")
	}

	if err := h.r.Kill(&gorun.Process{
		ID: req.Process.Id,
		Binary: &gorun.Binary{
			Path: req.Process.Binary.Path,
			Source: &gorun.Source{
				URL: req.Process.Binary.Source.Url,
				Dir: req.Process.Binary.Source.Dir,
			},
		},
	}); err != nil {
		return errors.InternalServerError(Name+".kill", err.Error())
	}

	return nil
}

func (h *runtimeHandler) Wait(ctx context.Context, req *proto.WaitRequest, stream proto.Runtime_WaitStream) error {
	if req.Process == nil {
		return errors.BadRequest(Name+".wait", "process is nil")
	}

	if req.Process.Binary == nil {
		return errors.BadRequest(Name+".wait", "process.Binary is nil")
	}

	if req.Process.Binary.Source == nil {
		return errors.BadRequest(Name+".wait", "process.Binary.Source is nil")
	}

	if err := h.r.Wait(&gorun.Process{
		ID: req.Process.Id,
		Binary: &gorun.Binary{
			Path: req.Process.Binary.Path,
			Source: &gorun.Source{
				URL: req.Process.Binary.Source.Url,
				Dir: req.Process.Binary.Source.Dir,
			},
		},
	}); err != nil {
		if serr := stream.Send(&proto.WaitResponse{
			Error: err.Error(),
		}); serr != nil {
			return errors.InternalServerError(Name+".wait", serr.Error())
		}
	}

	return nil
}

func (h *serviceHandler) Run(ctx context.Context, req *proto.RunRequest, rsp *proto.RunResponse) error {
	if len(req.Url) == 0 {
		return errors.BadRequest(Name+".run", "url is blank")
	}

	// TODO: should return stream?
	go h.m.Run(req.Url, req.Restart, req.Update)

	// publish run event
	h.prun.Publish(ctx, req)

	return nil
}

func (h *serviceHandler) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	if len(req.Url) == 0 {
		return errors.BadRequest(Name+".status", "url is blank")
	}

	info, err := h.m.Status(req.Url)
	if err != nil {
		return errors.InternalServerError(Name+".status", err.Error())
	}
	rsp.Info = info
	return nil
}

func (h *serviceHandler) Stop(ctx context.Context, req *proto.StopRequest, rsp *proto.StopResponse) error {
	if len(req.Url) == 0 {
		return errors.BadRequest(Name+".stop", "url is blank")
	}

	// TODO: return the error?
	h.m.Stop(req.Url)

	// publish stop event
	h.pstop.Publish(ctx, req)

	return nil
}
