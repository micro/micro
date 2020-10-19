package main

import (
	"fmt"
	"os"

	"github.com/micro/micro/v3/service/runtime/local"
)

// main prints out the entrypoint, e.g cmd/test/main.go for the current directory. If there is an
// error, e.g no main.go was found, then the error will be printed and the application will exit with
// an error code
func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory: ", err)
		os.Exit(1)
	}

	ep, err := local.Entrypoint(wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(ep)
}
