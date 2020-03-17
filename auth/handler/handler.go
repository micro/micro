package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/auth"

	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/go-micro/v2/auth/service/proto"
)

// New returns an instance of Handler
func New() *Handler {
	return &Handler{
		auth:  *cmd.DefaultOptions().Auth,
		store: *cmd.DefaultOptions().Store,
	}
}

var (
	// Duration the service account is valid for
	Duration = time.Hour * 24
)

// Handler processes RPC calls
type Handler struct {
	store store.Store
	auth  auth.Auth
}

// Generate creates a new  account in the store
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	if req.Account == nil {
		return errors.BadRequest("go.micro.auth", "account required")
	}
	if req.Account.Id == "" {
		return errors.BadRequest("go.micro.auth", "account id required")
	}

	opts := []auth.GenerateOption{}
	if req.Account.Metadata != nil {
		opts = append(opts, auth.Metadata(req.Account.Metadata))
	}
	if req.Account.Roles != nil {
		roles := make([]*auth.Role, len(req.Account.Roles))
		for i, r := range req.Account.Roles {
			var resouce *auth.Resource
			if r.Resource != nil {
				resouce = &auth.Resource{
					Name: r.Resource.Name,
					Type: r.Resource.Type,
				}
			}

			roles[i] = &auth.Role{
				Name:     r.Name,
				Resource: resouce,
			}
		}
		opts = append(opts, auth.Roles(roles))
	}

	acc, err := h.auth.Generate(req.Account.Id, opts...)
	if err != nil {
		return err
	}

	// encode the response
	rsp.Account = serializeAccount(acc)

	return nil
}

// Verify retrieves a  account from the store
func (h *Handler) Verify(ctx context.Context, req *pb.VerifyRequest, rsp *pb.VerifyResponse) error {
	if req.Token == "" {
		return errors.BadRequest("go.micro.auth", "token required")
	}

	acc, err := h.auth.Verify(req.Token)
	if err != nil {
		return err
	}

	rsp.Account = serializeAccount(acc)

	return nil
}

// Revoke deletes the  account
func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	if req.Token == "" {
		return errors.BadRequest("go.micro.auth", "token required")
	}

	return h.auth.Revoke(req.Token)
}

func serializeAccount(acc *auth.Account) *pb.Account {
	res := &pb.Account{
		Id:       acc.Id,
		Created:  acc.Created.Unix(),
		Expiry:   acc.Expiry.Unix(),
		Metadata: acc.Metadata,
		Token:    acc.Token,
		Roles:    make([]*pb.Role, len(acc.Roles)),
	}

	for i, r := range acc.Roles {
		res.Roles[i] = &pb.Role{Name: r.Name}

		if r.Resource != nil {
			res.Roles[i].Resource = &pb.Resource{
				Name: r.Resource.Name,
				Type: r.Resource.Type,
			}
		}
	}

	return res
}
