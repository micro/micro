// Package web is a web dashboard
package web

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/micro/v2/client"
	"github.com/micro/micro/v2/service/store"

	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/micro/cli/v2"
	res "github.com/micro/go-micro/v2/api/resolver"
	"github.com/micro/go-micro/v2/api/resolver/subdomain"
	"github.com/micro/go-micro/v2/api/server"
	"github.com/micro/go-micro/v2/api/server/acme"
	"github.com/micro/go-micro/v2/api/server/acme/autocert"
	"github.com/micro/go-micro/v2/api/server/acme/certmagic"
	"github.com/micro/go-micro/v2/api/server/cors"
	httpapi "github.com/micro/go-micro/v2/api/server/http"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/sync/memory"
	apiAuth "github.com/micro/micro/v2/client/api/auth"
	"github.com/micro/micro/v2/cmd"
	inauth "github.com/micro/micro/v2/internal/auth"
	"github.com/micro/micro/v2/internal/handler"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/internal/resolver/web"
	"github.com/micro/micro/v2/internal/stats"
	"github.com/micro/micro/v2/service"
	muauth "github.com/micro/micro/v2/service/auth"
	muregistry "github.com/micro/micro/v2/service/registry"
	"github.com/micro/micro/v2/service/router"
	"github.com/serenize/snaker"
)

//Meta Fields of micro web
var (
	// Default server name
	Name = "go.micro.web"
	// Default address to bind to
	Address = ":8082"
	// The namespace to serve
	// Example:
	// Namespace + /[Service]/foo/bar
	// Host: Namespace.Service Endpoint: /foo/bar
	Namespace = "go.micro"
	Type      = "web"
	Resolver  = "path"
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
	resolver res.Resolver
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

// ServeHTTP serves the web dashboard and proxies where appropriate
func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// set defaults on the request
	if len(r.URL.Host) == 0 {
		r.URL.Host = r.Host
	}
	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}

	// the auth wrapper will resolve the route so it can verify the callers access. To prevent the
	// resolution happening twice, we'll check to see if the endpont was set in the context before
	// trying to resolve it ourselves. if an endpoint was found, we'll proxy to it.
	if _, ok := (r.Context().Value(res.Endpoint{})).(*res.Endpoint); ok {
		s.prx.ServeHTTP(w, r)
		return
	}

	// no endpoint was set in the context, so we'll look it up. If the router returns an error we will
	// send the request to the mux which will render the web dashboard.
	endpoint, err := s.resolver.Resolve(r)
	if err != nil {
		s.Router.ServeHTTP(w, r)
		return
	}

	// set the endpoint in the request context and then proxy to that endpoint
	*r = *r.Clone(context.WithValue(r.Context(), res.Endpoint{}, endpoint))
	s.prx.ServeHTTP(w, r)
}

