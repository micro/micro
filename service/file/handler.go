// Package file contains handlers for file upload/download.
// The files are saved to and read from the store, but
// temporary local files are used to enable fast uploads and downloads
// (ie. the file is being used to write to and read from a specific offset).
// For this to work the client must call the same node during operations on the same file.
// (Or the server must only have once instance running in an environment).
// This package uses an other inelegant solution too for efficiency:
// files opened with `truncate: true` are considered uploads, others as downloads.
// For these reasons, make sure to use the `file.Client` `Upload` and `Download` methods
// which hide these implementation details to interact with these handlers.
package file

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	erro "errors"

	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/store"
	proto "github.com/micro/go-micro/v2/util/file/proto"
	"github.com/micro/micro/v2/internal/namespace"
	"golang.org/x/net/context"
)

const filePrefix = "files"

// NewHandler is a handler that can be registered with a micro Server
func NewHandler(readDir string, store store.Store) proto.FileHandler {
	return &handler{
		readDir: readDir,
		session: &session{
			files: make(map[int64]*fileInfo),
		},
		store: store,
	}
}

// RegisterHandler is a convenience method for registering a handler
func RegisterHandler(s server.Server, readDir string, store store.Store) {
	proto.RegisterFileHandler(s, NewHandler(readDir, store))
}

type handler struct {
	readDir string
	session *session
	store   store.Store
}

func (h *handler) Open(ctx context.Context, req *proto.OpenRequest, rsp *proto.OpenResponse) error {
	ns := namespace.FromContext(ctx)
	err := namespace.Authorize(ctx, ns)
	if err != nil {
		return err
	}

	// Create directory, ignore failure
	os.Mkdir(filepath.Join(h.readDir, ns), 0777)

	path := filepath.Join(h.readDir, ns, req.Filename)
	flags := os.O_CREATE | os.O_RDWR
	if req.GetTruncate() {
		flags = flags | os.O_TRUNC
	}
	file, err := os.OpenFile(path, flags, 0666)
	if err != nil {
		return errors.InternalServerError("go.micro.server", err.Error())
	}

	// Uploads are using truncate on open.
	// We use truncate as a (probably bad) proxy to know if a file is opened for reading or writing.
	isWrite := req.Truncate
	err = h.storeToDisk(ns, file, !isWrite)
	if err != nil {
		return err
	}
	rsp.Id = h.session.Add(&fileInfo{
		file:             file,
		openedForWriting: isWrite,
		namespace:        ns,
	})
	rsp.Result = true

	logger.Debugf("Open %s, sessionId=%d", req.Filename, rsp.Id)

	return nil
}

func (h *handler) Close(ctx context.Context, req *proto.CloseRequest, rsp *proto.CloseResponse) error {
	ns := namespace.FromContext(ctx)
	err := namespace.Authorize(ctx, ns)
	if err != nil {
		return err
	}

	fileInfo := h.session.Get(req.Id)
	if fileInfo == nil {
		return errors.BadRequest("go.micro.srv.file", "Closing unopened file")
	}
	if fileInfo.openedForWriting {
		err := h.diskToStore(fileInfo.namespace, fileInfo.file)
		if err != nil {
			return err
		}
	}
	h.session.Delete(req.Id)
	logger.Debugf("Close sessionId=%d", req.Id)
	return nil
}

func (h *handler) Stat(ctx context.Context, req *proto.StatRequest, rsp *proto.StatResponse) error {
	ns := namespace.FromContext(ctx)
	err := namespace.Authorize(ctx, ns)
	if err != nil {
		return err
	}

	path := filepath.Join(h.readDir, req.Filename)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.InternalServerError("go.micro.srv.file", err.Error())
	}

	if fi.IsDir() {
		rsp.Type = "Directory"
	} else {
		rsp.Type = "File"
		rsp.Size = fi.Size()
	}

	rsp.LastModified = fi.ModTime().Unix()
	logger.Debugf("Stat %s, %#v", req.Filename, rsp)

	return nil
}

func (h *handler) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	ns := namespace.FromContext(ctx)
	err := namespace.Authorize(ctx, ns)
	if err != nil {
		return err
	}

	fileInfo := h.session.Get(req.Id)
	if fileInfo == nil {
		return errors.InternalServerError("go.micro.srv.file", "You must call open first.")
	}

	rsp.Data = make([]byte, req.Size)
	n, err := fileInfo.file.ReadAt(rsp.Data, req.Offset)
	if err != nil && err != io.EOF {
		return errors.InternalServerError("go.micro.srv.file", err.Error())
	}

	if err == io.EOF {
		rsp.Eof = true
	}

	rsp.Size = int64(n)
	rsp.Data = rsp.Data[:n]

	logger.Debugf("Read sessionId=%d, Offset=%d, n=%d", req.Id, req.Offset, rsp.Size)

	return nil
}

func (h *handler) Write(ctx context.Context, req *proto.WriteRequest, rsp *proto.WriteResponse) error {
	ns := namespace.FromContext(ctx)
	err := namespace.Authorize(ctx, ns)
	if err != nil {
		return err
	}

	fileInfo := h.session.Get(req.Id)
	if fileInfo == nil {
		return errors.InternalServerError("go.micro.srv.file", "You must call open first.")
	}

	if _, err := fileInfo.file.WriteAt(req.GetData(), req.GetOffset()); err != nil {
		return err
	}

	logger.Debugf("Write sessionId=%d, Offset=%d, n=%d", req.Id, req.Offset)

	return nil
}

func (h *handler) diskToStore(namespace string, file *os.File) error {
	val, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%v/%v/%v", filePrefix, namespace, filepath.Base(file.Name()))
	return h.store.Write(&store.Record{
		Key:   key,
		Value: val,
	})
}

func (h *handler) storeToDisk(namespace string, file *os.File, failIfNotExist bool) error {
	key := fmt.Sprintf("%v/%v/%v", filePrefix, namespace, filepath.Base(file.Name()))
	recs, err := h.store.Read(key)
	if err != nil && err != store.ErrNotFound {
		return err
	}
	if (err != nil && err == store.ErrNotFound) || len(recs) == 0 {
		if failIfNotExist {
			return erro.New("File not found")
		}
		return nil
	}
	rec := recs[0]
	_, err = io.Copy(file, bytes.NewReader(rec.Value))
	return err
}

type fileInfo struct {
	file             *os.File
	openedForWriting bool
	namespace        string
}

type session struct {
	sync.Mutex
	files   map[int64]*fileInfo
	counter int64
}

func (s *session) Add(file *fileInfo) int64 {
	s.Lock()
	defer s.Unlock()

	s.counter++
	s.files[s.counter] = file

	return s.counter
}

func (s *session) Get(id int64) *fileInfo {
	s.Lock()
	defer s.Unlock()
	return s.files[id]
}

func (s *session) Delete(id int64) {
	s.Lock()
	defer s.Unlock()

	if fileInfo, exist := s.files[id]; exist {
		fileInfo.file.Close()
		delete(s.files, id)
	}
}

func (s *session) Len() int {
	return len(s.files)
}
