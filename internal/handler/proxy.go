package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"

	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/selector"
)

type proxy struct {
	Default  *httputil.ReverseProxy
	Director func(r *http.Request)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !isWebSocket(r) {
		// the usual path
		p.Default.ServeHTTP(w, r)
		return
	}

	// the websocket path
	req := new(http.Request)
	*req = *r
	p.Director(req)
	host := req.URL.Host

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

func Proxy(ns string) http.Handler {
	sel := selector.NewSelector(
		selector.Registry((*cmd.DefaultOptions().Registry)),
	)

	re := regexp.MustCompile("^[a-zA-Z0-9]+$")

	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			kill()
			return
		}
		if !re.MatchString(parts[1]) {
			kill()
			return
		}
		next, err := sel.Select(ns + "." + parts[1])
		if err != nil {
			kill()
			return
		}

		s, err := next()
		if err != nil {
			kill()
			return
		}

		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Scheme = "http"
	}

	return &proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}
