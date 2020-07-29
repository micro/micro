package test

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/token"
)

const (
	minPort = 8000
	maxPort = 60000
)

var (
	retryCount        = 2
	isParallel        = true
	ignoreThisError   = errors.New("Do not use this error")
	errFatal          = errors.New("Fatal error")
	testFilter        = []string{}
	maxTimeMultiplier = time.Duration(1)
)

type cmdFunc func() ([]byte, error)

type Server interface {
	Run() error
	Close()
	EnvFlag() string
	EnvName() string
}

// try is designed with command line executions in mind
// Error should be checked and a simple `return` from the test case should
// happen without calling `t.Fatal`. The error value should be disregarded.
func Try(blockName string, t *T, f cmdFunc, maxTime time.Duration) error {
	// hack. k8s can be slow locally
	maxTime = maxTimeMultiplier * maxTime

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

type ServerBase struct {
	cmd           *exec.Cmd
	t             *T
	envNm         string
	portNum       int
	containerName string
	opts          Options
	namespace     string
}

func (s *ServerBase) EnvName() string {
	return s.envNm
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

type Options struct {
	Login bool
}

type Option func(o *Options)

func WithLogin() Option {
	return func(o *Options) {
		o.Login = true
	}
}

func NewServer(t *T, opts ...Option) Server {
	fname := strings.Split(myCaller(), ".")[2]
	return newSrv(t, fname, opts...)
}

type NewServerFunc func(t *T, fname string, opts ...Option) Server

var newSrv NewServerFunc = newLocalServer

type ServerDefault struct {
	ServerBase
}

func newLocalServer(t *T, fname string, opts ...Option) Server {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	portnum := rand.Intn(maxPort-minPort) + minPort

	// kill container, ignore error because it might not exist,
	// we dont care about this that much
	exec.Command("docker", "kill", fname).CombinedOutput()
	exec.Command("docker", "rm", fname).CombinedOutput()

	// setup JWT keys
	base64 := "base64 -w0"
	if runtime.GOOS == "darwin" {
		base64 = "base64 -b0"
	}
	priv := "cat /tmp/sshkey | " + base64
	privKey, err := exec.Command("bash", "-c", priv).Output()
	if err != nil {
		panic(string(privKey))
	} else if len(strings.TrimSpace(string(privKey))) == 0 {
		panic("privKey has not been set")
	}

	pub := "cat /tmp/sshkey.pub | " + base64
	pubKey, err := exec.Command("bash", "-c", pub).Output()
	if err != nil {
		panic(string(pubKey))
	} else if len(strings.TrimSpace(string(pubKey))) == 0 {
		panic("pubKey has not been set")
	}

	// run the server
	cmd := exec.Command("docker", "run", "--name", fname,
		fmt.Sprintf("-p=%v:8081", portnum),
		"-e", "MICRO_AUTH_PRIVATE_KEY="+strings.Trim(string(privKey), "\n"),
		"-e", "MICRO_AUTH_PUBLIC_KEY="+strings.Trim(string(pubKey), "\n"),
		"-e", "MICRO_PROFILE=ci",
		"micro", "server")

	return &ServerDefault{ServerBase{
		cmd:           cmd,
		t:             t,
		envNm:         fname,
		containerName: fname,
		portNum:       portnum,
		opts:          options,
		namespace:     "micro",
	}}
}

// error value should not be used but caller should return in the test suite
// in case of error.
func (s *ServerBase) Run() error {
	go func() {
		if err := s.cmd.Start(); err != nil {
			s.t.Fatal(err)
		}
	}()

	// add the environment
	if err := Try("Adding micro env", s.t, func() ([]byte, error) {
		outp, err := exec.Command("micro", "env", "add", s.EnvName(), fmt.Sprintf("127.0.0.1:%v", s.portNum)).CombinedOutput()
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
		if !strings.Contains(string(outp), s.EnvName()) {
			return outp, errors.New("Not added")
		}

		return outp, nil
	}, 15 * time.Second); err != nil {
		return err
	}

	return nil
}

func (s *ServerDefault) Run() error {
	if err := s.ServerBase.Run(); err != nil {
		return err
	}

	if err := Try("Calling micro server", s.t, func() ([]byte, error) {
		outp, err := exec.Command("micro", s.EnvFlag(), "services").CombinedOutput()
		if !strings.Contains(string(outp), "runtime") ||
			!strings.Contains(string(outp), "registry") ||
			!strings.Contains(string(outp), "broker") ||
			!strings.Contains(string(outp), "config") ||
			!strings.Contains(string(outp), "debug") ||
			!strings.Contains(string(outp), "proxy") ||
			!strings.Contains(string(outp), "auth") ||
			!strings.Contains(string(outp), "store") {
			return outp, errors.New("Not ready")
		}

		return outp, err
	}, 60 * time.Second); err != nil {
		return err
	}

	// login to admin account
	if s.opts.Login {
		Login(s, s.t, "default", "password")
	}

	// // generate a new admin account for the env : user=ENV_NAME pass=password
	// req := fmt.Sprintf(`{"id":"%s", "secret":"password", "options":{"namespace":"%s"}}`, s.EnvName(), s.namespace)
	// outp, err := exec.Command("micro", s.EnvFlag(), "call", "go.micro.auth", "Auth.Generate", req).CombinedOutput()
	// if err != nil && !strings.Contains(string(outp), "already exists") { // until auth.Delete is implemented
	// 	s.t.Fatalf("Error generating auth: %s, %s", err, outp)
	// 	return err
	// }

	return nil
}

func (s *ServerBase) Close() {
	// remove the credentials so they aren't reused on next run
	token.Remove(s.EnvName())

	// reset back to the default namespace
	namespace.Set("micro", s.EnvName())

}

func (s *ServerDefault) Close() {
	s.ServerBase.Close()
	exec.Command("docker", "kill", s.containerName).CombinedOutput()
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGKILL)
	}
}

