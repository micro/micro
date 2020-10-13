package main

import (
	"io"
	"os"
	"os/exec"

	runtime "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// setup the client
	srv := service.New()
	cli := runtime.NewBuildService("runtime", srv.Client())

	// get the name and version of the service, these are injected by the runtime manager
	name := getEnv("MICRO_SERVICE_NAME")
	version := getEnv("MICRO_SERVICE_VERSION")

	// stream the binary from the runtime
	logger.Infof("Downloading service %v:%v", name, version)
	svc := &runtime.Service{Name: name, Version: version}
	stream, err := cli.Read(context.DefaultContext, svc, client.WithAuthToken())
	if err != nil {
		logger.Fatalf("Error starting stream: %v", err)
	}

	// create a file to write the result into
	file, err := os.Create("service")
	if err != nil {
		logger.Fatalf("Error creating output file: %v", err)
	}
	if err := os.Chmod("service", 744); err != nil {
		logger.Fatalf("Error setting output file permissions: %v", err)
	}

	// write the build to the local file
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Fatalf("Error reading from the stream: %v", err)
		}

		// write the bytes to the buffer
		if _, err := file.Write(req.Data); err != nil {
			logger.Fatalf("Error writing data to the file: %v", err)
		}
	}
	if err := file.Close(); err != nil {
		logger.Fatalf("Error closing the file: %v", err)
	}

	// execute the binary
	logger.Info("Starting service")
	cmd := exec.Command("./service")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logger.Fatalf("Error starting service: %v", err)
	}

	if err = cmd.Wait(); err != nil {
		logger.Fatalf("Service exited: %v", err)
	} else {
		logger.Fatalf("Service finished")
	}
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		logger.Fatalf("Missing required env var: %v", key)
	}
	return val
}
