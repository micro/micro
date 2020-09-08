package server

import (
	"context"
	"time"

	goauth "github.com/micro/go-micro/v3/auth"
	goevents "github.com/micro/go-micro/v3/events"
	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/namespace"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
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

func (r *Runtime) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("runtime.Runtime.Create", "blank service")
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
		return errors.Forbidden("runtime.Runtime.Create", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.Create", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.Create", err.Error())
	}

	// create the service
	service := toService(req.Service)
	setupServiceMeta(ctx, service)

	options := toCreateOptions(ctx, req.Options)

	log.Infof("Creating service %s version %s source %s", service.Name, service.Version, service.Source)
	if err := r.Runtime.Create(service, options...); err != nil {
		return errors.InternalServerError("runtime.Runtime.Create", err.Error())
	}

	// publish the create event
	ev := &runtime.EventPayload{
		Service:   service,
		Namespace: req.Options.Namespace,
		Type:      runtime.EventServiceCreated,
	}

	return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
		"type":      runtime.EventServiceCreated,
		"namespace": req.Options.Namespace,
	}))
}

func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("runtime.Runtime.Update", "blank service")
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
		return errors.Forbidden("runtime.Runtime.Update", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.Update", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.Update", err.Error())
	}

	service := toService(req.Service)
	setupServiceMeta(ctx, service)

	options := toUpdateOptions(ctx, req.Options)

	log.Infof("Updating service %s version %s source %s", service.Name, service.Version, service.Source)

	if err := r.Runtime.Update(service, options...); err != nil {
		return errors.InternalServerError("runtime.Runtime.Update", err.Error())
	}

	// publish the update event
	ev := &runtime.EventPayload{
		Service:   service,
		Namespace: req.Options.Namespace,
		Type:      runtime.EventServiceUpdated,
	}

	return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
		"type":      runtime.EventServiceUpdated,
		"namespace": req.Options.Namespace,
	}))
}

func setupServiceMeta(ctx context.Context, service *gorun.Service) {
	if service.Metadata == nil {
		service.Metadata = map[string]string{}
	}
	account, accOk := goauth.AccountFromContext(ctx)
	if accOk {
		service.Metadata["owner"] = account.ID
		// This is a hack - we don't want vanilla `micro server` users where the auth is noop
		// to have long uuid as owners, so we put micro here - not great, not terrible.
		if auth.DefaultAuth.String() == "noop" {
			service.Metadata["owner"] = "micro"
		}
		service.Metadata["group"] = account.Issuer
	}
	service.Metadata["started"] = time.Now().Format(time.RFC3339)
}

func (r *Runtime) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if req.Service == nil {
		return errors.BadRequest("runtime.Runtime.Delete", "blank service")
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
		return errors.Forbidden("runtime.Runtime.Delete", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.Delete", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.Delete", err.Error())
	}

	// delete the service
	service := toService(req.Service)
	options := toDeleteOptions(ctx, req.Options)

	log.Infof("Deleting service %s version %s source %s", service.Name, service.Version, service.Source)
	if err := r.Runtime.Delete(service, options...); err != nil {
		return errors.InternalServerError("runtime.Runtime.Delete", err.Error())
	}

	// publish the delete event
	ev := &runtime.EventPayload{
		Type:      runtime.EventServiceDeleted,
		Namespace: req.Options.Namespace,
		Service:   service,
	}

	return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
		"type":      runtime.EventServiceDeleted,
		"namespace": req.Options.Namespace,
	}))
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
		opts = append(opts, gorun.LogsCount(req.GetCount()))
	}
	if req.GetStream() {
		opts = append(opts, gorun.LogsStream(req.GetStream()))
	}

	logStream, err := r.Runtime.Logs(&gorun.Service{
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

func (r *Runtime) CreateNamespace(ctx context.Context, req *pb.CreateNamespaceRequest, rsp *pb.CreateNamespaceResponse) error {
	// authorize the request, only admins/core services should be able to call
	if err := namespace.Authorize(ctx, namespace.DefaultNamespace); err == namespace.ErrForbidden {
		return errors.Forbidden("runtime.Runtime.CreateNamespace", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.CreateNamespace", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.CreateNamespace", err.Error())
	}

	if err := r.Runtime.CreateNamespace(req.Namespace); err != nil {
		return err
	}

	ev := &runtime.EventNamespacePayload{
		Type:      runtime.EventNamespaceCreated,
		Namespace: req.Namespace,
	}
	return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
		"type":      runtime.EventNamespaceCreated,
		"namespace": req.Namespace,
	}))
}

func (r *Runtime) DeleteNamespace(ctx context.Context, req *pb.DeleteNamespaceRequest, rsp *pb.DeleteNamespaceResponse) error {
	// authorize the request, only admins/core services should be able to call
	if err := namespace.Authorize(ctx, namespace.DefaultNamespace); err == namespace.ErrForbidden {
		return errors.Forbidden("runtime.Runtime.DeleteNamespace", err.Error())
	} else if err == namespace.ErrUnauthorized {
		return errors.Unauthorized("runtime.Runtime.DeleteNamespace", err.Error())
	} else if err != nil {
		return errors.InternalServerError("runtime.Runtime.DeleteNamespace", err.Error())
	}

	if err := r.Runtime.DeleteNamespace(req.Namespace); err != nil {
		return err
	}

	ev := &runtime.EventNamespacePayload{
		Type:      runtime.EventNamespaceDeleted,
		Namespace: req.Namespace,
	}
	return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
		"type":      runtime.EventNamespaceDeleted,
		"namespace": req.Namespace,
	}))
}
