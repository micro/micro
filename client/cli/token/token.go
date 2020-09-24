// Package token contains CLI client token related helpers
package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/service/auth"
)

// Get tries a best effort read of auth token from user config.
// Might have missing `RefreshToken` or `Expiry` fields in case of
// incomplete or corrupted user config.
func Get(envName string) (*auth.AccountToken, error) {
	path := []string{"micro", "auth", envName}
	accessToken, _ := config.Get(config.Path(append(path, "token")...))

	refreshToken, err := config.Get(config.Path(append(path, "refresh-token")...))
	if err != nil {
		// Gracefully degrading here in case the user only has a temporary access token at hand.
		// The call will fail on the receiving end.
		return &auth.AccountToken{
			AccessToken: accessToken,
		}, nil
	}

	// See if the access token has expired
	expiry, _ := config.Get(config.Path(append(path, "expiry")...))
	if len(expiry) == 0 {
		return &auth.AccountToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
	expiryInt, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return nil, err
	}
	return &auth.AccountToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       time.Unix(expiryInt, 0),
	}, nil
}

// Save saves the auth token to the user's local config file
func Save(envName string, token *auth.AccountToken) error {
	if err := config.Set(config.Path("micro", "auth", envName, "token"), token.AccessToken); err != nil {
		return err
	}
	// Store the refresh token in micro config
	if err := config.Set(config.Path("micro", "auth", envName, "refresh-token"), token.RefreshToken); err != nil {
		return err
	}
	// Store the refresh token in micro config
	return config.Set(config.Path("micro", "auth", envName, "expiry"), fmt.Sprintf("%v", token.Expiry.Unix()))
}

// Remove deletes a token. Useful when trying to reset test
// for example at testing: not having a token is a different state
// than having an invalid token.
func Remove(envName string) error {
	if err := config.Set(config.Path("micro", "auth", envName, "token"), ""); err != nil {
		return err
	}
	// Store the refresh token in micro config
	if err := config.Set(config.Path("micro", "auth", envName, "refresh-token"), ""); err != nil {
		return err
	}
	// Store the refresh token in micro config
	return config.Set(config.Path("micro", "auth", envName, "expiry"), "")
}
