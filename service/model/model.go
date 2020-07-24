package model

import (
	"github.com/micro/go-micro/v2/model"
	"github.com/micro/go-micro/v2/model/mud"
)

var (
	// DefaultModel for the service
	DefaultModel model.Model = mud.NewModel()
)
