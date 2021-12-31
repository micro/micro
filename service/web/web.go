// Package web is a web dashboard
package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/camelcase"
	"github.com/go-acme/lego/v3/providers/dns/cloudflare"
	"github.com/gorilla/mux"
	"github.com/micro/micro/v3/cmd"
	"github.com/micro/micro/v3/service"
	server "github.com/micro/micro/v3/service/api"
	apiAuth "github.com/micro/micro/v3/service/api/auth"
	res "github.com/micro/micro/v3/service/api/resolver"
	"github.com/micro/micro/v3/service/api/resolver/subdomain"
	httpapi "github.com/micro/micro/v3/service/api/server/http"
	"github.com/micro/micro/v3/service/auth"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	muregistry "github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
	regRouter "github.com/micro/micro/v3/service/router/registry"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/util/acme"
	"github.com/micro/micro/v3/util/acme/autocert"
	"github.com/micro/micro/v3/util/acme/certmagic"
	"github.com/micro/micro/v3/util/helper"
	"github.com/micro/micro/v3/util/sync/memory"
	"github.com/serenize/snaker"
	"github.com/urfave/cli/v2"
)

//Meta Fields of micro web
var (
	Name                  = "web"
	Address               = ":8082"
	Namespace             = "micro"
	Resolver              = "path"
	LoginURL              = "/login"
	ACMEProvider          = "autocert"
	ACMEChallengeProvider = "cloudflare"
	ACMECA                = acme.LetsEncryptProductionCA

	// Host name the web dashboard is served on
	Host, _ = os.Hostname()
	// Token cookie name
	TokenCookieName = "micro-token"
)

type srv struct {
	*mux.Router
	// registry we use
	registry registry.Registry
	// the resolver
	resolver res.Resolver
}

type reg struct {
	registry.Registry

	sync.RWMutex
	lastPull time.Time
	services []*registry.Service
}

// ServeHTTP serves the web dashboard and proxies where appropriate
func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check if authenticated
	if r.URL.Path != LoginURL {
		c, err := r.Cookie(TokenCookieName)
		if err != nil || c == nil {
			http.Redirect(w, r, LoginURL, 302)
			return
		}

		// check the token is valid
		token := strings.TrimPrefix(c.Value, TokenCookieName+"=")
		if len(token) == 0 {
			http.Redirect(w, r, LoginURL, 302)
			return
		}
	}

	// set defaults on the request
	if len(r.URL.Host) == 0 {
		r.URL.Host = r.Host
	}
	if len(r.URL.Scheme) == 0 {
		r.URL.Scheme = "http"
	}

	// no endpoint was set in the context, so we'll look it up. If the router returns an error we will
	// send the request to the mux which will render the web dashboard.
	s.Router.ServeHTTP(w, r)
	//_, err := s.resolver.Resolve(r)
	//if err != nil {
	//	return
	//}
}

func split(v string) string {
	parts := camelcase.Split(strings.Replace(v, ".", "", 1))
	return strings.Join(parts, " ")
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
	httpapi.SetHeaders(w, r)

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

	var webServices []webService
	for _, srv := range services {
		name := srv.Name

		if len(srv.Endpoints) == 0 {
			continue
		}

		// in the case of 3 letter things e.g m3o convert to M3O
		if len(name) <= 3 && strings.ContainsAny(name, "012345789") {
			name = strings.ToUpper(name)
		}

		webServices = append(webServices, webService{Name: name, Link: fmt.Sprintf("/%v", name)})
	}

	sort.Slice(webServices, func(i, j int) bool { return webServices[i].Name < webServices[j].Name })

	type templateData struct {
		HasWebServices bool
		WebServices    []webService
	}

	data := templateData{len(webServices) > 0, webServices}
	s.render(w, r, indexTemplate, data)
}

func (s *srv) loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		s.generateTokenHandler(w, req)
		return
	}

	t, err := template.New("template").Parse(loginTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "basic", map[string]interface{}{
		"foo": "bar",
	}); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
	}
}

func (s *srv) logoutHandler(w http.ResponseWriter, req *http.Request) {
	var domain string
	if arr := strings.Split(req.Host, ":"); len(arr) > 0 {
		domain = arr[0]
	}

	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
		Value:   "",
		Expires: time.Unix(0, 0),
		Domain:  domain,
		Secure:  true,
	})

	http.Redirect(w, req, "/", 302)
}

