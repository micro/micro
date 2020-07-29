package manager

import (
	goauth "github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/logger"
)

func (m *manager) generateAccount(srv *runtime.Service, ns string) (*goauth.Account, error) {
	accName := srv.Name + "-" + srv.Version

	opts := []goauth.GenerateOption{
		goauth.WithIssuer(ns),
		goauth.WithScopes("service"),
		goauth.WithType("service"),
	}

	acc, err := auth.Generate(accName, opts...)
	if err != nil {
		if logger.V(logger.WarnLevel, logger.DefaultLogger) {
			logger.Warnf("Error generating account %v: %v", accName, err)
		}
		return nil, err
	}
	if logger.V(logger.DebugLevel, logger.DefaultLogger) {
		logger.Debugf("Generated auth account: %v, secret: [len: %v]", acc.ID, len(acc.Secret))
	}

	return acc, nil
}
