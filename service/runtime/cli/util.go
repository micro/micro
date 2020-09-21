package runtime

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/runtime/local/source/git"
	"github.com/micro/micro/v3/client/cli/namespace"
	"github.com/micro/micro/v3/client/cli/util"
	pb "github.com/micro/micro/v3/proto/runtime"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/urfave/cli/v2"
)

const bufferSize = 100

// upload source to the server. will return the source id, e.g. source://foo-bar and an error if
// one occured. The ID returned can be used as a source in runtime.Create.
func upload(ctx *cli.Context, source *git.Source) (string, error) {
	// if the source exists within a local git repository, archive the whole repository, otherwise
	// just archive the folder
	var tar io.Reader
	var err error
	if len(source.LocalRepoRoot) > 0 {
		tar, err = archive(source.LocalRepoRoot)
	} else {
		tar, err = archive(source.FullPath)
	}
	if err != nil {
		return "", err
	}

	// get the namespace of the client
	ns, err := namespace.Get(util.GetEnv(ctx).Name)
	if err != nil {
		return "", err
	}

	// create an upload stream
	cli := pb.NewSourceService("runtime", client.DefaultClient)
	stream, err := cli.Upload(context.WithNamespace(ns), goclient.WithAuthToken())
	if err != nil {
		return "", err
	}

	// read bytes from the tar and stream it to the server
	buffer := make([]byte, bufferSize)
	for {
		num, err := tar.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		if err := stream.Send(&pb.UploadRequest{Data: buffer[:num]}); err != nil {
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

// archive a local directory into a tar gzip, ready for streaming to a server.
func archive(dir string) (io.Reader, error) {
	// Create a tar writer and a buffer to store the archive
	tf := bytes.NewBuffer(nil)
	tw := tar.NewWriter(tf)
	defer tw.Close()

	// walkFn archives each file in the directory
	walkFn := func(path string, info os.FileInfo, err error) error {
		// get the relative path, e.g. cmd/main.go
		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// generate and write tar header
		header, err := tar.FileInfoHeader(info, relpath)
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// there is no body if it's a directory
		if info.IsDir() {
			return nil
		}

		// read the contents of the file
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// write the contents of the file to the tar
		_, err = tw.Write([]byte(bytes))
		return err
	}

	// Add the files to the archive
	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, err
	}

	// Return the archive
	return tf, nil
}
