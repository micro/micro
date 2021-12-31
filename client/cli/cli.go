// Package cli is a command line interface
package cli

import (
	"fmt"
	"os"
	osexec "os/exec"
	"strings"

	rl "github.com/chzyer/readline"
	"github.com/micro/micro/v3/client/cli/util"
	"github.com/micro/micro/v3/cmd"
	"github.com/urfave/cli/v2"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/token"
	"github.com/micro/micro/v3/internal/config"
	"github.com/micro/micro/v3/service/auth"

	_ "github.com/micro/micro/v3/client/cli/auth"
	_ "github.com/micro/micro/v3/client/cli/config"
	_ "github.com/micro/micro/v3/client/cli/gen"
	_ "github.com/micro/micro/v3/client/cli/init"
	_ "github.com/micro/micro/v3/client/cli/network"
	_ "github.com/micro/micro/v3/client/cli/new"
	_ "github.com/micro/micro/v3/client/cli/run"
	_ "github.com/micro/micro/v3/client/cli/store"
	_ "github.com/micro/micro/v3/client/cli/user"
)

var (
	commonF1 = []rl.PrefixCompleterInterface{
		rl.PcItem("--name"),
		rl.PcItem("--source"),
		rl.PcItem("--image"),
		rl.PcItem("--command"),
		rl.PcItem("--args"),
		rl.PcItem("--type"),
		rl.PcItem("--env_vars"),
	}

	commonF2 = []rl.PrefixCompleterInterface{
		rl.PcItem("--database"),
		rl.PcItem("--table"),
		rl.PcItem("--prefix"),
		rl.PcItem("--limit"),
		rl.PcItem("--offset"),
		rl.PcItem("--verbose"),
		rl.PcItem("--output"),
	}

	commonF3 = []rl.PrefixCompleterInterface{
		rl.PcItem("--nodes"),
		rl.PcItem("--database"),
		rl.PcItem("--table"),
	}

	commonF4 = []rl.PrefixCompleterInterface{
		rl.PcItem("--scope"),
		rl.PcItem("--resource"),
		rl.PcItem("--access"),
		rl.PcItem("--priority"),
	}

	completer = rl.NewPrefixCompleter(
		// need to check auth again
		rl.PcItem("auth",
			rl.PcItem("list",
				rl.PcItem("rules"),
				rl.PcItem("accounts"),
			),
			rl.PcItem("create",
				rl.PcItem("rule", commonF4...),
				rl.PcItem("account",
					rl.PcItem("--secret"),
					rl.PcItem("--scopes"),
					rl.PcItem("--namespace"),
				),
			),
			rl.PcItem("delete",
				rl.PcItem("rule", commonF4...),
				rl.PcItem("account",
					rl.PcItem("--secret"),
					rl.PcItem("--scopes"),
				),
			),
		),
		rl.PcItem("config",
			rl.PcItem("get",
				rl.PcItem("--secret"),
			),
			rl.PcItem("set",
				rl.PcItem("--secret"),
			),
			rl.PcItem("del"),
		),
		rl.PcItem("gen"),
		rl.PcItem("init",
			rl.PcItem("--package"),
			rl.PcItem("--profile"),
			rl.PcItem("--output"),
		),
		rl.PcItem("network",
			rl.PcItem("connect"),
			rl.PcItem("connections"),
			rl.PcItem("graph"),
			rl.PcItem("nodes"),
			rl.PcItem("services"),
			rl.PcItem("routes",
				rl.PcItem("--service"),
				rl.PcItem("--address"),
				rl.PcItem("--gateway"),
				rl.PcItem("--router"),
				rl.PcItem("--network"),
			),
			rl.PcItem("call",
				rl.PcItem("--address"),
				rl.PcItem("--output"),
				rl.PcItem("--metadata"),
			),
		),
		rl.PcItem("new"),
		rl.PcItem("run", commonF1...),
		rl.PcItem("update", commonF1...),
		rl.PcItem("kill", commonF1...),
		rl.PcItem("status", commonF1...),
		rl.PcItem("logs",
			rl.PcItem("--version"),
			rl.PcItem("--output"),
			rl.PcItem("--follow"),
			rl.PcItem("--since"),
			rl.PcItem("--lines"),
		),
		rl.PcItem("signup",
			rl.PcItem("--email"),
			rl.PcItem("--password"),
			rl.PcItem("--recover"),
		),
		rl.PcItem("store",
			rl.PcItem("read", commonF2...),
			rl.PcItem("list", removeIndex(commonF2, 5)...),
			rl.PcItem("write",
				rl.PcItem("--expiry"),
				rl.PcItem("--database"),
				rl.PcItem("--table"),
			),
			rl.PcItem("delete",
				rl.PcItem("--database"),
				rl.PcItem("--table"),
			),
			rl.PcItem("databases", rl.PcItem("--store")),
			rl.PcItem("tables",
				rl.PcItem("--store"),
				rl.PcItem("--database"),
			),
			rl.PcItem("snapshot", append(commonF3, rl.PcItem("--destination"))...),
			rl.PcItem("sync",
				rl.PcItem("--from-backend"),
				rl.PcItem("--from-nodes"),
				rl.PcItem("--from-database"),
				rl.PcItem("--from-table"),
				rl.PcItem("--to-backend"),
				rl.PcItem("--to-nodes"),
				rl.PcItem("--to-database"),
				rl.PcItem("--to-table"),
			),
			rl.PcItem("restore", append(commonF3, rl.PcItem("--source"))...),
		),
		rl.PcItem("user",
			rl.PcItem("config",
				rl.PcItem("get"),
				rl.PcItem("set"),
				rl.PcItem("delete"),
			),
			rl.PcItem("token"),
			rl.PcItem("namespace", rl.PcItem("set")),
			rl.PcItem("set",
				rl.PcItem("password",
					rl.PcItem("--email"),
					rl.PcItem("--old-password"),
					rl.PcItem("--new-password"),
				),
			),
		),
		rl.PcItem("call",
			rl.PcItem("--address"),
			rl.PcItem("--output"),
			rl.PcItem("--metadata"),
			rl.PcItem("--request_timeout"),
		),
		rl.PcItem("stream",
			rl.PcItem("--output"),
			rl.PcItem("--metadata"),
		),
		rl.PcItem("stats", rl.PcItem("--all")),
		rl.PcItem("env",
			rl.PcItem("get"),
			rl.PcItem("set"),
			rl.PcItem("add"),
			rl.PcItem("del"),
		),
		rl.PcItem("services"),
		rl.PcItem("exit"),
	)

	// TODO: only run fixed set of commands for security purposes
	commands = map[string]*command{}
)

