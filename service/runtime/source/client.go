package source

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/micro/v3/internal/git"
	"github.com/micro/go-micro/v3/client"
	pb "github.com/micro/micro/v3/service/runtime/source/proto"
)

// Client is the client for managing source code
type Client interface {
	// Package will package a directory as a tarball
	Package(source, pkg string) error
	// Extract a package into its source form
	Extract(pkg, source string) error
	// Download source to a local destination
	Download(source, dest string) error
	// Upload local source code to the runtime
	Upload(source, dest string) error
	// String indicates implementation
	String() string
}

// NewClient returns a new Client which uses a micro Client
func New(c client.Client) Client {
	return &src{pb.NewSourceService("runtime", c)}
}

const (
	blockSize = 512 * 1024
)

type src struct {
	c pb.SourceService
}

func grepMain(path string) error {
        files, err := ioutil.ReadDir(path)
        if err != nil {
                return err
        }

        for _, f := range files {
                if !strings.HasSuffix(f.Name(), ".go") {
                        continue
                }
                file := filepath.Join(path, f.Name())
                b, err := ioutil.ReadFile(file)
                if err != nil {
                        continue
                }
                if strings.Contains(string(b), "package main") {
                        return nil
                }
        }

        return fmt.Errorf("Directory does not contain a main package")
}

// Package the source
func (s *src) Package(source, dest string) error {
	// absolute path
	abs, err := filepath.Abs(source)
	if err != nil {
		return err
	}
	// look for main
        if err := grepMain(abs); err != nil {
                return err
        }

	// compress as a tarball
        return git.Compress(abs, dest)
}

// Download source from the destination
func (s *src) Download(source, dest string) error {
	if err := os.MkdirAll(path.Dir(dest)); err != nil {
		return err
	}

        flags := os.O_CREATE | os.O_RDWR | os.O_TRUNC

        // TODO: replace with write to backend store
        // create, open and truncate the file
        file, err := os.OpenFile(dest, flags, 0666)
        if err != nil {
                return err
        }
	defer file.Close()

	// start the download request
	stream, err := s.Download(context.Background(), &pb.DownloadRequest{
		Source: source,
		Format: "tar.gz",
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	// stream all the chunks of the file
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := file.Write(data); err != nil {
			return err
		}
	}

	return nil
}

// Upload source to the destination
func (s *src) Upload(source, dest string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Start the upload request
	stream, err := s.Upload(context.Background(), &pb.UploadRequest{
		Source: dest,
		Format: "tar.gz",
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	// stream all the chunks of the file
        for {
                data := new(pb.Data)
                // read a megabyte at a time
                data.Chunk = make([]byte, 1024*1024)

                n, err := file.Read(data.Chunk)
                if n > 0 {
                        if err := stream.Send(data); err != nil {
                                return err
                        }
                }
                if err == io.EOF {
                        return nil
                }
                if err != nil {
                        return err
                }
        }

	return nil
}

// Extract the package
func (s *src) Extract(pkg, source string) error {
	err := os.RemoveAll(source)
	if err != nil {
		return err
	}
	err = os.MkdirAll(source, 0777)
	if err != nil {
		return err
	}
	// uncompress the tarball into a directory
	return git.Uncompress(pkg, source)
}

func (s *src) String() string {
	return "tarball"
}
