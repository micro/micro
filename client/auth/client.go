package client

import (
	"strings"
	"sync"
	"time"

	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/auth/rules"
	"github.com/micro/micro/v3/util/auth/token"
	"github.com/micro/micro/v3/util/auth/token/jwt"
)

const (
	ruleCacheTTL = 2 * time.Minute
)

type rulesCache struct {
	sync.RWMutex
	ruleCache map[string]*cacheEntry
	ttl       time.Duration
}

func (r *rulesCache) get(key string) []*auth.Rule {
	r.RLock()
	entry := r.ruleCache[key]
	r.RUnlock()
	if entry != nil && time.Since(entry.t) < r.ttl {
		return entry.v
	}
	return nil
}

func (r *rulesCache) put(key string, v []*auth.Rule) {
	r.Lock()
	r.ruleCache[key] = &cacheEntry{t: time.Now(), v: v}
	r.Unlock()
}

type cacheEntry struct {
	t time.Time
	v []*auth.Rule
}

// srv is the service implementation of the Auth interface
type srv struct {
	options   auth.Options
	auth      pb.AuthService
	rules     pb.RulesService
	token     token.Provider
	ruleCache rulesCache
}

func (s *srv) String() string {
	return "service"
}

func (s *srv) Init(opts ...auth.Option) {
	for _, o := range opts {
		o(&s.options)
	}
	s.auth = pb.NewAuthService("auth", client.DefaultClient)
	s.rules = pb.NewRulesService("auth", client.DefaultClient)
	s.setupJWT()
	s.ruleCache = rulesCache{
		ruleCache: map[string]*cacheEntry{},
		ttl:       ruleCacheTTL,
	}
}

func (s *srv) Options() auth.Options {
	return s.options
}

// Generate a new account
func (s *srv) Generate(id string, opts ...auth.GenerateOption) (*auth.Account, error) {
	options := auth.NewGenerateOptions(opts...)
	if len(options.Issuer) == 0 {
		options.Issuer = s.options.Issuer
	}

	// we have the JWT private key and generate ourselves an account
	if len(s.options.PrivateKey) > 0 {
		acc := &auth.Account{
			ID:       id,
			Type:     options.Type,
			Scopes:   options.Scopes,
			Metadata: options.Metadata,
			Issuer:   options.Issuer,
		}

		tok, err := s.token.Generate(acc, token.WithExpiry(time.Hour*24*365))
		if err != nil {
			return nil, err
		}

		// when using JWTs, the account secret is the JWT's token. This
		// can be used as an argument in the Token method.
		acc.Secret = tok.Token
		return acc, nil
	}

	rsp, err := s.auth.Generate(context.DefaultContext, &pb.GenerateRequest{
		Id:       id,
		Type:     options.Type,
		Secret:   options.Secret,
		Scopes:   options.Scopes,
		Metadata: options.Metadata,
		Provider: options.Provider,
		Options: &pb.Options{
			Namespace: options.Issuer,
		},
		Name: options.Name,
	}, s.callOpts()...)
	if err != nil {
		return nil, err
	}

	return serializeAccount(rsp.Account), nil
}

// Grant access to a resource
func (s *srv) Grant(rule *auth.Rule) error {
	access := pb.Access_UNKNOWN
	if rule.Access == auth.AccessGranted {
		access = pb.Access_GRANTED
	} else if rule.Access == auth.AccessDenied {
		access = pb.Access_DENIED
	}

	_, err := s.rules.Create(context.DefaultContext, &pb.CreateRequest{
		Rule: &pb.Rule{
			Id:       rule.ID,
			Scope:    rule.Scope,
			Priority: rule.Priority,
			Access:   access,
			Resource: &pb.Resource{
				Type:     rule.Resource.Type,
				Name:     rule.Resource.Name,
				Endpoint: rule.Resource.Endpoint,
			},
		},
		Options: &pb.Options{
			Namespace: s.Options().Issuer,
		},
	}, s.callOpts()...)
	go s.refreshRulesCache(s.Options().Issuer)
	return err
}

// Revoke access to a resource
func (s *srv) Revoke(rule *auth.Rule) error {
	_, err := s.rules.Delete(context.DefaultContext, &pb.DeleteRequest{
		Id: rule.ID, Options: &pb.Options{
			Namespace: s.Options().Issuer,
		},
	}, s.callOpts()...)
	go s.refreshRulesCache(s.Options().Issuer)
	return err
}

func (s *srv) refreshRulesCache(ns string) error {
	rsp, err := s.rules.List(context.DefaultContext, &pb.ListRequest{
		Options: &pb.Options{Namespace: ns},
	}, s.callOpts()...)
	if err != nil {
		logger.Errorf("Error refreshing rules cache %s", err)
		return err
	}

	rules := make([]*auth.Rule, len(rsp.Rules))
	for i, r := range rsp.Rules {
		rules[i] = serializeRule(r)
	}
	s.ruleCache.put(ns, rules)
	return nil
}

func (s *srv) Rules(opts ...auth.RulesOption) ([]*auth.Rule, error) {
	var options auth.RulesOptions
	for _, o := range opts {
		o(&options)
	}
	if options.Context == nil {
		options.Context = context.DefaultContext
	}
	if len(options.Namespace) == 0 {
		options.Namespace = s.options.Issuer
	}

	if ret := s.ruleCache.get(options.Namespace); ret != nil {
		return ret, nil
	}
	if err := s.refreshRulesCache(options.Namespace); err != nil {
		return nil, err
	}

	return s.ruleCache.get(options.Namespace), nil
}

