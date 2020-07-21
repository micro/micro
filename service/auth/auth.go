package auth

import (
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/go-micro/v2/auth/service/proto"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"
	"github.com/micro/go-micro/v2/util/token"
	"github.com/micro/go-micro/v2/util/token/jwt"
	authHandler "github.com/micro/micro/v2/service/auth/handler/auth"
	rulesHandler "github.com/micro/micro/v2/service/auth/handler/rules"
)

const (
	name    = "go.micro.auth"
	address = ":8010"
)

// Run the auth service
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	srvOpts = append([]micro.Option{
		micro.Name(name),
		micro.Address(address),
	}, srvOpts...)

	srv := micro.NewService(srvOpts...)

	// setup the handlers
	ruleH := &rulesHandler.Rules{}
	authH := &authHandler.Auth{}

	// setup the auth handler to use JWTs
	pubKey := ctx.String("auth_public_key")
	privKey := ctx.String("auth_private_key")
	if len(pubKey) > 0 || len(privKey) > 0 {
		authH.TokenProvider = jwt.NewTokenProvider(
			token.WithPublicKey(pubKey),
			token.WithPrivateKey(privKey),
		)
	}

	// set the handlers store
	srv.Options().Store.Init(store.Table("auth"))
	authH.Init(auth.Store(srv.Options().Store))
	ruleH.Init(auth.Store(srv.Options().Store))

	// register handlers
	pb.RegisterAuthHandler(srv.Server(), authH)
	pb.RegisterRulesHandler(srv.Server(), ruleH)
	pb.RegisterAccountsHandler(srv.Server(), authH)

	// run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
