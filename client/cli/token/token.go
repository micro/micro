// Package token contains CLI client token related helpers
// tToken files consist of one line per token, each token having
// the structure of `micro://envAddress/namespace[/id]:token`, ie.
// micro://m3o.com/foo-bar-baz/asim@aslam.me:afsafasfasfaceevqcCEWVEWV
// or
// micro://m3o.com/foo-bar-baz:afsafasfasfaceevqcCEWVEWV
package token

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/service/auth"
	"github.com/urfave/cli/v2"
)

const tokenPath = "-tokens.json"

// Get tries a best effort read of auth token from user config.
// Might have missing `RefreshToken` or `Expiry` fields in case of
// incomplete or corrupted user config.
func Get(ctx *cli.Context) (*auth.AccountToken, error) {
	env, err := util.GetEnv(ctx)
	if err != nil {
		return nil, err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return nil, err
	}
	tok, err := getFromFile(env.Name, ns)
	if err == nil {
		return tok, nil
	}
	return getFromUserConfig(env.Name)
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
	if strings.HasSuffix(config.File, ".json") {
		return strings.ReplaceAll(config.File, ".json", "") + tokenPath
	}
	return config.File + tokenPath
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
	lines := strings.Split(string(dat), "\n")
	ret := map[string]token{}
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}
		key := strings.Join(parts[0:2], ":")
		base64Encoded := parts[3]
		jsonMarshalled, err := base64.StdEncoding.DecodeString(base64Encoded)
		if err != nil {
			return nil, err
		}
		tok := token{}
		err = json.Unmarshal(jsonMarshalled, &tok)
		if err != nil {
			return nil, err
		}
		ret[key] = tok
	}
	return ret, nil
}

func saveTokens(tokens map[string]token) error {
	buf := bytes.NewBuffer([]byte{})
	for key, t := range tokens {
		marshalledToken, err := json.Marshal(t)
		if err != nil {
			return err
		}
		base64Token := base64.StdEncoding.EncodeToString(marshalledToken)
		_, err = buf.WriteString(key + ":" + base64Token)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(tokensFilePath(), buf.Bytes(), 0700)
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
func Save(ctx *cli.Context, token *auth.AccountToken) error {
	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}
	return saveToFile(env.Name, ns, token)
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
func Remove(ctx *cli.Context) error {
	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return err
	}
	return removeFromFile(env.Name, ns)
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
