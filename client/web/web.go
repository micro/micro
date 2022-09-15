// Package web is a web dashboard
package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fatih/camelcase"
	"github.com/gorilla/mux"
	"github.com/micro/micro/v3/client/web/html"
	"github.com/micro/micro/v3/cmd"
	apiAuth "github.com/micro/micro/v3/service/api/auth"
	"github.com/micro/micro/v3/service/api/resolver"
	"github.com/micro/micro/v3/service/api/resolver/subdomain"
	httpapi "github.com/micro/micro/v3/service/api/server/http"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
	regRouter "github.com/micro/micro/v3/service/router/registry"
	"github.com/serenize/snaker"
	"github.com/urfave/cli/v2"
)

//Meta Fields of micro web
var (
	Name      = "web"
	API       = "http://localhost:8080"
	Address   = ":8082"
	Namespace = "micro"
	Resolver  = "path"
	LoginURL  = "/login"
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
	resolver resolver.Resolver
}

type reg struct {
	registry.Registry

	sync.RWMutex
	lastPull time.Time
	services []*registry.Service
}

func init() {
	cmd.Register(
		&cli.Command{
			Name:   "web",
			Usage:  "Run the micro web UI",
			Action: Run,
			Flags:  Flags,
		},
	)
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

	s.Router.ServeHTTP(w, r)
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
	s.render(w, r, html.NotFoundTemplate, nil)
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
		log.Printf("Error listing services: %v", err)
	}

	type webService struct {
		Name string
		Link string
		Icon string // TODO: lookup icon
	}

	var webServices []webService
	for _, srv := range services {
		name := srv.Name
		link := fmt.Sprintf("/%v", name)

		if len(srv.Endpoints) == 0 {
			continue
		}

		// in the case of 3 letter things e.g m3o convert to M3O
		if len(name) <= 3 && strings.ContainsAny(name, "012345789") {
			name = strings.ToUpper(name)
		}

		webServices = append(webServices, webService{Name: name, Link: link})
	}

	sort.Slice(webServices, func(i, j int) bool { return webServices[i].Name < webServices[j].Name })

	type templateData struct {
		HasWebServices bool
		WebServices    []webService
	}

	data := templateData{len(webServices) > 0, webServices}
	s.render(w, r, html.IndexTemplate, data)
}

func (s *srv) loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		s.generateTokenHandler(w, req)
		return
	}

	t, err := template.New("template").Parse(html.LoginTemplate)
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
		t, err := template.New("template").Parse(html.LoginTemplate)
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

		s.render(w, r, html.ServiceTemplate, sv)
		return
	}

	services, err := s.registry.ListServices(registry.ListDomain(domain))
	if err != nil {
		log.Printf("Error listing services: %v", err)
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

	s.render(w, r, html.RegistryTemplate, services)
}

func (s *srv) callHandler(w http.ResponseWriter, r *http.Request) {
	// if we're using the subdomain resolver, we want to use a custom domain
	domain := registry.DefaultDomain
	if res, ok := s.resolver.(*subdomain.Resolver); ok {
		domain = res.Domain(r)
	}

	services, err := s.registry.ListServices(registry.ListDomain(domain))
	if err != nil {
		log.Printf("Error listing services: %v", err)
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

	s.render(w, r, html.CallTemplate, serviceMap)
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
		log.Printf("Error getting service %s: %v", name, err)
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

	s.render(w, r, html.WebTemplate, serviceMap, templateValue{
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
		"Endpoint": func(ep string) string {
			return strings.Replace(ep, ".", "/", -1)
		},
	}).Parse(html.LayoutTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
	t, err = t.Parse(tmpl)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	apiURL := API
	u, err := url.Parse(apiURL)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	filepath.Join(u.Path, r.URL.Path)

	// If the user is logged in, render Account instead of Login
	loginTitle := "Login"
	loginLink := LoginURL
	user := ""
	token := r.Header.Get("Authorization")

	acc, ok := auth.AccountFromContext(r.Context())
	if ok {
		user = acc.ID
		loginTitle = "Logout"
		loginLink = "/logout"
	}

	templateData := map[string]interface{}{
		"ApiURL":     apiURL,
		"LoginTitle": loginTitle,
		"LoginURL":   loginLink,
		"Results":    data,
		"User":       user,
		"Token":      token,
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
	if len(ctx.String("api_address")) > 0 {
		API = ctx.String("api_address")
	}
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

	// Setup the web resolver
	var res resolver.Resolver

	// the default resolver
	res = &WebResolver{
		Router: regRouter.NewRouter(
			router.Registry(registry.DefaultRegistry),
		),
		Options: resolver.NewOptions(resolver.WithServicePrefix(
			Namespace,
		)),
	}

	// switch for subdomain resolver
	if Resolver == "subdomain" {
		res = subdomain.NewResolver(res)
	}

	srv := &srv{
		Router: mux.NewRouter(),
		registry: &reg{
			Registry: registry.DefaultRegistry,
		},
		resolver: res,
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
	srv.HandleFunc("/{service}", srv.serviceHandler)
	srv.HandleFunc("/", srv.indexHandler)

	// create the service and add the auth wrapper
	aw := apiAuth.Wrapper(srv.resolver, Namespace)
	server := httpapi.NewServer(Address)
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

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
	<-ch

	if err := server.Stop(); err != nil {
		log.Fatal(err)
	}

	return nil
}

var (
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "api_address",
			Usage:   "Set the api address to call e.g http://localhost:8080",
			EnvVars: []string{"MICRO_API_ADDRESS"},
		},
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
	}
)

func reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

type sortedServices struct {
	services []*registry.Service
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
