package server

import (
	pb "github.com/micro/micro/v5/proto/auth"
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/auth"
	"github.com/micro/micro/v5/service/auth/handler"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	mustore "github.com/micro/micro/v5/service/store"
	"github.com/micro/micro/v5/util/auth/token"
	"github.com/micro/micro/v5/util/auth/token/jwt"
	"github.com/urfave/cli/v2"
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
