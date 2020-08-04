package test

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
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

// Server is a micro server
type Server interface {
	// Run the server
	Run() error
	// Close shuts down the server
	Close()
	// Command provides a `micro` command for the server
	Command() *Command
	// Name of the environment
	Env() string
}

type Command struct {
	Env    string
	Config string

	sync.Mutex
	// in the event an async command is run
	cmd       *exec.Cmd
	cmdOutput bytes.Buffer

	// internal logging use
	t *T
}

func (c *Command) args(a ...string) []string {
	arguments := []string{}
	// add config flag
	arguments = append(arguments, []string{"-c", c.Config}...)
	// add env flag if not env command
	if v := len(a); v > 0 && a[0] != "env" {
		arguments = append(arguments, []string{"-e", c.Env}...)
	}
	return append(arguments, a...)
}

// Exec executes a command inline
func (c *Command) Exec(args ...string) ([]byte, error) {
	arguments := c.args(args...)
	// exec the command
	c.t.Logf("Executing command: micro %s\n", strings.Join(arguments, " "))
	return exec.Command("micro", arguments...).CombinedOutput()
}

// Starts a new command
func (c *Command) Start(args ...string) error {
	c.Lock()
	defer c.Unlock()

	if c.cmd != nil {
		return errors.New("command is already running")
	}

	arguments := c.args(args...)
	c.cmd = exec.Command("micro", arguments...)

	c.cmd.Stdout = &c.cmdOutput
	c.cmd.Stderr = &c.cmdOutput

	return c.cmd.Start()
}

// Stop a command thats running
func (c *Command) Stop() error {
	c.Lock()
	defer c.Unlock()

	if c.cmd != nil {
		err := c.cmd.Process.Kill()
		c.cmd = nil
		return err
	}

	return nil
}

// Output of a running command
func (c *Command) Output() ([]byte, error) {
	c.Lock()
	defer c.Unlock()
	if c.cmd == nil {
		return nil, errors.New("command is not running")
	}
	return c.cmdOutput.Bytes(), nil
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
	opts Options
	cmd  *exec.Cmd
	t    *T
	// path to config file
	config string
	// directory for test config
	dir string
	// name of the environment
	env string
	// port number
	port int
	// name of the container
	container string
	// namespace of server
	namespace string
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
	// Login specifies whether to login to the server
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

	userDir, _ := user.Current()
	dir := filepath.Join(userDir.HomeDir, ".micro/test")
	config := filepath.Join(dir, "config-"+fname+".json")

	return &ServerDefault{ServerBase{
		dir:       dir,
		config:    config,
		cmd:       cmd,
		t:         t,
		env:       fname,
		container: fname,
		port:      portnum,
		opts:      options,
		namespace: "micro",
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

	cmd := s.Command()

	// add the environment
	if err := Try("Adding micro env: "+s.env+" file: "+s.config, s.t, func() ([]byte, error) {
		out, err := cmd.Exec("env", "add", s.env, fmt.Sprintf("127.0.0.1:%d", s.port))
		if err != nil {
			return nil, err
		}

		if len(out) > 0 {
			return nil, errors.New("Not added")
		}

		out, err = cmd.Exec("env")
		if err != nil {
			return nil, err
		}

		if !strings.Contains(string(out), s.env) {
			return nil, errors.New("Not added")
		}

		return out, nil
	}, 15*time.Second); err != nil {
		return err
	}

	return nil
}

func (s *ServerDefault) Run() error {
	if err := s.ServerBase.Run(); err != nil {
		return err
	}

	if err := Try("Calling micro server", s.t, func() ([]byte, error) {
		out, err := s.Command().Exec("services")
		if !strings.Contains(string(out), "runtime") ||
			!strings.Contains(string(out), "registry") ||
			!strings.Contains(string(out), "broker") ||
			!strings.Contains(string(out), "config") ||
			!strings.Contains(string(out), "debug") ||
			!strings.Contains(string(out), "proxy") ||
			!strings.Contains(string(out), "auth") ||
			!strings.Contains(string(out), "store") {
			return out, errors.New("Not ready")
		}

		return out, err
	}, 60*time.Second); err != nil {
		return err
	}

	// login to admin account
	if s.opts.Login {
		Login(s, s.t, "default", "password")
	}

	// // generate a new admin account for the env : user=ENV_NAME pass=password
	// req := fmt.Sprintf(`{"id":"%s", "secret":"password", "options":{"namespace":"%s"}}`, s.env, s.namespace)
	// outp, err := exec.Command("micro", s.EnvFlag(), "call", "go.micro.auth", "Auth.Generate", req).CombinedOutput()
	// if err != nil && !strings.Contains(string(outp), "already exists") { // until auth.Delete is implemented
	// 	s.t.Fatalf("Error generating auth: %s, %s", err, outp)
	// 	return err
	// }

	return nil
}

func (s *ServerBase) Close() {
	// delete the config for this test
	os.Remove(s.config)

	// remove the credentials so they aren't reused on next run
	token.Remove(s.env)

	// reset back to the default namespace
	namespace.Set("micro", s.env)

}

func (s *ServerDefault) Close() {
	s.ServerBase.Close()
	exec.Command("docker", "kill", s.container).CombinedOutput()
	if s.cmd.Process != nil {
		s.cmd.Process.Signal(syscall.SIGKILL)
	}
}

func (s *ServerBase) Command() *Command {
	return &Command{
		Env:    s.env,
		Config: s.config,
		t:      s.t,
	}
}

func (s *ServerBase) Env() string {
	return s.env
}

type T struct {
	counter int
	failed  bool
	format  string
	values  []interface{}
	t       *testing.T
}

// Failed indicate whether the test failed
func (t *T) Failed() bool {
	return t.failed
}

// Expose testing.T
func (t *T) T() *testing.T {
	return t.t
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

func (t *T) Logf(format string, values ...interface{}) {
	t.t.Helper()
	t.t.Logf(format, values...)
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

// New returns a new test framework
func New(t *testing.T) *T {
	return &T{t: t}
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

	tee := New(t)
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
		out, err := serv.Command().Exec("login", "--email", email, "--password", password)
		if err != nil {
			return out, err
		}
		if !strings.Contains(string(out), "Success") {
			return out, errors.New("Login output does not contain 'Success'")
		}
		return out, err
	}, 4*time.Second)
}

func ChangeNamespace(cmd *Command, env, namespace string) error {
	if _, err := cmd.Exec("user", "config", "set", "namespaces."+env+".all", namespace); err != nil {
		return err
	}
	if _, err := cmd.Exec("user", "config", "set", "namespaces."+env+".current", namespace); err != nil {
		return err
	}
	if _, err := cmd.Exec("user", "config", "set", "micro.auth."+env+".token", ""); err != nil {
		return err
	}
	if _, err := cmd.Exec("user", "config", "set", "micro.auth."+env+".refresh-token", ""); err != nil {
		return err
	}
	if _, err := cmd.Exec("user", "config", "set", "micro.auth."+env+".expiry", ""); err != nil {
		return err
	}
	return nil
}
