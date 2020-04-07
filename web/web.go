// Package web is a web dashboard
package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"
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
	res "github.com/micro/go-micro/v2/api/resolver"
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
	nsResolver "github.com/micro/micro/v2/internal/namespace/resolver"
	"github.com/micro/micro/v2/internal/stats"
	"github.com/micro/micro/v2/plugin"
	"github.com/serenize/snaker"
	"golang.org/x/net/publicsuffix"
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
	Namespace     = "go.micro"
	Type          = "web"
	FullNamespace = "go.micro.web"
	// Resolver used to resolve services
	Resolver = "path"
	// Base path sent to web service.
	// This is stripped from the request path
	// Allows the web service to define absolute paths
	BasePathHeader        = "X-Micro-Web-Base-Path"
	statsURL              string
	loginURL              string
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA

	// Host name the web dashboard is served on
	Host, _ = os.Hostname()
)

type srv struct {
	*mux.Router
	// registry we use
	registry registry.Registry
	// the resolver
	resolver *resolver
	// the proxy server
	prx *proxy
	// auth service
	auth auth.Auth
}

type reg struct {
	registry.Registry

	sync.RWMutex
	lastPull time.Time
	services []*registry.Service
}

func (r *reg) watch() {
	// update once
	r.update()

	// periodically update the service cache
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()

		for _ = range t.C {
			r.update()
		}
	}()

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

func (r *reg) update() {
	// pull the services
	s, err := r.Registry.ListServices()
	if err != nil {
		return
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

	r.Lock()
	defer r.Unlock()

	// cache it
	r.services = services
	r.lastPull = time.Now()
}

func (r *reg) ListServices() ([]*registry.Service, error) {
	r.RLock()
	defer r.RUnlock()

	// return the cached list
	return r.services, nil
}

// ServeHTTP serves the web dashboard and proxies where appropriate
func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Host) == 0 {
		r.URL.Host = r.Host
	}

	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}

	// no host means dashboard
	host := r.URL.Hostname()
	if len(host) == 0 {
		h, _, err := net.SplitHostPort(r.Host)
		if err != nil && strings.Contains(err.Error(), "missing port in address") {
			host = r.Host
		} else if err == nil {
			host = h
		}
	}

	// check again
	if len(host) == 0 {
		s.Router.ServeHTTP(w, r)
		return
	}

	// check based on host set
	if len(Host) > 0 && Host == host {
		s.Router.ServeHTTP(w, r)
		return
	}

	// an ip instead of hostname means dashboard
	ip := net.ParseIP(host)
	if ip != nil {
		s.Router.ServeHTTP(w, r)
		return
	}

	// namespace matching host means dashboard
	parts := strings.Split(host, ".")
	reverse(parts)
	namespace := strings.Join(parts, ".")

	// replace mu since we know its ours
	if strings.HasPrefix(namespace, "mu.micro") {
		namespace = strings.Replace(namespace, "mu.micro", "go.micro", 1)
	}

	// web dashboard if namespace matches
	if namespace == FullNamespace {
		s.Router.ServeHTTP(w, r)
		return
	}

	// if a host has no subdomain serve dashboard
	v, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil || v == host {
		s.Router.ServeHTTP(w, r)
		return
	}

	// check if its a web request
	if _, _, isWeb := s.resolver.Info(r); isWeb {
		s.Router.ServeHTTP(w, r)
		return
	}

	// otherwise serve the proxy
	s.prx.ServeHTTP(w, r)
}

