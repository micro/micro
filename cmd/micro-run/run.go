package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/cmd"
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

func Run(c *cli.Context) error {
	dir := c.Args().Get(0)
	var tmpDir string
	if len(dir) == 0 {
		dir = "."
	} else if strings.HasPrefix(dir, "github.com/") || strings.HasPrefix(dir, "https://github.com/") {
		// Handle git URLs
		repo := dir
		if strings.HasPrefix(repo, "https://") {
			repo = strings.TrimPrefix(repo, "https://")
		}
		// Clone to a temp directory
		tmp, err := os.MkdirTemp("", "micro-run-")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		tmpDir = tmp
		cloneURL := repo
		if !strings.HasPrefix(cloneURL, "https://") {
			cloneURL = "https://" + repo
		}
		// Run git clone
		cmd := exec.Command("git", "clone", cloneURL, tmpDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo %s: %w", cloneURL, err)
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

	// Always run all services (find all main.go)
	var mainFiles []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		var serviceName string
		absDir, _ := filepath.Abs(dir)
		absServiceDir, _ := filepath.Abs(serviceDir)
		if absDir == absServiceDir {
			// If main.go is in the root dir being run, use the current working dir name
			cwd, _ := os.Getwd()
			serviceName = filepath.Base(cwd)
		} else {
			serviceName = filepath.Base(serviceDir)
		}
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
		go func(name string, color string, pr *io.PipeReader, logFile *os.File) {
			defer logFile.Close()
			scanner := bufio.NewScanner(pr)
			for scanner.Scan() {
				line := scanner.Text()
				// Write to terminal with color
				fmt.Printf("%s[%s]\033[0m %s\n", color, name, line)
				// Write to log file (without color)
				logFile.WriteString(line + "\n")
			}
		}(serviceName, color, pr, logFile)
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
			pw.Close()
			continue
		}
		procs = append(procs, cmd)
		pidFiles = append(pidFiles, pidFilePath)
		absServiceDir, _ = filepath.Abs(serviceDir)
		os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, absServiceDir)), 0644)
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

func init() {
	cmd.Register(&cli.Command{
		Name:   "run",
		Usage:  "Run all services in a directory",
		Action: Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Address to bind the micro web UI (default :8080)",
				Value:   ":8080",
			},
		},
	})
}
