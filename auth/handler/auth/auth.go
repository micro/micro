package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto/auth"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/basic"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
)

const (
	storePrefix = "accounts/"
)

// Auth processes RPC calls
type Auth struct {
	Options        auth.Options
	SecretProvider token.Provider
	TokenProvider  token.Provider
}

// Init the auth
func (a *Auth) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&a.Options)
	}

	// use the default store as a fallback
	if a.Options.Store == nil {
		a.Options.Store = store.DefaultStore
	}

	// noop will not work for auth
	if a.Options.Store.String() == "noop" {
		a.Options.Store = memStore.NewStore()
	}

	if a.TokenProvider == nil {
		a.TokenProvider = basic.NewTokenProvider(token.WithStore(a.Options.Store))
	}
	if a.SecretProvider == nil {
		a.SecretProvider = basic.NewTokenProvider(token.WithStore(a.Options.Store))
	}
}

// Generate an account
func (a *Auth) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// Generate a long-lived secret
	secretOpts := []token.GenerateOption{
		token.WithExpiry(time.Duration(req.SecretExpiry) * time.Second),
		token.WithNamespace(req.Namespace),
		token.WithMetadata(req.Metadata),
		token.WithRoles(req.Roles...),
	}
	secret, err := a.SecretProvider.Generate(req.Id, secretOpts...)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate secret: %v", err)
	}

	// construct the account
	acc := &pb.Account{
		Id:        req.Id,
		Metadata:  req.Metadata,
		Roles:     req.Roles,
		Namespace: req.Namespace,
		Secret:    serializeToken(secret),
	}

	// marshal to json
	bytes, err := json.Marshal(acc)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to marshal json: %v", err)
	}

	// write to the store
	key := storePrefix + acc.Id
	if err := a.Options.Store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to write account to store: %v", err)
	}

	// return the account
	rsp.Account = acc
	return nil
}

// Inspect a token and retrieve the account
func (a *Auth) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	tok, err := a.TokenProvider.Inspect(req.Token)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect token: %v", err)
	}

	rsp.Account = &pb.Account{
		Id:        tok.Subject,
		Roles:     tok.Roles,
		Metadata:  tok.Metadata,
		Namespace: tok.Namespace,
	}
	return nil
}

// Refresh a token using a secret
func (a *Auth) Refresh(ctx context.Context, req *pb.RefreshRequest, rsp *pb.RefreshResponse) error {
	sec, err := a.SecretProvider.Inspect(req.Secret)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect secret: %v", err)
	}

	tok, err := a.TokenProvider.Generate(sec.Subject,
		token.WithExpiry(time.Duration(req.TokenExpiry)*time.Second),
		token.WithMetadata(sec.Metadata),
		token.WithRoles(sec.Roles...),
		token.WithNamespace(sec.Namespace),
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
		Id:        a.ID,
		Roles:     a.Roles,
		Metadata:  a.Metadata,
		Namespace: a.Namespace,
		Secret:    secret,
	}
}

func serializeToken(t *auth.Token) *pb.Token {
	return &pb.Token{
		Token:     t.Token,
		Type:      t.Type,
		Created:   t.Created.Unix(),
		Expiry:    t.Expiry.Unix(),
		Subject:   t.Subject,
		Roles:     t.Roles,
		Metadata:  t.Metadata,
		Namespace: t.Namespace,
	}
}
