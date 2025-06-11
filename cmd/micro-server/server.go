package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
)

// Placeholder for dashboard/API logic
// Color codes for log output

func Run(c *cli.Context) error {
	addr := c.String("address")
	if addr == "" {
		addr = ":8080"
	}

	// HTML templates for micro web UI
	var htmlTemplate = `<!DOCTYPE html><html><head><meta charset="UTF-8" />
	<meta name="viewport" content="width=device-width" />
	<title>Micro Web</title>
	<script src="https://cdn.tailwindcss.com"></script>
	<style>
	.micro-link {@apply inline-block font-bold border-2 border-gray-400 rounded-lg px-4 py-2 bg-gray-50 mr-2 mb-2 transition-colors;border: 2px solid #888;border-radius: 8px;padding: 8px 18px;margin: 8px 8px 8px 0;}.micro-link:hover {@apply bg-gray-100;background: #f3f4f6;}#title { text-decoration: none; color: black; border: none; padding: 0; margin: 0; }#title:hover { background: none; }
	</style>
	</head>
	<body class="bg-gray-50 text-gray-900"><div id="head" class="relative m-6"><h1><a href="/" id="title" class="text-3xl font-bold">Micro</a></h1><a id="web-link" href="%s" class="micro-link absolute top-0 right-24">Web</a><a id="api-link" href="/api" class="micro-link absolute top-0 right-2">API</a></div><div class="container px-6 py-4 max-w-screen-xl mx-auto">%s</div>
	</body>
	</html>`

	// Helper functions
	normalize := func(v string) string {
		if len(v) == 0 {
			return v
		}
		return strings.ToUpper(v[:1]) + v[1:]
	}
	render := func(w http.ResponseWriter, v string, webLink string) error {
		html := fmt.Sprintf(htmlTemplate, webLink, v)
		_, err := w.Write([]byte(html))
		return err
	}
	rpcCall := func(service, endpoint string, request []byte) ([]byte, error) {
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

	// Serve web and API
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		absDir, _ := os.Getwd()
		webDir := filepath.Join(absDir, "web")
		parentDir := filepath.Base(absDir)
		webLink := "/web"
		if _, err := os.Stat(webDir); err != nil {
			webLink = "/"
		}

		// --- Handle /api prefix for micro-api functionality ---
		if r.URL.Path == "/api" || r.URL.Path == "/api/" {
			// Render API documentation page
			services, _ := registry.ListServices()
			var html string
			html += `<h2 class="text-2xl font-bold mb-4">API Endpoints</h2>`
			for _, srv := range services {
				srvs, err := registry.GetService(srv.Name)
				if err != nil || len(srvs) == 0 {
					continue
				}
				s := srvs[0]
				if len(s.Endpoints) == 0 {
					continue
				}
				html += fmt.Sprintf(`<h3 class="text-xl font-semibold mt-8 mb-2">%s</h3>`, s.Name)

				for _, ep := range s.Endpoints {
					parts := strings.Split(ep.Name, ".")
					if len(parts) != 2 {
						continue
					}
					apiPath := fmt.Sprintf("/api/%s/%s/%s", s.Name, parts[0], parts[1])
					var params string
					if ep.Request != nil && len(ep.Request.Values) > 0 {
						params += "<ul class=\"ml-4 mb-2\">"
						for _, v := range ep.Request.Values {
							params += fmt.Sprintf("<li><b>%s</b> <span class=\"text-gray-500\">%s</span></li>", v.Name, v.Type)
						}
						params += "</ul>"
					} else {
						params = "<i class=\"text-gray-500\">No parameters</i>"
					}
					var response string
					if ep.Response != nil && len(ep.Response.Values) > 0 {
						response += "<ul class=\"ml-4 mb-2\">"
						for _, v := range ep.Response.Values {
							response += fmt.Sprintf("<li><b>%s</b> <span class=\"text-gray-500\">%s</span></li>", v.Name, v.Type)
						}
						response += "</ul>"
					} else {
						response = "<i class=\"text-gray-500\">No response fields</i>"
					}
					html += fmt.Sprintf(
						`<div class="mb-10"><div class="text-lg font-bold mb-1">%s</div><hr class="mb-4 border-gray-300"><div class="mb-2"><span class="font-bold">HTTP Path:</span> <code class="font-mono">%s</code></div><div class="mb-2"><span class="font-bold">Request:</span> %s</div><div class="mb-2"><span class="font-bold">Response:</span> %s</div></div>`,
						ep.Name, apiPath, params, response,
					)
				}
			}
			render(w, html, webLink)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") {
			// /api/{service}/{endpointService}/{endpointMethod}
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 5 {
				http.Error(w, "Invalid API path. Use /api/{service}/{endpointService}/{endpointMethod}", 400)
				return
			}
			service := parts[2]
			endpointService := parts[3]
			endpointMethod := parts[4]
			endpointName := normalize(endpointService) + "." + normalize(endpointMethod)

			// Support GET params, POST form, and JSON body
			var reqBody map[string]interface{}

			// Prefer JSON body if present and Content-Type is application/json
			if r.Method == "POST" && strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				defer r.Body.Close()
				json.NewDecoder(r.Body).Decode(&reqBody)
			}

			// If not JSON, or for GET, use URL query/form values
			if reqBody == nil {
				reqBody = map[string]interface{}{}
				// Parse form for POST, or query for GET
				r.ParseForm()
				for k, v := range r.Form{
					if len(v) == 1 {
						if len(v[0]) == 0 {
							continue
						}
						reqBody[k] = v[0]
					} else {
						if len(v) == 0 {
							continue
						}
						reqBody[k] = v
					}
				}
			}

			b, _ := json.Marshal(reqBody)
			rsp, err := rpcCall(service, endpointName, b)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(rsp)
			return
		}

		// --- Serve /services and /service/{service} always, never remap ---
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 2 && parts[1] == "services" {
			services, _ := registry.ListServices()
			html := `<h2 class="text-2xl font-bold mb-4">Services</h2>`
			for _, service := range services {
				html += fmt.Sprintf(`<button onclick="location.href='/service/%s'" class="micro-link">%s</button>`, url.QueryEscape(service.Name), service.Name)
			}
			render(w, html, webLink)
			return
		}

		// /service/foo
		// /service/foo/bar

		if len(parts) >= 3 && parts[1] == "service" && parts[2] != "" {
			service := parts[2]
			service, _ = url.QueryUnescape(service)
			s, err := registry.GetService(service)
			if err != nil || len(s) == 0 {
				w.WriteHeader(404)
				w.Write([]byte(fmt.Sprintf("Service not found: %s", service)))
				return
			}

			if len(parts) <= 3 {
				var endpoints string
				endpoints += `<h4 class="font-semibold mb-2">Endpoints</h4>`
				if len(s[0].Endpoints) == 0 {
					endpoints += "<p>No endpoints registered</p>"
				}

				for _, ep := range s[0].Endpoints {
					p := strings.Split(ep.Name, ".")
					if len(p) != 2 {
						endpoints += "<p>" + ep.Name + "</p>"
						continue
					}
					uri := fmt.Sprintf("/service/%s/%s/%s", service, p[0], p[1])
					endpoints += fmt.Sprintf(`<button onclick="location.href='%s'" class="micro-link">%s</button>`, uri, ep.Name)
				}
				b, _ := json.MarshalIndent(s[0], "", "    ")
				serviceHTML := fmt.Sprintf(
					`<h2 class="text-xl font-bold mb-2">%s</h2>%s<h4 class="font-semibold mt-4 mb-2">Description</h4><pre class="bg-gray-100 rounded p-2">%s</pre>`,
					service, endpoints, string(b),
				)
				render(w, serviceHTML, webLink)
				return
			}

			endpoint := parts[3]
			if len(parts) == 5 {
				endpoint = normalize(endpoint) + "." + normalize(parts[4])
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
				w.WriteHeader(404)
				w.Write([]byte("Endpoint not found"))
				return
			}
			if r.Method == "GET" {
				var inputs string
				inputs += fmt.Sprintf(`<h3 class="text-lg font-bold mb-2">%s</h3>`, ep.Name)
				for _, input := range ep.Request.Values {
					inputs += fmt.Sprintf(`<label class="block font-semibold">%s</label><input id=%s name=%s placeholder=%s class="border rounded px-2 py-1 mb-2 w-full">`, input.Name, input.Name, input.Name, input.Name)
				}
				inputs += `<button class="micro-link mt-2" type="submit">Submit</button>`
				formHTML := fmt.Sprintf(`<h2>%s</h2><form action=%s method=POST>%s</form>`, service, r.URL.Path, inputs)
				render(w, formHTML, webLink)
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
					// Render form with error below
					var inputs string
					inputs += fmt.Sprintf(`<h3 class="text-lg font-bold mb-2">%s</h3>`, ep.Name)
					for _, input := range ep.Request.Values {
						inputs += fmt.Sprintf(`<label class="block font-semibold">%s</label><input id=%s name=%s placeholder=%s class="border rounded px-2 py-1 mb-2 w-full" value="%s">`, input.Name, input.Name, input.Name, input.Name, r.Form.Get(input.Name))
					}
					inputs += `<button class="micro-link mt-2" type="submit">Submit</button>`
					formHTML := fmt.Sprintf(`<h2>%s</h2><form action=%s method=POST>%s</form>`, service, r.URL.Path, inputs)
					errorHTML := fmt.Sprintf(`<div class="mt-4 text-red-600 font-bold">Error: %s</div>`, err.Error())
					render(w, formHTML+errorHTML, webLink)
					return
				}
				var response map[string]interface{}
				json.Unmarshal(rsp, &response)
				// Build response table
				var tableRows string
				for k, v := range response {
					tableRows += fmt.Sprintf(`<tr><td class="border px-2 py-1 font-semibold">%s</td><td class="border px-2 py-1">%v</td></tr>`, k, v)
				}
				tableHTML := `<table class="table-auto border-collapse border border-gray-300 mt-4 mb-2"><thead><tr><th class="border px-2 py-1">Field</th><th class="border px-2 py-1">Value</th></tr></thead><tbody>` + tableRows + `</tbody></table>`
				pretty, _ := json.MarshalIndent(response, "", "    ")
				jsonHTML := fmt.Sprintf(`<pre class="bg-gray-100 rounded p-2 mt-2">%s</pre>`, string(pretty))
				// Render form + response
				var inputs string
				inputs += fmt.Sprintf(`<h3 class="text-lg font-bold mb-2">%s</h3>`, ep.Name)
				for _, input := range ep.Request.Values {
					inputs += fmt.Sprintf(`<label class="block font-semibold">%s</label><input id=%s name=%s placeholder=%s class="border rounded px-2 py-1 mb-2 w-full" value="%s">`, input.Name, input.Name, input.Name, input.Name, r.Form.Get(input.Name))
				}
				inputs += `<button class="micro-link mt-2" type="submit">Submit</button>`
				formHTML := fmt.Sprintf(`<h2>%s</h2><form action=%s method=POST>%s</form>`, service, r.URL.Path, inputs)
				responseHTML := `<div class="mt-4"><h4 class="font-bold mb-2">Response</h4>` + tableHTML + jsonHTML + `</div>`
				render(w, formHTML+responseHTML, webLink)
				return
			}
		}

		// --- Serve /web and proxy all other requests to the web app if webDir exists ---
		if _, err := os.Stat(webDir); err == nil {
			if r.URL.Path == "/web" || r.URL.Path == "/web/" {
				html := `<h2 class="text-2xl font-bold mb-4">Web</h2><button onclick="location.href='/services'" class="micro-link">Services</button>`
				render(w, html, webLink)
				return
			}
			// Proxy everything else (except /api and /services/service) to the web app
			if !strings.HasPrefix(r.URL.Path, "/api") && !strings.HasPrefix(r.URL.Path, "/services") && !strings.HasPrefix(r.URL.Path, "/service") {
				srvs, err := registry.GetService(parentDir)
				if err == nil && len(srvs) > 0 && len(srvs[0].Nodes) > 0 {
					target := srvs[0].Nodes[0].Address
					u, _ := url.Parse("http://" + target)
					proxy := httputil.NewSingleHostReverseProxy(u)
					proxy.ServeHTTP(w, r)
					return
				}
			}
		}

		// --- Default index page ---
		if len(parts) < 2 || parts[1] == "" {
			html := `<h2 class="text-2xl font-bold mb-4">Web</h2><button onclick="location.href='/services'" class="micro-link">Services</button>`
			render(w, html, webLink)
			return
		}

		// --- Fallback: 404 ---
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
	})

	go func() {
		log.Printf("[micro-server] Web/API listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Web/API server error: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Shutting down micro server...")
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "server",
		Usage:  "Start the Micro server (dashboard/API/web UI)",
		Action: Run,
		Flags:  []cli.Flag{},
	})
}
