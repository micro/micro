package handler

import (
	"micro.dev/v4/service/api"
	"micro.dev/v4/service/client"
)

type Context interface {
	Client() client.Client
	Service() *api.Service
	Domain() string
}