type command struct {
	name  string
	usage string
	exec  util.Exec
}

// remove an element from a slice
func removeIndex(l []rl.PrefixCompleterInterface, i int) []rl.PrefixCompleterInterface {
	s := make([]rl.PrefixCompleterInterface, 0)
	s = append(s, l[:i]...)
	return append(s, l[i+1:]...)
}

func initPrompt(ctx *cli.Context) string {
	var u, env, ns, prompt string

	// get user
	token, _ := token.Get(ctx)

	acc, _ := auth.Inspect(token.AccessToken)

	if acc == nil {
		u = ""
	} else {
		// backward compatibility
		u = acc.Name
		if len(u) == 0 {
			u = acc.ID
		}
	}

	// get env
	env, _ = config.Get("env")

	// get ns
	ns, _ = namespace.Get(env)

	ns = fmt.Sprintf("\033[0;1m%s\033[0m", ns)
	env = fmt.Sprintf("\033[0;1m%s\033[0m", env)
	// prompt = ns + "|" + env + " " + "\033[30;107;1mMicro>\033[0m" + " "

	if u == "" {
		// prompt = "\033[31;1mlogin\033[0m" + "|" + prompt
		prompt = ns + "://" + "\033[31;1m$user\033[0m" + "@" + env + "\033[0;1m/micro> \033[0m"
	} else {
		u = fmt.Sprintf("\033[32;1m%s\033[0m", u)
		// prompt = u + "|" + prompt
		prompt = ns + "://" + u + "@" + env + "\033[0;1m/micro> \033[0m"
	}

	return prompt
}

