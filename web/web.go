// Package web is a web dashboard
package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/api/server"
	"github.com/micro/go-micro/v2/api/server/acme"
	"github.com/micro/go-micro/v2/api/server/acme/autocert"
	"github.com/micro/go-micro/v2/api/server/acme/certmagic"
	"github.com/micro/go-micro/v2/api/server/cors"
	httpapi "github.com/micro/go-micro/v2/api/server/http"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/config/cmd"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/cache"
	cfstore "github.com/micro/go-micro/v2/store/cloudflare"
	"github.com/micro/go-micro/v2/sync/lock/memory"
	"github.com/micro/micro/v2/internal/handler"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/internal/stats"
	"github.com/micro/micro/v2/plugin"
	"github.com/serenize/snaker"
)

var (
	re = regexp.MustCompile("^[a-zA-Z0-9]+([a-zA-Z0-9-]*[a-zA-Z0-9]*)?$")
	// Default server name
	Name = "go.micro.web"
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
	BasePathHeader        = "X-Micro-Web-Base-Path"
	statsURL              string
	loginURL              string
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA

	// A placeholder icon
	DefaultIcon = "https://micro.mu/circle.png"
)

type srv struct {
	*mux.Router
	// registry we use
	registry registry.Registry
}

type reg struct {
	registry.Registry

	sync.Mutex
	lastPull time.Time
	services []*registry.Service
}

func (r *reg) watch() {
Loop:
	for {
		// get a watcher
		w, err := r.Registry.Watch()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		// loop results
		for {
			_, err := w.Next()
			if err != nil {
				w.Stop()
				time.Sleep(time.Second)
				goto Loop
			}

			// next pull will be from the registry
			r.Lock()
			r.lastPull = time.Time{}
			r.Unlock()
		}
	}
}

func (r *reg) ListServices() ([]*registry.Service, error) {
	r.Lock()
	defer r.Unlock()

	if r.lastPull.IsZero() || time.Since(r.lastPull) > time.Minute {
		// pull the services
		s, err := r.Registry.ListServices()
		if err != nil {
			// return stale entries if they exist
			if len(r.services) > 0 {
				return r.services, nil
			}
			// otherwise return an error
			return nil, err
		}
		// collapse the list
		serviceMap := make(map[string]*registry.Service)
		for _, service := range s {
			serviceMap[service.Name] = service
		}
		var services []*registry.Service
		for _, service := range serviceMap {
			services = append(services, service)
		}
		// cache it
		r.services = services
		r.lastPull = time.Now()
		return s, nil
	}

	// return the cached list
	return r.services, nil
}

