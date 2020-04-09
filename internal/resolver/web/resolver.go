package web

import (
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"

	res "github.com/micro/go-micro/v2/api/resolver"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client/selector"
	"golang.org/x/net/publicsuffix"
)

var (
	re               = regexp.MustCompile("^[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9]*)?$")
	defaultNamespace = auth.DefaultNamespace + ".web"
)

type Resolver struct {
	// Type of resolver e.g path, domain
	Type string
	// a function which returns the namespace of the request
	Namespace func(*http.Request) string
	// selector to find services
	Selector selector.Selector
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func (r *Resolver) String() string {
	return "web/resolver"
}

// Info checks whether this is a web request.
// It returns host, namespace and whether its internal
func (r *Resolver) Info(req *http.Request) (string, string, bool) {
	// set to host
	host := req.URL.Hostname()

	// set as req.Host if blank
	if len(host) == 0 {
		host = req.Host
	}

	// split out ip
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// determine the namespace of the request
	namespace := r.Namespace(req)

	// isWeb sets if its a web.micro.mu request
	var isWeb bool

	// go.micro.web => go.micro.web use path based resolution if
	// explicitly set or the namespace matches the default namespace
	// (indicating we're on a micro.mu host or in dev)
	if r.Type == "path" || namespace != defaultNamespace {
		isWeb = true
	}

	return host, namespace, isWeb
}

// Resolve replaces the values of Host, Path, Scheme to calla backend service
// It accounts for subdomains for service names based on namespace
func (r *Resolver) Resolve(req *http.Request) (*res.Endpoint, error) {
	// get host, namespace and if its an internal request
	host, namespace, webReq := r.Info(req)

	// use path based resolution if its web dashboard related
	if webReq {
		parts := strings.Split(req.URL.Path, "/")
		if len(parts) < 2 {
			return nil, errors.New("unknown service")
		}

		if !re.MatchString(parts[1]) {
			return nil, res.ErrInvalidPath
		}

		next, err := r.Selector.Select(namespace + "." + parts[1])
		if err == selector.ErrNotFound {
			return nil, res.ErrNotFound
		} else if err != nil {
			return nil, err
		}

		// TODO: better retry strategy
		s, err := next()
		if err != nil {
			return nil, err
		}

		// we're done
		return &res.Endpoint{
			Name:   parts[1],
			Method: req.Method,
			Host:   s.Address,
			Path:   "/" + strings.Join(parts[2:], "/"),
		}, nil
	}

	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return nil, err
	}

	// get and reverse the subdomain
	subdomain := strings.TrimSuffix(host, "."+domain)
	parts := strings.Split(subdomain, ".")
	reverse(parts)

	// turn it into an alias
	alias := strings.Join(parts, ".")
	if len(alias) == 0 {
		return nil, errors.New("unknown host")
	}

	// set name to lookup
	name := defaultNamespace + "." + alias

	// check for go.micro.web (render dashboard)
	if namespace == defaultNamespace && alias == "web" {
		name = defaultNamespace
	}

	// get namespace + subdomain
	next, err := r.Selector.Select(name)
	if err == selector.ErrNotFound {
		return nil, res.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// TODO: better retry strategy
	s, err := next()
	if err != nil {
		return nil, err
	}

	// we're done
	return &res.Endpoint{
		Name:   alias,
		Method: req.Method,
		Host:   s.Address,
		Path:   req.URL.Path,
	}, nil
}
