package handler

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/auth"
	"github.com/micro/go-micro/util/log"

	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/store"

	pb "github.com/micro/go-micro/auth/service/proto"
)

// New returns an instance of Handler
func New() *Handler {
	return &Handler{
		store: *cmd.DefaultOptions().Store,
	}
}

var (
	// Duration is how long until the service account can no longer be used as auth
	Duration = time.Hour * 24 * 365
)

// Handler processes RPC calls
type Handler struct {
	store store.Store
}

// Generate creates a new service account in the store
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if req.ServiceAccount == nil {
		return errors.BadRequest("go.micro.auth", "service account required")
	}

	parent := req.ServiceAccount.Parent
	if parent == nil {
		return errors.BadRequest("go.micro.auth", "parent required")
	}
	if parent.Id == "" || parent.Type == "" {
		return errors.BadRequest("go.micro.auth", "invalid parent")
	}

	// generate the token
	token, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// key for the store
	key := fmt.Sprintf("%v/%v", prefixForResource(parent), token.String())

	// construct the service account
	sa := auth.ServiceAccount{
		Created:  time.Now(),
		Expiry:   time.Now().Add(Duration),
		Metadata: req.ServiceAccount.Metadata,
	}

	// add the roles
	sa.Roles = make([]*auth.Role, len(req.ServiceAccount.Roles))
	for i, r := range req.ServiceAccount.Roles {
		sa.Roles[i] = &auth.Role{Name: r.Name}

		if r.Resource != nil {
			sa.Roles[i].Resource = &auth.Resource{
				Id:   r.Resource.Id,
				Type: r.Resource.Type,
			}
		}
	}

	// encode the data to bytes
	buf := &bytes.Buffer{}
	e := gob.NewEncoder(buf)
	if err := e.Encode(sa); err != nil {
		return err
	}

	// write to the store
	err = h.store.Write(&store.Record{
		Key:    key,
		Value:  buf.Bytes(),
		Expiry: Duration,
	})
	if err != nil {
		return err
	}
	log.Infof("Created service account: %v", key)

	// encode the response
	rsp.ServiceAccount = &pb.ServiceAccount{
		Created:  sa.Created.Unix(),
		Expiry:   sa.Expiry.Unix(),
		Metadata: sa.Metadata,
		Token:    token.String(),
		Roles:    req.ServiceAccount.Roles,
	}

	return nil
}

// Validate retrieves a token from the store
func (h *Handler) Validate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	return nil
}

func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	return nil
}
func (h *Handler) AddRole(ctx context.Context, req *pb.AddRoleRequest, rsp *pb.AddRoleResponse) error {
	return nil
}
func (h *Handler) RemoveRole(ctx context.Context, req *pb.RemoveRoleRequest, rsp *pb.RemoveRoleResponse) error {
	return nil
}

// prefixForResource is used is the store's key name, e.g. user/asim@micro.mu || service/go.micro.srv.auth
func prefixForResource(r *pb.Resource) string {
	return fmt.Sprintf("%v/%v", r.Type, r.Id)
}

// // DEBUG
// records, err := h.store.List()
// if err != nil {
// 	return err
// }
// for _, r := range records {
// 	b := bytes.NewBuffer(r.Value)
// 	d := gob.NewDecoder(b)
// 	var f auth.ServiceAccount
// 	err = d.Decode(&f)
// 	if err == nil {
// 		fmt.Println(r.Key)
// 		fmt.Println(f)
// 	}
// }
