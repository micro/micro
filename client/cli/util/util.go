// Package cliutil contains methods used across all cli commands
// @todo: get rid of os.Exits and use errors instread
package util

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	merrors "github.com/micro/micro/v3/service/errors"

	"github.com/micro/micro/v3/internal/config"
	"github.com/urfave/cli/v2"
)

const (
	// EnvLocal is a builtin environment, it represents your local `micro server`
	EnvLocal = "local"
	// EnvDev is a builtin staging / dev environment in the cloud
	EnvDev = "dev"
	// EnvPlatform is a builtin highly available environment in the cloud,
	EnvPlatform = "platform"
)

const (
	// localProxyAddress is the default proxy address for environment server
	localProxyAddress = "127.0.0.1:8081"
	// deprecated dev env
	devProxyAddress = "proxy.m3o.dev"
	// platformProxyAddress is the default proxy address for environment platform
	platformProxyAddress = "proxy.m3o.com"
)

var (
	// list of services managed
	// TODO: make use server/server list
	services = []string{
		// runtime services
		"network",  // :8085 (peer), :8443 (proxy)
		"runtime",  // :8088
		"registry", // :8000
		"config",   // :8001
		"store",    // :8002
		"broker",   // :8003
		"router",   // :8084
		"auth",     // :8010
		"proxy",    // :8081
		"api",      // :8080
		"events",
	}
)

var defaultEnvs = map[string]Env{
	EnvLocal: {
		Name:         EnvLocal,
		ProxyAddress: localProxyAddress,
		Description:  "Local running Micro Server",
	},
	EnvDev: {
		Name:         EnvDev,
		ProxyAddress: devProxyAddress,
		Description:  "Deprecated: Please use platform environment",
	},
	EnvPlatform: {
		Name:         EnvPlatform,
		ProxyAddress: platformProxyAddress,
		Description:  "Cloud hosted Micro Platform",
	},
}

func IsBuiltInService(command string) bool {
	for _, service := range services {
		if command == service {
			return true
		}
	}
	return false
}

// CLIProxyAddress returns the proxy address which should be set for the client
func CLIProxyAddress(ctx *cli.Context) (string, error) {
	switch ctx.Args().First() {
	case "new", "server", "help", "env":
		return "", nil
	}

	// fix for "micro service [command]", e.g "micro service auth"
	if ctx.Args().First() == "service" && IsBuiltInService(ctx.Args().Get(1)) {
		return "", nil
	}

	// don't set the proxy address on the proxy
	if ctx.Args().First() == "proxy" {
		return "", nil
	}

	env, err := GetEnv(ctx)
	if err != nil {
		return "", err
	}
	addr := env.ProxyAddress
	if !strings.Contains(addr, ":") {
		return fmt.Sprintf("%v:443", addr), nil
	}
	return addr, nil
}

type Env struct {
	Name         string
	ProxyAddress string
	Description  string
}

func AddEnv(env Env) error {
	envs, err := getEnvs()
	if err != nil {
		return err
	}
	envs[env.Name] = env
	return setEnvs(envs)
}

