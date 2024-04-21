package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	"micro.dev/v4/cmd/client/util"
	"micro.dev/v4/service/auth"
	"micro.dev/v4/util/namespace"
)

func createToken(ctx *cli.Context) error {
	env, err := util.GetEnv(ctx)
	if err != nil {
		return err
	}
	ns, err := namespace.Get(env.Name)
	if err != nil {
		return fmt.Errorf("Error getting namespace: %v", err)
	}
	if len(ctx.String("namespace")) > 0 {
		ns = ctx.String("namespace")
	}

	id := ctx.String("id")
	secret := ctx.String("secret")
	expiry := ctx.Int("expiry")

	if len(id) == 0 {
		return errors.New("Missing account ID")
	}
	if len(secret) == 0 {
		return errors.New("Missing account secret")
	}

	options := []auth.TokenOption{auth.WithTokenIssuer(ns)}
	options = append(options, auth.WithCredentials(id, secret))

	if expiry > 0 {
		options = append(options, auth.WithExpiry(time.Second*time.Duration(expiry)))
	}

	token, err := auth.Token(options...)
	if err != nil {
		return fmt.Errorf("Error creating token: %v", err)
	}

	json, _ := json.Marshal(token)
	fmt.Printf("Token created: %v\n", string(json))
	return nil
}