func (s *srv) generateTokenHandler(w http.ResponseWriter, req *http.Request) {
	renderError := func(errMsg string) {
		t, err := template.New("template").Parse(loginTemplate)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, "basic", map[string]interface{}{
			"error": errMsg,
		}); err != nil {
			http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
		}
	}

	user := req.PostFormValue("username")
	if len(user) == 0 {
		renderError("Missing Username")
		return
	}

	pass := req.PostFormValue("password")
	if len(pass) == 0 {
		renderError("Missing Password")
		return
	}

	acc, err := auth.Token(
		auth.WithCredentials(user, pass),
		auth.WithTokenIssuer(Namespace),
		auth.WithExpiry(time.Hour*24*7),
	)
	if err != nil {
		renderError("Authentication failed: " + err.Error())
		return
	}

	var domain string
	if arr := strings.Split(req.Host, ":"); len(arr) > 0 {
		domain = arr[0]
	}

	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
		Value:   acc.AccessToken,
		Expires: acc.Expiry,
		Domain:  domain,
		Secure:  true,
	})

	http.Redirect(w, req, "/", http.StatusFound)
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

func (s *srv) serviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["service"]
	if len(name) == 0 {
		return
	}

	// if we're using the subdomain resolver, we want to use a custom domain
	domain := registry.DefaultDomain
	if res, ok := s.resolver.(*subdomain.Resolver); ok {
		domain = res.Domain(r)
	}

	services, err := s.registry.GetService(name, registry.GetDomain(domain))
	if err != nil {
		log.Errorf("Error getting service %s: %v", name, err)
	}

	sort.Sort(sortedServices{services})

	serviceMap := make(map[string][]*registry.Endpoint)

	for _, service := range services {
		if len(service.Endpoints) > 0 {
			serviceMap[service.Name] = service.Endpoints
			continue
		}
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

	s.render(w, r, webTemplate, serviceMap, templateValue{
		Key:   "Name",
		Value: name,
	})
}

type templateValue struct {
	Key   string
	Value interface{}
}

func (s *srv) render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, vals ...templateValue) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"Split":  split,
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
	loginLink := LoginURL
	user := ""

	acc, ok := auth.AccountFromContext(r.Context())
	if ok {
		user = acc.ID
		loginTitle = "Logout"
		loginLink = "/logout"
	}

	templateData := map[string]interface{}{
		"LoginTitle": loginTitle,
		"LoginURL":   loginLink,
		"Results":    data,
		"User":       user,
	}

	// add extra values
	for _, val := range vals {
		templateData[val.Key] = val.Value
	}

	if err := t.ExecuteTemplate(w, "layout",
		templateData,
	); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func Run(ctx *cli.Context) error {
	if len(ctx.String("server_name")) > 0 {
		Name = ctx.String("server_name")
	}
	if len(ctx.String("resolver")) > 0 {
		Resolver = ctx.String("resolver")
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
		Namespace = ctx.String("namespace")
	}

	// Initialize Server
	s := service.New(service.Name(Name))

	// Setup the web resolver
	var resolver res.Resolver

	// the default resolver
	resolver = &WebResolver{
		Router: regRouter.NewRouter(
			router.Registry(muregistry.DefaultRegistry),
		),
		Options: res.NewOptions(res.WithServicePrefix(
			Namespace,
		)),
	}

	// switch for subdomain resolver
	if Resolver == "subdomain" {
		resolver = subdomain.NewResolver(resolver)
	}

	srv := &srv{
		Router: mux.NewRouter(),
		registry: &reg{
			Registry: muregistry.DefaultRegistry,
		},
		resolver: resolver,
	}

	var h http.Handler
	// set as the server
	h = srv

	// the web handler itself
	srv.HandleFunc("/favicon.ico", faviconHandler)
	srv.HandleFunc("/404", srv.notFoundHandler)
	srv.HandleFunc("/login", srv.loginHandler)
	srv.HandleFunc("/logout", srv.logoutHandler)
	srv.HandleFunc("/client", srv.callHandler)
	srv.HandleFunc("/services", srv.registryHandler)
	srv.HandleFunc("/service/{name}", srv.registryHandler)
	srv.Handle("/rpc", NewRPCHandler(resolver, s.Client()))
	srv.HandleFunc("/{service}", srv.serviceHandler)
	srv.HandleFunc("/", srv.indexHandler)

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
	aw := apiAuth.Wrapper(srv.resolver, Namespace)
	server := httpapi.NewServer(Address)

	server.Init(opts...)
	server.Handle("/", aw(h))

	// Setup auth redirect
	if len(ctx.String("login_url")) > 0 {
		LoginURL = ctx.String("login_url")
	}

	// set the login url
	auth.DefaultAuth.Init(auth.LoginURL(LoginURL))

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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "web_address",
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
				Name:    "login_url",
				EnvVars: []string{"MICRO_WEB_LOGIN_URL"},
				Usage:   "The relative URL where a user can login",
			},
		},
	})
}

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

type sortedServices struct {
	services []*muregistry.Service
}

func (s sortedServices) Len() int {
	return len(s.services)
}

func (s sortedServices) Less(i, j int) bool {
	return s.services[i].Name < s.services[j].Name
}

func (s sortedServices) Swap(i, j int) {
	s.services[i], s.services[j] = s.services[j], s.services[i]
}
