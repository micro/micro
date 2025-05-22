package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

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
	daemon := c.Bool("daemon")
	dir := c.Args().Get(0)
	if len(dir) == 0 {
		dir = "."
	}

	// Detect git URL and clone if needed
	isGitURL := strings.HasPrefix(dir, "github.com/") || strings.HasPrefix(dir, "https://") || strings.HasPrefix(dir, "git@")
	if isGitURL {
		tmpDir, err := os.MkdirTemp("", "micro-git-")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		gitURL := dir
		if strings.HasPrefix(dir, "github.com/") {
			gitURL = "https://" + dir
		}
		cmd := exec.Command("git", "clone", gitURL, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo: %w", err)
		}
		dir = tmpDir
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	logsDir := filepath.Join(homeDir, "micro", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs dir: %w", err)
	}
	runDir := filepath.Join(homeDir, "micro", "run")
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create run dir: %w", err)
	}
	binDir := filepath.Join(homeDir, "micro", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin dir: %w", err)
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
		var procs []*exec.Cmd
		var pidFiles []string
		for _, mainFile := range mainFiles {
			serviceDir := filepath.Dir(mainFile)
			serviceName := filepath.Base(serviceDir)
			logFilePath := filepath.Join(logsDir, serviceName+".log")
			binPath := filepath.Join(binDir, serviceName)
			pidFilePath := filepath.Join(runDir, serviceName+".pid")
			logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open log file for %s: %v\n", serviceName, err)
				continue
			}
			// Build the binary
			buildCmd := exec.Command("go", "build", "-o", binPath, ".")
			buildCmd.Dir = serviceDir
			buildOut, buildErr := buildCmd.CombinedOutput()
			if buildErr != nil {
				logFile.WriteString(string(buildOut))
				lastLine := lastNonEmptyLine(string(buildOut))
				os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", serviceDir, lastLine)), 0644)
				fmt.Fprintf(os.Stderr, "failed to build %s: %v\n", serviceName, buildErr)
				logFile.Close()
				continue
			}
			// Run the binary
			cmd := exec.Command(binPath)
			cmd.Dir = serviceDir
			if daemon {
				cmd.Stdout = logFile
				cmd.Stderr = logFile
				if err := cmd.Start(); err != nil {
					lastLine := lastLogLine(logFilePath)
					os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", serviceDir, lastLine)), 0644)
					fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
					logFile.Close()
					continue
				}
				os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, serviceDir)), 0644)
				fmt.Printf("Started %s (pid %d), logging to %s\n", serviceName, cmd.Process.Pid, logFilePath)
				logFile.Close()
			} else {
				pr, pw := io.Pipe()
				cmd.Stdout = pw
				cmd.Stderr = pw
				go func(logFile *os.File, pr *io.PipeReader) {
					tee := io.MultiWriter(os.Stdout, logFile)
					io.Copy(tee, pr)
					logFile.Close()
				}(logFile, pr)
				if err := cmd.Start(); err != nil {
					lastLine := lastLogLine(logFilePath)
					os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", serviceDir, lastLine)), 0644)
					fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
					pw.Close()
					continue
				}
				os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, serviceDir)), 0644)
				fmt.Printf("Started %s (pid %d), logging to %s\n", serviceName, cmd.Process.Pid, logFilePath)
				procs = append(procs, cmd)
				pidFiles = append(pidFiles, pidFilePath)
			}
		}
		if !daemon {
			// Handle Ctrl+C cleanup
			waitAndCleanup(procs, pidFiles)
		}
		return nil
	}

	// single service mode
	serviceName := filepath.Base(dir)
	logFilePath := filepath.Join(logsDir, serviceName+".log")
	binPath := filepath.Join(binDir, serviceName)
	pidFilePath := filepath.Join(runDir, serviceName+".pid")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binPath, dir)
	buildCmd.Dir = dir
	buildOut, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		logFile.WriteString(string(buildOut))
		lastLine := lastNonEmptyLine(string(buildOut))
		os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", dir, lastLine)), 0644)
		return fmt.Errorf("failed to build %s: %v", serviceName, buildErr)
	}
	cmd := exec.Command(binPath)
	cmd.Dir = dir
	if daemon {
		cmd.Stdout = logFile
		cmd.Stderr = logFile
		if err := cmd.Start(); err != nil {
			lastLine := lastLogLine(logFilePath)
			os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", dir, lastLine)), 0644)
			return err
		}
		os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, dir)), 0644)
		return nil
	} else {
		pr, pw := io.Pipe()
		cmd.Stdout = pw
		cmd.Stderr = pw
		go func() {
			tee := io.MultiWriter(os.Stdout, logFile)
			io.Copy(tee, pr)
		}()
		if err := cmd.Start(); err != nil {
			lastLine := lastLogLine(logFilePath)
			os.WriteFile(pidFilePath, []byte(fmt.Sprintf("0\n%s\nreason: %s\n", dir, lastLine)), 0644)
			return err
		}
		os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, dir)), 0644)
		waitAndCleanup([]*exec.Cmd{cmd}, []string{pidFilePath})
		return nil
	}
}

