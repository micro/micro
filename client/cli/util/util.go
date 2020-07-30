// Package cliutil contains methods used across all cli commands
// @todo: get rid of os.Exits and use errors instread
package util

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/micro/cli/v2"
	"github.com/micro/micro/v3/internal/config"
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
		"config",   // ????
		"network",  // :8085
		"runtime",  // :8088
		"registry", // :8000
		"broker",   // :8001
		"store",    // :8002
		"router",   // :8084
		"debug",    // :????
		"proxy",    // :8081
		"api",      // :8080
		"auth",     // :8010
		"web",      // :8082
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

func isBuiltinService(command string) bool {
	for _, service := range services {
		if command == service {
			return true
		}
	}
	return false
}

// CLIProxyAddress returns the proxy address which should be set for the client
func CLIProxyAddress(ctx *cli.Context) string {
	// This makes `micro [command name] --help` work without a server
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			return ""
		}
	}
	switch ctx.Args().First() {
	case "new", "server", "help", "env":
		return ""
	}

	// fix for "micro service [command]", e.g "micro service auth"
	if ctx.Args().First() == "service" && isBuiltinService(ctx.Args().Get(1)) {
		return ""
	}

	// don't set the proxy address on the proxy
	if ctx.Args().First() == "proxy" {
		return ""
	}

	return GetEnv(ctx).ProxyAddress
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
	err = config.Set(string(envsJSON), "envs")
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

	if len(envir.ProxyAddress) == 0 {
		return envir
	}

	// default to :8081 (the proxy port)
	if _, port, _ := net.SplitHostPort(envir.ProxyAddress); len(port) == 0 {
		envir.ProxyAddress = net.JoinHostPort(envir.ProxyAddress, "8081")
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
	config.Set(envName, "env")
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