func (s *ServerBase) EnvFlag() string {
	return fmt.Sprintf("-env=%v", s.EnvName())
}

type T struct {
	counter int
	failed  bool
	format  string
	values  []interface{}
	t       *testing.T
}

// Fatal logs and exits immediately. Assumes it has come from a TrySuite() call. If called from within goroutine it does not immediately exit.
func (t *T) Fatal(values ...interface{}) {
	t.t.Helper()
	t.t.Log(values...)
	t.failed = true
	t.values = values
	doPanic()
}

func (t *T) Log(values ...interface{}) {
	t.t.Helper()
	t.t.Log(values...)
}

// Fatalf logs and exits immediately. Assumes it has come from a TrySuite() call. If called from within goroutine it does not immediately exit.
func (t *T) Fatalf(format string, values ...interface{}) {
	t.t.Helper()
	t.t.Log(fmt.Sprintf(format, values...))
	t.failed = true
	t.values = values
	t.format = format
	doPanic()
}

func doPanic() {
	stack := debug.Stack()
	// if we're not in TrySuite we're doing something funky in a goroutine (probably), don't panic because we won't recover
	if !strings.Contains(string(stack), "TrySuite(") {
		return
	}
	panic(errFatal)
}

func (t *T) Parallel() {
	if t.counter == 0 && isParallel {
		t.t.Parallel()
	}
	t.counter++
}

func newT(te *testing.T) *T {
	return &t{t: te}
}

// TrySuite is designed to retry a TestXX function
func TrySuite(t *testing.T, f func(t *T), times int) {
	t.Helper()
	if len(testFilter) > 0 {
		caller := strings.Split(getFrame(1).Function, ".")[2]
		runit := false
		for _, test := range testFilter {
			if test == caller {
				runit = true
				break
			}
		}
		if !runit {
			t.Skip()
		}
	}

	tee := newT(t)
	for i := 0; i < times; i++ {
		wrapF(tee, f)
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

func wrapF(t *T, f func(t *T)) {
	defer func() {
		if r := recover(); r != nil {
			if r != errFatal {
				panic(r)
			}
		}
	}()
	f(t)
}

func Login(serv Server, t *T, email, password string) error {
	return Try("Logging in with "+email, t, func() ([]byte, error) {
		readCmd := exec.Command("micro", serv.EnvFlag(), "login", "--email", email, "--password", password)
		outp, err := readCmd.CombinedOutput()
		if err != nil {
			return outp, err
		}
		if !strings.Contains(string(outp), "Success") {
			return outp, errors.New("Login output does not contain 'Success'")
		}
		return outp, err
	}, 4 * time.Second)
}
