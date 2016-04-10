package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/selector"
)

type proxy struct {
	Selector  selector.Selector
	Namespace string

	regex     *regexp.Regexp
	wsEnabled bool
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serviceHost, err := p.serviceHostForRequest(r)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if len(serviceHost) == 0 {
		w.WriteHeader(404)
		return
	}

	if isWebSocket(r) && p.wsEnabled {
		p.serveWebSocket(serviceHost, w, r)
		return
	}

	rpURL, err := url.Parse(serviceHost)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP(w, r)
}

func (p *proxy) serviceHostForRequest(r *http.Request) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 {
		return "", nil
	}
	if !p.regex.MatchString(parts[1]) {
		return "", nil
	}
	next, err := p.Selector.Select(p.Namespace + "." + parts[1])
	if err != nil {
		return "", nil
	}

	s, err := next()
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("http://%s:%d", s.Address, s.Port), nil
}

func (p *proxy) serveWebSocket(host string, w http.ResponseWriter, r *http.Request) {
	// the websocket path
	req := new(http.Request)
	*req = *r

	if len(host) == 0 {
		http.Error(w, "invalid host", 500)
		return
	}

	// connect to the backend host
	conn, err := net.Dial("tcp", host)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// hijack the connection
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "failed to connect", 500)
		return
	}

	nc, _, err := hj.Hijack()
	if err != nil {
		return
	}

	defer nc.Close()
	defer conn.Close()

	if err = req.Write(conn); err != nil {
		return
	}

	errCh := make(chan error, 2)

	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errCh <- err
	}

	go cp(conn, nc)
	go cp(nc, conn)

	<-errCh
}

func isWebSocket(r *http.Request) bool {
	contains := func(key, val string) bool {
		vv := strings.Split(r.Header.Get(key), ",")
		for _, v := range vv {
			if val == strings.ToLower(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}

	if contains("Connection", "upgrade") && contains("Upgrade", "websocket") {
		return true
	}

	return false
}

func Proxy(ns string, ws bool) http.Handler {
	return &proxy{
		Namespace: ns,
		Selector: selector.NewSelector(
			selector.Registry((*cmd.DefaultOptions().Registry)),
		),

		regex:     regexp.MustCompile("^[a-zA-Z0-9]+$"),
		wsEnabled: ws,
	}
}
