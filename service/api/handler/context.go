package handler

import (
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/client"
)

type Context interface {
	Client() client.Client
	Service() *api.Service
	Domain() string
}
