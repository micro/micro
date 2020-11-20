// +build integration

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	branch  = ""
	version = "latest"
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

	outp, err := cmd.Exec("run", "./services/helloworld")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find helloworld service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 90*time.Second); err != nil {
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

	outp, err := cmd.Exec("run", "./services/helloworld")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 30*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "helloworld") {
			outp1, _ := cmd.Exec("logs", "helloworld")
			return append(outp, outp1...), errors.New("Can't find helloworld service in list")
		}
		return outp, err
	}, 90*time.Second); err != nil {
		return
	}

	outp, err = cmd.Exec("kill", "helloworld")
	if err != nil {
		t.Fatalf("micro kill failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Should not find example service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find helloworld in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if strings.Contains(string(outp), "helloworld") {
			return outp, errors.New("Should not find helloworld service in list")
		}
		return outp, err
	}, 20*time.Second); err != nil {
		return
	}
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
	outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "github.com/micro/services/helloworld@"+branch)
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", version, outp) {
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

	cmd.Exec("kill", "helloworld")

	// test it works for a branch with a funny name
	outp, err = cmd.Exec("run", "--image", "localhost:5000/cells:v3", "github.com/micro/services/helloworld@integrationtest/branch_name")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}
		//
		if !statusRunning("helloworld", "integrationtest/branch_name", outp) {
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

	if err := Try("Find basic-micro-service in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("basic-micro-service", version, outp) {
			return outp, errors.New("Output should contain basic-micro-service")
		}
		return outp, nil
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find example in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "example") {
			return outp, errors.New("Does not example")
		}
		return outp, err
	}, 120*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "basic-micro-service")
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

	outp, err := cmd.Exec("run", "subfolder-test"+branch)
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", version, outp) {
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

	if err := Try("Find subfolder-test in runtime", t, func() ([]byte, error) {
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

	if err := Try("Find example in registry", t, func() ([]byte, error) {
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
	outp, err := cmd.Exec("run", "./services/helloworld")
	if err != nil {
		t.Fatalf("micro run failure, output: %v", string(outp))
		return
	}

	if err := Try("Finding helloworld service with micro status", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "example".
		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("can't find service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Call helloworld service", t, func() ([]byte, error) {
		outp, err := cmd.Exec("helloworld", "--name=Joe")
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

	replaceStringInFile(t, "./services/helloworld/handler/helloworld.go", `"Hello "`, `"Hi "`)
	defer func() {
		// Change file back
		replaceStringInFile(t, "./services/helloworld/handler/helloworld.go", `"Hi "`, `"Hello "`)
	}()

	outp, err = cmd.Exec("update", "./services/helloworld")
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := Try("Finding helloworld service with micro status", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "example".
		if !statusRunning("helloworld", "latest", outp) {
			outp1, _ := cmd.Exec("logs", "helloworld")
			return append(outp, outp1...), errors.New("can't find service in runtime")
		}
		return outp, err
	}, 45*time.Second); err != nil {
		return
	}

	if err := Try("Call helloworld service after modification", t, func() ([]byte, error) {
		outp, err := cmd.Exec("helloworld", "--name=Joe")
		if err != nil {
			outp1, _ := cmd.Exec("logs", "helloworld")
			return append(outp, outp1...), err
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
	}, 45*time.Second); err != nil {
		return
	}
}

func TestRunCurrentFolder(t *testing.T) {
	TrySuite(t, testRunCurrentFolder, retryCount)
}

func testRunCurrentFolder(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	usr, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Dir = usr.HomeDir
	err = os.RemoveAll(filepath.Join(usr.HomeDir, "helloworld"))
	if err != nil {
		t.Fatal(err)
	}
	//if err != nil {
	//	t.Fatal(string(outp))
	//}

	outp, err := cmd.Exec("new", "helloworld")
	if err != nil {
		t.Fatal(string(outp))
	}
	makeProt := exec.Command("make", "proto")
	makeProt.Dir = filepath.Join(usr.HomeDir, "helloworld")
	outp, err = makeProt.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	cmd.Dir = filepath.Join(usr.HomeDir, "helloworld")
	outp, err = cmd.Exec("run", ".")
	if err != nil {
		t.Fatal(outp)
	}

	Try("Find helloworld", t, func() ([]byte, error) {
		outp, err = cmd.Exec("status")
		if !statusRunning("helloworld", "latest", outp) {
			return outp, errors.New("Can't find helloworld")
		}
		return outp, err
	}, 20*time.Second)
}

func TestRunParentFolder(t *testing.T) {
	TrySuite(t, testRunParentFolder, retryCount)
}

func testRunParentFolder(t *T) {
	defer func() {
		os.RemoveAll("../test-top-level")
	}()
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	cmd.Dir = ".."
	outp, err := cmd.Exec("new", "test-top-level")
	if err != nil {
		t.Fatal(string(outp))
	}
	makeProt := exec.Command("make", "proto")
	makeProt.Dir = "../test-top-level"
	outp, err = makeProt.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	gomod := exec.Command("go", "mod", "edit", "-replace", "github.com/micro/micro/v3=github.com/micro/micro/v3@v3.0.0")
	gomod.Dir = "../test-top-level"
	if outp, err := gomod.CombinedOutput(); err != nil {
		t.Fatal(string(outp))
	}

	err = os.MkdirAll("../parent/folder/test", 0777)
	if err != nil {
		t.Fatal(err)
	}

	cmd.Dir = "../parent/folder/test"
	outp, err = cmd.Exec("run", "../../../test-top-level")
	if err != nil {
		t.Fatal(string(outp))
	}

	if err := Try("Find example", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("test-top-level", "latest", outp) {
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
		if !strings.Contains(string(outp), "test-top-level") {
			l, _ := cmd.Exec("logs", "test-top-level")
			return outp, fmt.Errorf("Can't find example service in list. \nLogs: %v", string(l))
		}
		return outp, err
	}, 90*time.Second); err != nil {
		return
	}
}

func TestRunNewWithGit(t *testing.T) {
	TrySuite(t, testRunNewWithGit, retryCount)
}

func testRunNewWithGit(t *T) {
	defer func() {
		os.RemoveAll("/tmp/new-with-git")
	}()
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	cmd.Dir = "/tmp"

	outp, err := cmd.Exec("new", "new-with-git")
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	makeProt := exec.Command("make", "proto")
	makeProt.Dir = "/tmp/new-with-git"

	outp, err = makeProt.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
		return
	}

	// for tests, update the micro import to use the current version of the code.
	fname := fmt.Sprintf(makeProt.Dir + "/go.mod")
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(string(outp))
		return
	}
	if _, err := f.WriteString("\nreplace github.com/micro/micro/v3 => github.com/micro/micro/v3 master"); err != nil {
		t.Fatal(string(outp))
		return
	}
	// This should point to master, but GOPROXY is not on in the runtime. Remove later.
	if _, err := f.WriteString("\nreplace github.com/micro/go-micro/v3 => github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200922112322-927d4f8eced6"); err != nil {
		t.Fatal(string(outp))
		return
	}
	f.Close()

	gitInit := exec.Command("git", "init")
	gitInit.Dir = "/tmp/new-with-git"
	outp, err = gitInit.CombinedOutput()
	if err != nil {
		t.Fatal(string(outp))
	}

	cmd = serv.Command()
	cmd.Dir = "/tmp/new-with-git"
	outp, err = cmd.Exec("run", ".")
	if err != nil {
		t.Fatal(outp)
	}

	if err := Try("Find example", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		outp1, _ := cmd.Exec("logs", "new-with-git")
		outp = append(outp, outp1...)
		if err != nil {
			return outp, err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("new-with-git", "latest", outp) {
			return outp, errors.New("Can't find example service in runtime")
		}
		return outp, err
	}, 15*time.Second); err != nil {
		return
	}

	if err := Try("Find example in list", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		outp1, _ := cmd.Exec("logs", "new-with-git")
		outp = append(outp, outp1...)
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "new-with-git") {
			return outp, errors.New("Can't find example service in list")
		}
		return outp, err
	}, 90*time.Second); err != nil {
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

	outp, err := cmd.Exec("run", "./services/test/logger")
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
	}, 90*time.Second); err != nil {
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

	outp, err := cmd.Exec("run", "./services/test/logger")
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
	}, 60*time.Second); err != nil {
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

	outp, err := cmd.Exec("run", "./services/test/logger")
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
	}, 90*time.Second); err != nil {
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
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "github.com/micro/test/helloworld"); err != nil {
		t.Fatalf("Expected no run error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find helloworld in runtime", t, func() ([]byte, error) {
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
	}, 90*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "helloworld")
		t.Logf("logs %s", string(outp))
		return
	}

	// call the service
	if err := Try("Calling helloworld", t, func() ([]byte, error) {
		return cmd.Exec("helloworld", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}
}

func TestRunCustomCredentials(t *testing.T) {
	TrySuite(t, testRunCustomCredentials, retryCount)
}

func testRunCustomCredentials(t *T) {
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
	if outp, err := cmd.Exec("user", "config", "set", "git.credentials.url", "github.com"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if outp, err := cmd.Exec("user", "config", "set", "git.credentials.token", pat); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// run the service
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "github.com/micro/test/helloworld"+branch); err != nil {
		t.Fatalf("Expected no run error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("helloworld", version, outp) {
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
		return cmd.Exec("helloworld", "--name=John")
	}, 30*time.Second); err != nil {
		return
	}
}

func TestGitSourceUpdateByShortName(t *testing.T) {
	TrySuite(t, testGitSourceUpdateByShortName, retryCount)
}

func testGitSourceUpdateByShortName(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()

	// run the service
	if outp, err := cmd.Exec("run", "github.com/m3o/services/invite"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find invite in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("invite", version, outp) {
			return outp, errors.New("Can't find subfolder-test service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find invite in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "invite") {
			return outp, errors.New("Does not contain invite")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "invite")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling invite", t, func() ([]byte, error) {
		outp, _ := cmd.Exec("logs", "-n=1000", "invite")
		outp1, err := cmd.Exec("invite", "--help")

		return append(outp1, outp...), err
	}, 70*time.Second); err != nil {
		return
	}

	// update service

	// run the service
	if outp, err := cmd.Exec("update", "invite"+branch); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	// look for disconnect
	if err := Try("Find invite disconnect proof in logs", t, func() ([]byte, error) {
		outp, err := cmd.Exec("logs", "invite")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "Disconnect") {
			return outp, errors.New("Does not contain Disconnect")
		}
		return outp, err
	}, 40*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "invite")
		t.Log(string(outp))
		return
	}

	if err := Try("Find invite in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("invite", version, outp) {
			return outp, errors.New("Can't find invite service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find invite in registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("services")
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "invite") {
			return outp, errors.New("Does not contain example")
		}
		return outp, err
	}, 300*time.Second); err != nil {
		outp, _ := cmd.Exec("logs", "invite")
		t.Log(string(outp))
		return
	}

	// call the service
	if err := Try("Calling invite", t, func() ([]byte, error) {
		outp, _ := cmd.Exec("logs", "-n=1000", "invite")
		outp1, err := cmd.Exec("invite", "--help")

		return append(outp1, outp...), err
	}, 70*time.Second); err != nil {
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
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "gitlab.com/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find subfolder-test in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", version, outp) {
			return outp, errors.New("Can't find subfolder-test service in runtime")
		}
		return outp, err
	}, 60*time.Second); err != nil {
		return
	}

	if err := Try("Find example in registry", t, func() ([]byte, error) {
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
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "bitbucket.org/micro-test/private-monorepo-test/subfolder-test"); err != nil {
		t.Fatalf("Expected no error, got %v %v", err, string(outp))
		return
	}

	if err := Try("Find micro/test/helloworld in runtime", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		if err != nil {
			return outp, err
		}

		if !statusRunning("subfolder-test", version, outp) {
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
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", "gitlab.com/micro-test/private-monorepo-test/subfolder-test"+branch); err != nil {
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

		if !statusRunning("subfolder-test", version, outp) {
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

func TestIdiomaticFolderStructure(t *testing.T) {
	TrySuite(t, testIdiomaticFolderStructure, retryCount)
}

func testIdiomaticFolderStructure(t *T) {
	t.Parallel()
	serv := NewServer(t, WithLogin())
	defer serv.Close()
	if err := serv.Run(); err != nil {
		return
	}

	cmd := serv.Command()
	src := "./services/test/template"
	if outp, err := cmd.Exec("run", "--image", "localhost:5000/cells:v3", src); err != nil {
		t.Fatalf("Error running service: %v, %v", err, string(outp))
		return
	}

	if err := Try("Find template service in the registry", t, func() ([]byte, error) {
		outp, err := cmd.Exec("status")
		outp1, _ := cmd.Exec("logs", "template")
		if err != nil {
			return append(outp, outp1...), err
		}

		// The started service should have the runtime name of "service/example",
		// as the runtime name is the relative path inside a repo.
		if !statusRunning("template", "latest", outp) {
			return outp, errors.New("Can't find template service in runtime")
		}
		return outp, err
	}, 120*time.Second); err != nil {
		return
	}
}
