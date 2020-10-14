package server

import (
	"context"
	"time"

	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/runtime"
)

type Runtime struct {
	Runtime gorun.Runtime
}

func (r *Runtime) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.ReadOptions{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("runtime.Runtime.Read", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.Read", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.Read", err.Error())
	}

	// lookup the services
	options := toReadOptions(ctx, req.Options)
	services, err := r.Runtime.Read(options...)
	if err != nil {
		return errors.InternalServerError("runtime.Runtime.Read", err.Error())
	}

	// serialize the response
	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func setupServiceMeta(ctx context.Context, service *runtime.Service) {
	if service.Metadata == nil {
		service.Metadata = map[string]string{}
	}
	account, accOk := auth.AccountFromContext(ctx)
	if accOk {
		// Try to use the account name as it's more user friendly. If none, fall back to ID
		owner := account.Name
		if len(owner) == 0 {
			owner = account.ID
		}
		service.Metadata["owner"] = owner
		// This is a hack - we don't want vanilla `micro server` users where the auth is noop
		// to have long uuid as owners, so we put micro here - not great, not terrible.
		if auth.DefaultAuth.String() == "noop" {
			service.Metadata["owner"] = "micro"
		}
		service.Metadata["group"] = account.Issuer
	}
	service.Metadata["started"] = time.Now().Format(time.RFC3339)
}

func (r *Runtime) Logs(ctx context.Context, req *pb.LogsRequest, stream pb.Runtime_LogsStream) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.LogsOptions{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("runtime.Runtime.Logs", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.Logs", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.Logs", err.Error())
	}

	opts := toLogsOptions(ctx, req.Options)

	// options passed in the request
	if req.GetCount() > 0 {
		opts = append(opts, runtime.LogsCount(req.GetCount()))
	}
	if req.GetStream() {
		opts = append(opts, runtime.LogsStream(req.GetStream()))
	}

	logStream, err := r.Runtime.Logs(&runtime.Service{
		Name: req.GetService(),
	}, opts...)
	if err != nil {
		return err
	}
	defer logStream.Stop()
	defer stream.Close()

	recordChan := logStream.Chan()
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				return logStream.Error()
			}
			// send record
			if err := stream.Send(&pb.LogRecord{
				//Timestamp: record.Timestamp.Unix(),
				Metadata: record.Metadata,
				Message:  record.Message,
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
