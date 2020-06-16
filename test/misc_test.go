// +build integration

package test

import (
	"fmt"
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
	// @TODO this is obviously bad that we have to start a server for this. Why?
	// What happens is in `cmd/cmd.go` `/service/store/cli/util.go`.SetupCommand is called
	// which does not run for builtin services and help etc but there is no such exception for
	// missing/unrecognized commands, so the behaviour below will only happen if a `micro server`
	// is running. This is most likely because some config/auth wrapper in the background failing.
	// Fix this later.
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	t.Parallel()

	comm := exec.Command("micro", serv.envFlag())
	outp, err := comm.CombinedOutput()
	if err == nil {
		t.Fatal("Missing command should error")
	}

	if !strings.Contains(string(outp), "No command") {
		t.Fatalf("Unexpected output for no command: %v", string(outp))
	}

	comm = exec.Command("micro", serv.envFlag(), "asdasd")
	outp, err = comm.CombinedOutput()
	if err == nil {
		t.Fatal("Wrong command should error")
	}

	if !strings.Contains(string(outp), "Unrecognized micro command") {
		t.Fatalf("Unexpected output for unrecognized command: %v", string(outp))
	}

	comm = exec.Command("micro", serv.envFlag(), "config", "asdasd")
	outp, err = comm.CombinedOutput()
	if err == nil {
		t.Fatal("Wrong subcommand should error")
	}

	// @todod for some reason this one returns multiple lines so we don't check for line count now
	if !strings.Contains(string(outp), "Unrecognized subcommand for micro config") {
		t.Fatalf("Unexpected output for unrecognized subcommand: %v", string(outp))
	}
}

// TestHelps ensures all `micro [command name] help` && `micro [command name] --help` commands are working.
func TestHelps(t *testing.T) {
	trySuite(t, testHelps, retryCount)
}

func testHelps(t *t) {
	comm := exec.Command("micro", "help")
	outp, err := comm.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	commands := strings.Split(strings.Split(string(outp), "COMMANDS:")[1], "GLOBAL OPTIONS:")[0]
	for _, line := range strings.Split(commands, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		// no help for help ;)
		if strings.Contains(trimmed, "help") {
			continue
		}
		commandName := strings.Split(trimmed, " ")[0]
		comm = exec.Command("micro", commandName, "--help")
		outp, err = comm.CombinedOutput()

		if err != nil {
			t.Fatal(fmt.Errorf("Command %v output is wrong: %v", commandName, string(outp)))
			break
		}
		if !strings.Contains(string(outp), "micro "+commandName+" -") {
			t.Fatal(commandName + " output is wrong: " + string(outp))
			break
		}
	}
}
