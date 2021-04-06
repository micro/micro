// +build integration

package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	TrySuite(t, testNew, retryCount)
}

func testNew(t *T) {
	t.Parallel()

	tcs := []struct {
		svcName    string
		skipBuild  bool
		skipProtoc bool
	}{
		{svcName: "foobarsvc"},
		{svcName: "foo-bar"},
		{svcName: "foo-bar-baz"},
	}

	for _, tc := range tcs {
		t.t.Run(tc.svcName, func(t *testing.T) {
			defer func() {
				exec.Command("rm", "-r", "./"+tc.svcName).CombinedOutput()
			}()
			outp, err := exec.Command("micro", "new", tc.svcName).CombinedOutput()
			if err != nil {
				t.Fatal(err)
				return
			}
			if !tc.skipProtoc && !strings.Contains(string(outp), "protoc") {
				t.Fatalf("micro new lacks protobuf install instructions %v", string(outp))
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
				if !tc.skipProtoc && strings.HasPrefix(line, "make proto") {
					mp := strings.Split(line, " ")
					protocCmd := exec.Command(mp[0], mp[1:]...)
					protocCmd.Dir = "./" + tc.svcName
					pOutp, pErr := protocCmd.CombinedOutput()
					if pErr != nil {
						t.Log("That didn't work ", pErr)
						t.Fatal(string(pOutp))
						return
					}
				}
			}
			if tc.skipBuild {
				return
			}

			// for tests, update the micro import to use the current version of the code.
			fname := fmt.Sprintf("./%v/go.mod", tc.svcName)
			f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.Fatal(string(outp))
				return
			}
			if _, err := f.WriteString("\nreplace github.com/micro/micro/v3 => ../.."); err != nil {
				t.Fatal(string(outp))
				return
			}
			f.Close()

			buildCommand := exec.Command("go", "build")
			buildCommand.Dir = "./" + tc.svcName
			outp, err = buildCommand.CombinedOutput()
			if err != nil {
				t.Fatal(string(outp))
				return
			}

		})
	}

}

func TestWrongCommands(t *testing.T) {
	TrySuite(t, testWrongCommands, retryCount)
}

func testWrongCommands(t *T) {
	// @TODO this is obviously bad that we have to start a server for this. Why?
	// What happens is in `cmd/cmd.go` `/service/store/cli/util.go`.SetupCommand is called
	// which does not run for builtin services and help etc but there is no such exception for
	// missing/unrecognized commands, so the behaviour below will only happen if a `micro server`
	// is running. This is most likely because some config/auth wrapper in the background failing.
	// Fix this later.
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	t.Parallel()

	cmd := serv.Command()

	outp, err := cmd.Exec()
	if err == nil {
		t.Fatal("Missing command should error")
	}

	if !strings.Contains(string(outp), "No command") {
		t.Fatalf("Unexpected output for no command: %v", string(outp))
	}

	outp, err = cmd.Exec("asdasd")
	if err == nil {
		t.Fatal("Wrong command should error")
	}

	if !strings.Contains(string(outp), "Unrecognized micro command: asdasd. Please refer to 'micro --help'") {
		t.Fatalf("Unexpected output for unrecognized command: %v", string(outp))
	}

	outp, err = cmd.Exec("config", "asdasd")
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
	TrySuite(t, testHelps, retryCount)
}

func testHelps(t *T) {
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

		outp, err = exec.Command("micro", commandName, "--help").CombinedOutput()
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
