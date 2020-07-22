package auth

import (
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/micro/v2/service/auth/client"
)

// DefaultAuth implementation
var DefaultAuth auth.Auth = client.NewAuth()