// Verify an account has access to a resource
func (s *srv) Verify(acc *auth.Account, res *auth.Resource, opts ...auth.VerifyOption) error {
	var options auth.VerifyOptions
	for _, o := range opts {
		o(&options)
	}

	rs, err := s.Rules(
		auth.RulesContext(options.Context),
		auth.RulesNamespace(options.Namespace),
	)
	if err != nil {
		return err
	}
	return rules.VerifyAccess(rs, acc, res, opts...)
}

// Inspect a token
func (s *srv) Inspect(token string) (*auth.Account, error) {
	// validate the request
	if len(token) == 0 {
		return nil, auth.ErrInvalidToken
	}

	// optimisation - is the key the right format for jwt auth?
	if s.token.String() == "jwt" && strings.Count(token, ".") != 2 {
		return nil, auth.ErrInvalidToken
	}

	// try to decode JWT locally and fall back to srv if an error occurs
	if len(strings.Split(token, ".")) == 3 && len(s.options.PublicKey) > 0 {
		return s.token.Inspect(token)
	}

	// the token is not a JWT or we do not have the keys to decode it,
	// fall back to the auth service
	rsp, err := s.auth.Inspect(context.DefaultContext, &pb.InspectRequest{
		Token: token, Options: &pb.Options{Namespace: s.Options().Issuer},
	}, s.callOpts()...)
	if err != nil {
		return nil, err
	}
	return serializeAccount(rsp.Account), nil
}

// Token generation using an account ID and secret
func (s *srv) Token(opts ...auth.TokenOption) (*auth.AccountToken, error) {
	options := auth.NewTokenOptions(opts...)
	if len(options.Issuer) == 0 {
		options.Issuer = s.options.Issuer
	}

	tok := options.RefreshToken
	if len(options.Secret) > 0 {
		tok = options.Secret
	}

	// we have the JWT private key and refresh accounts locally
	if len(s.options.PrivateKey) > 0 && len(strings.Split(tok, ".")) == 3 {
		acc, err := s.token.Inspect(tok)
		if err != nil {
			return nil, err
		}

		token, err := s.token.Generate(acc, token.WithExpiry(options.Expiry))
		if err != nil {
			return nil, err
		}

		return &auth.AccountToken{
			Expiry:       token.Expiry,
			AccessToken:  token.Token,
			RefreshToken: tok,
		}, nil
	}

	rsp, err := s.auth.Token(context.DefaultContext, &pb.TokenRequest{
		Id:           options.ID,
		Secret:       options.Secret,
		RefreshToken: options.RefreshToken,
		TokenExpiry:  int64(options.Expiry.Seconds()),
		Options: &pb.Options{
			Namespace: options.Issuer,
		},
	}, s.callOpts()...)
	if err != nil && errors.FromError(err).Detail == auth.ErrInvalidToken.Error() {
		return nil, auth.ErrInvalidToken
	} else if err != nil {
		return nil, err
	}

	return serializeToken(rsp.Token), nil
}

func serializeToken(t *pb.Token) *auth.AccountToken {
	return &auth.AccountToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		Created:      time.Unix(t.Created, 0),
		Expiry:       time.Unix(t.Expiry, 0),
	}
}

func serializeAccount(a *pb.Account) *auth.Account {
	return &auth.Account{
		ID:       a.Id,
		Secret:   a.Secret,
		Issuer:   a.Issuer,
		Metadata: a.Metadata,
		Scopes:   a.Scopes,
		Name:     a.Name,
		Type:     a.Type,
	}
}

func serializeRule(r *pb.Rule) *auth.Rule {
	var access auth.Access
	if r.Access == pb.Access_GRANTED {
		access = auth.AccessGranted
	} else {
		access = auth.AccessDenied
	}

	return &auth.Rule{
		ID:       r.Id,
		Scope:    r.Scope,
		Access:   access,
		Priority: r.Priority,
		Resource: &auth.Resource{
			Type:     r.Resource.Type,
			Name:     r.Resource.Name,
			Endpoint: r.Resource.Endpoint,
		},
	}
}

func (s *srv) callOpts() []client.CallOption {
	return []client.CallOption{
		client.WithAddress(s.options.Addrs...),
		client.WithAuthToken(),
	}
}

// NewAuth returns a new instance of the Auth service
func NewAuth(opts ...auth.Option) auth.Auth {
	service := &srv{
		auth:    pb.NewAuthService("auth", client.DefaultClient),
		rules:   pb.NewRulesService("auth", client.DefaultClient),
		options: auth.NewOptions(opts...),
	}

	service.setupJWT()

	return service
}

func (s *srv) setupJWT() {
	tokenOpts := []token.Option{}

	// if we have a JWT public key passed as an option,
	// we can decode tokens with the type "JWT" locally
	// and not have to make an RPC call
	if key := s.options.PublicKey; len(key) > 0 {
		tokenOpts = append(tokenOpts, token.WithPublicKey(key))
	}

	// if we have a JWT private key passed as an option,
	// we can generate accounts locally and not have to make
	// an RPC call, this is used for micro clients such as
	// api, web, proxy.
	if key := s.options.PrivateKey; len(key) > 0 {
		tokenOpts = append(tokenOpts, token.WithPrivateKey(key))
	}

	s.token = jwt.NewTokenProvider(tokenOpts...)
}
