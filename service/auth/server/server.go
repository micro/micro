package server

import (
	"github.com/micro/micro/v3/internal/auth/token"
	"github.com/micro/micro/v3/internal/auth/token/jwt"
	pb "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	authHandler "github.com/micro/micro/v3/service/auth/server/auth"
	rulesHandler "github.com/micro/micro/v3/service/auth/server/rules"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	mustore "github.com/micro/micro/v3/service/store"
	"github.com/urfave/cli/v2"
)

// Flags specific to the router
var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "disable_admin",
		EnvVars: []string{"MICRO_AUTH_DISABLE_ADMIN"},
		Usage:   "Prevent generation of default accounts in namespaces",
	},
}

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
	authH.TokenProvider = jwt.NewTokenProvider(
		token.WithPublicKey(auth.DefaultAuth.Options().PublicKey),
		token.WithPrivateKey(auth.DefaultAuth.Options().PrivateKey),
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