// lastNonEmptyLine returns the last non-empty line from a string
func lastNonEmptyLine(s string) string {
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			return lines[i]
		}
	}
	return ""
}

// lastLogLine returns the last non-empty line from a file
func lastLogLine(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	var last string
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		if strings.TrimSpace(scan.Text()) != "" {
			last = scan.Text()
		}
	}
	return last
}

// waitAndCleanup waits for all procs and removes pid files on exit or Ctrl+C
func waitAndCleanup(procs []*exec.Cmd, pidFiles []string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		for _, proc := range procs {
			if proc.Process != nil {
				_ = proc.Process.Kill()
			}
		}
		for _, pf := range pidFiles {
			_ = os.Remove(pf)
		}
		os.Exit(1)
	}()
	for i, proc := range procs {
		_ = proc.Wait()
		if proc.Process != nil {
			_ = os.Remove(pidFiles[i])
		}
	}
}

func main() {
	cmd.Register([]*cli.Command{
		{
			Name:  "run",
			Usage: "Run a service",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "all",
					Usage: "Run all services (find all main.go)",
				},
				&cli.BoolFlag{
					Name:    "daemon",
					Aliases: []string{"d"},
					Usage:   "Daemonize (detach and only log to file)",
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
		{
			Name:  "status",
			Usage: "Check status of running services",
			Action: func(ctx *cli.Context) error {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home dir: %w", err)
				}
				runDir := filepath.Join(homeDir, "micro", "run")
				files, err := os.ReadDir(runDir)
				if err != nil {
					return fmt.Errorf("failed to read run dir: %w", err)
				}
				fmt.Printf("%-20s %-8s %-8s %s\n", "SERVICE", "PID", "STATUS", "DIRECTORY")
				for _, f := range files {
					if f.IsDir() || !strings.HasSuffix(f.Name(), ".pid") {
						continue
					}
					service := f.Name()[:len(f.Name())-4]
					pidFilePath := filepath.Join(runDir, f.Name())
					pidFile, err := os.Open(pidFilePath)
					if err != nil {
						continue
					}
					var pid int
					var dir, reason string
					fmt.Fscanf(pidFile, "%d\n%s\nreason: [%s]\n", &pid, &dir, &reason)
					pidFile.Close()
					status := "stopped"
					if pid > 0 {
						proc, err := os.FindProcess(pid)
						if err == nil {
							// On unix, sending syscall.Signal(0) checks if running
							// import "syscall" at the top
							if err := proc.Signal(syscall.Signal(0)); err == nil {
								status = "running"
							}
						}
					}
					if reason != "" && status != "running" {
						fmt.Printf("%-20s %-8d %-8s %-40s %s\n", service, pid, status, reason, dir)
					} else {
						fmt.Printf("%-20s %-8d %-8s %-40s %s\n", service, pid, status, "", dir)
					}
				}
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Stop a running service",
			Action: func(ctx *cli.Context) error {
				if ctx.Args().Len() != 1 {
					return fmt.Errorf("Usage: micro stop [service]")
				}
				service := ctx.Args().Get(0)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home dir: %w", err)
				}
				runDir := filepath.Join(homeDir, "micro", "run")
				pidFilePath := filepath.Join(runDir, service+".pid")
				pidFile, err := os.Open(pidFilePath)
				if err != nil {
					return fmt.Errorf("no pid file for service %s", service)
				}
				var pid int
				var dir, reason string
				fmt.Fscanf(pidFile, "%d\n%s\nreason: [%s]\n", &pid, &dir, &reason)
				pidFile.Close()
				if pid <= 0 {
					_ = os.Remove(pidFilePath)
					return fmt.Errorf("service %s is not running", service)
				}
				proc, err := os.FindProcess(pid)
				if err != nil {
					_ = os.Remove(pidFilePath)
					return fmt.Errorf("could not find process for %s", service)
				}
				if err := proc.Signal(syscall.SIGTERM); err != nil {
					_ = os.Remove(pidFilePath)
					return fmt.Errorf("failed to stop service %s: %v", service, err)
				}
				_ = os.Remove(pidFilePath)
				fmt.Printf("Stopped service %s (pid %d)\n", service, pid)
				return nil
			},
		},
		{
			Name:  "logs",
			Usage: "Show logs for a service",
			Action: func(ctx *cli.Context) error {
				if ctx.Args().Len() != 1 {
					return fmt.Errorf("Usage: micro logs [service]")
				}
				service := ctx.Args().Get(0)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home dir: %w", err)
				}
				logsDir := filepath.Join(homeDir, "micro", "logs")
				logFilePath := filepath.Join(logsDir, service+".log")
				f, err := os.Open(logFilePath)
				if err != nil {
					return fmt.Errorf("could not open log file for service %s: %v", service, err)
				}
				defer f.Close()
				scan := bufio.NewScanner(f)
				for scan.Scan() {
					fmt.Println(scan.Text())
				}
				return scan.Err()
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
