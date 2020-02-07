package main

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/util/log"

	pb "github.com/micro/micro/v2/auth/api/proto"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.auth"),
		micro.Version("latest"),
	)
	service.Init()

	pb.RegisterAuthHandler(service.Server(), NewHandler(service))

	if err := service.Run(); err != nil {
		log.Error(err)
	}
}

// Handler is an impementation of the auth api
type Handler struct {
	auth auth.Auth
}

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{auth: auth.DefaultAuth}
}

// Validate gets a token and verifies it with the auth package
func (h *Handler) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if len(req.Token) == 0 {
		return errors.BadRequest("go.micro.api.auth", "token required")
	}

	_, err := h.auth.Validate(req.Token)
	return err
}
