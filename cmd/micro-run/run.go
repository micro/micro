package run

import (
	"fmt"
	"bufio"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/registry"
)

// Color codes for log output
var colors = []string{
	"\033[31m", // red
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[34m", // blue
	"\033[35m", // magenta
	"\033[36m", // cyan
}

func colorFor(idx int) string {
	return colors[idx%len(colors)]
}

func serveMicroWeb(dir string) {
	webDir := filepath.Join(dir, "web")
	parentDir := filepath.Base(dir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(webDir); err == nil {
			// web subdir exists, look for service by parent dir name
			srvs, err := registry.GetService(parentDir)
			if err == nil && len(srvs) > 0 && len(srvs[0].Nodes) > 0 {
				// reverse proxy to first node
				target := srvs[0].Nodes[0].Address
				u, _ := url.Parse("http://" + target)
				proxy := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					proxyReq, _ := http.NewRequest(req.Method, u.String()+req.RequestURI, req.Body)
					for k, v := range req.Header {
						proxyReq.Header[k] = v
					}
					resp, err := http.DefaultClient.Do(proxyReq)
					if err != nil {
						http.Error(w, "Proxy error", 502)
						return
					}
					defer resp.Body.Close()
					for k, v := range resp.Header {
						w.Header()[k] = v
					}
					w.WriteHeader(resp.StatusCode)
					io.Copy(w, resp.Body)
				})
				proxy.ServeHTTP(w, r)
				return
			}
		}
		// else: serve index page listing services
		services, _ := registry.ListServices()
		html := "<h1>Available Services</h1>"
		for _, s := range services {
			html += "<p>" + s.Name + "</p>"
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
	go http.ListenAndServe(":8080", nil)
}

func Run(c *cli.Context) error {
	all := c.Bool("all")
	dir := c.Args().Get(0)
	if len(dir) == 0 {
		dir = "."
	}
	serveMicroWeb(dir)

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
		for i, mainFile := range mainFiles {
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
			buildCmd := exec.Command("go", "build", "-o", binPath, ".")
			buildCmd.Dir = serviceDir
			buildOut, buildErr := buildCmd.CombinedOutput()
			if buildErr != nil {
				logFile.WriteString(string(buildOut))
				logFile.Close()
				fmt.Fprintf(os.Stderr, "failed to build %s: %v\n", serviceName, buildErr)
				continue
			}
			cmd := exec.Command(binPath)
			cmd.Dir = serviceDir
			pr, pw := io.Pipe()
			cmd.Stdout = pw
			cmd.Stderr = pw
			color := colorFor(i)
			go func(name string, color string, pr *io.PipeReader) {
				scanner := bufio.NewScanner(pr)
				for scanner.Scan() {
					fmt.Printf("%s[%s]\033[0m %s\n", color, name, scanner.Text())
				}
			}(serviceName, color, pr)
			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
				pw.Close()
				continue
			}
			procs = append(procs, cmd)
			pidFiles = append(pidFiles, pidFilePath)
			os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, serviceDir)), 0644)
		}
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
		for _, proc := range procs {
			_ = proc.Wait()
		}
		return nil
	}

	// single service mode (no color needed)
	serviceName := filepath.Base(dir)
	logFilePath := filepath.Join(logsDir, serviceName+".log")
	binPath := filepath.Join(binDir, serviceName)
	pidFilePath := filepath.Join(runDir, serviceName+".pid")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()
	buildCmd := exec.Command("go", "build", "-o", binPath, dir)
	buildCmd.Dir = dir
	buildOut, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		logFile.WriteString(string(buildOut))
		return fmt.Errorf("failed to build %s: %v", serviceName, buildErr)
	}
	cmd := exec.Command(binPath)
	cmd.Dir = dir
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			fmt.Printf("[%s] %s\n", serviceName, scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, dir)), 0644)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = os.Remove(pidFilePath)
		os.Exit(1)
	}()
	_ = cmd.Wait()
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "run",
		Usage:  "Run a service or all services in a directory",
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
		Action: Run,
	})
}
