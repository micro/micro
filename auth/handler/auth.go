package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto/auth"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/errors"
)

// Generate an account
func (h *Handler) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// Generate a long-lived secret
	secretOpts := []token.GenerateOption{
		token.WithExpiry(time.Duration(req.SecretExpiry) * time.Second),
		token.WithMetadata(req.Metadata),
		token.WithRoles(req.Roles...),
	}
	secret, err := h.SecretProvider.Generate(req.Id, secretOpts...)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate secret: %v", err)
	}

	// set the account
	rsp.Account = &pb.Account{
		Id:       req.Id,
		Metadata: req.Metadata,
		Roles:    req.Roles,
		Secret:   serializeToken(secret),
	}

	return nil
}

// Inspect a token and retrieve the account
func (h *Handler) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	tok, err := h.TokenProvider.Inspect(req.Token)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect token: %v", err)
	}

	rsp.Account = &pb.Account{
		Id:       tok.Subject,
		Roles:    tok.Roles,
		Metadata: tok.Metadata,
	}
	return nil
}

// Refresh a token using a secret
func (h *Handler) Refresh(ctx context.Context, req *pb.RefreshRequest, rsp *pb.RefreshResponse) error {
	sec, err := h.SecretProvider.Inspect(req.Secret)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect secret: %v", err)
	}

	tok, err := h.TokenProvider.Generate(sec.Subject,
		token.WithExpiry(time.Duration(req.TokenExpiry)*time.Second),
		token.WithMetadata(sec.Metadata),
		token.WithRoles(sec.Roles...),
	)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate token: %v", err)
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
