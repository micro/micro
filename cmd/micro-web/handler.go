package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"go-micro.dev/v5/client"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
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
		margin: 0;
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

func WebHandler() http.Handler {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// MCP server on /mcp
		if r.URL.Path == "/mcp" {
			// Import here to avoid circular import at top
			apiPkg, err := importAPI()
			if err == nil && apiPkg != nil {
				apiPkg.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/web/") {
			// Remove /web prefix
			path := strings.TrimPrefix(r.URL.Path, "/web")
			if path == "" || path == "/" {
				services, _ := registry.ListServices()
				var html string
				for _, service := range services {
					html += fmt.Sprintf(`<p><a href="/web/%s">%s</a></p>`, url.QueryEscape(service.Name), service.Name)
				}
				html = fmt.Sprintf(htmlTemplate, html)
				w.Write([]byte(html))
				return
			}
			parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
			if len(parts) < 1 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			service := parts[0]
			s, err := registry.GetService(service)
			if err != nil || len(s) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if len(parts) < 2 {
				var endpoints string
				for _, ep := range s[0].Endpoints {
					parts := strings.Split(ep.Name, ".")
					uri := fmt.Sprintf("/web/%s/%s/%s", service, parts[0], parts[1])
					endpoints += fmt.Sprintf(`<div><a href="%s">%s</a>`, uri, ep.Name)
				}
				b, _ := json.MarshalIndent(s[0], "", "    ")
				serviceHTML := fmt.Sprintf(serviceTemplate, service, endpoints, string(b))
				render(w, serviceHTML)
				return
			}
			endpoint := parts[1]
			if len(parts) == 3 {
				endpoint = normalize(endpoint) + "." + normalize(parts[2])
			} else {
				endpoint = normalize(service) + "." + normalize(endpoint)
			}
			var ep *registry.Endpoint
			for _, eps := range s[0].Endpoints {
				if eps.Name == endpoint {
					ep = eps
					break
				}
			}
			if ep == nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
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
				pretty, _ := json.MarshalIndent(response, "", "    ")
				output += fmt.Sprintf(`<pre>%s</pre>`, string(pretty))
				render(w, fmt.Sprintf(responseTemplate, output))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// fallback: root web UI
		services, _ := registry.ListServices()
		var html string
		for _, service := range services {
			html += fmt.Sprintf(`<p><a href="/web/%s">%s</a></p>`, url.QueryEscape(service.Name), service.Name)
		}
		html = fmt.Sprintf(htmlTemplate, html)
		w.Write([]byte(html))
// importAPI returns the /mcp handler from the api package, or nil if not available
func importAPI() (http.Handler, error) {
	// Import locally to avoid circular import at top
	apiPkg, err := importAPIPkg()
	if err != nil {
		return nil, err
	}
	return apiPkg, nil
}

// importAPIPkg is a helper to import the APIHandler's /mcp endpoint
func importAPIPkg() (http.Handler, error) {
	// Import the handler from the api package
	// This import path is correct for this monorepo
	// If you move files, update this import
	api "github.com/micro/micro/v5/cmd/micro-api"
	return api.APIHandler(), nil
}
	})
	return h
}

func render(w http.ResponseWriter, v string) error {
	html := fmt.Sprintf(htmlTemplate, v)
	_, err := w.Write([]byte(html))
	return err
}

func normalize(v string) string {
	return strings.Title(v)
}

func rpcCall(service, endpoint string, request []byte) ([]byte, error) {
	data := []byte(`{}`)
	if len(request) > 0 {
		data = request
	}
	req := client.NewRequest(service, endpoint, &bytes.Frame{Data: data})
	var rsp bytes.Frame
	err := client.Call(nil, req, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp.Data, nil
}
