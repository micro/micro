package handler

import (
	"fmt"
	"strings"

	"github.com/micro/go-micro/v2/runtime"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
)

var (
	validSources  = []string{"github.com"}
	defaultSource = "github.com/micro/services"
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
	var srcProvided bool
	for _, src := range validSources {
		if strings.HasPrefix(s.Source, src) {
			srcProvided = true
			break
		}
	}
	if !srcProvided {
		s.Source = fmt.Sprintf("%v/%v", defaultSource, s.Source)
	}

	return &runtime.Service{
		Name:     s.Name,
		Version:  s.Version,
		Source:   s.Source,
		Metadata: s.Metadata,
	}
}

func toCreateOptions(opts *pb.CreateOptions) []runtime.CreateOption {
	options := []runtime.CreateOption{}
	// command options
	if len(opts.Command) > 0 {
		options = append(options, runtime.WithCommand(opts.Command...))
	}
	// env options
	if len(opts.Env) > 0 {
		options = append(options, runtime.WithEnv(opts.Env))
	}

	// TODO: output options

	return options
}

func toReadOptions(opts *pb.ReadOptions) []runtime.ReadOption {
	options := []runtime.ReadOption{}
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
