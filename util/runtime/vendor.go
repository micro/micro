package runtime

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// VendorDependencies will use `go mod vendor` to generate a vendor directory containing all of a
// services deps. This is then uploaded to the server along with the source code to be built into
// a binary.
func VendorDependencies(dir string) error {
	// find the go command
	gopath, err := locateGo()
	if err != nil {
		return err
	}

	// construct the command
	cmd := exec.Command(gopath, "mod", "vendor")
	cmd.Env = append(os.Environ(), "GO111MODULE=auto")
	cmd.Dir = dir

	// execute the command and get the error output
	outp := bytes.NewBuffer(nil)
	cmd.Stderr = outp
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %v", err, outp.String())
	}

	return nil
}

// locateGo locates the go command
func locateGo() (string, error) {
	if gr := os.Getenv("GOROOT"); len(gr) > 0 {
		return filepath.Join(gr, "bin", "go"), nil
	}

	// check path
	for _, p := range filepath.SplitList(os.Getenv("PATH")) {
		bin := filepath.Join(p, "go")
		if _, err := os.Stat(bin); err == nil {
			return bin, nil
		}
	}

	return exec.LookPath("go")
}
