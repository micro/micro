package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/micro/v2/service/config/client"
)

// DefaultConfig implementation
var DefaultConfig, _ = config.NewConfig(
	config.WithSource(client.NewSource()),
)
