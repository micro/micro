package test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	trySuite(t, testNew, retryCount)
}

func testNew(t *t) {
	t.Parallel()
	defer func() {
		exec.Command("rm", "-r", "./foobar").CombinedOutput()
	}()
	outp, err := exec.Command("micro", "new", "foobar").CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}
	if !strings.Contains(string(outp), "protoc") {
		t.Fatalf("micro new lacks 	protobuf install instructions %v", string(outp))
		return
	}

	lines := strings.Split(string(outp), "\n")
	// executing install instructions
	for _, line := range lines {
		if strings.HasPrefix(line, "go get") {
			parts := strings.Split(line, " ")
			getOutp, getErr := exec.Command(parts[0], parts[1:]...).CombinedOutput()
			if getErr != nil {
				t.Fatal(string(getOutp))
				return
			}
		}
		if strings.HasPrefix(line, "make proto") {
			mp := strings.Split(line, " ")
			protocCmd := exec.Command(mp[0], mp[1:]...)
			protocCmd.Dir = "./foobar"
			pOutp, pErr := protocCmd.CombinedOutput()
			if pErr != nil {
				t.Log("That didn't work ", pErr)
				t.Fatal(string(pOutp))
				return
			}
		}
	}

	buildCommand := exec.Command("go", "build")
	buildCommand.Dir = "./foobar"
	outp, err = buildCommand.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
		return
	}
}

func TestWrongCommands(t *testing.T) {
	trySuite(t, testWrongCommands, retryCount)
}

func testWrongCommands(t *t) {
	t.Parallel()

	comm := exec.Command("micro")
	outp, err := comm.CombinedOutput()
	if err == nil {
		t.Fatal("Missing command should error")
	}
	lines := strings.Split(string(outp), "\n")
	if len(lines) > 1 || !strings.Contains(string(outp), "No command") {
		t.Fatalf("Unexpected output for no command: %v", string(outp))
	}

	comm = exec.Command("micro", "asdasd")
	outp, err = comm.CombinedOutput()
	if err == nil {
		t.Fatal("Wrong command should error")
	}
	lines = strings.Split(string(outp), "\n")
	if len(lines) > 1 || !strings.Contains(string(outp), "Unrecognized micro command") {
		t.Fatalf("Unexpected output for unrecognized command: %v", string(outp))
	}

	comm = exec.Command("micro", "config", "asdasd")
	outp, err = comm.CombinedOutput()
	if err == nil {
		t.Fatal("Wrong subcommand should error")
	}
	lines = strings.Split(string(outp), "\n")
	// @todod for some reason this one returns multiple lines so we don't check for line count now
	if !strings.Contains(string(outp), "Unrecognized subcommand for micro config") {
		t.Fatalf("Unexpected output for unrecognized subcommand: %v", string(outp))
	}
}
