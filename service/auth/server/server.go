package server

import (
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/go-micro/v3/util/token"
	"github.com/micro/go-micro/v3/util/token/jwt"
	"github.com/micro/micro/v3/internal/user"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service"
	authHandler "github.com/micro/micro/v3/service/auth/server/auth"
	rulesHandler "github.com/micro/micro/v3/service/auth/server/rules"
	"github.com/micro/micro/v3/service/logger"
	log "github.com/micro/micro/v3/service/logger"
	mustore "github.com/micro/micro/v3/service/store"
	"github.com/urfave/cli/v2"
)

const (
	name    = "auth"
	address = ":8010"
)

// Run the auth service
func Run(ctx *cli.Context) error {
	srv := service.New(
		service.Name(name),
		service.Address(address),
	)

	// setup the handlers
	ruleH := &rulesHandler.Rules{}
	authH := &authHandler.Auth{
		DisableAdmin: ctx.Bool("disable_admin"),
	}

	// setup the auth handler to use JWTs
	pubKey := ctx.String("auth_public_key")
	privKey := ctx.String("auth_private_key")
	if len(privKey) == 0 || len(pubKey) == 0 {
		privB, pubB, err := user.GetJWTCerts()
		if err != nil {
			logger.Fatalf("Error getting keys; %v", err)
		}
		privKey = string(privB)
		pubKey = string(pubB)
	}

	authH.TokenProvider = jwt.NewTokenProvider(
		token.WithPublicKey(pubKey),
		token.WithPrivateKey(privKey),
	)

	// set the handlers store
	mustore.DefaultStore.Init(store.Table("auth"))
	authH.Init(auth.Store(mustore.DefaultStore))
	ruleH.Init(auth.Store(mustore.DefaultStore))

	// register handlers
	pb.RegisterAuthHandler(srv.Server(), authH)
	pb.RegisterRulesHandler(srv.Server(), ruleH)
	pb.RegisterAccountsHandler(srv.Server(), authH)

	// run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}
