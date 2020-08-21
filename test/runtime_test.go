// +build integration

package test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestServerModeCall(t *testing.T) {
	TrySuite(t, ServerModeCall, retryCount)
}

func ServerModeCall(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())

	cmd := serv.Command()

	outp, err := cmd.Exec("call", "runtime", "Runtime.Read", "{}")
	if err == nil {
		t.Fatalf("Call to server should fail, got no error, output: %v", string(outp))
		return
	}

	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	if err := Try("Calling Runtime.Read", t, func() ([]byte, error) {
		outp, err := cmd.Exec("call", "runtime", "Runtime.Read", "{}")
		if err != nil {
			return outp, errors.New("Call to runtime read should succeed")
		}
		return outp, err
	}, 5*time.Second); err != nil {
		return
	}
}

func TestRunLocalSource(t *testing.T) {
	TrySuite(t, testRunLocalSource, retryCount)
}

func testRunLocalSource(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "./service/example")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find test/example", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("example", "latest", outp) {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find example in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 50*time.Second); err != nil {
		return
	}
}

func TestRunAndKill(t *testing.T) {
	TrySuite(t, testRunAndKill, retryCount)
}

func testRunAndKill(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "./service/example")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find test/example", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("example", "latest", outp) {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find example in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 50*time.Second); err != nil {
		return
	}

	outp, err = cmd.Exec("kill", "service/example")
	if err != nil {
		t.Fatalf("micro kill failure, output: %v", string(outp))
		return
	}

	if err := Try("Find test/example", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if strings.Contains(string(outp), "service/example") {
			return outp, errors.New("Should not find example service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find example in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if strings.Contains(string(outp), "example") {
			return outp, errors.New("Should not find example service in list")
		}
		return outp, err
	}, 20*time.Second); err != nil {
		return
	}
}

func statusRunning(service, branch string, statusOutput []byte) bool {
	reg, _ := regexp.Compile(service + "\\s+" + branch + "\\s+\\S+\\s+running")
	return reg.Match(statusOutput)
}

func TestRunGithubSource(t *testing.T) {
	TrySuite(t, testRunGithubSource, retryCount)
}

func testRunGithubSource(t *T) {
	t.Parallel()

	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "github.com/micro/services/helloworld@master")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", "master", outp) {
			return outp, errors.New("Output should contain helloworld")
		}
		if !strings.Contains(string(outp), "owner=admin") || !(strings.Contains(string(outp), "group=micro") || strings.Contains(string(outp), "group="+serv.Env())) {
			return outp, errors.New("micro status does not have correct owner or group")
		}
		if strings.Contains(string(outp), "unknown") {
			return outp, errors.New("there should be no unknown in the micro status output")
		}
		return outp, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err = cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("helloworld is not running")
		}
		return outp, nil
	}, 180*time.Second); err != nil {
		return
	}

	if err := Try("Call helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("call", "helloworld", "Helloworld.Call", "{\"name\":\"John\"}")
		if err != nil {
			return outp, err
		}
		rsp := map[string]string{}
		err = json.Unmarshal(outp, &rsp)
		if err != nil {
			return outp, err
		}
		if rsp["msg"] != "Hello John" {
			return outp, errors.New("Helloworld resonse is unexpected")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

}

// Note: @todo this method should truly be the same as TestGithubSource.
func TestRunGitlabSource(t *testing.T) {
	TrySuite(t, testRunGitlabSource, retryCount)
}

func testRunGitlabSource(t *T) {
	t.Parallel()

	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	cmd.Exec("user", "config", "set", "git."+serv.Env()+".baseurl", "gitlab.com/micro-test")

	outp, err := cmd.Exec("run", "basic-micro-service")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("basic-micro-service", "latest", outp) {
			return outp, errors.New("Output should contain basic-micro-service")
		}
		return outp, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not example")
		}
		return outp, err
	}, 120*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "basic-micro-serviced")
		t.Log(string(outp))
		return
	}
}

func TestRunGitlabSourceMonoRepo(t *testing.T) {
	TrySuite(t, testRunGitlabSourceMonoRepo, retryCount)
}

func testRunGitlabSourceMonoRepo(t *T) {
	t.Parallel()

	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	cmd.Exec("user", "config", "set", "git."+serv.Env()+".baseurl", "gitlab.com/micro-test/monorepo-test")

	outp, err := cmd.Exec("run", "subfolder-test")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Output should contain subfolder-test")
		}
		return outp, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 120*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}
}

// This test exists to test the path of "generic git checkout", not just bitbucket
func TestRunGenericRemote(t *testing.T) {
	TrySuite(t, testRunGenericRemote, retryCount)
}

func testRunGenericRemote(t *T) {
	t.Parallel()

	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	cmd.Exec("user", "config", "set", "git."+serv.Env()+".baseurl", "bitbucket.org/micro-test/monorepo-test")

	outp, err := cmd.Exec("run", "subfolder-test")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Output should contain subfolder-test")
		}
		return outp, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 120*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}
}

func TestRunLocalUpdateAndCall(t *testing.T) {
	TrySuite(t, testRunLocalUpdateAndCall, retryCount)
}

