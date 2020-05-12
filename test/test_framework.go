package test

import (
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
	retryCount = 1
	isParallel = true
)

type cmdFunc func() ([]byte, error)

// try is designed with command line executions in mind
func try(blockName string, t *t, f cmdFunc, maxTime time.Duration) {
	start := time.Now()
	var outp []byte
	var err error

	for {
		if time.Since(start) > maxTime {
			_, file, line, _ := runtime.Caller(1)
			fname := filepath.Base(file)
			if err != nil {
				t.Fatalf("%v:%v, %v (failed after %v with '%v'), output: '%v'", fname, line, blockName, time.Since(start), err, string(outp))
			}
			return
		}
		outp, err = f()
		if err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func once(blockName string, t *testing.T, f cmdFunc) {
	outp, err := f()
	if err != nil {
		t.Fatalf("%v with '%v', output: %v", blockName, err, string(outp))
	}
}

type server struct {
	cmd       *exec.Cmd
	t         *t
	proxyPort int
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

func newServer(t *t) server {
	min := 8000
	max := 60000
	portnum := rand.Intn(max-min) + min
	fname := strings.Split(myCaller(), ".")[2]
	exec.Command("docker", "rm", fname).CombinedOutput()

	return server{
		cmd: exec.Command("docker", "run", "--network", "alpine-net", "--name", fname,
			fmt.Sprintf("-p=%v:8081", portnum), "micro", "server"),
		t:         t,
		proxyPort: portnum,
	}
}

func (s server) launch() {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.Fatal(err)
		}
	}()
	// @todo find a way to know everything is up and running
	try("Calling micro server", s.t, func() ([]byte, error) {
		return exec.Command("micro", s.envFlag(), "call", "go.micro.runtime", "Runtime.Read", "{}").CombinedOutput()
	}, 10000*time.Millisecond)
	time.Sleep(5 * time.Second)
}

func (s server) close() {
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGTERM)
	}
}

func (s server) envFlag() string {
	return fmt.Sprintf("-env=127.0.0.1:%v", s.proxyPort)
}

func (s server) trySuite() {

}

type t struct {
	counter int
	failed  bool
	format  string
	values  []interface{}
	t       *testing.T
}

func (t *t) Fatal(values ...interface{}) {
	//t.Log(values...)
	t.failed = true
	t.values = values
}

func (t *t) Log(values ...interface{}) {
	t.t.Log(values...)
}

func (t *t) Fatalf(format string, values ...interface{}) {
	//t.Log(fmt.Sprintf(format, values...))
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
	tee := newT(t)
	for i := 0; i < times; i++ {
		f(tee)
		if !tee.failed {
			return
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
