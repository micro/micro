package run

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/client"
	"go-micro.dev/v5/cmd"
	"go-micro.dev/v5/codec/bytes"
	"go-micro.dev/v5/registry"
)

// Color codes for log output
var colors = []string{
	"\033[31m", // red
	"\033[32m", // green
	"\033[33m", // yellow
	"\033[34m", // blue
	"\033[35m", // magenta
	"\033[36m", // cyan
}

func colorFor(idx int) string {
	return colors[idx%len(colors)]
}

// HTML templates for micro web UI
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
        position: relative;
      }
      #head a { color: black; text-decoration: none; }
      #api-link, #web-link {
        position: absolute;
        top: 25px;
        font-size: 18px;
        font-weight: bold;
      }
      #api-link { right: 25px; }
      #web-link { right: 90px; }
      .container {
         padding: 25px;
         max-width: 1400px;
         margin: 0 auto;
      }
      .micro-link {
        color: black;
        text-decoration: none;
        font-weight: bold;
        margin-bottom: 10px;
        margin-right: 10px;
        margin-top: 10px;
        margin-left: 0;
        border: 2px solid #888;
        border-radius: 8px;
        padding: 5px 14px;
        display: inline-block;
        transition: background 0.15s;
        background: #f9f9f9;
      }
      .micro-link:hover {
        background: whitesmoke;
      }
      #title { text-decoration: none; color: black; border: none; padding: 0; margin: 0; }
      #title:hover { background: none; }
      pre { background: #f5f5f5; border-radius: 5px; padding: 10px; overflow: scroll;}
      input, button { border-radius: 5px; padding: 10px; display: block; margin-bottom: 5px; }
      button:hover { cursor: pointer; }
    </style>
  </head>
  <body>
     <div id="head">
       <h1><a href="/" id="title">Micro</a></h1>
       <a id="web-link" href="/web" class="micro-link">Web</a>
       <a id="api-link" href="/api" class="micro-link">API</a>
     </div>
     <div class="container">
	%s
     </div>
  </body>
</html>
`

var serviceTemplate = `
<h2>%s</h2>
<h3>Endpoints</h3>
<div>%s</div>
<h3>Description</h3>
<pre>%s</pre>
`

var endpointTemplate = `
<h2>%s</h2>
<form action=%s method=POST>%s</form>
`
var responseTemplate = `<h2>Response</h2><div>%s</div>`

func normalize(v string) string {
	if len(v) == 0 {
		return v
	}
	return strings.ToUpper(v[:1]) + v[1:]
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

func serveMicroWeb(dir string, addr string) {
	// Always resolve to absolute path for dir
	absDir, err := filepath.Abs(dir)
	if err != nil {
		absDir = dir // fallback
	}
	webDir := filepath.Join(absDir, "web")
	parentDir := filepath.Base(absDir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// --- Handle /api prefix for micro-api functionality ---
		if r.URL.Path == "/api" || r.URL.Path == "/api/" {
			// Render API documentation page
			services, _ := registry.ListServices()
			var html string
			html += "<h2>API Endpoints</h2>"
			for _, srv := range services {
				srvs, err := registry.GetService(srv.Name)
				if err != nil || len(srvs) == 0 {
					continue
				}
				s := srvs[0]
				html += fmt.Sprintf(`<h3>%s</h3>`, s.Name)
				for _, ep := range s.Endpoints {
					// Parse endpoint name
					parts := strings.Split(ep.Name, ".")
					if len(parts) != 2 {
						continue
					}
					// Build API path
					apiPath := fmt.Sprintf("/api/%s/%s/%s", s.Name, parts[0], parts[1])
					// Params
					var params string
					if ep.Request != nil && len(ep.Request.Values) > 0 {
						params += "<ul>"
						for _, v := range ep.Request.Values {
							params += fmt.Sprintf("<li><b>%s</b> (%s)</li>", v.Name, v.Type)
						}
						params += "</ul>"
					} else {
						params = "<i>No parameters</i>"
					}
					// Response
					var response string
					if ep.Response != nil && len(ep.Response.Values) > 0 {
						response += "<ul>"
						for _, v := range ep.Response.Values {
							response += fmt.Sprintf("<li><b>%s</b> (%s)</li>", v.Name, v.Type)
						}
						response += "</ul>"
					} else {
						response = "<i>No response fields</i>"
					}
					html += fmt.Sprintf(
						`<div><code>%s</code></div>
						  <hr>
						  <div style="margin: 1em 1em 2em 1em;">
							<div><b>HTTP Path:</b> <code>%s</code></div>
							<br>
							<div><b>Parameters:</b> %s</div>
							<div><b>Response:</b> %s</div>
						  </div>`,
						ep.Name, apiPath, params, response,
					)
				}
			}
			render(w, html)
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
				for k, v := range r.Form {
					if len(v) == 1 {
						reqBody[k] = v[0]
					} else {
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

		// --- Restore web reverse proxy logic ---
		if _, err := os.Stat(webDir); err == nil {
			// If /web or /web/ is requested, serve the micro web index
			if r.URL.Path == "/web" || r.URL.Path == "/web/" {
				// Render the micro web index page (same as old "/")
				html := `<h2>Web</h2>
            <p><a href="/web/services" class="micro-link">Services</a></p>`
				render(w, html)
				return
			}
			// Remap /services and /{service}... to /web/services and /web/{service}...
			if r.URL.Path == "/services" || strings.HasPrefix(r.URL.Path, "/services/") {
				http.Redirect(w, r, "/web"+r.URL.Path, http.StatusFound)
				return
			}
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) > 1 && parts[1] != "" && parts[1] != "web" && parts[1] != "api" {
				// e.g. /foo, /foo/bar
				http.Redirect(w, r, "/web"+r.URL.Path, http.StatusFound)
				return
			}
			// web subdir exists, look for service by parent dir name
			srvs, err := registry.GetService(parentDir)
			if err == nil && len(srvs) > 0 && len(srvs[0].Nodes) > 0 {
				// reverse proxy to first node
				target := srvs[0].Nodes[0].Address
				u, _ := url.Parse("http://" + target)
				proxy := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					proxyReq, _ := http.NewRequest(req.Method, u.String()+req.RequestURI, req.Body)
					for k, v := range req.Header {
						proxyReq.Header[k] = v
					}
					resp, err := http.DefaultClient.Do(proxyReq)
					if err != nil {
						http.Error(w, "Proxy error", 502)
						return
					}
					defer resp.Body.Close()
					for k, v := range resp.Header {
						w.Header()[k] = v
					}
					w.WriteHeader(resp.StatusCode)
					io.Copy(w, resp.Body)
				})
				// Only proxy if not /web or /web/
				if !(r.URL.Path == "/web" || r.URL.Path == "/web/") {
					proxy.ServeHTTP(w, r)
					return
				}
			}
		}

		// --- Custom routing for / and /services ---
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) == 2 && parts[1] == "services" {
			// List all services on /services
			services, _ := registry.ListServices()
			html := `<h2>Services</h2>`
			for _, service := range services {
				html += fmt.Sprintf(`<a class="micro-link" href="/%s">%s</a>`, url.QueryEscape(service.Name), service.Name)
			}
			render(w, html)
			return
		}
		if len(parts) < 2 || parts[1] == "" {
			// If webDir exists, redirect "/" to "/web" for the index page
			if _, err := os.Stat(webDir); err == nil {
				http.Redirect(w, r, "/web", http.StatusFound)
				return
			}
			// Otherwise, show the default index
			html := `<h2>Web</h2>
            <p><a href="/services" class="micro-link">Services</a></p>`
			render(w, html)
			return
		}
		service := parts[1]
		// Try to decode URL-encoded service names (for /services links)
		service, _ = url.QueryUnescape(service)
		s, err := registry.GetService(service)
		if err != nil || len(s) == 0 {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf("Service not found: %s", service)))
			return
		}
		if len(parts) < 3 || parts[2] == "" {
			// List endpoints for the service
			var endpoints string
			for _, ep := range s[0].Endpoints {
				p := strings.Split(ep.Name, ".")
				uri := fmt.Sprintf("/%s/%s/%s", service, p[0], p[1])
				endpoints += fmt.Sprintf(`<a class="micro-link" href="%s">%s</a>`, uri, ep.Name)
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
		if ep == nil {
			w.WriteHeader(404)
			w.Write([]byte("Endpoint not found"))
			return
		}
		if r.Method == "GET" {
			var inputs string
			inputs += fmt.Sprintf("<h3>%s</h3>", ep.Name)
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
	})
	go http.ListenAndServe(addr, nil)
}

func Run(c *cli.Context) error {
	dir := c.Args().Get(0)
	if len(dir) == 0 {
		dir = "."
	}
	addr := c.String("address")
	if addr == "" {
		addr = ":8080"
	}
	serveMicroWeb(dir, addr)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home dir: %w", err)
	}
	logsDir := filepath.Join(homeDir, "micro", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs dir: %w", err)
	}
	runDir := filepath.Join(homeDir, "micro", "run")
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create run dir: %w", err)
	}
	binDir := filepath.Join(homeDir, "micro", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin dir: %w", err)
	}

	// Always run all services (find all main.go)
	var mainFiles []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "main.go" {
			mainFiles = append(mainFiles, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path: %w", err)
	}
	if len(mainFiles) == 0 {
		return fmt.Errorf("no main.go files found in %s", dir)
	}
	var procs []*exec.Cmd
	var pidFiles []string
	for i, mainFile := range mainFiles {
		serviceDir := filepath.Dir(mainFile)
		var serviceName string
		absDir, _ := filepath.Abs(dir)
		absServiceDir, _ := filepath.Abs(serviceDir)
		if absDir == absServiceDir {
			// If main.go is in the root dir being run, use the current working dir name
			cwd, _ := os.Getwd()
			serviceName = filepath.Base(cwd)
		} else {
			serviceName = filepath.Base(serviceDir)
		}
		logFilePath := filepath.Join(logsDir, serviceName+".log")
		binPath := filepath.Join(binDir, serviceName)
		pidFilePath := filepath.Join(runDir, serviceName+".pid")

		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open log file for %s: %v\n", serviceName, err)
			continue
		}
		buildCmd := exec.Command("go", "build", "-o", binPath, ".")
		buildCmd.Dir = serviceDir
		buildOut, buildErr := buildCmd.CombinedOutput()
		if buildErr != nil {
			logFile.WriteString(string(buildOut))
			logFile.Close()
			fmt.Fprintf(os.Stderr, "failed to build %s: %v\n", serviceName, buildErr)
			continue
		}
		cmd := exec.Command(binPath)
		cmd.Dir = serviceDir
		pr, pw := io.Pipe()
		cmd.Stdout = pw
		cmd.Stderr = pw
		color := colorFor(i)
		go func(name string, color string, pr *io.PipeReader) {
			scanner := bufio.NewScanner(pr)
			for scanner.Scan() {
				fmt.Printf("%s[%s]\033[0m %s\n", color, name, scanner.Text())
			}
		}(serviceName, color, pr)
		if err := cmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start service %s: %v\n", serviceName, err)
			pw.Close()
			continue
		}
		procs = append(procs, cmd)
		pidFiles = append(pidFiles, pidFilePath)
		os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d\n%s\n", cmd.Process.Pid, serviceDir)), 0644)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		for _, proc := range procs {
			if proc.Process != nil {
				_ = proc.Process.Kill()
			}
		}
		for _, pf := range pidFiles {
			_ = os.Remove(pf)
		}
		os.Exit(1)
	}()
	for _, proc := range procs {
		_ = proc.Wait()
	}
	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "run",
		Usage:  "Run all services in a directory",
		Action: Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Address to bind the micro web UI (default :8080)",
				Value:   ":8080",
			},
		},
	})
}
