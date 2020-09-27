package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	usage := fmt.Sprintf("%s {install|uninstall} {dev|staging|platform}", os.Args[0])

	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}

	switch os.Args[1] {
	case "install", "uninstall":
		action := os.Args[1] + ".sh"
		args := []string{action}
		args = append(args, os.Args[2:]...)
		cmd := exec.Command("bash", args...)
		cmd.Dir = "./kubernetes"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println(usage)
		return
	}
}
