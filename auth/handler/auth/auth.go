package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
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
	Options       auth.Options
	TokenProvider token.Provider
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

	// setup a token provider
	if a.TokenProvider == nil {
		a.TokenProvider = basic.NewTokenProvider(token.WithStore(a.Options.Store))
	}
}

// Generate an account
func (a *Auth) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// construct the account
	acc := &pb.Account{
		Id:        req.Id,
		Metadata:  req.Metadata,
		Roles:     req.Roles,
		Namespace: req.Namespace,
		Secret:    uuid.New().String(),
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

// Token generation using an account ID and secret
func (a *Auth) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.auth", "ID required")
	}
	if len(req.Secret) == 0 {
		return errors.BadRequest("go.micro.auth", "Secret required")
	}

	// Lookup the account in the store
	key := storePrefix + req.Id
	recs, err := a.Options.Store.Read(key)
	if err == store.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Account not found with this ID")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to read from store: %v", err)
	}

	// Unmarshal the record
	var acc *auth.Account
	if err := json.Unmarshal(recs[0].Value, &acc); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to unmarshal account: %v", err)
	}

	// Check the secret
	if acc.Secret != req.Secret {
		return errors.BadRequest("go.micro.auth", "Secret not correct")
	}

	// Generate a new token
	tok, err := a.TokenProvider.Generate(acc.ID,
		token.WithExpiry(time.Duration(req.TokenExpiry)*time.Second),
		token.WithMetadata(acc.Metadata),
		token.WithRoles(acc.Roles...),
		token.WithNamespace(acc.Namespace),
	)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate token: %v", err)
	}

	rsp.Token = serializeToken(tok)
	return nil
}

func serializeAccount(a *auth.Account) *pb.Account {
	return &pb.Account{
		Id:        a.ID,
		Roles:     a.Roles,
		Metadata:  a.Metadata,
		Namespace: a.Namespace,
		Secret:    a.Secret,
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
