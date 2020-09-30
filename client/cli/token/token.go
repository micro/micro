// Package token contains CLI client token related helpers
package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/service/auth"
)

// Get tries a best effort read of auth token from user config.
// Might have missing `RefreshToken` or `Expiry` fields in case of
// incomplete or corrupted user config.
func Get(envName, namespace string) (*auth.AccountToken, error) {
	tok, err := getFromFile(envName, namespace)
	if err == nil {
		return tok, nil
	}
	return getFromUserConfig(envName)
}

type token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// unix timestamp
	Created int64 `json:"created"`
	// unix timestamp
	Expiry int64 `json:"expiry"`
}

func tokensFilePath() string {
	return config.File + "-tokens.json"
}

func getFromFile(envName, namespace string) (*auth.AccountToken, error) {
	tokens, err := getTokens()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%v:%v", envName, namespace)
	tok, ok := tokens[key]
	if !ok {
		return nil, fmt.Errorf("Token not found under %v in file %v", key, tokensFilePath())
	}
	return &auth.AccountToken{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		Created:      time.Unix(tok.Created, 0),
		Expiry:       time.Unix(tok.Expiry, 0),
	}, nil
}

func getTokens() (map[string]token, error) {
	// @todo work on the path as `~/.micro/config.json-tokens` is not nice enough
	dat, err := ioutil.ReadFile(tokensFilePath())
	if err != nil {
		return nil, err
	}
	m := map[string]token{}
	return m, json.Unmarshal(dat, &m)
}

func saveTokens(m map[string]token) error {
	dat, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(tokensFilePath(), dat, 0700)
}

func getFromUserConfig(envName string) (*auth.AccountToken, error) {
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
func Save(envName, namespace string, token *auth.AccountToken) error {
	return saveToFile(envName, namespace, token)
}

func saveToFile(envName, namespace string, authToken *auth.AccountToken) error {
	tokens, err := getTokens()
	if err != nil {
		return err
	}
	tokens[fmt.Sprintf("%v:%v", envName, namespace)] = token{
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
		Created:      authToken.Created.Unix(),
		Expiry:       authToken.Expiry.Unix(),
	}
	return saveTokens(tokens)
}

func saveToUserConfig(envName string, token *auth.AccountToken) error {
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
func Remove(envName, namespace string) error {
	return removeFromFile(envName, namespace)
}

func removeFromFile(envName, namespace string) error {
	tokens, err := getTokens()
	if err != nil {
		return err
	}
	delete(tokens, fmt.Sprintf("%v:%v", envName, namespace))
	return saveTokens(tokens)
}

func removeFromUserConfig(envName string) error {
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
