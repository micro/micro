package server

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/fs"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"syscall"
	"path/filepath"

	goMicroClient "go-micro.dev/v5/client"
	goMicroBytes "go-micro.dev/v5/codec/bytes"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v5/auth"
	jwtAuth "go-micro.dev/v5/auth/jwt"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/store"
)

// HTML is the embedded filesystem for templates and static files, set by main.go
var HTML fs.FS

var (
	apiCache struct {
		sync.Mutex
		data map[string]any
		time time.Time
	}
)

type templates struct {
	api      *template.Template
	service  *template.Template
	form     *template.Template
	home     *template.Template
	logs     *template.Template
	log      *template.Template
	status   *template.Template
	authTokens *template.Template
	authLogin *template.Template
}

func parseTemplates() *templates {
	return &templates{
		api:      template.Must(template.ParseFS(HTML, "html/base.html", "html/api.html")),
		service:  template.Must(template.ParseFS(HTML, "html/base.html", "html/service.html")),
		form:     template.Must(template.ParseFS(HTML, "html/base.html", "html/form.html")),
		home:     template.Must(template.ParseFS(HTML, "html/base.html", "html/home.html")),
		logs:     template.Must(template.ParseFS(HTML, "html/base.html", "html/logs.html")),
		log:      template.Must(template.ParseFS(HTML, "html/base.html", "html/log.html")),
		status:   template.Must(template.ParseFS(HTML, "html/base.html", "html/status.html")),
		authTokens: template.Must(template.ParseFS(HTML, "html/base.html", "html/auth_tokens.html")),
		authLogin: template.Must(template.ParseFS(HTML, "html/base.html", "html/login.html")),
	}
}

// Helper to render templates
func render(w http.ResponseWriter, tmpl *template.Template, data any) error {
	return tmpl.Execute(w, data)
}

// Helper to extract user info from JWT cookie
func getUser(r *http.Request) string {
	cookie, err := r.Cookie("micro_token")
	if err != nil || cookie.Value == "" {
		return ""
	}
	// Parse JWT claims (just decode, don't verify)
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := decodeSegment(parts[1])
	if err != nil {
		return ""
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	if sub, ok := claims["sub"].(string); ok {
		return sub
	}
	if id, ok := claims["id"].(string); ok {
		return id
	}
	return ""
}

// Helper to decode JWT base64url segment
func decodeSegment(seg string) ([]byte, error) {
	// JWT uses base64url, no padding
	missing := len(seg) % 4
	if missing != 0 {
		seg += strings.Repeat("=", 4-missing)
	}
	return decodeBase64Url(seg)
}

func decodeBase64Url(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Fix authRequired to check JWT expiry from claims, not acc.Expiry
func authRequired(authSrv auth.Auth) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("micro_token")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			// Parse JWT expiry
			parts := strings.Split(cookie.Value, ".")
			if len(parts) != 3 {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			payload, err := decodeSegment(parts[1])
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			var claims map[string]any
			if err := json.Unmarshal(payload, &claims); err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}
			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					http.Redirect(w, r, "/auth/login", http.StatusFound)
					return
				}
			}
			// Optionally: verify token with authSrv (optional, since we just check expiry here)
			next(w, r)
		}
	}
}

func wrapAuth(authRequired func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth/login") || strings.HasPrefix(r.URL.Path, "/auth/logout") {
				h(w, r)
				return
			}
			authRequired(h)(w, r)
		}
	}
}

func getDashboardData() (serviceCount, runningCount, stoppedCount int, statusDot string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	pidDir := homeDir + "/micro/run"
	dirEntries, err := os.ReadDir(pidDir)
	if err != nil {
		return
	}
	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".pid") || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		pidFile := pidDir + "/" + entry.Name()
		pidBytes, err := os.ReadFile(pidFile)
		if err != nil {
			continue
		}
		lines := strings.Split(string(pidBytes), "\n")
		pid := "-"
		if len(lines) > 0 && len(lines[0]) > 0 {
			pid = lines[0]
		}
		serviceCount++
		if pid != "-" {
			if _, err := os.FindProcess(parsePid(pid)); err == nil {
				if processRunning(pid) {
					runningCount++
				} else {
					stoppedCount++
				}
			} else {
				stoppedCount++
			}
		} else {
			stoppedCount++
		}
	}
	if serviceCount > 0 && runningCount == serviceCount {
		statusDot = "green"
	} else if serviceCount > 0 && runningCount > 0 {
		statusDot = "yellow"
	} else {
		statusDot = "red"
	}
	return
}

