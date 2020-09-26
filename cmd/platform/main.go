package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	usage := fmt.Sprintf("%s {install|uninstall}", os.Args[0])

	if len(os.Args) == 1 {
		fmt.Println(usage)
		return
	}

	switch os.Args[1] {
	case "install", "uninstall":
		cmd := exec.Command("bash", os.Args[1]+".sh")
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
