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

// Delete a resource
func (r *Runtime) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {

	// validate the request
	if req.Resource == nil || (req.Resource.Namespace == nil && req.Resource.Networkpolicy == nil && req.Resource.Service == nil) {
		return errors.BadRequest("runtime.Runtime.Delete", "blank resource")
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

	// Handle the different possible types of resource
	switch {
	case req.Resource.Namespace != nil:
		ns, err := gorun.NewNamespace(req.Resource.Namespace.Name)
		if err != nil {
			return err
		}

		if err := r.Runtime.Delete(ns); err != nil {
			return err
		}

		ev := &runtime.EventNamespacePayload{
			Type:      runtime.EventNamespaceDeleted,
			Namespace: ns.Name,
		}

		return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
			"type":      runtime.EventNamespaceDeleted,
			"namespace": ns.Name,
		}))

	case req.Resource.Networkpolicy != nil:
		np, err := gorun.NewNetworkPolicy(req.Resource.Networkpolicy.Name, req.Resource.Networkpolicy.Namespace, req.Resource.Networkpolicy.Allowedlabels)
		if err != nil {
			return err
		}

		if err := r.Runtime.Delete(np); err != nil {
			return err
		}

		ev := &runtime.EventNetworkPolicyPayload{
			Type:      runtime.EventNetworkPolicyDeleted,
			Name:      np.Name,
			Namespace: np.Namespace,
		}

		return events.Publish(runtime.EventTopic, ev, goevents.WithMetadata(map[string]string{
			"type":      ev.Type,
			"namespace": ev.Namespace,
		}))

	case req.Resource.Service != nil:

		// delete the service
		service := toService(req.Resource.Service)
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

	default:
		return nil
	}
}