// proxy is a http reverse proxy
func (s *srv) proxy() *proxy {
	director := func(r *http.Request) {
		// the endpoint would have been set by either the auth wrapper or by the servers ServeHTTP method
		// which invokes the proxy. If no endpoint is found, the router couldn't resolve the endpoint
		// and the mux didn't match any routes.
		endpoint, ok := (r.Context().Value(res.Endpoint{})).(*res.Endpoint)
		if !ok {
			r.URL.Path = "/404"
			return
		}

		// rewrite the request to go to this endpoint
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

func (s *srv) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.render(w, r, notFoundTemplate, nil)
}

func (s *srv) indexHandler(w http.ResponseWriter, r *http.Request) {
	cors.SetHeaders(w, r)

	if r.Method == "OPTIONS" {
		return
	}

	// if we're using the subdomain resolver, we want to use a custom domain
	domain := registry.DefaultDomain
	if res, ok := s.resolver.(*subdomain.Resolver); ok {
		domain = res.Domain(r)
	}

	services, err := s.registry.ListServices(registry.ListDomain(domain))
	if err != nil {
		log.Errorf("Error listing services: %v", err)
	}

	type webService struct {
		Name string
		Link string
		Icon string // TODO: lookup icon
	}

	prefix := Namespace + "." + Type + "."

	var webServices []webService
	for _, srv := range services {
		if !strings.HasPrefix(srv.Name, prefix) {
			continue
		}

		name := strings.TrimPrefix(srv.Name, prefix)

		// in the case of 3 letter things e.g m3o convert to M3O
		if len(name) <= 3 && strings.ContainsAny(name, "012345789") {
			name = strings.ToUpper(name)
		}

		webServices = append(webServices, webService{Name: name, Link: fmt.Sprintf("/%v/", name)})
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

	// if we're using the subdomain resolver, we want to use a custom domain
	domain := registry.DefaultDomain
	if res, ok := s.resolver.(*subdomain.Resolver); ok {
		domain = res.Domain(r)
	}

	if len(svc) > 0 {
		sv, err := s.registry.GetService(svc, registry.GetDomain(domain))
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

	services, err := s.registry.ListServices(registry.ListDomain(domain))
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
	// if we're using the subdomain resolver, we want to use a custom domain
	domain := registry.DefaultDomain
	if res, ok := s.resolver.(*subdomain.Resolver); ok {
		domain = res.Domain(r)
	}

	services, err := s.registry.ListServices(registry.ListDomain(domain))
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
		s, err := s.registry.GetService(service.Name, registry.GetDomain(domain))
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

	if c, err := r.Cookie(inauth.TokenCookieName); err == nil && c != nil {
		token := strings.TrimPrefix(c.Value, inauth.TokenCookieName+"=")
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

func Run(ctx *cli.Context) error {
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
	if len(ctx.String("web_address")) > 0 {
		Address = ctx.String("web_address")
	}
	if len(ctx.String("web_namespace")) > 0 {
		Namespace = ctx.String("web_namespace")
	}
	if len(ctx.String("web_host")) > 0 {
		Host = ctx.String("web_host")
	}
	if len(ctx.String("namespace")) > 0 {
		// remove the service type from the namespace to allow for
		// backwards compatability
		Namespace = strings.TrimSuffix(ctx.String("namespace"), "."+Type)
	}

	// Initialize Server
	s := service.New(service.Name(Name))

	// Setup the web resolver
	var resolver res.Resolver
	resolver = &web.Resolver{
		Router: router.DefaultRouter,
		Options: res.NewOptions(res.WithServicePrefix(
			Namespace + "." + Type,
		)),
	}
	if Resolver == "subdomain" {
		resolver = subdomain.NewResolver(resolver)
	}

	srv := &srv{
		Router: mux.NewRouter(),
		registry: &reg{
			Registry: muregistry.DefaultRegistry,
		},
		resolver: resolver,
		auth:     muauth.DefaultAuth,
	}

	var h http.Handler
	// set as the server
	h = srv

	if ctx.Bool("enable_stats") {
		statsURL = "/stats"
		st := stats.New()
		srv.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(srv)
		st.Start()
		defer st.Stop()
	}

	// create the proxy
	p := srv.proxy()

	// the web handler itself
	srv.HandleFunc("/favicon.ico", faviconHandler)
	srv.HandleFunc("/404", srv.notFoundHandler)
	srv.HandleFunc("/client", srv.callHandler)
	srv.HandleFunc("/services", srv.registryHandler)
	srv.HandleFunc("/service/{name}", srv.registryHandler)
	srv.Handle("/rpc", handler.NewRPCHandler(resolver))
	srv.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(p)
	srv.HandleFunc("/", srv.indexHandler)

	// insert the proxy
	srv.prx = p

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
			// TODO: support multiple providers in internal/acme as a map
			if ACMEChallengeProvider != "cloudflare" {
				log.Fatal("The only implemented DNS challenge provider is cloudflare")
			}

			apiToken := os.Getenv("CF_API_TOKEN")
			if len(apiToken) == 0 {
				log.Fatal("env variables CF_API_TOKEN and CF_ACCOUNT_ID must be set")
			}

			// create the store
			storage := certmagic.NewStorage(memory.NewSync(), store.DefaultStore)

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
			return err
		}

		opts = append(opts, server.EnableTLS(true))
		opts = append(opts, server.TLSConfig(config))
	}

	// create the service and add the auth wrapper
	aw := apiAuth.Wrapper(srv.resolver, Namespace+"."+Type)
	server := httpapi.NewServer(Address, server.WrapHandler(aw))

	server.Init(opts...)
	server.Handle("/", h)

	// Setup auth redirect
	if len(ctx.String("auth_login_url")) > 0 {
		loginURL = ctx.String("auth_login_url")
		muauth.DefaultAuth.Init(auth.LoginURL(loginURL))
	}

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

	if err := server.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func init() {
	cmd.Register(&cli.Command{
		Name:   "web",
		Usage:  "Run the web dashboard",
		Action: Run,
		Flags: append(client.Flags,
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
		),
	})
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
