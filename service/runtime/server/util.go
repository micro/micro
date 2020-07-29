package server

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/runtime"
	pb "github.com/micro/micro/v3/service/runtime/proto"
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

func toCreateOptions(ctx context.Context, opts *pb.CreateOptions) []runtime.CreateOption {
	options := []runtime.CreateOption{
		runtime.CreateNamespace(opts.Namespace),
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
	options := []runtime.ReadOption{
		runtime.ReadNamespace(opts.Namespace),
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
	return []runtime.UpdateOption{
		runtime.UpdateNamespace(opts.Namespace),
	}
}

func toDeleteOptions(ctx context.Context, opts *pb.DeleteOptions) []runtime.DeleteOption {
	return []runtime.DeleteOption{
		runtime.DeleteNamespace(opts.Namespace),
	}
}

func toLogsOptions(ctx context.Context, opts *pb.LogsOptions) []runtime.LogsOption {
	return []runtime.LogsOption{
		runtime.LogsNamespace(opts.Namespace),
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

		srcWithSlash := src
		if !strings.HasSuffix(src, string(filepath.Separator)) {
			srcWithSlash = src + string(filepath.Separator)
		}
		header.Name = strings.ReplaceAll(file, srcWithSlash, "")
		if header.Name == src || len(strings.TrimSpace(header.Name)) == 0 {
			return nil
		}

		// @todo This is a quick hack to speed up whole repo uploads
		// https://github.com/micro/micro/pull/956
		if !fi.IsDir() && !strings.HasSuffix(header.Name, ".go") &&
			!strings.HasSuffix(header.Name, ".mod") &&
			!strings.HasSuffix(header.Name, ".sum") {
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
