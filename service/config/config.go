package config

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/service"
)

// DefaultConfig implementation
var DefaultConfig, _ = config.NewConfig(
	config.WithSource(service.NewSource()),
)