func registerHandlers(tmpls *templates, authSrv auth.Auth, storeInst store.Store) {
	authMw := authRequired(authSrv)
	wrap := wrapAuth(authMw)

	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.FS(HTML))))

	http.HandleFunc("/", wrap(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/auth/") {
			// Let the dedicated /auth/* handlers process this
			return
		}
		if path == "/" {
			serviceCount, runningCount, stoppedCount, statusDot := getDashboardData()
			_ = tmpls.home.Execute(w, map[string]any{
				"Title": "Micro Dashboard",
				"WebLink": "/",
				"ServiceCount": serviceCount,
				"RunningCount": runningCount,
				"StoppedCount": stoppedCount,
				"StatusDot": statusDot,
				"User": getUser(r),
			})
			return
		}
		if path == "/api" || path == "/api/" {
			apiCache.Lock()
			useCache := false
			if apiCache.data != nil && time.Since(apiCache.time) < 30*time.Second {
				useCache = true
			}
			var apiData map[string]any
			var sidebarEndpoints []map[string]string
			if useCache {
				apiData = apiCache.data
				if v, ok := apiData["SidebarEndpoints"]; ok {
					sidebarEndpoints, _ = v.([]map[string]string)
				}
			} else {
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
						// Build params and response HTML from endpoint values
						var params, response string
						if ep.Request != nil && len(ep.Request.Values) > 0 {
							params += "<ul class=no-bullets>"
							for _, v := range ep.Request.Values {
								params += fmt.Sprintf("<li><b>%s</b> <span style='color:#888;'>%s</span></li>", v.Name, v.Type)
							}
							params += "</ul>"
						} else {
							params = "<i style='color:#888;'>No parameters</i>"
						}
						if ep.Response != nil && len(ep.Response.Values) > 0 {
							response += "<ul class=no-bullets>"
							for _, v := range ep.Response.Values {
								response += fmt.Sprintf("<li><b>%s</b> <span style='color:#888;'>%s</span></li>", v.Name, v.Type)
							}
							response += "</ul>"
						} else {
							response = "<i style='color:#888;'>No response fields</i>"
						}
						endpoints = append(endpoints, map[string]any{
							"Name": ep.Name,
							"Path": apiPath,
							"Params": params,
							"Response": response,
						})
					}
					anchor := strings.ReplaceAll(s.Name, ".", "-")
					apiServices = append(apiServices, map[string]any{
						"Name": s.Name,
						"Anchor": anchor,
						"Endpoints": endpoints,
					})
					sidebarEndpoints = append(sidebarEndpoints, map[string]string{"Name": s.Name, "Anchor": anchor})
				}
				// Sort sidebarEndpoints by Name
				sort.Slice(sidebarEndpoints, func(i, j int) bool {
					return sidebarEndpoints[i]["Name"] < sidebarEndpoints[j]["Name"]
				})
				apiData = map[string]any{"Title": "API", "WebLink": "/", "Services": apiServices, "SidebarEndpoints": sidebarEndpoints, "SidebarEndpointsEnabled": true}
				apiCache.data = apiData
				apiCache.time = time.Now()
			}
			apiCache.Unlock()
			_ = render(w, tmpls.api, apiData)
			return
		}
		if path == "/services" {
			services, _ := registry.ListServices()
			var serviceNames []string
			for _, service := range services {
				serviceNames = append(serviceNames, service.Name)
			}
			sort.Strings(serviceNames)
			_ = render(w, tmpls.service, map[string]any{"Title": "Services", "WebLink": "/", "Services": serviceNames})
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
			_ = render(w, tmpls.logs, map[string]any{"Title": "Logs", "WebLink": "/", "Services": serviceNames})
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
			_ = render(w, tmpls.log, map[string]any{"Title": "Logs for " + service, "WebLink": "/logs", "Service": service, "Log": logText})
			return
		}
		if path == "/status" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not get home directory"))
				return
			}
			pidDir := homeDir + "/micro/run"
			dirEntries, err := os.ReadDir(pidDir)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Could not list pid directory: " + err.Error()))
				return
			}
			statuses := []map[string]string{}
			for _, entry := range dirEntries {
				if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".pid") || strings.HasPrefix(entry.Name(), ".") {
					continue
				}
				pidFile := pidDir + "/" + entry.Name()
				pidBytes, err := os.ReadFile(pidFile)
				if err != nil {
					statuses = append(statuses, map[string]string{
						"Service": entry.Name(),
						"Dir": "-",
						"Status": "unknown",
						"PID": "-",
						"Uptime": "-",
						"ID": strings.TrimSuffix(entry.Name(), ".pid"),
					})
					continue
				}
				lines := strings.Split(string(pidBytes), "\n")
				pid := "-"
				dir := "-"
				service := "-"
				start := "-"
				if len(lines) > 0 && len(lines[0]) > 0 {
					pid = lines[0]
				}
				if len(lines) > 1 && len(lines[1]) > 0 {
					dir = lines[1]
				}
				if len(lines) > 2 && len(lines[2]) > 0 {
					service = lines[2]
				}
				if len(lines) > 3 && len(lines[3]) > 0 {
					start = lines[3]
				}
				status := "stopped"
				if pid != "-" {
					if _, err := os.FindProcess(parsePid(pid)); err == nil {
						if processRunning(pid) {
							status = "running"
						}
					}
				}
				uptime := "-"
				if start != "-" {
					if t, err := parseStartTime(start); err == nil {
						uptime = time.Since(t).Truncate(time.Second).String()
					}
				}
				statuses = append(statuses, map[string]string{
					"Service": service,
					"Dir": dir,
					"Status": status,
					"PID": pid,
					"Uptime": uptime,
					"ID": strings.TrimSuffix(entry.Name(), ".pid"),
				})
			}
			_ = render(w, tmpls.status, map[string]any{"Title": "Service Status", "WebLink": "/", "Statuses": statuses})
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
				_ = render(w, tmpls.service, map[string]any{
					"Title": "Service: " + service,
					"WebLink": "/",
					"ServiceName": service,
					"Endpoints": endpoints,
					"Description": string(b),
					"User": getUser(r),
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
					_ = render(w, tmpls.form, map[string]any{
						"Title":       "Service: " + service,
						"WebLink":     "/",
						"ServiceName": service,
						"EndpointName": ep.Name,
						"Inputs":      inputs,
						"Action":      service + "/" + endpoint,
						"User": getUser(r),
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
	}))

	http.HandleFunc("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "micro_token", Value: "", Path: "/", Expires: time.Now().Add(-1 * time.Hour), HttpOnly: true})
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
	})

	http.HandleFunc("/auth/tokens", authMw(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id := r.FormValue("id")
			typeStr := r.FormValue("type")
			scopesStr := r.FormValue("scopes")
			accType := "user"
			if typeStr == "admin" {
				accType = "admin"
			} else if typeStr == "service" {
				accType = "service"
			}
			scopes := []string{"*"}
			if scopesStr != "" {
				scopes = strings.Split(scopesStr, ",")
				for i := range scopes {
					scopes[i] = strings.TrimSpace(scopes[i])
				}
			}
			acc := &auth.Account{
				ID: id,
				Type: accType,
				Scopes: scopes,
				Metadata: map[string]string{"created": time.Now().Format(time.RFC3339)},
			}
			// Service tokens do not require a password, generate a JWT directly
			tok, _ := authSrv.Generate(acc.ID, auth.WithType(accType), auth.WithScopes(acc.Scopes...))
			acc.Metadata["token"] = tok.Secret
			b, _ := json.Marshal(acc)
			storeInst.Write(&store.Record{Key: "auth/" + id, Value: b})
			http.Redirect(w, r, "/auth/tokens", http.StatusSeeOther)
			return
		}
		recs, _ := storeInst.Read("auth/", store.ReadPrefix())
		var tokens []map[string]any
		for _, rec := range recs {
			var acc auth.Account
			if err := json.Unmarshal(rec.Value, &acc); err == nil {
				tok := ""
				if t, ok := acc.Metadata["token"]; ok {
					tok = t
				}
				tokens = append(tokens, map[string]any{
					"ID": acc.ID,
					"Type": acc.Type,
					"Scopes": acc.Scopes,
					"Metadata": acc.Metadata,
					"Token": tok,
				})
			}
		}
		_ = tmpls.authTokens.Execute(w, map[string]any{"Title": "Auth Tokens", "Tokens": tokens, "User": getUser(r)})
	}))

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Render login form
			loginTmpl, err := template.ParseFS(HTML, "html/base.html", "html/auth_login.html")
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("Template error: " + err.Error()))
				return
			}
			_ = loginTmpl.Execute(w, map[string]any{"Title": "Login", "Error": "", "User": getUser(r)})
			return
		}
		if r.Method == "POST" {
			id := r.FormValue("id")
			pass := r.FormValue("password")
			recKey := "auth/" + id
			recs, _ := storeInst.Read(recKey)
			if len(recs) == 0 {
				loginTmpl, _ := template.ParseFS(HTML, "html/base.html", "html/auth_login.html")
				_ = loginTmpl.Execute(w, map[string]any{"Title": "Login", "Error": "Invalid credentials", "User": ""})
				return
			}
			var acc auth.Account
			if err := json.Unmarshal(recs[0].Value, &acc); err != nil || acc.Secret != pass {
				loginTmpl, _ := template.ParseFS(HTML, "html/base.html", "html/auth_login.html")
				_ = loginTmpl.Execute(w, map[string]any{"Title": "Login", "Error": "Invalid credentials", "User": ""})
				return
			}
			// Generate JWT token
			tok, err := authSrv.Generate(acc.ID, auth.WithType(acc.Type), auth.WithScopes(acc.Scopes...))
			if err != nil {
				loginTmpl, _ := template.ParseFS(HTML, "html/base.html", "html/auth_login.html")
				_ = loginTmpl.Execute(w, map[string]any{"Title": "Login", "Error": "Token error", "User": ""})
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "micro_token",
				Value:    tok.Secret,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
	})
}

