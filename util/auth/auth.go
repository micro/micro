package auth

import (
	"time"

	"micro.dev/v4/service/auth"
	"micro.dev/v4/service/logger"
)

const (
	// BearerScheme used for Authorization header
	BearerScheme = "Bearer "
	// TokenCookieName is the name of the cookie which stores the auth token
	TokenCookieName = "micro-token"
)

// SystemRules are the default rules which are applied to the runtime services
var SystemRules = []*auth.Rule{
	&auth.Rule{
		ID:       "default",
		Scope:    auth.ScopePublic,
		Access:   auth.AccessGranted,
		Resource: &auth.Resource{Type: "*", Name: "*", Endpoint: "*"},
	},
}

func RefreshToken() {
	// can't refresh a token we don't have
	if auth.DefaultAuth.Options().Token == nil {
		return
	}

	t := time.NewTicker(time.Minute)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// don't refresh the token if it's not close to expiring
			tok := auth.DefaultAuth.Options().Token
			if tok.Expiry.Unix() > time.Now().Add(time.Hour).Unix() {
				continue
			}

			// generate the first token
			tok, err := auth.Token(
				auth.WithToken(tok.RefreshToken),
				auth.WithExpiry(time.Minute*10),
			)
			if err == auth.ErrInvalidToken {
				logger.Warnf("[Auth] Refresh token expired, regenerating using account credentials")

				tok, err = auth.Token(
					auth.WithCredentials(
						auth.DefaultAuth.Options().ID,
						auth.DefaultAuth.Options().Secret,
					),
					auth.WithExpiry(time.Hour*24),
				)
			} else if err != nil {
				logger.Warnf("[Auth] Error refreshing token: %v", err)
				continue
			}

			// set the token
			logger.Debugf("Auth token refreshed, expires at %v", tok.Expiry.Format(time.UnixDate))
			auth.DefaultAuth.Init(auth.ClientToken(tok))
		}
	}
}
