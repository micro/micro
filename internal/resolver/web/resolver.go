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

	// overide host if the namespace is go.micro.web, since
	// this will also catch localhost & 127.0.0.1, resulting
	// in a more consistent dev experience
	if host == "localhost" || host == "127.0.0.1" {
		host = "web.micro.mu"
	}

	// if the type is path, always resolve using the path
	if r.Type == "path" {
		return host, namespace, true
	}

	// if the namespace is not the default (go.micro.web),
	// we always resolve using path
	if namespace != defaultNamespace {
		return host, namespace, true
	}

	// check to see if this request is for a micro.mu subdomain
	if host != "web.micro.mu" {
		return host, namespace, false
	}

	// Check if the request is a top level path
	isWeb := strings.Count(req.URL.Path, "/") == 1
	return host, namespace, isWeb
}

// Resolve replaces the values of Host, Path, Scheme to calla backend service
// It accounts for subdomains for service names based on namespace
func (r *Resolver) Resolve(req *http.Request) (*res.Endpoint, error) {
	// get host, namespace and if its an internal request
	host, _, webReq := r.Info(req)

	// use path based resolution if its web dashboard related.
	if webReq {
		return r.resolveWithPath(req)
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

	// get namespace + subdomain
	next, err := r.Selector.Select(name)
	if err == selector.ErrNotFound {
		// fallback to path based
		return r.resolveWithPath(req)
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

func (r *Resolver) resolveWithPath(req *http.Request) (*res.Endpoint, error) {
	parts := strings.Split(req.URL.Path, "/")
	if len(parts) < 2 {
		return nil, errors.New("unknown service")
	}

	if !re.MatchString(parts[1]) {
		return nil, res.ErrInvalidPath
	}

	_, namespace, _ := r.Info(req)
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
