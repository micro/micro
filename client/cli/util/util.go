// Package cliutil contains methods used across all cli commands
// @todo: get rid of os.Exits and use errors instread
package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	// devProxyAddress is the address for the proxy server in the dev environment
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
	}
)

var defaultEnvs = map[string]Env{
	EnvLocal: {
		Name:         EnvLocal,
		ProxyAddress: localProxyAddress,
	},
	EnvDev: {
		Name:         EnvDev,
		ProxyAddress: devProxyAddress,
	},
	EnvPlatform: {
		Name:         EnvPlatform,
		ProxyAddress: platformProxyAddress,
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
		return fmt.Errorf("Environment '%v' does not exist\n", envName)
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
		return fmt.Errorf("Environment '%v' does not exist\n", envName)
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
			fmt.Println(err)
			os.Exit(1)
		}
		if len(rsp) > 0 {
			fmt.Printf("%s\n", string(rsp))
		}
		return nil
	}
}
