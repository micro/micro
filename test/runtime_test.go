// +build integration

package test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
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

func TestServerModeCall(t *testing.T) {
	trySuite(t, testServerModeCall, retryCount)
}

func testServerModeCall(t *t) {
	t.Parallel()
	serv := newServer(t)

	callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}")
	outp, err := callCmd.CombinedOutput()
	if err == nil {
		t.Fatalf("Call to server should fail, got no error, output: %v", string(outp))
		return
	}

	serv.launch()
	defer serv.close()

	try("Calling Runtime.Read", t, func() ([]byte, error) {
		outp, err = exec.Command("micro", serv.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
		if err != nil {
			return outp, errors.New("Call to runtime read should succeed")
		}
		return outp, err
	}, 5*time.Second)
}

func TestRunLocalSource(t *testing.T) {
	trySuite(t, testRunLocalSource, retryCount)
}

func testRunLocalSource(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("Find test/example", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if !strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 15*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 50*time.Second)
}

func TestRunAndKill(t *testing.T) {
	trySuite(t, testRunAndKill, retryCount)
}

func testRunAndKill(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("Find test/example", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if !strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 15*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 50*time.Second)

	outp, err = exec.Command("micro", serv.envFlag(), "kill", "test/example-service").CombinedOutput()
	if err != nil {
		t.Fatalf("micro kill failure, output: %v", string(outp))
		return
	}

	try("Find test/example", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("Should not find example service in runtime")
		}
		return outp, err
	}, 15*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Should not find example service in list")
		}
		return outp, err
	}, 20*time.Second)
}

func TestLocalOutsideRepo(t *testing.T) {
	trySuite(t, testLocalOutsideRepo, retryCount)
}

func testLocalOutsideRepo(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	dirname := "last-dir-of-path"
	folderPath := filepath.Join(os.TempDir(), dirname)

	err := os.MkdirAll(folderPath, 0777)
	if err != nil {
		t.Fatal(err)
		return
	}

	// since copying a whole folder is rather involved and only Linux sources
	// are available, see https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
	// we fall back to `cp`
	outp, err := exec.Command("cp", "-r", "example-service/.", folderPath).CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	runCmd := exec.Command("micro", serv.envFlag(), "run", ".")
	runCmd.Dir = folderPath
	outp, err = runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("Find "+dirname, t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		lines := strings.Split(string(outp), "\n")
		found := false
		for _, line := range lines {
			if strings.HasPrefix(line, dirname) {
				found = true
			}
		}
		if !found {
			return outp, errors.New("Can't find '" + dirname + "' in runtime")
		}
		return outp, err
	}, 12*time.Second)

	try("Find go.micro.service.example in list", t, func() ([]byte, error) {
		outp, err := exec.Command("micro", serv.envFlag(), "list", "services").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "go.micro.service.example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 75*time.Second)
}

func TestLocalEnvRunGithubSource(t *testing.T) {
	//trySuite(t, testLocalEnvRunGithubSource, retryCount)
}

func testLocalEnvRunGithubSource(t *t) {
	t.Parallel()
	outp, err := exec.Command("micro", "env", "set", "local").CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to set env to local, err: %v, output: %v", err, string(outp))
		return
	}
	var cmd *exec.Cmd
	go func() {
		cmd = exec.Command("micro", "run", "location")
		// fire and forget as this will run forever
		cmd.CombinedOutput()
	}()
	time.Sleep(100 * time.Millisecond)
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	try("Find location", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", "status")
		outp, err := psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "location") {
			return outp, errors.New("Output should contain location")
		}
		return outp, nil
	}, 30*time.Second)
}

func TestRunGithubSource(t *testing.T) {
	trySuite(t, testRunGithubSource, retryCount)
}

func testRunGithubSource(t *t) {
	t.Parallel()
	p, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
		return
	}
	if len(p) == 0 {
		t.Fatal("Git is not available")
		return
	}
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "github.com/micro/examples/helloworld")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("Find hello world", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Output should contain hello world")
		}
		return outp, nil
	}, 60*time.Second)

	try("Call hello world", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.helloworld", "Helloworld.Call", `{"name": "Joe"}`)
		outp, err := callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello Joe" {
			return outp, errors.New("Helloworld resonse is unexpected")
		}
		return outp, err
	}, 60*time.Second)

}