func Run(c *cli.Context) error {
	// take the first arg as the binary
	binary := os.Args[0]

	r, err := rl.NewEx(&rl.Config{
		Prompt:       initPrompt(c),
		AutoComplete: completer,
	})

	if err != nil {
		return err
	}

	defer r.Close()

	for {

		args, err := r.Readline()
		if err != nil {
			fmt.Fprint(os.Stdout, err)
			return err
		}

		args = strings.TrimSpace(args)

		// exit from cli
		if args == "exit" {
			goto exit
		}

		// skip no args
		if len(args) == 0 {
			continue
		}

		parts := strings.Split(args, " ")
		if len(parts) == 0 {
			continue
		}

		cmd := osexec.Command(binary, parts...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(string(err.(*osexec.ExitError).Stderr))
		}

		r.SetPrompt(initPrompt(c))
	}

exit:

	return nil
}

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "cli",
			Usage:  "Run the interactive CLI",
			Action: Run,
		},
		&cli.Command{
			Name:   "call",
			Usage:  `Call a service e.g micro call greeter Say.Hello '{"name": "John"}'`,
			Action: util.Print(CallService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "address",
					Usage:   "Set the address of the service instance to call",
					EnvVars: []string{"MICRO_ADDRESS"},
				},
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
				&cli.StringFlag{
					Name:  "request_timeout",
					Usage: "timeout duration",
				},
			},
		},
		&cli.Command{
			Name:  "get",
			Usage: `Get resources from micro`,
			Subcommands: []*cli.Command{
				{
					Name:   "service",
					Usage:  "Get a specific service from the registry",
					Action: util.Print(GetService),
				},
			},
		},
		&cli.Command{
			Name:   "health",
			Usage:  `Get the service health`,
			Action: util.Print(QueryHealth),
		},
		&cli.Command{
			Name:   "stream",
			Usage:  `Create a service stream e.g. micro stream foo Bar.Baz '{"key": "value"}'`,
			Action: util.Print(streamService),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "output, o",
					Usage:   "Set the output format; json (default), raw",
					EnvVars: []string{"MICRO_OUTPUT"},
				},
				&cli.StringSliceFlag{
					Name:    "metadata",
					Usage:   "A list of key-value pairs to be forwarded as metadata",
					EnvVars: []string{"MICRO_METADATA"},
				},
			},
		},
		&cli.Command{
			Name:   "stats",
			Usage:  "Query the stats of specified service(s), e.g micro stats srv1 srv2 srv3",
			Action: util.Print(queryStats),
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "all",
					Usage: "to list all builtin services use --all builtin, for user's services use --all custom",
				},
			},
		},
		&cli.Command{
			Name:   "env",
			Usage:  "Get/set micro cli environment",
			Action: util.Print(listEnvs),
			Subcommands: []*cli.Command{
				{
					Name:   "get",
					Usage:  "Get the currently selected environment",
					Action: util.Print(getEnv),
				},
				{
					Name:   "set",
					Usage:  "Set the environment to use for subsequent commands e.g. micro env set dev",
					Action: util.Print(setEnv),
				},
				{
					Name:   "add",
					Usage:  "Add a new environment e.g. micro env add foo 127.0.0.1:8081",
					Action: util.Print(addEnv),
				},
				{
					Name:   "del",
					Usage:  "Delete an environment from your list e.g. micro env del foo",
					Action: util.Print(delEnv),
				},
			},
		},
		&cli.Command{
			Name:   "services",
			Usage:  "List services in the registry",
			Action: util.Print(ListServices),
		},
	)
}
