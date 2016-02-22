package web

import (
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
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

	if !contains("Connection", "upgrade") {
		return false
	}

	if !contains("Upgrade", "websocket") {
		return false
	}

	return true
}
