package web

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"regexp"
	"sort"
	"strings"
	"time"

	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
	"github.com/serenize/snaker"
)

var (
	re = regexp.MustCompile("^[a-zA-Z0-9]+$")
	// Default address to bind to
	Address = ":8082"
	// The namespace to serve
	// Example:
	// Namespace + /[Service]/foo/bar
	// Host: Namespace.Service Endpoint: /foo/bar
	Namespace = "go.micro.web"
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	BasePathHeader = "X-Micro-Web-Base-Path"

	CORS = map[string]bool{"*": true}
)

type srv struct {
	*mux.Router
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); CORS[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && CORS["*"] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	s.Router.ServeHTTP(w, r)
}

func (s *srv) proxy() http.Handler {
	sel := selector.NewSelector(
		selector.Registry((*cmd.DefaultOptions().Registry)),
	)

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
		next, err := sel.Select(Namespace + "." + parts[1])
		if err != nil {
			kill()
			return
		}

		s, err := next()
		if err != nil {
			kill()
			return
		}

		r.Header.Set(BasePathHeader, "/"+parts[1])
		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Path = "/" + strings.Join(parts[2:], "/")
		r.URL.Scheme = "http"
	}

	return &proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}

func format(v *registry.Value) string {
	if v == nil || len(v.Values) == 0 {
		return "{}"
	}
	var f []string
	for _, k := range v.Values {
		f = append(f, formatEndpoint(k, 0))
	}
	return fmt.Sprintf("{\n%s}", strings.Join(f, ""))
}

func formatEndpoint(v *registry.Value, r int) string {
	// default format is tabbed plus the value plus new line
	fparts := []string{"", "%s %s", "\n"}
	for i := 0; i < r+1; i++ {
		fparts[0] += "\t"
	}
	// its just a primitive of sorts so return
	if len(v.Values) == 0 {
		return fmt.Sprintf(strings.Join(fparts, ""), snaker.CamelToSnake(v.Name), v.Type)
	}

	// this thing has more things, it's complex
	fparts[1] += " {"

	vals := []interface{}{snaker.CamelToSnake(v.Name), v.Type}

	for _, val := range v.Values {
		fparts = append(fparts, "%s")
		vals = append(vals, formatEndpoint(val, r+1))
	}

	// at the end
	l := len(fparts) - 1
	for i := 0; i < r+1; i++ {
		fparts[l] += "\t"
	}
	fparts = append(fparts, "}\n")

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	var webServices []string
	for _, s := range services {
		if strings.Index(s.Name, Namespace) == 0 && len(strings.TrimPrefix(s.Name, Namespace)) > 0 {
			webServices = append(webServices, strings.Replace(s.Name, Namespace+".", "", 1))
		}
	}

	sort.Strings(webServices)

	type templateData struct {
		HasWebServices bool
		WebServices    []string
	}

	data := templateData{len(webServices) > 0, webServices}
	render(w, r, indexTemplate, data)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	svc := r.Form.Get("service")

	if len(svc) > 0 {
		s, err := (*cmd.DefaultOptions().Registry).GetService(svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(s) == 0 {
			http.Error(w, "Not found", 404)
			return
		}

		if r.Header.Get("Content-Type") == "application/json" {
			b, err := json.Marshal(map[string]interface{}{
				"services": s,
			})
			if err != nil {
				http.Error(w, "Error occurred:"+err.Error(), 500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}

		render(w, r, serviceTemplate, s)
		return
	}

	services, err := (*cmd.DefaultOptions().Registry).ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	sort.Sort(sortedServices{services})

	if r.Header.Get("Content-Type") == "application/json" {
		b, err := json.Marshal(map[string]interface{}{
			"services": services,
		})
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}

	render(w, r, registryTemplate, services)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, queryTemplate, nil)
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
	}).Parse(layoutTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	t, err = t.Parse(tmpl)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func run(ctx *cli.Context) {
	var h http.Handler
	r := mux.NewRouter()
	s := &srv{r}
	h = s

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		s.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	s.HandleFunc("/registry", registryHandler)
	s.HandleFunc("/rpc", handler.RPC)
	s.HandleFunc("/query", queryHandler)
	s.HandleFunc("/favicon.ico", faviconHandler)
	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(s.proxy())
	s.HandleFunc("/", indexHandler)

	var opts []server.Option

	if ctx.GlobalBool("enable_tls") {
		cert := ctx.GlobalString("tls_cert_file")
		key := ctx.GlobalString("tls_key_file")

		if len(cert) > 0 && len(key) > 0 {
			certs, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			config := &tls.Config{
				Certificates: []tls.Certificate{certs},
			}
			opts = append(opts, server.EnableTLS(true))
			opts = append(opts, server.TLSConfig(config))
		} else {
			fmt.Println("Enable TLS specified without certificate and key files")
			return
		}
	}

	srv := server.NewServer(Address)
	srv.Init(opts...)
	srv.Handle("/", h)

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.web"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Stop(); err != nil {
		log.Fatal(err)
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "web",
			Usage: "Run the micro web app",
			Action: func(c *cli.Context) {
				run(c)
			},
		},
	}
}
