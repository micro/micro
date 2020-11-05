package golang

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/build/util/tar"
	"github.com/micro/micro/v3/service/build/util/zip"
)

// NewBuilder returns a golang build which can build a go binary given some source
func NewBuilder() (build.Builder, error) {
	path, err := locateGo()
	if err != nil {
		return nil, fmt.Errorf("Error locating go binary: %v", err)
	}

	return &golang{
		cmdPath: path,
		tmpDir:  os.TempDir(),
	}, nil
}

type golang struct {
	cmdPath string
	tmpDir  string
}

// Build a binary using source. If an archive was used, e.g. tar, this should be specified in the
// options. If no archive option is passed, the build will treat the source as a single file.
func (g *golang) Build(src io.Reader, opts ...build.Option) (io.Reader, error) {
	// parse the options
	var options build.Options
	for _, o := range opts {
		o(&options)
	}

	// create a tmp dir to contain the source
	dir, err := ioutil.TempDir(g.tmpDir, "src")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	// decode the source and write to the tmp directory
	switch options.Archive {
	case "":
		err = writeFile(src, dir)
	case "tar":
		err = tar.Unarchive(src, dir)
	case "zip":
		err = zip.Unarchive(src, dir)
	default:
		return nil, errors.New("Invalid Archive")
	}
	if err != nil {
		return nil, err
	}

	// check for vendor directory before setting mod to vendor. the vendor directory wasn't uploaded
	// in early v3 betas so this enables backwards compatability
	args := []string{"build", "-o", "micro_build"}
	if _, err := os.Stat(filepath.Join(dir, "vendor")); err == nil {
		args = append(args, "-mod", "vendor")
	}

	// build the binary
	cmd := exec.Command(g.cmdPath, append(args, ".")...)
	cmd.Env = append(os.Environ(), "GO111MODULE=auto")
	cmd.Dir = filepath.Join(dir, options.Entrypoint)

	outp := bytes.NewBuffer(nil)
	cmd.Stderr = outp

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v: %v", err, outp.String())
	}

	// read the bytes from the file
	dst, err := ioutil.ReadFile(filepath.Join(cmd.Dir, "micro_build"))
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(dst), nil
}

// writeFile takes a single file to a directory
func writeFile(src io.Reader, dir string) error {
	// copy the source to the temp file
	bytes, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	// write the file, note: in order for the golang build to access the file, it cannot be
	// os.ModeTemp. This is okay because we delete all the files in the tmp dir at the end of this
	// function.
	return ioutil.WriteFile(filepath.Join(dir, "main.go"), bytes, os.ModePerm)
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
