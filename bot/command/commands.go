package command

import (
	"strings"

	"github.com/micro/cli"
	"github.com/micro/micro/internal/command"
)

// Hello returns a greeting
func Hello(ctx *cli.Context) Command {
	usage := "hello"
	desc := "Returns a greeting"

	return NewCommand("hello", usage, desc, func(args ...string) ([]byte, error) {
		return []byte("hey what's up?"), nil
	})
}

// Ping returns pong
func Ping(ctx *cli.Context) Command {
	usage := "ping"
	desc := "Returns pong"

	return NewCommand("ping", usage, desc, func(args ...string) ([]byte, error) {
		return []byte("pong"), nil
	})
}

// Get service returns a service
func Get(ctx *cli.Context) Command {
	usage := "get service [name]"
	desc := "Returns a registered service"

	return NewCommand("get", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("get what?"), nil
		}
		switch args[1] {
		case "service":
			if len(args) < 3 {
				return []byte("require service name"), nil
			}
			rsp, err := command.GetService(ctx, args[2:])
			if err != nil {
				return nil, err
			}
			return rsp, nil
		default:
			return []byte("unknown command...\nsupported commands: \nget service [name]"), nil
		}
	})
}

// Health returns the health of a service
func Health(ctx *cli.Context) Command {
	usage := "health [service]"
	desc := "Returns health of a service"

	return NewCommand("health", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("health of what?"), nil
		}
		rsp, err := command.QueryHealth(ctx, args[1:])
		if err != nil {
			return nil, err
		}
		return rsp, nil
	})
}

// List returns a list of services
func List(ctx *cli.Context) Command {
	usage := "list services"
	desc := "Returns a list of registered services"

	return NewCommand("list", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("list what?"), nil
		}
		switch args[1] {
		case "services":
			rsp, err := command.ListServices(ctx)
			if err != nil {
				return nil, err
			}
			return rsp, nil
		default:
			return []byte("unknown command...\nsupported commands: \nlist services"), nil
		}
	})
}

// Query returns a service query
func Query(ctx *cli.Context) Command {
	usage := "query [service] [method] [request]"
	desc := "Returns the response for a service query"

	return NewCommand("query", usage, desc, func(args ...string) ([]byte, error) {
		var cargs []string

		for _, arg := range args {
			if len(strings.TrimSpace(arg)) == 0 {
				continue
			}
			cargs = append(cargs, arg)
		}

		if len(cargs) < 2 {
			return []byte("query what?"), nil
		}

		rsp, err := command.QueryService(ctx, cargs[1:])
		if err != nil {
			return nil, err
		}
		return rsp, nil
	})
}
