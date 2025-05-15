package main

import (
	"context"
	"encoding/json"
	"fmt"
	//"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
	"tailscale.com/tsnet"
)

var htmlTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width" />
    <title>Micro Web</title>
    <style>
      html, body {
        font-family: Arial;
        font-size: 16px;
        margin: 0:
	padding: 0;
      }
      #head {
        margin: 25px;
      }
      #head a { color: black; text-decoration: none; }
      .container {
	 padding: 25px;
         max-width: 1400px;
         margin: 0 auto;
      }
      a { color: black; text-decoration: none; font-weight: bold; margin-bottom: 10px;}
      pre { background: #f5f5f5; border-radius: 5px; padding: 10px; overflow: scroll;}
      input, button { border-radius: 5px; padding: 10px; display: block; margin-bottom: 5px; }
      button:hover { cursor: pointer; }
    </style>
  </head>
  <body>
     <div id="head">
       <h1><a href="/">Micro</a></h1>
     </div>
     <div class="container">
	%s
     </div>
  </body>
</html>

`

var serviceTemplate = `
<h2>%s</h2>
<div>%s</div>
<pre>%s</pre>
`

var endpointTemplate = `
<h2>%s</h2>
<form action=%s method=POST>%s</form>
`
var responseTemplate = `<h2>Response</h2><div>%s</div>`

func normalize(v string) string {
	return strings.Title(v)
}

func render(w http.ResponseWriter, v string) error {
	html := fmt.Sprintf(htmlTemplate, v)
	_, err := w.Write([]byte(html))
	return err
}

func rpcCall(service, endpoint string, request []byte) ([]byte, error) {
	data := []byte(`{}`)

	if len(request) > 0 {
		data = request
	}

	req := client.NewRequest(service, endpoint, &bytes.Frame{Data: data})
	var rsp bytes.Frame
	err := client.Call(context.TODO(), req, &rsp)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// list serfvices
		if r.URL.Path == "/" {
			services, _ := registry.ListServices()

			var html string

			for _, service := range services {
				html += fmt.Sprintf(`<p><a href="/%s">%s</a></p>`, url.QueryEscape(service.Name), service.Name)
			}

			html = fmt.Sprintf(htmlTemplate, html)
			w.Write([]byte(html))

			return
		}

		// got more e.g /helloworld
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) < 2 {
			return
		}

		service := parts[1]

		// get service
		s, err := registry.GetService(service)
		if err != nil {
			return
		}

		// no service
		if len(s) == 0 {
			return
		}

		// service definition for /helloworld
		if len(parts) < 3 {
			var endpoints string
			for _, ep := range s[0].Endpoints {
				parts := strings.Split(ep.Name, ".")
				uri := fmt.Sprintf("/%s/%s/%s", service, parts[0], parts[1])
				endpoints += fmt.Sprintf(`<div><a href="%s">%s</a>`, uri, ep.Name)
			}

			b, _ := json.MarshalIndent(s[0], "", "    ")

			serviceHTML := fmt.Sprintf(serviceTemplate, service, endpoints, string(b))
			render(w, serviceHTML)
			return
		}

		endpoint := parts[2]

		if len(parts) == 4 {
			endpoint = normalize(endpoint) + "." + normalize(parts[3])
		} else {
			endpoint = normalize(service) + "." + normalize(endpoint)
		}

		// get the endpoint
		var ep *registry.Endpoint
		for _, eps := range s[0].Endpoints {
			if eps.Name == endpoint {
				ep = eps
				break
			}
		}

		// no endpoint match
		if ep == nil {
			return
		}

		// render form
		if r.Method == "GET" {
			var inputs string
			if ep != nil {
				for _, input := range ep.Request.Values {
					inputs += fmt.Sprintf(`<input id=%s name=%s placeholder=%s>`, input.Name, input.Name, input.Name)
				}
			} else {
				inputs = `<textarea></textarea>`
			}

			inputs += `<button>Submit</button>`

			formHTML := fmt.Sprintf(endpointTemplate, service, r.URL.Path, inputs)
			render(w, formHTML)
			return
		}

		if r.Method == "POST" {
			r.ParseForm()

			request := map[string]interface{}{}

			for k, v := range r.Form {
				request[k] = strings.Join(v, ",")
			}

			b, _ := json.Marshal(request)

			rsp, err := rpcCall(service, endpoint, b)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			var response map[string]interface{}
			json.Unmarshal(rsp, &response)

			var output string
			for k, v := range response {
				output += fmt.Sprintf(`<div>%s: %s</div>`, k, v)
			}

			render(w, fmt.Sprintf(responseTemplate, output))
			return
		}
	})

	app := cmd.App()

	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "network", Value: "", Usage: "Set the network e.g --network=tailscale requires TS_AUTHKEY"},
	}

	app.Action = func(c *cli.Context) error {
		var network string
		var key string

		if c.IsSet("network") {
			network = c.Value("network").(string)
		}

		if network == "tailscale" {
			// check for TS_AUTHKEY
			key = os.Getenv("TS_AUTHKEY")
			if len(key) == 0 {
				return fmt.Errorf("missing TS_AUTHKEY")
			}

			srv := new(tsnet.Server)
			srv.AuthKey = key
			srv.Hostname = "micro"

			ln, err := srv.Listen("tcp", ":8082")
			if err != nil {
				return err
			}

			return http.Serve(ln, h)
		}

		return http.ListenAndServe(":8082", h)
	}

	cmd.Init(
		cmd.Name("micro-web"),
	)
}
