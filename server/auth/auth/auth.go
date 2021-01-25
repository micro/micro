package auth

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	authns "github.com/micro/micro/v3/internal/auth/namespace"
	"github.com/micro/micro/v3/internal/auth/token"
	"github.com/micro/micro/v3/internal/auth/token/basic"
	"github.com/micro/micro/v3/internal/namespace"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"golang.org/x/crypto/bcrypt"
)

const (
	joinKey                  = "/"
	storePrefixAccounts      = "account"
	storePrefixRefreshTokens = "refresh"

	// used to enable login with username rather than ID (username can change e.g. email, id is stable)
	storePrefixAccountsByName = "accountByName"
)

var defaultAccount = auth.Account{
	ID:     "admin",
	Type:   "user",
	Scopes: []string{"admin"},
	Secret: "micro",
}

// Auth processes RPC calls
type Auth struct {
	Options       auth.Options
	TokenProvider token.Provider

	namespaces map[string]bool
	sync.Mutex
	// Prevent the generation of default accounts
	DisableAdmin bool
}

// Init the auth
func (a *Auth) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&a.Options)
	}

	// setup a token provider
	if a.TokenProvider == nil {
		a.TokenProvider = basic.NewTokenProvider(token.WithStore(store.DefaultStore))
	}
}

func (a *Auth) setupDefaultAccount(ns string) error {
	if ns != namespace.DefaultNamespace {
		return nil
	}
	if a.DisableAdmin {
		return nil
	}
	a.Lock()
	defer a.Unlock()

	// setup the namespace cache if not yet done
	if a.namespaces == nil {
		a.namespaces = make(map[string]bool)
	}

	// check to see if the default account has already been verified
	if _, ok := a.namespaces[ns]; ok {
		return nil
	}

	// check to see if we need to create the default account
	prefix := strings.Join([]string{storePrefixAccounts, ns, ""}, joinKey)
	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return err
	}

	hasUser := false
	for _, rec := range recs {
		acc := &auth.Account{}
		err := json.Unmarshal(rec.Value, acc)
		if err != nil {
			return err
		}
		if acc.Type == "user" {
			hasUser = true
			break
		}
	}

	// create the account if none exist in the namespace
	if !hasUser {
		acc := defaultAccount
		acc.Issuer = ns
		if err := a.createAccount(&acc); err != nil {
			return err
		}
	}

	// set the namespace in the cache
	a.namespaces[ns] = true
	return nil
}

// Generate an account
func (a *Auth) Generate(ctx context.Context, req *pb.GenerateRequest, rsp *pb.GenerateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("auth.Auth.Generate", "ID required")
	}

	// set the defaults
	if len(req.Type) == 0 {
		req.Type = "user"
	}
	if len(req.Secret) == 0 {
		req.Secret = uuid.New().String()
	}
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.FromContext(ctx)
	}

	// authorize the request
	if err := authns.AuthorizeAdmin(ctx, req.Options.Namespace, "auth.Auth.Generate"); err != nil {
		return err
	}

	// check the user does not already exists
	key := strings.Join([]string{storePrefixAccounts, req.Options.Namespace, req.Id}, joinKey)
	if _, err := store.Read(key); err != store.ErrNotFound {
		return errors.BadRequest("auth", "Account with this ID already exists")
	}

	// construct the account
	acc := &auth.Account{
		ID:       req.Id,
		Type:     req.Type,
		Scopes:   req.Scopes,
		Metadata: req.Metadata,
		Issuer:   req.Options.Namespace,
		Secret:   req.Secret,
		Name:     req.Name,
	}

	// create the account
	if err := a.createAccount(acc); err != nil {
		return err
	}

	// return the account
	rsp.Account = serializeAccount(acc)
	rsp.Account.Secret = req.Secret // return unhashed secret
	return nil
}
func (a *Auth) createAccount(acc *auth.Account) error {
	// check the user does not already exists
	key := strings.Join([]string{storePrefixAccounts, acc.Issuer, acc.ID}, joinKey)
	if _, err := store.Read(key); err != store.ErrNotFound {
		return errors.BadRequest("auth.Auth.Generate", "Account with this ID already exists")
	}
	if acc.Metadata == nil {
		acc.Metadata = map[string]string{}
	}
	if acc.Name == "" {
		acc.Name = acc.ID
	}
	usernameKey := strings.Join([]string{storePrefixAccountsByName, acc.Issuer, acc.Name}, joinKey)
	if _, err := store.Read(usernameKey); err != store.ErrNotFound {
		return errors.BadRequest("auth.Auth.Generate", "Account with this Name already exists")
	}

	// hash the secret
	secret, err := hashSecret(acc.Secret)
	if err != nil {
		return errors.InternalServerError("auth.Auth.Generate", "Unable to hash password: %v", err)
	}
	acc.Secret = secret

	// marshal to json
	bytes, err := json.Marshal(acc)
	if err != nil {
		return errors.InternalServerError("auth.Auth.Generate", "Unable to marshal json: %v", err)
	}

	// write to the store
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.InternalServerError("auth.Auth.Generate", "Unable to write account to store: %v", err)
	}

	if err := store.Write(&store.Record{Key: usernameKey, Value: bytes}); err != nil {
		return errors.InternalServerError("auth.Auth.Generate", "Unable to write account to store: %v", err)
	}

	// set a refresh token
	if err := a.setRefreshToken(acc.Issuer, acc.ID, uuid.New().String()); err != nil {
		return errors.InternalServerError("auth.Auth.Generate", "Unable to set a refresh token: %v", err)
	}

	return nil
}

