package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	goMicroClient "go-micro.dev/v5/client"
	goMicroBytes "go-micro.dev/v5/codec/bytes"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/registry"
	htmltemplate "html/template"
)

// HTML is the embedded filesystem for templates and static files, set by main.go
var HTML fs.FS

func Run(c *cli.Context) error {
	addr := c.String("address")
	if addr == "" {
		addr = ":8080"
	}

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
			_ = render(w, webTmpl, map[string]any{"Title": "Micro Web", "WebLink": "/", "Content": nil})
			return
		}
		if path == "/api" || path == "/api/" {
			services, _ := registry.ListServices()
			var apiServices []map[string]any
			for _, srv := range services {
				srvs, err := registry.GetService(srv.Name)
				if err != nil || len(srvs) == 0 {
					continue
				}
				s := srvs[0]
				if len(s.Endpoints) == 0 {
					continue
				}
				endpoints := []map[string]any{}
				for _, ep := range s.Endpoints {
					parts := strings.Split(ep.Name, ".")
					if len(parts) != 2 {
						continue
					}
					apiPath := fmt.Sprintf("/api/%s/%s/%s", s.Name, parts[0], parts[1])
					var params, response string
					if ep.Request != nil && len(ep.Request.Values) > 0 {
						params += "<ul class=\"ml-4 mb-2\">"
						for _, v := range ep.Request.Values {
							params += fmt.Sprintf("<li><b>%s</b> <span class=\"text-gray-500\">%s</span></li>", v.Name, v.Type)
						}
						params += "</ul>"
					} else {
						params = "<i class=\"text-gray-500\">No parameters</i>"
					}
					if ep.Response != nil && len(ep.Response.Values) > 0 {
						response += "<ul class=\"ml-4 mb-2\">"
						for _, v := range ep.Response.Values {
							response += fmt.Sprintf("<li><b>%s</b> <span class=\"text-gray-500\">%s</span></li>", v.Name, v.Type)
						}
						response += "</ul>"
					} else {
						response = "<i class=\"text-gray-500\">No response fields</i>"
					}
					endpoints = append(endpoints, map[string]any{
						"Name": ep.Name,
						"Path": apiPath,
						"Params": htmltemplate.HTML(params),
						"Response": htmltemplate.HTML(response),
					})
				}
				apiServices = append(apiServices, map[string]any{
					"Name": s.Name,
					"Endpoints": endpoints,
				})
			}
			_ = render(w, apiTmpl, map[string]any{"Title": "API", "WebLink": "/", "Services": apiServices})
			return
		}
		if path == "/services" {
			services, _ := registry.ListServices()
			var serviceNames []string
			for _, service := range services {
				serviceNames = append(serviceNames, service.Name)
			}
			_ = render(w, serviceTmpl, map[string]any{"Title": "Services", "WebLink": "/", "Services": serviceNames})
			return
		}
		if path == "/logs" || path == "/logs/" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not get home directory"))
				return
			}
			logsDir := homeDir + "/micro/logs"
			dirEntries, err := os.ReadDir(logsDir)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not list logs directory: " + err.Error()))
				return
			}
			serviceNames := []string{}
			for _, entry := range dirEntries {
				name := entry.Name()
				if !entry.IsDir() && strings.HasSuffix(name, ".log") && !strings.HasPrefix(name, ".") {
					serviceNames = append(serviceNames, strings.TrimSuffix(name, ".log"))
				}
			}
			_ = render(w, logsTmpl, map[string]any{"Title": "Logs", "WebLink": "/", "Services": serviceNames})
			return
		}
		if strings.HasPrefix(path, "/logs/") {
			service := strings.TrimPrefix(path, "/logs/")
			if service == "" {
				w.WriteHeader(404)
				w.Write([]byte("Service not specified"))
				return
			}
			homeDir, err := os.UserHomeDir()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not get home directory"))
				return
			}
			logFilePath := homeDir + "/micro/logs/" + service + ".log"
			f, err := os.Open(logFilePath)
			if err != nil {
				w.WriteHeader(404)
				w.Write([]byte("Could not open log file for service: " + service))
				return
			}
			defer f.Close()
			logBytes, err := io.ReadAll(f)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not read log file for service: " + service))
				return
			}
			logText := string(logBytes)
			_ = render(w, logTmpl, map[string]any{"Title": "Logs for " + service, "WebLink": "/logs", "Service": service, "Log": logText})
			return
		}
		// Match /{service} and /{service}/{endpoint}
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) >= 1 && parts[0] != "api" && parts[0] != "html" && parts[0] != "services" {
			service := parts[0]
			if len(parts) == 1 {
				s, err := registry.GetService(service)
				if err != nil || len(s) == 0 {
					w.WriteHeader(404)
					w.Write([]byte(fmt.Sprintf("Service not found: %s", service)))
					return
				}
				endpoints := []map[string]string{}
				for _, ep := range s[0].Endpoints {
					endpoints = append(endpoints, map[string]string{
						"Name": ep.Name,
						"Path": fmt.Sprintf("/%s/%s", service, ep.Name),
					})
				}
				b, _ := json.MarshalIndent(s[0], "", "    ")
				_ = render(w, serviceTmpl, map[string]any{
					"Title": "Service: " + service,
					"WebLink": "/",
					"ServiceName": service,
					"Endpoints": endpoints,
					"Description": string(b),
				})
				return
			}
			if len(parts) == 2 {
				service := parts[0]
				endpoint := parts[1] // Use the actual endpoint name from the URL, e.g. Foo.Bar
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
					// Build form fields from endpoint request values
					var inputs []map[string]string
					if ep.Request != nil && len(ep.Request.Values) > 0 {
						for _, input := range ep.Request.Values {
							inputs = append(inputs, map[string]string{
								"Label":      input.Name,
								"Name":       input.Name,
								"Placeholder": input.Name,
								"Value":      "",
							})
						}
					}
					_ = render(w, formTmpl, map[string]any{
						"Title":       "Service: " + service,
						"WebLink":     "/",
						"ServiceName": service,
						"EndpointName": ep.Name,
						"Inputs":      inputs,
						"Action":      service + "/" + endpoint,
					})
					return
				}
				if r.Method == "POST" {
					var reqBody map[string]interface{}
					if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
						defer r.Body.Close()
						json.NewDecoder(r.Body).Decode(&reqBody)
					} else {
						reqBody = map[string]interface{}{}
						r.ParseForm()
						for k, v := range r.Form {
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
					req := goMicroClient.NewRequest(service, endpoint, &goMicroBytes.Frame{Data: b})
					var rsp goMicroBytes.Frame
					err := goMicroClient.Call(r.Context(), req, &rsp)
					if err != nil {
						w.WriteHeader(500)
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(rsp.Data)
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

	ch := make(chan struct{})
	<-ch
	return nil
}

func normalize(v string) string {
	if len(v) == 0 {
		return v
	}
	return strings.ToUpper(v[:1]) + v[1:]
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "server",
		Usage:  "Start the Micro server (dashboard/API/web UI)",
		Action: Run,
		Flags:  []cli.Flag{},
	})
}