func getEnvs() (map[string]Env, error) {
	envsJSON, err := config.Get("envs")
	if err != nil {
		return nil, fmt.Errorf("Error getting environment: %v", err)
	}
	envs := map[string]Env{}
	if len(envsJSON) > 0 {
		err := json.Unmarshal([]byte(envsJSON), &envs)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range defaultEnvs {
		envs[k] = v
	}
	return envs, nil
}

func setEnvs(envs map[string]Env) error {
	envsJSON, err := json.Marshal(envs)
	if err != nil {
		return err
	}
	return config.Set("envs", string(envsJSON))
}

// GetEnv returns the current selected environment
// Does not take
func GetEnv(ctx *cli.Context) (Env, error) {
	var envName string
	if len(ctx.String("env")) > 0 {
		envName = ctx.String("env")
	} else {
		env, err := config.Get("env")
		if err != nil {
			return Env{}, err
		}
		if env == "" {
			env = EnvLocal
		}
		envName = env
	}

	return GetEnvByName(envName)
}

func GetEnvByName(env string) (Env, error) {
	envs, err := getEnvs()
	if err != nil {
		return Env{}, err
	}
	envir, ok := envs[env]
	if !ok {
		return Env{}, fmt.Errorf("Env \"%s\" not found. See `micro env` for available environments.", env)
	}
	return envir, nil
}

func GetEnvs() ([]Env, error) {
	envs, err := getEnvs()
	if err != nil {
		return nil, err
	}

	var ret []Env

	// populate the default environments
	for _, env := range defaultEnvs {
		ret = append(ret, env)
	}

	var nonDefaults []Env

	for _, env := range envs {
		if _, isDefault := defaultEnvs[env.Name]; !isDefault {
			nonDefaults = append(nonDefaults, env)
		}
	}

	// @todo order nondefault envs alphabetically
	ret = append(ret, nonDefaults...)

	return ret, nil
}

// SetEnv selects an environment to be used.
func SetEnv(envName string) error {
	envs, err := getEnvs()
	if err != nil {
		return err
	}
	_, ok := envs[envName]
	if !ok {
		return fmt.Errorf("Environment '%v' does not exist", envName)
	}
	return config.Set("env", envName)
}

// DelEnv deletes an env from config
func DelEnv(envName string) error {
	envs, err := getEnvs()
	if err != nil {
		return err
	}
	_, ok := envs[envName]
	if !ok {
		return fmt.Errorf("Environment '%v' does not exist", envName)
	}
	delete(envs, envName)
	return setEnvs(envs)
}

func IsPlatform(ctx *cli.Context) bool {
	env, err := GetEnv(ctx)
	if err == nil && env.Name == EnvPlatform {
		return true
	}
	return false
}

type Exec func(*cli.Context, []string) ([]byte, error)

func Print(e Exec) func(*cli.Context) error {
	return func(c *cli.Context) error {
		rsp, err := e(c, c.Args().Slice())
		if err != nil {
			return CliError(err)
		}
		if len(rsp) > 0 {
			fmt.Printf("%s\n", string(rsp))
		}
		return nil
	}
}

// CliError returns a user friendly message from error. If we can't determine a good one returns an error with code 128
func CliError(err error) cli.ExitCoder {
	if err == nil {
		return nil
	}
	// if it's already a cli.ExitCoder we use this
	cerr, ok := err.(cli.ExitCoder)
	if ok {
		return cerr
	}

	// grpc errors
	if mname := regexp.MustCompile(`malformed method name: \\?"(\w+)\\?"`).FindStringSubmatch(err.Error()); len(mname) > 0 {
		return cli.Exit(fmt.Sprintf(`Method name "%s" invalid format. Expecting service.endpoint`, mname[1]), 3)
	}
	if service := regexp.MustCompile(`service ([\w\.]+): route not found`).FindStringSubmatch(err.Error()); len(service) > 0 {
		return cli.Exit(fmt.Sprintf(`Service "%s" not found`, service[1]), 4)
	}
	if service := regexp.MustCompile(`unknown service ([\w\.]+)`).FindStringSubmatch(err.Error()); len(service) > 0 {
		if strings.Contains(service[0], ".") {
			return cli.Exit(fmt.Sprintf(`Service method "%s" not found`, service[1]), 5)
		}
		return cli.Exit(fmt.Sprintf(`Service "%s" not found`, service[1]), 5)
	}
	if address := regexp.MustCompile(`Error while dialing dial tcp.*?([\w]+\.[\w:\.]+): `).FindStringSubmatch(err.Error()); len(address) > 0 {
		return cli.Exit(fmt.Sprintf(`Failed to connect to micro server at %s`, address[1]), 4)
	}

	merr, ok := err.(*merrors.Error)
	if !ok {
		return cli.Exit(err, 128)
	}

	switch merr.Code {
	case 408:
		return cli.Exit("Request timed out", 1)
	case 401:
		// TODO check if not signed in, prompt to sign in
		return cli.Exit("Not authorized to perform this request", 2)
	}

	// fallback to using the detail from the merr
	return cli.Exit(merr.Detail, 127)
}
