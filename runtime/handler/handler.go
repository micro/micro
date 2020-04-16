package handler

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
)

type Runtime struct {
	// The runtime used to manage services
	Runtime runtime.Runtime
	// The client used to publish events
	Client micro.Publisher
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

func (r *Runtime) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	var options []runtime.CreateOption
	if req.Options != nil {
		options = toCreateOptions(req.Options)
	}

	service := toService(req.Service)

	sourceInfo, err := extractSource(service.Source)
	if err != nil {
		return err
	}
	service.Name = sourceInfo.serviceName
	service.Version = sourceInfo.serviceVersion

	// non local source
	if sourceInfo.githubURL != nil {
		service.Source = sourceInfo.githubURL.folder
	}
	// This is needed to support local `micro server` execution of git urls
	if r.Runtime.String() == "local" && sourceInfo.githubURL != nil {
		service.Source = filepath.Join(sourceInfo.repoRoot, sourceInfo.githubURL.folder)
	}

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

func (r *Runtime) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	if req.Service == nil {
		return errors.BadRequest("go.micro.runtime", "blank service")
	}

	// TODO: add opts
	service := toService(req.Service)
	sourceInfo, err := extractSource(service.Source)
	if err != nil {
		return err
	}
	service.Name = sourceInfo.serviceName
	service.Version = sourceInfo.serviceVersion
	// non local source
	if sourceInfo.githubURL != nil {
		service.Source = sourceInfo.githubURL.folder
	}
	// This is needed to support local `micro server` execution of git urls
	if r.Runtime.String() == "local" && sourceInfo.githubURL != nil {
		service.Source = filepath.Join(sourceInfo.repoRoot, sourceInfo.githubURL.folder)
	}

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
	sourceInfo, err := extractSource(service.Source)
	if err != nil {
		return err
	}
	service.Name = sourceInfo.serviceName
	service.Version = sourceInfo.serviceVersion
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

// exists returns whether the given file or directory exists
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// parsedGithubURL represent
// the strutured information we care about when
// extracting the full provided github URL source
type parsedGithubURL struct {
	// for cloning purposes
	repoAddress string
	// path of folder to repo root
	folder string
	// github ref
	ref string
}

// sourceInfo contains all information
// that was extracted from a source.
// for source examples see `micro run --help`
type sourceInfo struct {
	githubURL *parsedGithubURL
	// repo root in local filesystem
	repoRoot string
	// local urls ie. micro run helloworld/web
	// will not have github URLs but we need to pass
	// the relative folder to the micro/services image still.
	relativePath string
	// name of service
	serviceName string
	// service version
	serviceVersion string
}

// extractSource tries to gather as much information
// as possible just from the source.
// Also clones and checks out remote services from github
// for the local runtime implementation which has no code available.
// for source examples see `micro run --help`
func extractSource(source string) (*sourceInfo, error) {
	sinf := &sourceInfo{}
	var mainFilePath string

	if local, err := pathExists(source); err == nil && local {
		// Local directories to be deployed are not expected
		// to be in source control. @todo we could still try
		// to detect source control if exists and take the commit hash
		// from there.
		repoRoot, err := getRepoRoot(source)
		if err != nil {
			return nil, err
		}
		// is source controlled
		if repoRoot != "" {
			sinf.repoRoot = repoRoot
			sinf.relativePath = strings.ReplaceAll(source, repoRoot+string(filepath.Separator), "")
			repo, err := git.PlainOpen(repoRoot)
			if err != nil {
				return nil, err
			}
			var head *plumbing.Reference
			head, err = repo.Head()
			if err != nil {
				return nil, err
			}
			sinf.serviceVersion = head.String()
		}
		// @todo think about non source controlled
		// deploys. They will miss the needed relative path
		sinf.serviceVersion = "latest"
		mainFilePath = filepath.Join(source, "main.go")
	} else {
		parsed, err := parseGithubURL(source)
		if err != nil {
			return nil, err
		}
		sinf.githubURL = parsed

		dirifiedURL := dirifyRepo(parsed.repoAddress)
		repoDir := filepath.Join(os.TempDir(), dirifiedURL)

		// Only clone if doesn't exist already.
		// @todo implement pull and check out of correct version
		// by parsing commit hash from the git URL.
		if exists, err := pathExists(repoDir); err == nil && !exists {
			_, err = git.PlainClone(repoDir, false, &git.CloneOptions{
				URL:      parsed.repoAddress,
				Progress: os.Stdout,
			})
			if err != nil {
				return nil, err
			}
		}
		repo, err := git.PlainOpen(repoDir)
		if err != nil {
			return nil, err
		}
		remotes, err := repo.Remotes()
		if err != nil {
			return nil, err
		}
		err = remotes[0].Fetch(&git.FetchOptions{})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, err
		}
		worktree, err := repo.Worktree()
		if err != nil {
			return nil, err
		}
		reference := plumbing.NewReferenceFromStrings("refname", parsed.ref)
		if reference.Type() == plumbing.HashReference {
			err = worktree.Checkout(&git.CheckoutOptions{
				Hash:  plumbing.NewHash(parsed.ref),
				Force: true,
			})
			if err != nil {
				return nil, err
			}
		} else {
			branch := parsed.ref
			if parsed.ref == "latest" {
				branch = "master"
			}
			worktree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(branch),
				Force:  true,
			})
		}
		sinf.repoRoot = repoDir
		sinf.serviceVersion = parsed.ref
		mainFilePath = filepath.Join(repoDir, parsed.folder, "main.go")
	}

	fileContent, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return nil, err
	}
	sinf.serviceName = extractServiceName(fileContent)
	return sinf, nil
}

// get repo root from full path
// returns empty string and no error if not found
func getRepoRoot(fullPath string) (string, error) {
	// traverse parent directories
	prev := fullPath
	for {
		current := prev
		log.Infof("++", current)
		exists, err := pathExists(filepath.Join(current, ".git"))
		if err != nil {
			return "", err
		}
		if exists {
			return current, nil
		}
		prev = filepath.Dir(current)
		// reached top level, see:
		// https://play.golang.org/p/rDgVdk3suzb
		if current == prev {
			break
		}
	}
	return "", nil
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

// @todo rename, source is not an actual URL
// but more like `go get`.
func parseGithubURL(url string) (*parsedGithubURL, error) {
	// If github is not present, we got a shorthand for `micro/services`
	if !strings.Contains(url, "github.com") {
		url = "github.com/micro/services/" + url
	}
	if !strings.Contains(url, "@") {
		url += "@latest"
	}
	ret := &parsedGithubURL{}
	refs := strings.Split(url, "@")
	ret.ref = refs[1]
	parts := strings.Split(refs[0], "/")
	ret.repoAddress = "https://" + strings.Join(parts[0:3], "/")
	if len(parts) > 1 {
		ret.folder = strings.Join(parts[3:], "/")
	}

	return ret, nil
}

func dirifyRepo(s string) string {
	s = strings.ReplaceAll(s, "https://", "")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}
