package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/service/logger"
)

// generate an auth account for the service
func generateCreds(issuer string) error {
	acc, err := Generate(
		uuid.New().String(),
		auth.WithType("service"),
		auth.WithScopes("service"),
		auth.WithIssuer(issuer),
	)
	if err != nil {
		return err
	}
	if logger.V(logger.DebugLevel, logger.DefaultLogger) {
		logger.Debugf("Auth [%v] Generated an auth account", DefaultAuth.String())
	}

	// generate the first token
	token, err := Token(
		auth.WithCredentials(acc.ID, acc.Secret),
		auth.WithExpiry(time.Minute*10),
	)
	if err != nil {
		return err
	}

	// set the credentials and token in auth options
	DefaultAuth.Init(
		auth.ClientToken(token),
		auth.Credentials(acc.ID, acc.Secret),
	)
	return nil
}

// authUsingCreds generates auth credentials for the service
func authUsingCreds(issuer, id, secret string) error {
	// generate the first token
	token, err := Token(
		auth.WithCredentials(id, secret),
		auth.WithExpiry(time.Minute*10),
		auth.WithTokenIssuer(issuer),
	)
	if err != nil {
		return err
	}

	// set the credentials and token in options
	DefaultAuth.Init(
		auth.ClientToken(token),
		auth.Credentials(id, secret),
	)
	return nil
}

// refreshAccessToken if it is close to expiring
func refreshAccessToken() {
	t := time.NewTicker(time.Second * 15)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// don't refresh the token if it's not close to expiring
			tok := DefaultAuth.Options().Token
			if tok == nil {
				return
			}
			if tok.Expiry.Unix() > time.Now().Add(time.Minute).Unix() {
				continue
			}

			// generate the first token
			tok, err := Token(
				auth.WithToken(tok.RefreshToken),
				auth.WithExpiry(time.Minute*10),
			)
			if err != nil {
				logger.Warnf("[Auth] Error refreshing token: %v", err)
				continue
			}

			// set the token
			DefaultAuth.Init(auth.ClientToken(tok))
		}
	}
}
