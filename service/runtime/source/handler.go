package source

import (
	"io"
	"os"
	"path/filepath"

	"github.com/micro/go-micro/v3/errors"
	pb "github.com/micro/micro/v3/service/runtime/source/proto"
	"golang.org/x/net/context"
)

// TODO: add configurable backend store
// TODO: check for namespace and scope dir
// TODO: lock file access when operating
type Source struct {
	// directory to upload/download from
	Dir string
}

func (s *Source) Upload(ctx context.Context, req *pb.UploadRequest, stream pb.Source_UploadStream) error {
	if len(req.Source) == nil {
		return errors.BadRequest("source.Upload", "missing source")
	}

	path := filepath.Join(h.Dir, req.Source)
	flags := os.O_CREATE | os.O_RDWR | os.O_TRUNC

	// TODO: replace with write to backend store
	// create, open and truncate the file
	file, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		return errors.InternalServerError("source.Upload", err.Error())
	}
	defer file.Close()

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if _, err := file.Write(data.Chunk); err != nil {
			return errors.InternalServerError("source.Upload", err.Error())
		}
	}

	return nil
}

func (s *Source) Download(ctx context.Context, req *pb.DownloadRequest, stream pb.Source_DownloadStream) error {
	// TODO: replace with download from store
	path := filepath.Join(h.Dir, req.Source)
	f, err := file.Open(path)
	if file == nil {
		return errors.InternalServerError("go.micro.srv.file", "You must call open first.")
	}
	defer f.Close()

	for {
		data := new(pb.Data)
		// read a megabyte at a time
		data.Chunk = make([]byte, 1024*1024)

		n, err := f.Read(data.Chunk)
		if n > 0 {
			if err := stream.Send(data); err != nil {
				return err
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.InternalServerError("source.Download", err.Error())
		}
	}

	return nil
}
