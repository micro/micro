package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/errors"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/config/cmd"
)

// New returns an instance of Handler
func New() *Handler {
	return &Handler{
		auth: *cmd.DefaultOptions().Auth,
	}
}

// Handler processes RPC calls
type Handler struct {
	auth auth.Auth
}

// Generate an account
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	acc, err := h.auth.Generate(req.Id,
		auth.WithRoles(req.Roles),
		auth.WithMetadata(req.Metadata),
		auth.WithSecretExpiry(time.Duration(req.SecretExpiry)),
	)
	if err != nil {
		return err
	}

	rsp.Account = serializeAccount(acc)
	return nil
}

// Grant a role access to a resource
func (h *Handler) Grant(ctx context.Context, req *pb.GrantRequest, rsp *pb.GrantResponse) error {
	if req.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}

	return h.auth.Grant(req.Role, &auth.Resource{
		Type:     req.Resource.Type,
		Name:     req.Resource.Name,
		Endpoint: req.Resource.Endpoint,
	})
}

// Revoke a roles access to a resource
func (h *Handler) Revoke(ctx context.Context, req *pb.RevokeRequest, rsp *pb.RevokeResponse) error {
	if req.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}

	return h.auth.Revoke(req.Role, &auth.Resource{
		Type:     req.Resource.Type,
		Name:     req.Resource.Name,
		Endpoint: req.Resource.Endpoint,
	})
}

// Verify an account has access to a resource
func (h *Handler) Verify(ctx context.Context, req *pb.VerifyRequest, rsp *pb.VerifyResponse) error {
	if req.Account == nil {
		return errors.BadRequest("go.micro.auth", "Account missing")
	}
	if req.Resource == nil {
		return errors.BadRequest("go.micro.auth", "Resource missing")
	}

	return h.auth.Verify(
		&auth.Account{
			ID:    req.Account.Id,
			Roles: req.Account.Roles,
		},
		&auth.Resource{
			Type:     req.Resource.Type,
			Name:     req.Resource.Name,
			Endpoint: req.Resource.Endpoint,
		},
	)
}

// Inspect a token and retrieve the account
func (h *Handler) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	acc, err := h.auth.Inspect(req.Token)
	if err != nil {
		return err
	}
	rsp.Account = serializeAccount(acc)
	return nil
}

// Refresh a token using a secret
func (h *Handler) Refresh(ctx context.Context, req *pb.RefreshRequest, rsp *pb.RefreshResponse) error {
	tok, err := h.auth.Refresh(req.Secret)
	if err != nil {
		return err
	}
	rsp.Token = serializeToken(tok)
	return nil
}

func serializeAccount(a *auth.Account) *pb.Account {
	var secret *pb.Token
	if a.Secret != nil {
		secret = serializeToken(a.Secret)
	}

	return &pb.Account{
		Id:       a.ID,
		Roles:    a.Roles,
		Metadata: a.Metadata,
		Secret:   secret,
	}
}

func serializeToken(t *auth.Token) *pb.Token {
	return &pb.Token{
		Token:    t.Token,
		Type:     t.Type,
		Created:  t.Created.Unix(),
		Expiry:   t.Expiry.Unix(),
		Subject:  t.Subject,
		Roles:    t.Roles,
		Metadata: t.Metadata,
	}
}