func TestRunLocalUpdateAndCall(t *testing.T) {
	trySuite(t, testRunLocalUpdateAndCall, retryCount)
}

func testRunLocalUpdateAndCall(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	// Run the example service
	runCmd := exec.Command("micro", serv.envFlag(), "run", "./example-service")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("Finding example service with micro status", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "status")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "test/example-service",
		// as the runtime name is the relative path inside a repo.
		if !strings.Contains(string(outp), "test/example-service") {
			return outp, errors.New("can't find service in runtime")
		}
		return outp, err
	}, 15*time.Second)

	try("Call example service", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
		outp, err := callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello Joe" {
			return outp, errors.New("Response is unexpected")
		}
		return outp, err
	}, 50*time.Second)

	replaceStringInFile(t, "./example-service/handler/handler.go", "Hello", "Hi")
	defer func() {
		// Change file back
		replaceStringInFile(t, "./example-service/handler/handler.go", "Hi", "Hello")
	}()

	updateCmd := exec.Command("micro", serv.envFlag(), "update", "./example-service")
	outp, err = updateCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
		return
	}

	try("Call example service after modification", t, func() ([]byte, error) {
		callCmd := exec.Command("micro", serv.envFlag(), "call", "go.micro.service.example", "Example.Call", `{"name": "Joe"}`)
		outp, err = callCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hi Joe" {
			return outp, errors.New("Response is not what's expected")
		}
		return outp, err
	}, 15*time.Second)
}

func TestExistingLogs(t *testing.T) {
	trySuite(t, testExistingLogs, retryCount)
}

func testExistingLogs(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "github.com/crufter/micro-services/logspammer")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("logspammer logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "5", "crufter/micro-services/logspammer")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "never stopping") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 50*time.Second)
}

func TestStreamLogsAndThirdPartyRepo(t *testing.T) {
	trySuite(t, testStreamLogsAndThirdPartyRepo, retryCount)
}

func testStreamLogsAndThirdPartyRepo(t *t) {
	t.Parallel()
	serv := newServer(t)
	serv.launch()
	defer serv.close()

	runCmd := exec.Command("micro", serv.envFlag(), "run", "github.com/crufter/micro-services/logspammer")
	outp, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	try("logspammer logs", t, func() ([]byte, error) {
		psCmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "5", "crufter/micro-services/logspammer")
		outp, err = psCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "never stopping") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 50*time.Second)

	// Test streaming logs
	cmd := exec.Command("micro", serv.envFlag(), "logs", "-n", "1", "-f", "crufter-micro-services-logspammer")

	time.Sleep(7 * time.Second)

	go func() {
		outp, err := cmd.CombinedOutput()
		if err != nil {
			t.Log(err)
		}
		if len(outp) == 0 {
			t.Fatal("No log lines streamed")
			return
		}
		if !strings.Contains(string(outp), "never stopping") {
			t.Fatalf("Unexpected logs: %v", string(outp))
			return
		}
		// Logspammer logs every 2 seconds, so we need 2 different
		now := time.Now()
		// leaving the hour here to fix a docker issue
		// when the containers clock is a few hours behind
		stampA := now.Add(-2 * time.Second).Format("04:05")
		stampB := now.Add(-1 * time.Second).Format("04:05")
		if !strings.Contains(string(outp), stampA) && !strings.Contains(string(outp), stampB) {
			t.Fatalf("Timestamp %v or %v not found in logs: %v", stampA, stampB, string(outp))
			return
		}
	}()

	time.Sleep(7 * time.Second)

	err = cmd.Process.Kill()
	if err != nil {
		t.Fatal(err)
		return
	}
	time.Sleep(2 * time.Second)
}

func replaceStringInFile(t *t, filepath string, original, newone string) {
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
		return
	}

	output := strings.ReplaceAll(string(input), original, newone)
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		t.Fatal(err)
		return
	}
}
