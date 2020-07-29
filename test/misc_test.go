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
		sType      string
		skipBuild  bool
		skipProtoc bool
	}{
		{svcName: "foobarsvc", sType: "service"},
		{svcName: "foobarweb", sType: "web", skipProtoc: true, skipBuild: true}, // web service has no proto generated
		{svcName: "foobarapi", sType: "api", skipBuild: true},                   // api service actually fails build out of the box because it's supposed to point to a service proto
		{svcName: "foo-bar", sType: "service"},
		{svcName: "foo-barfn", sType: "function"},
		{svcName: "foo-barweb", sType: "web", skipProtoc: true, skipBuild: true}, // web service has no proto generated
		{svcName: "foo-barapi", sType: "api", skipBuild: true},                   // api service actually fails build out of the box because it's supposed to point to a service proto
		{svcName: "foo-bar-baz", sType: "service"},
	}

	for _, tc := range tcs {
		t.t.Run(tc.svcName, func(t *testing.T) {
			defer func() {
				exec.Command("rm", "-r", "./"+tc.svcName).CombinedOutput()
			}()
			outp, err := exec.Command("micro", "new", "--type", tc.sType, tc.svcName).CombinedOutput()
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

	comm := exec.Command("micro", serv.EnvFlag())
	outp, err := comm.CombinedOutput()
	if err == nil {
		t.Fatal("Missing command should error")
	}

	if !strings.Contains(string(outp), "No command") {
		t.Fatalf("Unexpected output for no command: %v", string(outp))
	}

	comm = exec.Command("micro", serv.EnvFlag(), "asdasd")
	outp, err = comm.CombinedOutput()
	if err == nil {
		t.Fatal("Wrong command should error")
	}

	if !strings.Contains(string(outp), "No command provided to micro") {
		t.Fatalf("Unexpected output for unrecognized command: %v", string(outp))
	}

	comm = exec.Command("micro", serv.EnvFlag(), "config", "asdasd")
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

func TestUnrecognisedCommand(t *testing.T) {
	TrySuite(t, testUnrecognisedCommand, retryCount)
}

func testUnrecognisedCommand(t *T) {
	serv := NewServer(t)
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	t.Parallel()
	outp, _ := exec.Command("micro", serv.EnvFlag(), "foobar").CombinedOutput()
	if !strings.Contains(string(outp), "No command provided to micro. Please refer to 'micro --help'") {
		t.Fatalf("micro foobar does not return correct error %v", string(outp))
		return
	}
}

func TestPlatformErrorLocalSource(t *testing.T) {
	t.Parallel()
	// @todo reintroduce this test as a change after the creation of this test broke it
	return
	outp, _ := exec.Command("micro", "-env=platform", "run", "./service/example").CombinedOutput()
	if !strings.Contains(string(outp), "Local sources are not yet supported on m3o") {
		t.Fatalf("Local source does not return expected error %v", string(outp))
		return
	}
}
