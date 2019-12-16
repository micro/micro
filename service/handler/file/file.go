package file

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/proxy"
	"github.com/micro/go-micro/server"
)

//Proxy for a proxy instance
type Proxy struct {
	options proxy.Options

	// The file or directory to read from
	Endpoint string
}

func filePath(eps ...string) string {
	p := filepath.Join(eps...)
	return strings.Replace(p, "../", "", -1)
}

func getMethod(hdr map[string]string) string {
	switch hdr["Micro-Method"] {
	case "read", "write":
		return hdr["Micro-Method"]
	default:
		return "read"
	}
}

func getEndpoint(hdr map[string]string) string {
	ep := hdr["Micro-Endpoint"]
	if len(ep) > 0 && ep[0] == '/' {
		return ep
	}
	return ""
}

func (p *Proxy) ProcessMessage(ctx context.Context, msg server.Message) error {
	return nil
}

// ServeRequest honours the server.Router interface
func (p *Proxy) ServeRequest(ctx context.Context, req server.Request, rsp server.Response) error {
	if p.Endpoint == "" {
		exe, err := os.Executable()
		if err != nil {
			return err
		}
		// set the endpoint to the current path
		p.Endpoint = filepath.Dir(exe)
	}

	for {
		// get data
		// Read the body if we're writing the file
		_, err := req.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// get the header
		hdr := req.Header()

		// get method
		//method := getMethod(hdr)

		// get endpoint
		endpoint := getEndpoint(hdr)

		// filepath
		file := filePath(p.Endpoint, endpoint)

		// lookup the file
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}

		// write back the header
		rsp.WriteHeader(hdr)
		// write the body
		err = rsp.Write(b)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.InternalServerError(req.Service(), err.Error())
		}
	}

}

func (p *Proxy) String() string {
	return "file"
}

//NewSingleHostProxy returns a Proxy which stand for a endpoint.
func NewSingleHostProxy(url string) proxy.Proxy {
	return &Proxy{
		Endpoint: url,
	}
}

// NewProxy returns a new proxy which will route using a http client
func NewProxy(opts ...proxy.Option) proxy.Proxy {
	var options proxy.Options
	for _, o := range opts {
		o(&options)
	}

	p := new(Proxy)
	p.Endpoint = options.Endpoint

	return p
}
