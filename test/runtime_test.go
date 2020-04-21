// +build integration

package test

import (
	"log"
	"os/exec"
	"testing"
	"time"
)

func TestMicroServer(t *testing.T) {
	serverCmd := exec.Command("micro ", "server")
	go func() {
		if err := serverCmd.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	defer serverCmd.Process.Kill()
	time.Sleep(3 * time.Second)
}
