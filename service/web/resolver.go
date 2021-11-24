package web

import (
	"errors"
	"math/rand"
	"net/http"
	"regexp"
	"strings"

	"github.com/micro/micro/v3/service/api/resolver"
	res "github.com/micro/micro/v3/service/api/resolver"
	"github.com/micro/micro/v3/service/router"
)

var re = regexp.MustCompile("^[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9]*)?$")

type WebResolver struct {
	// Options
	Options resolver.Options
	// selector to choose from a pool of nodes
	// Selector selector.Selector
	// router to lookup routes
	Router router.Router
}

func (r *WebResolver) String() string {
	return "web/resolver"
}

// Resolve replaces the values of Host, Path, Scheme to calla backend service
// It accounts for subdomains for service names based on namespace
func (r *WebResolver) Resolve(req *http.Request, opts ...res.ResolveOption) (*res.Endpoint, error) {
	// parse the options
	options := resolver.NewResolveOptions(opts...)

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 2 {
		return nil, errors.New("unknown service")
	}

	if !re.MatchString(parts[1]) {
		return nil, res.ErrInvalidPath
	}

	// lookup the routes for the service
	query := []router.LookupOption{
		router.LookupNetwork(options.Domain),
	}

	routes, err := r.Router.Lookup(parts[1], query...)
	if err == router.ErrRouteNotFound {
		return nil, res.ErrNotFound
	} else if err != nil {
		return nil, err
	} else if len(routes) == 0 {
		return nil, res.ErrNotFound
	}

	// select a random route to use
	// todo: update to use selector once go-micro has updated the interface
	// route, err := r.Selector.Select(routes...)
	// if err != nil {
	// 	return nil, err
	// }
	route := routes[rand.Intn(len(routes))]

	// we're done
	return &res.Endpoint{
		Name:   parts[1],
		Method: req.Method,
		Host:   route.Address,
		Path:   "/" + strings.Join(parts[2:], "/"),
		Domain: options.Domain,
	}, nil
}
