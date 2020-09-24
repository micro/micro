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
	// EnvPlatform is a builtin environment, the One True Micro Live(tm) environment.
	EnvPlatform = "platform"
)

const (
	// localProxyAddress is the default proxy address for environment server
	localProxyAddress = "127.0.0.1:8081"
	// platformProxyAddress is teh default proxy address for environment platform
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
func CLIProxyAddress(ctx *cli.Context) string {
	switch ctx.Args().First() {
	case "new", "server", "help", "env":
		return ""
	}

	// fix for "micro service [command]", e.g "micro service auth"
	if ctx.Args().First() == "service" && IsBuiltInService(ctx.Args().Get(1)) {
		return ""
	}

	// don't set the proxy address on the proxy
	if ctx.Args().First() == "proxy" {
		return ""
	}

	addr := GetEnv(ctx).ProxyAddress
	if !strings.Contains(addr, ":") {
		return fmt.Sprintf("%v:443", addr)
	}
	return addr
}

type Env struct {
	Name         string
	ProxyAddress string
}

func AddEnv(env Env) {
	envs := getEnvs()
	envs[env.Name] = env
	setEnvs(envs)
}

func getEnvs() map[string]Env {
	envsJSON, err := config.Get("envs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	envs := map[string]Env{}
	if len(envsJSON) > 0 {
		err := json.Unmarshal([]byte(envsJSON), &envs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	for k, v := range defaultEnvs {
		envs[k] = v
	}
	return envs
}

func setEnvs(envs map[string]Env) {
	envsJSON, err := json.Marshal(envs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Set("envs", string(envsJSON))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetEnv returns the current selected environment
// Does not take
func GetEnv(ctx *cli.Context) Env {
	var envName string
	if len(ctx.String("env")) > 0 {
		envName = ctx.String("env")
	} else {
		env, err := config.Get("env")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if env == "" {
			env = EnvLocal
		}
		envName = env
	}

	return GetEnvByName(envName)
}

func GetEnvByName(env string) Env {
	envs := getEnvs()

	envir, ok := envs[env]
	if !ok {
		fmt.Println(fmt.Sprintf("Env \"%s\" not found. See `micro env` for available environments.", env))
		os.Exit(1)
	}
	return envir
}

func GetEnvs() []Env {
	envs := getEnvs()
	ret := []Env{defaultEnvs[EnvLocal], defaultEnvs[EnvPlatform]}
	nonDefaults := []Env{}
	for _, env := range envs {
		if _, isDefault := defaultEnvs[env.Name]; !isDefault {
			nonDefaults = append(nonDefaults, env)
		}
	}
	// @todo order nondefault envs alphabetically
	ret = append(ret, nonDefaults...)
	return ret
}

// SetEnv selects an environment to be used.
func SetEnv(envName string) {
	envs := getEnvs()
	_, ok := envs[envName]
	if !ok {
		fmt.Printf("Environment '%v' does not exist\n", envName)
		os.Exit(1)
	}
	config.Set("env", envName)
}

// DelEnv deletes an env from config
func DelEnv(envName string) {
	envs := getEnvs()
	_, ok := envs[envName]
	if !ok {
		fmt.Printf("Environment '%v' does not exist\n", envName)
		os.Exit(1)
	}
	delete(envs, envName)
	setEnvs(envs)
}

func IsPlatform(ctx *cli.Context) bool {
	return GetEnv(ctx).Name == EnvPlatform
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
