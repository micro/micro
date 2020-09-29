package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/micro/micro/v3/cmd"

	// load packages so they can register commands
	_ "github.com/micro/micro/v3/client/cli"
	_ "github.com/micro/micro/v3/server"
	_ "github.com/micro/micro/v3/service/cli"

	// include the platform profile
	_ "github.com/micro/micro/profile/platform/v3"
)

var (
	image   = "micro/platform"
	profile = "platform"
)

func main() {
	usage := fmt.Sprintf("%s {install|uninstall|update}", os.Args[0])

	switch os.Args[1] {
	case "install", "uninstall":
		usage = fmt.Sprintf("%s {install|uninstall} {dev|staging|platform}", os.Args[0])

		if len(os.Args) < 3 {
			fmt.Println(usage)
			return
		}

		// set the install/uninstall script
		action := os.Args[1] + ".sh"
		// create the args
		args := []string{action}
		args = append(args, os.Args[2:]...)
		// exec the command
		cmd := exec.Command("bash", args...)
		cmd.Dir = "./kubernetes"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	case "update":
		usage = fmt.Sprintf("%s update [tag]", os.Args[0])

		if len(os.Args) < 3 {
			fmt.Println(usage)
			return
		}

		tag := os.Args[2]

		// set the tag for the micro deployment
		cmd := exec.Command("kubectl", "set", "image", "deployments", "micro="+image+":"+tag, "-l", "micro=runtime")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	default:
		// set the profile
		os.Setenv("MICRO_PROFILE", profile)

		// run micro by default
		cmd.Run()
		return
	}
}