func Run(c *cli.Context) error {
	if err := initAuth(); err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	homeDir, _ := os.UserHomeDir()
	keyDir := filepath.Join(homeDir, "micro", "keys")
	privPath := filepath.Join(keyDir, "private.pem")
	pubPath := filepath.Join(keyDir, "public.pem")
	privPem, _ := os.ReadFile(privPath)
	pubPem, _ := os.ReadFile(pubPath)
	authSrv := jwtAuth.NewAuth()
	authSrv.Init(
		auth.PublicKey(string(pubPem)),
		auth.PrivateKey(string(privPem)),
	)
	storeInst := store.DefaultStore

	tmpls := parseTemplates()
	registerHandlers(tmpls, authSrv, storeInst)

	addr := c.String("address")
	if addr == "" {
		addr = ":8080"
	}

	go func() {
		log.Printf("[micro-server] Web/API listening on %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Web/API server error: %v", err)
		}
	}()

	return nil
}

// --- PID FILES ---

func parsePid(pidStr string) int {
	pid, _ := strconv.Atoi(pidStr)
	return pid
}

func processRunning(pid string) bool {
	proc, err := os.FindProcess(parsePid(pid))
	if err != nil {
		return false
	}
	// On unix, sending syscall.Signal(0) checks if process exists
	return proc.Signal(syscall.Signal(0)) == nil
}

