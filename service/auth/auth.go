package auth

import (
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/micro/v3/service/auth/client"
)

// DefaultAuth implementation
var DefaultAuth auth.Auth = client.NewAuth()
