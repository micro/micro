package manager

import (
	"strings"

	gorun "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/runtime"
)

// createServiceInRuntime will add all the required env vars and secrets and then create the service
func (m *manager) createServiceInRuntime(srv *service) error {
	// generate an auth account for the service to use
	acc, err := m.generateAccount(srv)
	if err != nil {
		return err
	}

	// construct the options
	options := []gorun.CreateOption{
		gorun.CreateImage(srv.Options.Image),
		gorun.CreateType(srv.Options.Type),
		gorun.CreateNamespace(srv.Options.Namespace),
		gorun.WithArgs(srv.Options.Args...),
		gorun.WithCommand(srv.Options.Command...),
		gorun.WithEnv(m.runtimeEnv(srv.Service, srv.Options)),
	}

	// add the secrets
	for key, value := range srv.Options.Secrets {
		options = append(options, gorun.WithSecret(key, value))
	}

	// inject the credentials into the service if present
	if len(acc.ID) > 0 && len(acc.Secret) > 0 {
		options = append(options, gorun.WithSecret("MICRO_AUTH_ID", acc.ID))
		options = append(options, gorun.WithSecret("MICRO_AUTH_SECRET", acc.Secret))
	}

	// create the service
	return m.Runtime.Create(srv.Service, options...)
}

// runtimeEnv returns the environment variables which should  be used when creating a service.
func (m *manager) runtimeEnv(srv *gorun.Service, options *gorun.CreateOptions) []string {
	setEnv := func(p []string, env map[string]string) {
		for _, v := range p {
			parts := strings.Split(v, "=")
			if len(parts) <= 1 {
				continue
			}
			env[parts[0]] = strings.Join(parts[1:], "=")
		}
	}

	// overwrite any values
	env := map[string]string{
		// ensure a profile for the services isn't set, they
		// should use the default RPC clients
		"MICRO_PROFILE": "",
		// pass the service's name and version
		"MICRO_SERVICE_NAME":    srv.Name,
		"MICRO_SERVICE_VERSION": srv.Version,
		// set the proxy for the service to use (e.g. micro network)
		// using the proxy which has been configured for the runtime
		"MICRO_PROXY": client.DefaultClient.Options().Proxy,
	}

	// bind to port 8080, this is what the k8s tcp readiness check will use
	if runtime.DefaultRuntime.String() == "kubernetes" {
		env["MICRO_SERVICE_ADDRESS"] = ":8080"
	}

	// set the env vars provided
	setEnv(options.Env, env)

	// set the service namespace
	if len(options.Namespace) > 0 {
		env["MICRO_NAMESPACE"] = options.Namespace
	}

	// create a new env
	var vars []string
	for k, v := range env {
		vars = append(vars, k+"="+v)
	}

	// setup the runtime env
	return vars
}

func (m *manager) generateAccount(srv *service) (*auth.Account, error) {
	accName := srv.Service.Name + "-" + srv.Service.Version

	opts := []auth.GenerateOption{
		auth.WithIssuer(srv.Options.Namespace),
		auth.WithScopes("service"),
		auth.WithType("service"),
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
