package server

import (
	"context"
	"time"

	"github.com/micro/go-micro/v3/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/namespace"
	"github.com/micro/micro/v3/service"
	pb "github.com/micro/micro/v3/service/runtime/proto"
)

type Runtime struct {
	Runtime runtime.Runtime
	Event   *service.Event
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
		return errors.Forbidden("go.micro.runtime", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.runtime", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// lookup the services
	options := toReadOptions(ctx, req.Options)
	services, err := r.Runtime.Read(options...)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// serialize the response
	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func (r *Runtime) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.CreateOptions{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.runtime", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.runtime", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// create the service
	service := toService(req.Service)
	options := toCreateOptions(ctx, req.Options)

	log.Infof("Creating service %s version %s source %s", service.Name, service.Version, service.Source)
	if err := r.Runtime.Create(service, options...); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the create event
	r.Event.Publish(ctx, &pb.Event{
		Type:      "create",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.UpdateOptions{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.runtime", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.runtime", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	service := toService(req.Service)
	options := toUpdateOptions(ctx, req.Options)

	log.Infof("Updating service %s version %s source %s", service.Name, service.Version, service.Source)

	if err := r.Runtime.Update(service, options...); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the update event
	r.Event.Publish(ctx, &pb.Event{
		Type:      "update",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// set defaults
	if req.Options == nil {
		req.Options = &pb.DeleteOptions{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// authorize the request
	if err := namespace.Authorize(ctx, req.Options.Namespace); err == namespace.ErrForbidden {
		return errors.Forbidden("go.micro.runtime", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.runtime", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// delete the service
	service := toService(req.Service)
	options := toDeleteOptions(ctx, req.Options)

	log.Infof("Deleting service %s version %s source %s", service.Name, service.Version, service.Source)
	if err := r.Runtime.Delete(service, options...); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the delete event
	r.Event.Publish(ctx, &pb.Event{
		Type:      "delete",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
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
		return errors.Forbidden("go.micro.runtime", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("go.micro.runtime", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
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
