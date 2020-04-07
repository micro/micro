package resolver

import (
	"net"
	"net/http"
	"strings"

	"github.com/micro/go-micro/v2/api/resolver"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/logger"
	"golang.org/x/net/publicsuffix"
)

func NewNamespaceResolver(srvType, namespace string) resolver.NamespaceResolver {
	return &Resolver{srvType, namespace}
}

type Resolver struct {
	srvType   string
	namespace string
}

func (r Resolver) String() string {
	return "internal/namespace"
}

func (r Resolver) Resolve(req *http.Request) string {
	withTypeSuffix := func(ns string) string {
		return ns + "." + r.srvType
	}

	// check to see what the provided namespace is, we only do
	// domain mapping if the namespace is set to 'domain'
	if r.namespace != "domain" {
		return withTypeSuffix(r.namespace)
	}

	// determine the host, e.g. dev.micro.mu:8080
	host := req.URL.Hostname()
	if len(host) == 0 {
		if h, _, err := net.SplitHostPort(req.Host); err == nil {
			host = h // host does contain a port
		} else if strings.Contains(err.Error(), "missing port in address") {
			host = req.Host // host does not contain a port
		}
	}

	// check for an ip address
	if net.ParseIP(host) != nil {
		return withTypeSuffix(auth.DefaultNamespace)
	}

	// check for dev enviroment
	if host == "localhost" || host == "127.0.0.1" {
		return withTypeSuffix(auth.DefaultNamespace)
	}

	// extract the top level domain plus one (e.g. 'myapp.com')
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		logger.Debugf("Unable to extract domain from %v", host)
		return withTypeSuffix(auth.DefaultNamespace)
	}

	// check to see if the domain matches the host of micro.mu, in
	// these cases we return the default namespace
	if domain == host || domain == "micro.mu" {
		return withTypeSuffix(auth.DefaultNamespace)
	}

	// remove the domain from the host, leaving the subdomain
	subdomain := strings.TrimSuffix(host, "."+domain)

	// return the reversed subdomain as the namespace
	comps := strings.Split(subdomain, ".")
	for i := len(comps)/2 - 1; i >= 0; i-- {
		opp := len(comps) - 1 - i
		comps[i], comps[opp] = comps[opp], comps[i]
	}
	return withTypeSuffix(strings.Join(comps, "."))
}
