package command

import (
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/micro/internal/command"
)

// Echo returns the same message
func Echo(ctx *cli.Context) Command {
	usage := "echo [text]"
	desc := "Returns the [text]"

	return NewCommand("echo", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("echo what?"), nil
		}
		return []byte(strings.Join(args[1:], " ")), nil
	})
}

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

// Register registers a service
func Register(ctx *cli.Context) Command {
	usage := "register service [definition]"
	desc := "Registers a service"

	return NewCommand("register", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("register what?"), nil
		}
		switch args[1] {
		case "service":
			if len(args) < 3 {
				return []byte("require service definition"), nil
			}
			rsp, err := command.RegisterService(ctx, args[2:])
			if err != nil {
				return nil, err
			}
			return rsp, nil
		default:
			return []byte("unknown command...\nsupported commands: \nregister service [definition]"), nil
		}
	})
}

// Deregister registers a service
func Deregister(ctx *cli.Context) Command {
	usage := "deregister service [definition]"
	desc := "Deregisters a service"

	return NewCommand("deregister", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("deregister what?"), nil
		}
		switch args[1] {
		case "service":
			if len(args) < 3 {
				return []byte("require service definition"), nil
			}
			rsp, err := command.DeregisterService(ctx, args[2:])
			if err != nil {
				return nil, err
			}
			return rsp, nil
		default:
			return []byte("unknown command...\nsupported commands: \nderegister service [definition]"), nil
		}
	})
}

// Laws of robotics
func ThreeLaws(ctx *cli.Context) Command {
	usage := "the three laws"
	desc := "Returns the three laws of robotics"

	return NewCommand("the three laws", usage, desc, func(args ...string) ([]byte, error) {
		laws := []string{
			"1. A robot may not injure a human being or, through inaction, allow a human being to come to harm.",
			"2. A robot must obey the orders given it by human beings except where such orders would conflict with the First Law.",
			"3. A robot must protect its own existence as long as such protection does not conflict with the First or Second Laws.",
		}
		return []byte("\n" + strings.Join(laws, "\n")), nil
	})
}

// Time returns the time
func Time(ctx *cli.Context) Command {
	usage := "time"
	desc := "Returns the server time"

	return NewCommand("time", usage, desc, func(args ...string) ([]byte, error) {
		t := time.Now().Format(time.RFC1123)
		return []byte("Server time is: " + t), nil
	})
}