func testRunLocalUpdateAndCall(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// Run the example service
	outp, err := cmd.Exec("run", "./service/example")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Finding example service with micro status", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "example".
		if !statusRunning("example", "latest", outp) {
			return outp, errors.New("can't find service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Call example service", t, func() ([]byte, error) {
		outp, err := cmd.Exec("call", "example", "Example.Call", `{"name": "Joe"}`)
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
	}, 50*time.Second); err != nil {
		return
	}

	replaceStringInFile(t, "./service/example/handler/handler.go", "Hello", "Hi")
	defer func() {
		// Change file back
		replaceStringInFile(t, "./service/example/handler/handler.go", "Hi", "Hello")
	}()

	outp, err = cmd.Exec("update", "./service/example")
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := Try("Call example service after modification", t, func() ([]byte, error) {
		outp, err := cmd.Exec("call", "example", "Example.Call", `{"name": "Joe"}`)
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
	}, 15*time.Second); err != nil {
		return
	}
}

func TestExistingLogs(t *testing.T) {
	TrySuite(t, testExistingLogs, retryCount)
}

func testExistingLogs(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "./service/logger")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("logger logs", t, func() ([]byte, error) {
		outp, err = cmd.Exec("logs", "logger")
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "This is a log line") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 50*time.Second); err != nil {
		return
	}
}

func TestBranchCheckout(t *testing.T) {
	TrySuite(t, testBranchCheckout, retryCount)
}

func testBranchCheckout(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "github.com/micro/micro/test/service/logger@master")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("logger logs", t, func() ([]byte, error) {
		outp, err = cmd.Exec("logs", "logger")
		if err != nil {
			return outp, err
		}

		// The log that this branch outputs is different from master, that's what we look for
		if !strings.Contains(string(outp), "Listening on") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 30*time.Second); err != nil {
		return
	}
}

func TestStreamLogsAndThirdPartyRepo(t *testing.T) {
	TrySuite(t, testStreamLogsAndThirdPartyRepo, retryCount)
}

func testStreamLogsAndThirdPartyRepo(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "github.com/micro/micro/test/service/logger")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("logger logs", t, func() ([]byte, error) {
		outp, err = cmd.Exec("logs", "logger")
		if err != nil {
			return outp, err
		}

		if !strings.Contains(string(outp), "Listening on") || !strings.Contains(string(outp), "This is a log line") {
			return outp, errors.New("Output does not contain expected")
		}
		return outp, nil
	}, 50*time.Second); err != nil {
		return
	}

	cmd.Start("logs", "-n", "1", "-f", "logger")

	time.Sleep(7 * time.Second)

	go func() {
		outp, err := cmd.Output()
		if err != nil {
			t.Log(err)
		}
		if len(outp) == 0 {
			t.Fatal("No log lines streamed")
			return
		}
		if !strings.Contains(string(outp), "This is a log line") {
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

	if err := cmd.Stop(); err != nil {
		t.Fatal(err)
		return
	}

	time.Sleep(2 * time.Second)
}

func replaceStringInFile(t *T, filepath string, original, newone string) {
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

func TestParentDependency(t *testing.T) {
	TrySuite(t, testParentDependency, retryCount)
}

func testParentDependency(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	outp, err := cmd.Exec("run", "./dep-test/dep-test-service")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find dep-test-service", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("dep-test-service", "latest", outp) {
			return outp, errors.New("Output should contain dep-test-service")
		}
		return outp, nil
	}, 30*time.Second); err != nil {
		return
	}
}

func TestRunPrivateSource(t *testing.T) {
	TrySuite(t, testRunPrivateSource, retryCount)
}

func testRunPrivateSource(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// get the git credentials, injected by the k8s integration test
	pat := os.Getenv("GITHUB_PAT")
	if len(pat) == 0 {
		t.Logf("Skipping test, missing GITHUB_PAT")
		return
	}

	// set the pat in the users config
	if outp, err := cmd.Exec("user", "config", "set", "git.credentials.github", pat); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:go", "github.com/micro/test/helloworld"); err != nil {
		t.Fatalf("Expected no run error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Does not contain helloworld")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "helloworld")
		t.Logf("logs %s", string(outp))
		return
	}

	// call the service
	if err := Try("Calling helloworld", t, func() ([]byte, error) {
		outp, _ := cmd.Exec("logs", "helloworld")
		t.Logf("logs %s", string(outp))
		return cmd.Exec("helloworld", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}
}

func TestRunPrivateGitlabSource(t *testing.T) {
	TrySuite(t, testRunPrivateGitlabSource, retryCount)
}

func testRunPrivateGitlabSource(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// get the git credentials, injected by the k8s integration test
	pat := os.Getenv("GITLAB_PAT")
	if len(pat) == 0 {
		t.Logf("Skipping test, missing GITLAB_PAT")
		return
	}

	// set the pat in the users config
	if outp, err := cmd.Exec("user", "config", "set", "git.credentials.gitlab", pat); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:go", "gitlab.com/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling example", t, func() ([]byte, error) {
		return cmd.Exec("example", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}

	// update service

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:go", "gitlab.com/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// We need a better way of knowing the server restarted than a sleep...
	time.Sleep(5 * time.Second)

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling example", t, func() ([]byte, error) {
		return cmd.Exec("example", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}
}

func TestRunPrivateGenericRemote(t *testing.T) {
	TrySuite(t, testRunPrivateGenericRemote, retryCount)
}

func testRunPrivateGenericRemote(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// get the git credentials, injected by the k8s integration test
	pat := os.Getenv("BITBUCKET_PAT")
	if len(pat) == 0 {
		t.Logf("Skipping test, missing BITBUCKET_PAT")
		return
	}

	// set the pat in the users config
	if outp, err := cmd.Exec("user", "config", "set", "git.credentials.bitbucket", pat); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:go", "bitbucket.org/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling example", t, func() ([]byte, error) {
		return cmd.Exec("example", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}

	// update service

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:go", "gitlab.com/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// We need a better way of knowing the server restarted than a sleep...
	time.Sleep(5 * time.Second)

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "subfolder-test")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling example", t, func() ([]byte, error) {
		return cmd.Exec("example", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}
}
