package handler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
	"github.com/micro/micro/v2/internal/namespace"
)

func toProto(s *runtime.Service) *pb.Service {
	return &pb.Service{
		Name:     s.Name,
		Version:  s.Version,
		Source:   s.Source,
		Metadata: s.Metadata,
	}
}

func toService(s *pb.Service) *runtime.Service {
	return &runtime.Service{
		Name:     s.Name,
		Version:  s.Version,
		Source:   s.Source,
		Metadata: s.Metadata,
	}
}

// getNamespace replaces the default auth namespace until we move
// we wil replace go.micro with micro and move our default things there
func getNamespace(ctx context.Context) string {
	return namespace.FromContext(ctx)
}

func toCreateOptions(ctx context.Context, opts *pb.CreateOptions) []runtime.CreateOption {
	if opts == nil {
		opts = &pb.CreateOptions{}
	}

	options := []runtime.CreateOption{
		runtime.CreateNamespace(getNamespace(ctx)),
	}

	// command options
	if len(opts.Command) > 0 {
		options = append(options, runtime.WithCommand(opts.Command...))
	}

	// args for command
	if len(opts.Args) > 0 {
		options = append(options, runtime.WithArgs(opts.Args...))
	}

	// env options
	if len(opts.Env) > 0 {
		options = append(options, runtime.WithEnv(opts.Env))
	}

	// create specific type of service
	if len(opts.Type) > 0 {
		options = append(options, runtime.CreateType(opts.Type))
	}

	// use specific image
	if len(opts.Image) > 0 {
		options = append(options, runtime.CreateImage(opts.Image))
	}

	// TODO: output options

	return options
}

func toReadOptions(ctx context.Context, opts *pb.ReadOptions) []runtime.ReadOption {
	if opts == nil {
		opts = &pb.ReadOptions{}
	}

	options := []runtime.ReadOption{
		runtime.ReadNamespace(getNamespace(ctx)),
	}

	if len(opts.Service) > 0 {
		options = append(options, runtime.ReadService(opts.Service))
	}
	if len(opts.Version) > 0 {
		options = append(options, runtime.ReadVersion(opts.Version))
	}
	if len(opts.Type) > 0 {
		options = append(options, runtime.ReadType(opts.Type))
	}

	return options
}

func toUpdateOptions(ctx context.Context, opts *pb.UpdateOptions) []runtime.UpdateOption {
	if opts == nil {
		opts = &pb.UpdateOptions{}
	}

	return []runtime.UpdateOption{
		runtime.UpdateNamespace(getNamespace(ctx)),
	}
}

func toDeleteOptions(ctx context.Context, opts *pb.DeleteOptions) []runtime.DeleteOption {
	if opts == nil {
		opts = &pb.DeleteOptions{}
	}

	return []runtime.DeleteOption{
		runtime.DeleteNamespace(getNamespace(ctx)),
	}
}

func toLogsOptions(ctx context.Context, opts *pb.LogsOptions) []runtime.LogsOption {
	if opts == nil {
		opts = &pb.LogsOptions{}
	}

	return []runtime.LogsOption{
		runtime.LogsNamespace(getNamespace(ctx)),
	}
}

// taken from https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
// likely should be in go micro

func Compress(sourceFolderPath, destinationFilePath string) error {
	// tar + gzip
	var buf bytes.Buffer
	_ = compress(sourceFolderPath, &buf)

	// write the .tar.gzip
	fileToWrite, err := os.OpenFile(destinationFilePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(fileToWrite, &buf)
	return err
}

func compress(src string, buf io.Writer) error {
	// tar > gzip > buf
	zr := gzip.NewWriter(buf)
	tw := tar.NewWriter(zr)

	// walk through every file in the folder
	filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// generate tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)

		header.Name = filepath.ToSlash(strings.ReplaceAll(file, src+string(filepath.Separator), ""))
		if header.Name == src {
			return nil
		}

		// write header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		// if not a dir, write file content

		data, err := os.Open(file)
		if err != nil {
			return err
		}
		if _, err := io.Copy(tw, data); err != nil {
			return err
		}

		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return err
	}
	//
	return nil
}
