package runtime

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/build/util/tar"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/runtime"
	"github.com/micro/micro/v3/service/runtime/source/git"
	"github.com/urfave/cli/v2"
)

const bufferSize = 1024

// upload source to the server. will return the source id, e.g. source://foo-bar and an error if
// one occured. The ID returned can be used as a source in runtime.Create.
func upload(ctx *cli.Context, srv *runtime.Service, source *git.Source) (string, error) {
	// if the source exists within a local git repository, archive the whole repository, otherwise
	// just archive the folder
	var data io.Reader
	var err error
	if len(source.LocalRepoRoot) > 0 {
		data, err = tar.Archive(source.LocalRepoRoot)
	} else {
		data, err = tar.Archive(source.FullPath)
	}
	if err != nil {
		return "", err
	}

	// create an upload stream
	cli := pb.NewSourceService("runtime", client.DefaultClient)
	stream, err := cli.Upload(context.DefaultContext, client.WithAuthToken())
	if err != nil {
		return "", err
	}

	// read bytes from the tar and stream it to the server
	buffer := make([]byte, bufferSize)
	var sentService bool
	for {
		num, err := data.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		req := &pb.UploadRequest{Data: buffer[:num]}

		// construct the service object, we'll send this on the first message only to reduce the amount of
		// data needed to be streamed
		if !sentService {
			req.Service = &pb.Service{Name: srv.Name, Version: srv.Version}
			sentService = true
		}

		if err := stream.Send(req); err != nil {
			return "", err
		}
	}

	// wait for the server to process the source
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}
	return rsp.Id, nil
}

// vendorDependencies will use `go mod vendor` to generate a vendor directory containing all of a
// services deps. This is then uploaded to the server along with the source code to be built into
// a binary.
func vendorDependencies(dir string) error {
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