func generateKeyPair(bits int) (*rsa.PrivateKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func exportPrivateKeyAsPEM(priv *rsa.PrivateKey) ([]byte, error) {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(priv)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	var buf bytes.Buffer
	err := pem.Encode(&buf, block)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func exportPublicKeyAsPEM(pub *rsa.PublicKey) ([]byte, error) {
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pub)
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	var buf bytes.Buffer
	err := pem.Encode(&buf, block)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func importPrivateKeyFromPEM(privKeyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func importPublicKeyFromPEM(pubKeyPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func initAuth() error {
	// --- AUTH SETUP ---
	homeDir, _ := os.UserHomeDir()
	keyDir := filepath.Join(homeDir, "micro", "keys")
	privPath := filepath.Join(keyDir, "private.pem")
	pubPath := filepath.Join(keyDir, "public.pem")
	os.MkdirAll(keyDir, 0700)
	// Generate keypair if not exist
	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		priv, _ := rsa.GenerateKey(rand.Reader, 2048)
		privBytes := x509.MarshalPKCS1PrivateKey(priv)
		privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
		os.WriteFile(privPath, privPem, 0600)
		pubBytes := x509.MarshalPKCS1PublicKey(&priv.PublicKey)
		pubPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})
		os.WriteFile(pubPath, pubPem, 0644)
	}
	_, _ = os.ReadFile(privPath)
	_, _ = os.ReadFile(pubPath)
	storeInst := store.DefaultStore
	// --- Ensure default admin account exists ---
	adminID := "admin"
	adminPass := "micro"
	adminKey := "auth/" + adminID
	if recs, _ := storeInst.Read(adminKey); len(recs) == 0 {
		acc := &auth.Account{
			ID: adminID,
			Type: "admin",
			Scopes: []string{"*"},
			Metadata: map[string]string{"created": "true"},
		}
		acc.Secret = adminPass
		b, _ := json.Marshal(acc)
		storeInst.Write(&store.Record{Key: adminKey, Value: b})
	}
	return nil
}

// parseStartTime parses a string as RFC3339 time
func parseStartTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "server",
		Usage:  "Run the micro server",
		Action: Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Address to listen on",
				EnvVars: []string{"MICRO_SERVER_ADDRESS"},
				Value:   ":8080",
			},
		},
	})
}
