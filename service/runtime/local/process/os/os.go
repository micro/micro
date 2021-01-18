// +build !windows

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Original source: github.com/micro/go-micro/v3/runtime/local/process/os/os.go

// Package os runs processes locally
package os

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/micro/micro/v3/service/runtime/local/process"
)

func (p *Process) Exec(exe *process.Binary) error {
	cmd := exec.Command(exe.Package.Path)
	cmd.Dir = exe.Dir
	return cmd.Run()
}

func (p *Process) Fork(exe *process.Binary) (*process.PID, error) {

	// create command
	cmd := exec.Command(exe.Package.Path, exe.Args...)

	cmd.Dir = exe.Dir
	// set env vars
	for _, e := range os.Environ() {
		// HACK - MICRO_AUTH_* env vars will cause weird behaviour with jwt auth so make sure we don't
		// pass these through to services that we run
		if strings.HasPrefix(e, "MICRO_AUTH_") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	cmd.Env = append(cmd.Env, exe.Env...)

	// create process group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	er, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	// start the process
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &process.PID{
		ID:     fmt.Sprintf("%d", cmd.Process.Pid),
		Input:  in,
		Output: out,
		Error:  er,
	}, nil
}

func (p *Process) Kill(pid *process.PID) error {
	id, err := strconv.Atoi(pid.ID)
	if err != nil {
		return err
	}
	if _, err := os.FindProcess(id); err != nil {
		return err
	}

	// now kill it
	// using -ve PID kills the process group which we created in Fork()
	return syscall.Kill(-id, syscall.SIGTERM)
}

func (p *Process) Wait(pid *process.PID) error {
	id, err := strconv.Atoi(pid.ID)
	if err != nil {
		return err
	}

	pr, err := os.FindProcess(id)
	if err != nil {
		return err
	}

	ps, err := pr.Wait()
	if err != nil {
		return err
	}

	if ps.Success() {
		return nil
	}

	return fmt.Errorf(ps.String())
}
