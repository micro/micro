package test

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

const (
	retryCount = 2
	isParallel = true
)

var ignoreThisError = errors.New("Do not use this error")

type cmdFunc func() ([]byte, error)

// try is designed with command line executions in mind
// Error should be checked and a simple `return` from the test case should
// happen without calling `t.Fatal`. The error value should be disregarded.
func try(blockName string, t *t, f cmdFunc, maxTime time.Duration) error {
	start := time.Now()
	var outp []byte
	var err error

	for {
		if t.failed {
			return ignoreThisError
		}
		if time.Since(start) > maxTime {
			_, file, line, _ := runtime.Caller(1)
			fname := filepath.Base(file)
			if err != nil {
				t.Fatalf("%v:%v, %v (failed after %v with '%v'), output: '%v'", fname, line, blockName, time.Since(start), err, string(outp))
				return ignoreThisError
			}
			return nil
		}
		outp, err = f()
		if err == nil {
			return nil
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func once(blockName string, t *testing.T, f cmdFunc) {
	outp, err := f()
	if err != nil {
		t.Fatalf("%v with '%v', output: %v", blockName, err, string(outp))
	}
}

type server struct {
	cmd           *exec.Cmd
	t             *t
	envName       string
	portNum       int
	containerName string
	opts          options
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// taken from https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
// MyCaller returns the caller of the function that called it :)
func myCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return getFrame(2).Function
}

type options struct {
	auth string // eg. jwt
}

func newServer(t *t, opts ...options) server {
	min := 8000
	max := 60000
	portnum := rand.Intn(max-min) + min
	fname := strings.Split(myCaller(), ".")[2]

	// kill container, ignore error because it might not exist,
	// we dont care about this that much
	exec.Command("docker", "kill", fname).CombinedOutput()
	exec.Command("docker", "rm", fname).CombinedOutput()

	cmd := exec.Command("docker", "run", "--name", fname,
		fmt.Sprintf("-p=%v:8081", portnum), "micro", "server")
	if len(opts) == 1 && opts[0].auth == "jwt" {

		base64 := "base64 -w0"
		if runtime.GOOS == "darwin" {
			base64 = "base64 -b0"
		}
		priv := "cat /tmp/sshkey | " + base64
		privKey, err := exec.Command("bash", "-c", priv).Output()
		if err != nil {
			panic(string(privKey))
		}

		pub := "cat /tmp/sshkey.pub | " + base64
		pubKey, err := exec.Command("bash", "-c", pub).Output()
		if err != nil {
			panic(string(pubKey))
		}
		cmd = exec.Command("docker", "run", "--name", fname,
			fmt.Sprintf("-p=%v:8081", portnum),
			"-e", "MICRO_AUTH=jwt",
			"-e", "MICRO_AUTH_PRIVATE_KEY="+strings.Trim(string(privKey), "\n"),
			"-e", "MICRO_AUTH_PUBLIC_KEY="+strings.Trim(string(pubKey), "\n"),
			"micro", "server")
	}
	//fmt.Println("docker", "run", "--name", fname, fmt.Sprintf("-p=%v:8081", portnum), "micro", "server")
	opt := options{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	return server{
		cmd:           cmd,
		t:             t,
		envName:       fname,
		containerName: fname,
		portNum:       portnum,
		opts:          opt,
	}
}

// error value should not be used but caller should return in the test suite
// in case of error.
func (s server) launch() error {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.t.Fatal(err)
		}
	}()

	// add the environment
	if err := try("Adding micro env", s.t, func() ([]byte, error) {
		outp, err := exec.Command("micro", "env", "add", s.envName, fmt.Sprintf("127.0.0.1:%v", s.portNum)).CombinedOutput()
		if err != nil {
			return outp, err
		}
		if len(outp) > 0 {
			return outp, errors.New("Not added")
		}

		outp, err = exec.Command("micro", "env").CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), s.envName) {
			return outp, errors.New("Not added")
		}

		return outp, nil
	}, 15*time.Second); err != nil {
		return err
	}

	if err := try("Calling micro server", s.t, func() ([]byte, error) {
		outp, err := exec.Command("micro", s.envFlag(), "list", "services").CombinedOutput()
		if !strings.Contains(string(outp), "runtime") ||
			!strings.Contains(string(outp), "registry") ||
			!strings.Contains(string(outp), "api") ||
			!strings.Contains(string(outp), "broker") ||
			!strings.Contains(string(outp), "config") ||
			!strings.Contains(string(outp), "debug") ||
			!strings.Contains(string(outp), "proxy") ||
			!strings.Contains(string(outp), "auth") ||
			!strings.Contains(string(outp), "store") {
			return outp, errors.New("Not ready")
		}

		return outp, err
	}, 60*time.Second); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return nil
}

func (s server) close() {
	exec.Command("docker", "kill", s.containerName).CombinedOutput()
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGKILL)
	}
}

func (s server) envFlag() string {
	return fmt.Sprintf("-env=%v", s.envName)
}

type t struct {
	counter int
	failed  bool
	format  string
	values  []interface{}
	t       *testing.T
}

func (t *t) Fatal(values ...interface{}) {
	t.t.Helper()
	t.t.Log(values...)
	t.failed = true
	t.values = values
}

func (t *t) Log(values ...interface{}) {
	t.t.Helper()
	t.t.Log(values...)
}

func (t *t) Fatalf(format string, values ...interface{}) {
	t.t.Helper()
	t.t.Log(fmt.Sprintf(format, values...))
	t.failed = true
	t.values = values
	t.format = format
}

func (t *t) Parallel() {
	if t.counter == 0 && isParallel {
		t.t.Parallel()
	}
	t.counter++
}

func newT(te *testing.T) *t {
	return &t{t: te}
}

// trySuite is designed to retry a TestXX function
func trySuite(t *testing.T, f func(t *t), times int) {
	t.Helper()
	tee := newT(t)
	for i := 0; i < times; i++ {
		f(tee)
		if !tee.failed {
			return
		}
		if i != times-1 {
			tee.failed = false
		}
		time.Sleep(200 * time.Millisecond)
	}
	if tee.failed {
		if len(tee.format) > 0 {
			t.Fatalf(tee.format, tee.values...)
		} else {
			t.Fatal(tee.values...)
		}
	}
}