// Inspect a token and retrieve the account
func (a *Auth) Inspect(ctx context.Context, req *pb.InspectRequest, rsp *pb.InspectResponse) error {
	acc, err := a.TokenProvider.Inspect(req.Token)
	if err == token.ErrInvalidToken || err == token.ErrNotFound {
		return errors.BadRequest("auth.Auth.Inspect", err.Error())
	} else if err != nil {
		return errors.InternalServerError("auth.Auth.Inspect", "Unable to inspect token: %v", err)
	}

	rsp.Account = serializeAccount(acc)
	return nil
}

// Token generation using an account ID and secret
func (a *Auth) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	// set defaults
	if req.Options == nil {
		req.Options = &pb.Options{}
	}
	if len(req.Options.Namespace) == 0 {
		req.Options.Namespace = namespace.DefaultNamespace
	}

	// setup the defaults incase none exist
	err := a.setupDefaultAccount(req.Options.Namespace)
	if err != nil {
		// failing gracefully here
		logger.Errorf("Error setting up default accounts: %v", err)
	}

	// validate the request
	if (len(req.Id) == 0 || len(req.Secret) == 0) && len(req.RefreshToken) == 0 {
		return errors.BadRequest("auth.Auth.Token", "Credentials or a refresh token required")
	}

	// check to see if the secret is a JWT. this is a workaround to allow accounts issued
	// by the runtime to be refreshed whilst keeping the private key in the server.
	if a.TokenProvider.String() == "jwt" {
		jwt := req.Secret
		if len(req.RefreshToken) > 0 {
			jwt = req.RefreshToken
		}

		if acc, err := a.TokenProvider.Inspect(jwt); err == nil {
			expiry := time.Duration(int64(time.Second) * req.TokenExpiry)
			tok, _ := a.TokenProvider.Generate(acc, token.WithExpiry(expiry))
			rsp.Token = serializeToken(tok, tok.Token)
			return nil
		}
	}

	// Declare the account id and refresh token
	accountID := req.Id
	refreshToken := req.RefreshToken

	// If the refresh token is set, check this
	if len(req.RefreshToken) > 0 {
		accID, err := a.accountIDForRefreshToken(req.Options.Namespace, req.RefreshToken)
		if err == store.ErrNotFound {
			return errors.BadRequest("auth.Auth.Token", auth.ErrInvalidToken.Error())
		} else if err != nil {
			return errors.InternalServerError("auth.Auth.Token", "Unable to lookup token: %v", err)
		}
		accountID = accID
	}

	acc, err := a.getAccountForID(accountID, req.Options.Namespace, "auth.Auth.Token")
	if err != nil {
		return err
	}

	// If the refresh token was not used, validate the secrets match and then set the refresh token
	// so it can be returned to the user
	if len(req.RefreshToken) == 0 {
		if !secretsMatch(acc.Secret, req.Secret) {
			return errors.BadRequest("auth.Auth.Token", "Secret not correct")
		}

		refreshToken, err = a.refreshTokenForAccount(req.Options.Namespace, acc.ID)
		if err != nil {
			return errors.InternalServerError("auth.Auth.Token", "Unable to get refresh token: %v", err)
		}
	}

	// Generate a new access token
	duration := time.Duration(req.TokenExpiry) * time.Second
	tok, err := a.TokenProvider.Generate(acc, token.WithExpiry(duration))
	if err != nil {
		return errors.InternalServerError("auth.Auth.Token", "Unable to generate token: %v", err)
	}

	rsp.Token = serializeToken(tok, refreshToken)
	return nil
}

func (a *Auth) getAccountForID(id, namespace, errCode string) (*auth.Account, error) {
	// Lookup the account in the store
	key := strings.Join([]string{storePrefixAccounts, namespace, id}, joinKey)
	recs, err := store.Read(key)
	if err != nil && err != store.ErrNotFound {
		return nil, errors.InternalServerError(errCode, "Unable to read from store: %v", err)
	}
	if err == store.ErrNotFound {
		// maybe id is the username and not the actual ID
		key = strings.Join([]string{storePrefixAccountsByName, namespace, id}, joinKey)
		recs, err = store.Read(key)
		if err == store.ErrNotFound {
			return nil, errors.BadRequest(errCode, "Account not found with this ID")
		}
		if err != nil {
			return nil, errors.InternalServerError(errCode, "Unable to read from store: %v", err)
		}
	}

	// Unmarshal the record
	var acc *auth.Account
	if err := json.Unmarshal(recs[0].Value, &acc); err != nil {
		return nil, errors.InternalServerError(errCode, "Unable to unmarshal account: %v", err)
	}
	return acc, nil
}

// set the refresh token for an account
func (a *Auth) setRefreshToken(ns, id, token string) error {
	key := strings.Join([]string{storePrefixRefreshTokens, ns, id, token}, joinKey)
	return store.Write(&store.Record{Key: key})
}

// get the refresh token for an account
func (a *Auth) refreshTokenForAccount(ns, id string) (string, error) {
	prefix := strings.Join([]string{storePrefixRefreshTokens, ns, id, ""}, joinKey)
	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return "", err
	} else if len(recs) == 0 {
		return "", store.ErrNotFound
	}

	comps := strings.Split(recs[0].Key, "/")
	if len(comps) != 4 {
		return "", store.ErrNotFound
	}
	return comps[3], nil
}

// get the account ID for the given refresh token
func (a *Auth) accountIDForRefreshToken(ns, token string) (string, error) {
	prefix := strings.Join([]string{storePrefixRefreshTokens, ns}, joinKey)
	keys, err := store.List(store.ListPrefix(prefix))
	if err != nil {
		return "", err
	}

	for _, k := range keys {
		if strings.HasSuffix(k, "/"+token) {
			comps := strings.Split(k, "/")
			if len(comps) != 4 {
				return "", store.ErrNotFound
			}
			return comps[2], nil
		}
	}

	return "", store.ErrNotFound
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
