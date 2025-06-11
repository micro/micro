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
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
	"io/fs"
)

// HTML is the embedded filesystem for templates and static files, set by main.go
var HTML fs.FS

func Run(c *cli.Context) error {
	addr := c.String("address")
	if addr == "" {
		addr = ":8080"
	}

	// Use embedded html directory for templates and static files
	staticFS, _ := fs.Sub(HTML, "html")

	parseTmpl := func(name string) *template.Template {
		tmpl, err := template.ParseFS(HTML, "html/base.html", "html/"+name)
		if err != nil {
			panic(err)
		}
		return tmpl
	}

	apiTmpl := parseTmpl("api.html")
	serviceTmpl := parseTmpl("service.html")
	formTmpl := parseTmpl("form.html")
	webTmpl := parseTmpl("web.html")
	logsTmpl := parseTmpl("logs.html")
	logTmpl := parseTmpl("log.html")

	render := func(w http.ResponseWriter, tmpl *template.Template, data any) error {
		return tmpl.Execute(w, data)
	}

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.FS(staticFS))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			// Home page
			_ = render(w, webTmpl, map[string]any{"Title": "Micro Web", "WebLink": "/", "Content": nil})
			return
		}
		if path == "/api" || path == "/api/" {
			// API page
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
			_ = render(w, apiTmpl, map[string]any{"Title": "API", "WebLink": "/", "Content": html})
			return
		}
		if path == "/services" {
			// List services
			services, _ := registry.ListServices()
			html := `<h2 class="text-2xl font-bold mb-4">Services</h2>`
			for _, service := range services {
				html += fmt.Sprintf(`<button onclick="location.href='/service/%s'" class="micro-link">%s</button>`, url.QueryEscape(service.Name), service.Name)
			}
			_ = render(w, serviceTmpl, map[string]any{"Title": "Services", "WebLink": "/", "Content": html})
			return
		}
		if path == "/logs" || path == "/logs/" {
			// List all services for logs
			services, _ := registry.ListServices()
			serviceNames := []string{}
			for _, srv := range services {
				serviceNames = append(serviceNames, srv.Name)
			}
			_ = render(w, logsTmpl, map[string]any{"Title": "Logs", "WebLink": "/", "Services": serviceNames})
			return
		}
		if strings.HasPrefix(path, "/logs/") {
			// Show logs for a specific service
			service := strings.TrimPrefix(path, "/logs/")
			if service == "" {
				w.WriteHeader(404)
				w.Write([]byte("Service not specified"))
				return
			}
			// Run 'micro logs <service>' and capture output
			cmd := exec.Command("micro", "logs", service)
			output, err := cmd.CombinedOutput()
			logText := string(output)
			if err != nil && logText == "" {
				logText = err.Error()
			}
			_ = render(w, logTmpl, map[string]any{"Title": "Logs for " + service, "WebLink": "/logs", "Service": service, "Log": logText})
			return
		}
		// Match /{service} and /{service}/{endpoint}
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 1 && parts[0] != "api" && parts[0] != "html" && parts[0] != "services" {
			service := parts[0]
			if len(parts) == 1 {
				// Service page
				s, err := registry.GetService(service)
				if err != nil || len(s) == 0 {
					w.WriteHeader(404)
					w.Write([]byte(fmt.Sprintf("Service not found: %s", service)))
					return
				}
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
				_ = render(w, serviceTmpl, map[string]any{"Title": "Service: " + service, "WebLink": "/", "Content": serviceHTML})
				return
			}
			if len(parts) == 2 {
				// Endpoint form
				service := parts[0]
				endpoint := normalize(service) + "." + normalize(parts[1])
				s, err := registry.GetService(service)
				if err != nil || len(s) == 0 {
					w.WriteHeader(404)
					w.Write([]byte(fmt.Sprintf("Service not found: %s", service)))
					return
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
					_ = render(w, formTmpl, map[string]any{"Title": "Service: " + service, "WebLink": "/", "Content": formHTML})
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
						_ = render(w, formTmpl, map[string]any{"Title": "Service: " + service, "WebLink": "/", "Content": formHTML + errorHTML})
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
					_ = render(w, formTmpl, map[string]any{"Title": "Service: " + service, "WebLink": "/", "Content": formHTML + responseHTML})
					return
				}
			}
		}
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
