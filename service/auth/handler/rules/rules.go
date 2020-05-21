package roles

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/internal/namespace"
)

const (
	storePrefix = "rules"
	joinKey     = "/"
)

var defaultRule = &pb.Rule{
	Id:       "default",
	Role:     "", // a blank role  allows public access
	Priority: 0,
	Resource: &pb.Resource{
		Name:     "*",
		Type:     "*",
		Endpoint: "*",
	},
	Access: pb.Access_GRANTED,
}

// Rules processes RPC calls
type Rules struct {
	Options auth.Options
}

// Init the auth
func (r *Rules) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&r.Options)
	}

	// use the default store as a fallback
	if r.Options.Store == nil {
		r.Options.Store = store.DefaultStore
	}

	// noop will not work for auth
	if r.Options.Store.String() == "noop" {
		r.Options.Store = memStore.NewStore()
	}

	resp := &pb.ListResponse{}
	err := r.List(context.Background(), &pb.ListRequest{}, resp)
	if err != nil {
		log.Errorf("Error listing rules in init: %v", err)
		return
	}
	if len(resp.GetRules()) > 0 {
		log.Info("Rules exists. Skipping rule injection.")
		return
	}
	log.Info("Generating default rules")
	err = r.Create(context.Background(), &pb.CreateRequest{
		Rule: defaultRule,
	}, &pb.CreateResponse{})
	if err != nil {
		log.Errorf("Error creating default rule in init: %v", err)
	}
}

// Create a role access to a resource
func (r *Rules) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Validate the request
	if req.Rule == nil {
		return errors.BadRequest("go.micro.auth", "Rule missing")
	}
	if len(req.Rule.Id) == 0 {
		return errors.BadRequest("go.micro.auth", "ID missing")
	}
	if req.Rule.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}
	if req.Rule.Access == pb.Access_UNKNOWN {
		return errors.BadRequest("go.micro.auth", "Access missing")
	}

	// Chck the rule doesn't exist
	ns := namespace.FromContext(ctx)
	key := strings.Join([]string{storePrefix, ns, req.Rule.Id}, joinKey)
	if _, err := r.Options.Store.Read(key); err == nil {
		return errors.BadRequest("go.micro.auth", "A rule with this ID already exists")
	}

	// Encode the rule
	bytes, err := json.Marshal(req.Rule)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to marshal rule: %v", err)
	}

	// Write to the store
	if err := r.Options.Store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to write to the store: %v", err)
	}

	return nil
}

// Delete a roles access to a resource
func (r *Rules) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Validate the request
	if req.Rule == nil {
		return errors.BadRequest("go.micro.auth", "Rule missing")
	}
	if len(req.Rule.Id) == 0 {
		return errors.BadRequest("go.micro.auth", "ID missing")
	}

	// Delete the rule
	ns := namespace.FromContext(ctx)
	key := strings.Join([]string{storePrefix, ns, req.Rule.Id}, joinKey)
	err := r.Options.Store.Delete(key)
	if err == store.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Rule not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to delete key from store: %v", err)
	}

	return nil
}

// List returns all the rules
func (r *Rules) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// get the records from the store
	ns := namespace.FromContext(ctx)
	prefix := strings.Join([]string{storePrefix, ns, ""}, joinKey)
	fmt.Println(prefix)
	recs, err := r.Options.Store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to read from store: %v", err)
	}

	// unmarshal the records
	rsp.Rules = make([]*pb.Rule, 0, len(recs))
	for _, rec := range recs {
		var r *pb.Rule
		if err := json.Unmarshal(rec.Value, &r); err != nil {
			return errors.InternalServerError("go.micro.auth", "Error to unmarshaling json: %v. Value: %v", err, string(rec.Value))
		}
		rsp.Rules = append(rsp.Rules, r)
	}

	return nil
}