func (s *srv) proxy() http.Handler {
	// our internal resolver
	res := &resolver{
		Namespace: Namespace,
		Selector: selector.NewSelector(
			selector.Registry(s.registry),
		),
	}

	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		// TODO: better error handling
		if err := res.Resolve(r); err != nil {
			fmt.Printf("Failed to resolve url: %v: %v\n", r.URL, err)
			kill()
			return
		}
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

func (s *srv) indexHandler(w http.ResponseWriter, r *http.Request) {
	cors.SetHeaders(w, r)

	if r.Method == "OPTIONS" {
		return
	}

	services, err := s.registry.ListServices()
	if err != nil {
		log.Errorf("Error listing services: %v", err)
	}

	type webService struct {
		Name string
		Icon string
	}

	var webServices []webService
	for _, srv := range services {
		if len(srv.Nodes) == 0 {
			srvs, _ := s.registry.GetService(srv.Name)
			if len(srvs) == 0 {
				continue
			}
			srv = srvs[0]
		}

		if len(srv.Nodes) == 0 {
			continue
		}

		if strings.Index(srv.Name, Namespace) == 0 && len(strings.TrimPrefix(srv.Name, Namespace)) > 0 {
			icon := DefaultIcon

			if ico := srv.Nodes[0].Metadata["icon"]; len(ico) > 0 {
				icon = ico
			}

			name := strings.Replace(srv.Name, Namespace+".", "", 1)

			webServices = append(webServices, webService{
				Name: name,
				Icon: icon,
			})
		}
	}

	sort.Slice(webServices, func(i, j int) bool { return webServices[i].Name < webServices[j].Name })

	type templateData struct {
		HasWebServices bool
		WebServices    []webService
	}

	data := templateData{len(webServices) > 0, webServices}
	render(w, r, indexTemplate, data)
}

func (s *srv) registryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	svc := vars["name"]

	if len(svc) > 0 {
		s, err := s.registry.GetService(svc)
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

	services, err := s.registry.ListServices()
	if err != nil {
		log.Errorf("Error listing services: %v", err)
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

func (s *srv) callHandler(w http.ResponseWriter, r *http.Request) {
	services, err := s.registry.ListServices()
	if err != nil {
		log.Errorf("Error listing services: %v", err)
	}

	sort.Sort(sortedServices{services})

	serviceMap := make(map[string][]*registry.Endpoint)
	for _, service := range services {
		if len(service.Endpoints) > 0 {
			serviceMap[service.Name] = service.Endpoints
			continue
		}
		// lookup the endpoints otherwise
		s, err := s.registry.GetService(service.Name)
		if err != nil {
			continue
		}
		if len(s) == 0 {
			continue
		}
		serviceMap[service.Name] = s[0].Endpoints
	}

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

	render(w, r, callTemplate, serviceMap)
}

func render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
		"Title":  strings.Title,
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

	if err := t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"LoginURL": loginURL,
		"StatsURL": statsURL,
		"Results":  data,
	}); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "web"}))

	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// use the caching registry
	cache := cache.New((*cmd.DefaultOptions().Registry))
	reg := &reg{Registry: cache}

	// start the watcher
	go reg.watch()

	var h http.Handler
	s := &srv{
		Router:   mux.NewRouter(),
		registry: reg,
	}
	h = s

	if ctx.Bool("enable_stats") {
		statsURL = "/stats"
		st := stats.New()
		s.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	s.HandleFunc("/client", s.callHandler)
	s.HandleFunc("/services", s.registryHandler)
	s.HandleFunc("/service/{name}", s.registryHandler)
	s.HandleFunc("/rpc", handler.RPC)
	s.HandleFunc("/favicon.ico", faviconHandler)
	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(s.proxy())
	s.HandleFunc("/", s.indexHandler)

	var opts []server.Option

	if len(ctx.String("acme_provider")) > 0 {
		ACMEProvider = ctx.String("acme_provider")
	}
	if ctx.Bool("enable_acme") {
		hosts := helper.ACMEHosts(ctx)
		opts = append(opts, server.EnableACME(true))
		opts = append(opts, server.ACMEHosts(hosts...))
		switch ACMEProvider {
		case "autocert":
			opts = append(opts, server.ACMEProvider(autocert.NewProvider()))
		case "certmagic":
			if ACMEChallengeProvider != "cloudflare" {
				log.Fatal("The only implemented DNS challenge provider is cloudflare")
			}
			apiToken, accountID := os.Getenv("CF_API_TOKEN"), os.Getenv("CF_ACCOUNT_ID")
			kvID := os.Getenv("KV_NAMESPACE_ID")
			if len(apiToken) == 0 || len(accountID) == 0 {
				log.Fatal("env variables CF_API_TOKEN and CF_ACCOUNT_ID must be set")
			}
			if len(kvID) == 0 {
				log.Fatal("env var KV_NAMESPACE_ID must be set to your cloudflare workers KV namespace ID")
			}

			cloudflareStore := cfstore.NewStore(
				cfstore.Token(apiToken),
				cfstore.Account(accountID),
				cfstore.Namespace(kvID),
				cfstore.CacheTTL(time.Minute),
			)
			storage := certmagic.NewStorage(
				memory.NewLock(),
				cloudflareStore,
			)
			config := cloudflare.NewDefaultConfig()
			config.AuthToken = apiToken
			config.ZoneToken = apiToken
			challengeProvider, err := cloudflare.NewDNSProviderConfig(config)
			if err != nil {
				log.Fatal(err.Error())
			}

			opts = append(opts,
				server.ACMEProvider(
					certmagic.NewProvider(
						acme.AcceptToS(true),
						acme.CA(ACMECA),
						acme.Cache(storage),
						acme.ChallengeProvider(challengeProvider),
						acme.OnDemand(false),
					),
				),
			)
		default:
			log.Fatalf("%s is not a valid ACME provider\n", ACMEProvider)
		}
	} else if ctx.Bool("enable_tls") {
		config, err := helper.TLSConfig(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		opts = append(opts, server.EnableTLS(true))
		opts = append(opts, server.TLSConfig(config))
	}

	// reverse wrap handler
	plugins := append(Plugins(), plugin.Plugins()...)
	for i := len(plugins); i > 0; i-- {
		h = plugins[i-1].Handler()(h)
	}

	srv := httpapi.NewServer(Address)
	srv.Init(opts...)
	srv.Handle("/", h)

	// service opts
	srvOpts = append(srvOpts, micro.Name(Name))
	if i := time.Duration(ctx.Int("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.Int("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	// Initialise Server
	service := micro.NewService(srvOpts...)

	// Setup auth redirect
	if len(ctx.String("auth_login_url")) > 0 {
		loginURL = ctx.String("auth_login_url")
		service.Options().Auth.Init(auth.LoginURL(loginURL))
	}

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

func Commands(options ...micro.Option) []*cli.Command {
	command := &cli.Command{
		Name:  "web",
		Usage: "Run the web dashboard",
		Action: func(c *cli.Context) error {
			run(c, options...)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Usage:   "Set the web UI address e.g 0.0.0.0:8082",
				EnvVars: []string{"MICRO_WEB_ADDRESS"},
			},
			&cli.StringFlag{
				Name:    "namespace",
				Usage:   "Set the namespace used by the Web proxy e.g. com.example.web",
				EnvVars: []string{"MICRO_WEB_NAMESPACE"},
			},
			&cli.StringFlag{
				Name:    "auth_login_url",
				EnvVars: []string{"MICRO_AUTH_LOGIN_URL"},
				Usage:   "The relative URL where a user can login",
			},
		},
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []*cli.Command{command}
}
