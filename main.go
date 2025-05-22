package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"

	_ "github.com/micro/micro/v5/cmd/micro-api"
	_ "github.com/micro/micro/v5/cmd/micro-cli"
	_ "github.com/micro/micro/v5/cmd/micro-mcp"
	_ "github.com/micro/micro/v5/cmd/micro-web"
	"github.com/micro/micro/v5/util"
)

var (
	// version is set by the release action
	// this is the default for local builds
	version = "5.0.0-dev"
)

func genProtoHandler(c *cli.Context) error {
	cmd := exec.Command("find", ".", "-name", "*.proto", "-exec", "protoc", "--proto_path=.", "--micro_out=.", "--go_out=.", `{}`, `;`)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runHandler(c *cli.Context) error {
	all := c.Bool("all")
	dir := c.Args().Get(0)
	if len(dir) == 0 {
		dir = "."
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	logsDir := filepath.Join(homeDir, "micro", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs dir: %w", err)
	}

	if all {
		var mainFiles []string
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if info.Name() == "main.go" {
				mainFiles = append(mainFiles, path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error walking the path: %w", err)
		}
		if len(mainFiles) == 0 {
			return fmt.Errorf("no main.go files found in %s", dir)
		}
		for _, mainFile := range mainFiles {
			serviceDir := filepath.Dir(mainFile)
			serviceName := filepath.Base(serviceDir)
			logFilePath := filepath.Join(logsDir, serviceName+".log")
			logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open log file for %s: %v\n", serviceName, err)
				continue
			}
			cmd := exec.Command("go", "run", mainFile)
			cmd.Stdout = logFile
			cmd.Stderr = logFile
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
				logFile.Close()
				continue
			}
			fmt.Printf("Started %s (pid %d), logging to %s\n", serviceName, cmd.Process.Pid, logFilePath)
			logFile.Close()
		}
		return nil
	}

	// single service mode
	serviceName := filepath.Base(dir)
	logFilePath := filepath.Join(logsDir, serviceName+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()
	cmd := exec.Command("go", "run", dir)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	return cmd.Run()
}

func main() {
	cmd.Register([]*cli.Command{
		{
			Name:   "run",
			Usage:  "Run a service",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "all",
					Usage:   "Run all services (find all main.go)",
				},
			},
			Action: runHandler,
		},
		{
			Name:  "gen",
			Usage: "Generate various things",
			Subcommands: []*cli.Command{
				{
					Name:   "proto",
					Usage:  "Generate proto requires protoc and protoc-gen-micro",
					Action: genProtoHandler,
				},
			},
		},
		{
			Name:  "services",
			Usage: "List available services",
			Action: func(ctx *cli.Context) error {
				services, err := registry.ListServices()
				if err != nil {
					return err
				}
				for _, service := range services {
					fmt.Println(service.Name)
				}
				return nil
			},
		},
		{
			Name:  "call",
			Usage: "Call a service",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if args.Len() < 2 {
					return fmt.Errorf("Usage: [service] [endpoint] [request]")
				}

				service := args.Get(0)
				endpoint := args.Get(1)
				request := `{}`

				// get the request if present
				if args.Len() == 3 {
					request = args.Get(2)
				}

				req := client.NewRequest(service, endpoint, &bytes.Frame{Data: []byte(request)})
				var rsp bytes.Frame
				err := client.Call(context.TODO(), req, &rsp)
				if err != nil {
					return err
				}

				fmt.Print(string(rsp.Data))
				return nil
			},
		},
		{
			Name:  "describe",
			Usage: "Describe a service",
			Action: func(ctx *cli.Context) error {
				args := ctx.Args()

				if args.Len() != 1 {
					return fmt.Errorf("Usage: [service]")
				}

				service := args.Get(0)
				services, err := registry.GetService(service)
				if err != nil {
					return err
				}
				if len(services) == 0 {
					return nil
				}
				b, _ := json.MarshalIndent(services[0], "", "    ")
				fmt.Println(string(b))
				return nil
			},
		},
	}...)

	cmd.App().Action = func(c *cli.Context) error {
		if c.Args().Len() == 0 {
			return nil
		}

		// if an executable is available with the name of
		// the command, execute it with the arguments from
		// index 1 on.
		v, err := exec.LookPath("micro-" + c.Args().First())
		if err == nil {
			ce := exec.Command(v, c.Args().Slice()[1:]...)
			ce.Stdout = os.Stdout
			ce.Stderr = os.Stderr
			return ce.Run()
		}

		command := c.Args().Get(0)
		args := c.Args().Slice()

		if srv, err := util.LookupService(command); err != nil {
			return util.CliError(err)
		} else if srv != nil && util.ShouldRenderHelp(args) {
			return cli.Exit(util.FormatServiceUsage(srv, c), 0)
		} else if srv != nil {
			err := util.CallService(srv, args)
			return util.CliError(err)
		}

		return nil
	}

	cmd.Init(
		cmd.Name("micro"),
		cmd.Version(version),
	)
}
