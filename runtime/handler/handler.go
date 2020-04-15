package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
	uuid "github.com/satori/go.uuid"
)

type Runtime struct {
	// The runtime used to manage services
	Runtime runtime.Runtime
	// The client used to publish events
	Client micro.Publisher
}

func (r *Runtime) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	var options []runtime.CreateOption
	if req.Options != nil {
		options = toCreateOptions(req.Options)
	}

	service := toService(req.Service)
	// local directory path of either local or checked out
	// service source code
	var dirPath string
	uid := uuid.NewV4().String()
	if local, err := isLocal(service.Source); err != nil && local {
		dirPath = service.Source
	} else {
		repoDir := filepath.Join(os.TempDir(), uid)
		parsed, err := parseGithubURL(service.Source)
		if err != nil {
			return err
		}
		_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
			URL:      parsed.repoAddress,
			Progress: os.Stdout,
		})
		if err != nil {
			return err
		}
		dirPath = filepath.Join(repoDir, parsed.folder)
	}
	name, version, err := extractNameAndVersion(dirPath)
	if err != nil {
		return err
	}
	fmt.Println("name and version", name, version)

	log.Infof("Creating service %s version %s source %s", service.Name, service.Version, service.Source)

	if err := r.Runtime.Create(service, options...); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the create event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "create",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

// exists returns whether the given file or directory exists
func isLocal(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type parsedGithubURL struct {
	// for cloning purposes
	repoAddress string
	// path of folder to repo root
	folder string
}

func extractNameAndVersion(dirPath string) (name, version string, err error) {
	repo, err := git.PlainOpen(dirPath)
	if err != nil {
		return
	}
	head, err := repo.Head()
	if err != nil {
		return
	}
	version = head.Hash().String()
	_, err = ioutil.ReadFile("/tmp/dat")
	if err != nil {
		return
	}
	return
}

var nameExtractRegexp = regexp.MustCompile(`(micro\.Name\(")(.*)("\))`)

func extractServiceName(fileContent []byte) string {
	hits := nameExtractRegexp.FindAll(fileContent, 1)
	if len(hits) == 0 {
		return ""
	}
	hit := string(hits[0])
	return strings.Split(hit, "\"")[1]
}

func parseGithubURL(url string) (*parsedGithubURL, error) {
	// If github is not present, we got a shorthand for `micro/services`
	if !strings.Contains(url, "github.com") {
		url = "https://github.com/micro/services/tree/master/" + url
	}
	parts := strings.Split(url, "tree/master")
	ret := &parsedGithubURL{
		repoAddress: parts[0][0 : len(parts[0])-1],
	}
	if len(parts) > 1 {
		ret.folder = parts[1][1:]
	}
	return ret, nil
}

func (r *Runtime) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	var options []runtime.ReadOption

	if req.Options != nil {
		options = toReadOptions(req.Options)
	}

	services, err := r.Runtime.Read(options...)
	if err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// TODO: add opts
	service := toService(req.Service)

	log.Infof("Updating service %s version %s source %s", service.Name, service.Version, service.Source)

	if err := r.Runtime.Update(service); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the update event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "update",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// TODO: add opts
	service := toService(req.Service)

	log.Infof("Deleting service %s version %s source %s", service.Name, service.Version, service.Source)

	if err := r.Runtime.Delete(service); err != nil {
		return errors.InternalServerError("go.micro.runtime", err.Error())
	}

	// publish the delete event
	r.Client.Publish(ctx, &pb.Event{
		Type:      "delete",
		Timestamp: time.Now().Unix(),
		Service:   req.Service.Name,
		Version:   req.Service.Version,
	})

	return nil
}

func (r *Runtime) Logs(ctx context.Context, req *pb.LogsRequest, stream pb.Runtime_LogsStream) error {
	opts := []runtime.LogsOption{}
	if req.GetCount() > 0 {
		opts = append(opts, runtime.LogsCount(req.GetCount()))
	}
	if req.GetStream() {
		opts = append(opts, runtime.LogsStream(req.GetStream()))
	}
	logStream, err := r.Runtime.Logs(&runtime.Service{
		Name: req.GetService(),
	}, opts...)
	if err != nil {
		return err
	}
	defer logStream.Stop()
	defer stream.Close()

	recordChan := logStream.Chan()
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				return logStream.Error()
			}
			// send record
			if err := stream.Send(&pb.LogRecord{
				//Timestamp: record.Timestamp.Unix(),
				Message: record.Message,
			}); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}
