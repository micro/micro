package server

import (
	"github.com/urfave/cli/v2"
	pb "micro.dev/v4/proto/auth"
	"micro.dev/v4/service"
	"micro.dev/v4/service/auth"
	"micro.dev/v4/service/auth/handler"
	log "micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"
	mustore "micro.dev/v4/service/store"
	"micro.dev/v4/util/auth/token"
	"micro.dev/v4/util/auth/token/jwt"
)

const (
	address = ":8010"
)

// Run the auth service
func Run(ctx *cli.Context) error {
	srv := service.New(
		service.Name("auth"),
		service.Address(address),
	)

	// setup the handlers
	ruleH := &handler.Rules{}
	authH := &handler.Auth{}

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
