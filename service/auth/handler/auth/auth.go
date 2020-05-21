package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	"github.com/micro/go-micro/v2/auth/token"
	"github.com/micro/go-micro/v2/auth/token/basic"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	memStore "github.com/micro/go-micro/v2/store/memory"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/micro/v2/internal/namespace"
	"golang.org/x/crypto/bcrypt"
)

const (
	storePrefixAccounts      = "account/"
	storePrefixRefreshTokens = "refresh/"
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

	keys, err := a.Options.Store.List(store.ListPrefix(storePrefixAccounts), store.ListLimit(2))
	if err != nil {
		log.Errorf("Error listing accounts in init: %v", err)
		return
	}
	if len(keys) > 0 {
		log.Info("Accounts exists. Skipping account injection.")
		return
	}
	log.Info("Generating default account")
	resp := &pb.GenerateResponse{}
	err = a.Generate(context.Background(), &pb.GenerateRequest{
		Id:     "admin",
		Secret: "Password1",
	}, resp)
	if err != nil {
		log.Errorf("Error creating default account in init: %v", err)
	}
}

// Generate an account
func (a *Auth) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.auth", "ID required")
	}

	// set the defaults
	if len(req.Type) == 0 {
		req.Type = "user"
	}
	if len(req.Secret) == 0 {
		req.Secret = uuid.New().String()
	}

	// check the user does not already exist
	key := storePrefixAccounts + req.Id
	if _, err := a.Options.Store.Read(key); err != store.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Account with this ID already exists")
	}

	// hash the secret
	secret, err := hashSecret(req.Secret)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to hash password: %v", err)
	}

	// Default to the current namespace as the scope. Once we add identity we can auto-generate
	// these scopes and prevent users from generating accounts with any scope.
	if len(req.Scopes) == 0 {
		req.Scopes = []string{"namespace." + namespace.FromContext(ctx)}
	}

	// construct the account
	acc := &auth.Account{
		ID:       req.Id,
		Type:     req.Type,
		Scopes:   req.Scopes,
		Metadata: req.Metadata,
		Issuer:   namespace.FromContext(ctx),
		Secret:   secret,
	}

	// marshal to json
	bytes, err := json.Marshal(acc)
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to marshal json: %v", err)
	}

	// write to the store
	if err := a.Options.Store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to write account to store: %v", err)
	}

	// set a refresh token
	if err := a.setRefreshToken(acc.ID, uuid.New().String()); err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to set a refresh token: %v", err)
	}

	// return the account
	rsp.Account = serializeAccount(acc)
	rsp.Account.Secret = req.Secret // return unhashed secret
	return nil
}

// Inspect a token and retrieve the account
func (a *Auth) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	acc, err := a.TokenProvider.Inspect(req.Token)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("go.micro.auth", "Invalid token")
	} else if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to inspect token: %v", err)
	}

	rsp.Account = serializeAccount(acc)
	return nil
}

// Token generation using an account ID and secret
func (a *Auth) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	// validate the request
	if (len(req.Id) == 0 || len(req.Secret) == 0) && len(req.RefreshToken) == 0 {
		return errors.BadRequest("go.micro.auth", "Credentials or a refresh token required")
	}

	// Declare the account id and refresh token
	accountID := req.Id
	refreshToken := req.RefreshToken

	// If the refresh token is set, check this
	if len(req.RefreshToken) > 0 {
		accID, err := a.accountIDForRefreshToken(req.RefreshToken)
		if err == store.ErrNotFound {
			return errors.BadRequest("go.micro.auth", "Invalid token")
		} else if err != nil {
			return errors.InternalServerError("go.micro.auth", "Unable to lookup token: %v", err)
		}
		accountID = accID
	}

	// Lookup the account in the store
	key := storePrefixAccounts + accountID
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

	// If the refresh token was not used, validate the secrets match and then set the refresh token
	// so it can be returned to the user
	if len(req.RefreshToken) == 0 {
		if !secretsMatch(acc.Secret, req.Secret) {
			return errors.BadRequest("go.micro.auth", "Secret not correct")
		}

		refreshToken, err = a.refreshTokenForAccount(acc.ID)
		if err != nil {
			return errors.InternalServerError("go.micro.auth", "Unable to get refresh token: %v", err)
		}
	}

	// Generate a new access token
	duration := time.Duration(req.TokenExpiry) * time.Second
	tok, err := a.TokenProvider.Generate(acc, token.WithExpiry(duration))
	if err != nil {
		return errors.InternalServerError("go.micro.auth", "Unable to generate token: %v", err)
	}

	rsp.Token = serializeToken(tok, refreshToken)
	return nil
}

// set the refresh token for an account
func (a *Auth) setRefreshToken(id, token string) error {
	key := storePrefixRefreshTokens + id + "/" + token
	return a.Options.Store.Write(&store.Record{Key: key})
}

// get the refresh token for an accutn
func (a *Auth) refreshTokenForAccount(id string) (string, error) {
	recs, err := a.Options.Store.Read(storePrefixRefreshTokens+id+"/", store.ReadPrefix())
	if err != nil {
		return "", err
	} else if len(recs) != 1 {
		return "", store.ErrNotFound
	}

	comps := strings.Split(recs[0].Key, "/")
	if len(comps) != 3 {
		return "", store.ErrNotFound
	}
	return comps[2], nil
}

// get the account ID for the given refresh token
func (a *Auth) accountIDForRefreshToken(token string) (string, error) {
	keys, err := a.Options.Store.List(store.ListPrefix(storePrefixRefreshTokens))
	if err != nil {
		return "", err
	}
	for _, k := range keys {
		if strings.HasSuffix(k, "/"+token) {
			comps := strings.Split(k, "/")
			if len(comps) != 3 {
				return "", store.ErrNotFound
			}
			return comps[1], nil
		}
	}
	return "", store.ErrNotFound
}

func serializeAccount(a *auth.Account) *pb.Account {
	return &pb.Account{
		Id:       a.ID,
		Type:     a.Type,
		Scopes:   a.Scopes,
		Issuer:   a.Issuer,
		Metadata: a.Metadata,
	}
}

func serializeToken(t *token.Token, refresh string) *pb.Token {
	return &pb.Token{
		Created:      t.Created.Unix(),
		Expiry:       t.Expiry.Unix(),
		AccessToken:  t.Token,
		RefreshToken: refresh,
	}
}

func hashSecret(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func secretsMatch(hash string, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming) == nil
}
