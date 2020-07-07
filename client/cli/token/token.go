// Package token contains user side config (eg. `~/.micro`) token helpers
package token

import (
	"fmt"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/micro/v2/internal/config"
)


// SaveToken saves the auth token to the user's local config file
func SaveToken(envName string, token *auth.Token) error {
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