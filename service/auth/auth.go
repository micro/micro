package auth

import (
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/auth/service"
)

// DefaultAuth implementation
var DefaultAuth auth.Auth = service.NewAuth()
