package main

import (
	"context"
	"errors"
	"strings"

	"go-micro.dev/v5"
	"go-micro.dev/v5/auth"
	"go-micro.dev/v5/metadata"
	"go-micro.dev/v5/server"
)

var (
	catchallResource = &auth.Resource{
		Type:     "*",
		Name:     "*",
		Endpoint: "*",
	}

	callResource = &auth.Resource{
		Type:     "service",
		Name:     "my.service.name",
		Endpoint: "Helloworld.MyMethod",
	}

	// Rules to validate against
	rules = []*auth.Rule{
		// Enforce auth on all endpoints.
		{Scope: "*", Resource: catchallResource},

		// Enforce auth on one specific endpoint.
		{Scope: "*", Resource: callResource},
	}
)

func NewAuthWrapper(service micro.Service) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// Fetch metadata from context (request headers).
			md, b := metadata.FromContext(ctx)
			if !b {
				return errors.New("no metadata found")
			}

			// Get auth header.
			authHeader, ok := md["Authorization"]
			if !ok || !strings.HasPrefix(authHeader, auth.BearerScheme) {
				return errors.New("no auth token provided")
			}

			// Extract auth token.
			token := strings.TrimPrefix(authHeader, auth.BearerScheme)

			// Extract account from token.
			a := service.Options().Auth
			acc, err := a.Inspect(token)
			if err != nil {
				return errors.New("auth token invalid")
			}

			// Create resource for current endpoint from request headers.
			currentResource := auth.Resource{
				Type:     "service",
				Name:     md["Micro-Service"],
				Endpoint: md["Micro-Endpoint"],
			}

			// Verify if account has access.
			if err := auth.Verify(rules, acc, &currentResource); err != nil {
				return errors.New("no access")
			}

			return h(ctx, req, rsp)
		}
	}
}