// proxy is a http reverse proxy
func (s *srv) proxy() *proxy {
	director := func(r *http.Request) {
		kill := func() {
			r.URL.Host = ""
			r.URL.Path = ""
			r.URL.Scheme = ""
			r.Host = ""
			r.RequestURI = ""
		}

		// check to see if the endpoint was encoded in the request context
		// by the auth wrapper
		var endpoint *res.Endpoint
		if val, ok := (r.Context().Value(res.Endpoint{})).(*res.Endpoint); ok {
			endpoint = val
		}

		// TODO: better error handling
		var err error
		if endpoint == nil {
			if endpoint, err = s.resolver.Resolve(r); err != nil {
				fmt.Printf("Failed to resolve url: %v: %v\n", r.URL, err)
				kill()
				return
			}
		}

		r.Header.Set(BasePathHeader, "/"+endpoint.Name)
		r.URL.Host = endpoint.Host
		r.URL.Path = endpoint.Path
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &proxy{
		Router:   &httputil.ReverseProxy{Director: director},
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

	// TODO: lookup icon

	var webServices []webService
	for _, srv := range services {
		if strings.Index(srv.Name, FullNamespace) == 0 && len(strings.TrimPrefix(srv.Name, FullNamespace)) > 0 {
			webServices = append(webServices, webService{
				Name: strings.Replace(srv.Name, FullNamespace+".", "", 1),
			})
		}
	}

	sort.Slice(webServices, func(i, j int) bool { return webServices[i].Name < webServices[j].Name })

	type templateData struct {
		HasWebServices bool
		WebServices    []webService
	}

	data := templateData{len(webServices) > 0, webServices}
	s.render(w, r, indexTemplate, data)
}

func (s *srv) registryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	svc := vars["name"]

	if len(svc) > 0 {
		sv, err := s.registry.GetService(svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if len(sv) == 0 {
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

		s.render(w, r, serviceTemplate, sv)
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

	s.render(w, r, registryTemplate, services)
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

	s.render(w, r, callTemplate, serviceMap)
}

func (s *srv) render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
		"Title":  strings.Title,
		"First": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.Title(string(s[0]))
		},
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

	// If the user is logged in, render Account instead of Login
	loginTitle := "Login"
	user := ""

	if c, err := r.Cookie(auth.TokenCookieName); err == nil && c != nil {
		token := strings.TrimPrefix(c.Value, auth.TokenCookieName+"=")
		if acc, err := s.auth.Inspect(token); err == nil {
			loginTitle = "Account"
			user = acc.ID
		}
	}

	if err := t.ExecuteTemplate(w, "layout", map[string]interface{}{
		"LoginTitle": loginTitle,
		"LoginURL":   loginURL,
		"StatsURL":   statsURL,
		"Results":    data,
		"User":       user,
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
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
	}
	if len(ctx.String("type")) > 0 {
		Type = ctx.String("type")
	}
	if len(ctx.String("namespace")) > 0 {
		// remove the service type from the namespace to allow for
		// backwards compatability
		Namespace = strings.TrimSuffix(ctx.String("namespace"), "."+Type)
	}

	// FullNamespace has the format: "go.micro.web"
	FullNamespace = Namespace + "." + Type

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	// use the caching registry
	cache := cache.New((*cmd.DefaultOptions().Registry))
	reg := &reg{Registry: cache}

	// start the watcher
	go reg.watch()

	s := &srv{
		Router:   mux.NewRouter(),
		registry: reg,
		// our internal resolver
		resolver: &resolver{
			// Default to type path
			Type:      Resolver,
			Namespace: FullNamespace,
			Selector: selector.NewSelector(
				selector.Registry(reg),
			),
		},
		auth: *cmd.DefaultOptions().Auth,
	}

	var h http.Handler
	// set as the server
	h = s

	if ctx.Bool("enable_stats") {
		statsURL = "/stats"
		st := stats.New()
		s.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	// create the proxy
	p := s.proxy()

	// the web handler itself
	s.HandleFunc("/favicon.ico", faviconHandler)
	s.HandleFunc("/client", s.callHandler)
	s.HandleFunc("/services", s.registryHandler)
	s.HandleFunc("/service/{name}", s.registryHandler)
	s.HandleFunc("/rpc", handler.RPC)
	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)
	s.HandleFunc("/", s.indexHandler)

	// insert the proxy
	s.prx = p

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

	// create the namespace resolver
	nsResolver := nsResolver.NewNamespaceResolver(Type, Namespace)

	// create the service with the resolver and namespace resolver
	srv := httpapi.NewServer(Address,
		server.Resolver(s.resolver),
		server.NamespaceResolver(nsResolver),
	)

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
				Name:    "resolver",
				Usage:   "Set the resolver to route to services e.g path, domain",
				EnvVars: []string{"MICRO_WEB_RESOLVER"},
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
