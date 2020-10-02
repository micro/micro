package server

import (
	"context"

	goevents "github.com/micro/go-micro/v3/events"
	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/internal/auth/namespace"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// Update a resource
func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {

	// validate the request
	if req.Resource == nil || (req.Resource.Namespace == nil && req.Resource.Networkpolicy == nil && req.Resource.Service == nil) {
		return errors.BadRequest("runtime.Runtime.Update", "blank resource")
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

	// Handle the different possible types of resource
	switch {
	case req.Resource.Namespace != nil:
		// No updates to namespace
		return nil

	case req.Resource.Networkpolicy != nil:
		np, err := gorun.NewNetworkPolicy(req.Resource.Networkpolicy.Name, req.Resource.Networkpolicy.Namespace, req.Resource.Networkpolicy.Allowedlabels)
		if err != nil {
			return err
		}

		if err := r.Runtime.Update(np); err != nil {
			return err
		}

		ev := &runtime.EventNetworkPolicyPayload{
			Type:      runtime.EventNetworkPolicyUpdated,
			Name:      np.Name,
			Namespace: np.Namespace,
		}

		return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
			"type":      ev.Type,
			"namespace": ev.Namespace,
		}))

	case req.Resource.Service != nil:

		service := toService(req.Resource.Service)
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

	default:
		return nil
	}
}
