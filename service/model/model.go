package model

import (
	"github.com/micro/go-micro/v3/model"
	"github.com/micro/go-micro/v3/model/mud"
)

var (
	// DefaultModel for the service
	DefaultModel model.Model = mud.NewModel()
)
