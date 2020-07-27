// Package token contains CLI client token related helpers
package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v2/internal/config"
)

// Get tries a best effort read of auth token from user config.
// Might have missing `RefreshToken` or `Expiry` fields in case of
// incomplete or corrupted user config.
func Get(envName string) (*auth.Token, error) {
	path := []string{"micro", "auth", envName}
	accessToken, _ := config.Get(append(path, "token")...)

	refreshToken, err := config.Get(append(path, "refresh-token")...)
	if err != nil {
		// Gracefully degrading here in case the user only has a temporary access token at hand.
		// The call will fail on the receiving end.
		return &auth.Token{
			AccessToken: accessToken,
		}, nil
	}

	// See if the access token has expired
	expiry, _ := config.Get(append(path, "expiry")...)
	if len(expiry) == 0 {
		return &auth.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
	expiryInt, err := strconv.ParseInt(expiry, 10, 64)
	if err != nil {
		return nil, err
	}
	return &auth.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiry:       time.Unix(expiryInt, 0),
	}, nil
}

// Save saves the auth token to the user's local config file
func Save(envName string, token *auth.Token) error {
	if err := config.Set(token.AccessToken, "micro", "auth", envName, "token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	if err := config.Set(token.RefreshToken, "micro", "auth", envName, "refresh-token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	return config.Set(fmt.Sprintf("%v", token.Expiry.Unix()), "micro", "auth", envName, "expiry")
}

// Remove deletes a token. Useful when trying to reset test
// for example at testing: not having a token is a different state
// than having an invalid token.
func Remove(envName string) error {
	if err := config.Set("", "micro", "auth", envName, "token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	if err := config.Set("", "micro", "auth", envName, "refresh-token"); err != nil {
		return err
	}
	// Store the refresh token in micro config
	return config.Set("", "micro", "auth", envName, "expiry")
}
