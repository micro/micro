package run


import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/cmd"

	api "github.com/micro/micro/v5/cmd/micro-api"
	web "github.com/micro/micro/v5/cmd/micro-web"
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
	// Unified server on :8080
	mux := http.NewServeMux()
	mux.Handle("/api/", api.APIHandler())
	mux.Handle("/web/", web.WebHandler())
	// Root (/) serves the web UI dashboard
	mux.Handle("/", web.WebHandler())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Micro server running on http://localhost:8080")
		fmt.Println("- API: http://localhost:8080/api/")
		fmt.Println("- Web UI: http://localhost:8080/web/")
		fmt.Println("- Dashboard: http://localhost:8080/")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	<-stop
	fmt.Println("Shutting down micro server...")
	return srv.Close()
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
