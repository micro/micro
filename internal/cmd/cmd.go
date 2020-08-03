package cmd

import (
	"github.com/micro/cli/v2"
)

var InitFuncs []cli.ActionFunc

func Init(fn cli.ActionFunc) {
	InitFuncs = append(InitFuncs, fn)
}
