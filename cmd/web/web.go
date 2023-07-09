// Package web is a web dashboard
package web

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"micro.dev/v4/cmd"
	"micro.dev/v4/cmd/web/html"
)

// Meta Fields of micro web
var (
	Name      = "web"
	API       = "http://localhost:8080"
	Address   = ":8082"
	Namespace = "micro"
	LoginURL  = "/login"
	// Host name the web dashboard is served on
	Host, _ = os.Hostname()
	// Token cookie name
	TokenCookieName = "micro-token"

	// create a session store
	mtx      sync.RWMutex
	sessions = map[string]*session{}
)

type srv struct {
	*mux.Router
}

type session struct {
	// token used for the session
	Token string
}

//go:embed html/* html/assets/*
var content embed.FS

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
	if strings.HasPrefix(r.URL.Path, "/assets/") {
		s.Router.ServeHTTP(w, r)
		return
	}

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

		// if we have a session retrieve it
		mtx.RLock()
		sess, ok := sessions[token]
		mtx.RUnlock()

		// no session, go get the account
		if !ok {
			// save the session
			mtx.Lock()
			sess = &session{
				Token: token,
			}
			sessions[token] = sess
			mtx.Unlock()
		}

		// create a new context
		ctx := context.WithValue(r.Context(), session{}, sess)

		// redefine request with context
		r = r.Clone(ctx)
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

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func (s *srv) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.render(w, r, html.NotFoundTemplate, nil)
}

func (s *srv) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	s.render(w, r, html.IndexTemplate, struct{}{})
}

func (s *srv) loginHandler(w http.ResponseWriter, req *http.Request) {
	s.render(w, req, html.LoginTemplate, struct{}{})
}

func (s *srv) logoutHandler(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    TokenCookieName,
		Value:   "",
		Expires: time.Unix(0, 0),
		Secure:  true,
	})

	http.Redirect(w, req, "/", 302)
}

type templateValue struct {
	Key   string
	Value interface{}
}

func (s *srv) render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}, vals ...templateValue) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"Title": strings.Title,
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

	// set api from the hdear if available
	if v := r.Header.Get("Micro-API"); len(v) > 0 {
		apiURL = v
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	filepath.Join(u.Path, r.URL.Path)

	var token string

	sess, ok := r.Context().Value(session{}).(*session)
	if ok {
		token = sess.Token
	}

	templateData := map[string]interface{}{
		"ApiURL":    template.URL(apiURL),
		"Token":     token,
		"Results":   data,
		"Namespace": Namespace,
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
	// Setup auth redirect
	if len(ctx.String("login_url")) > 0 {
		LoginURL = ctx.String("login_url")
	}

	srv := &srv{
		Router: mux.NewRouter(),
	}

	htmlContent, err := fs.Sub(content, "html")
	if err != nil {
		log.Fatal(err)
	}

	// the web handler itself
	srv.HandleFunc("/favicon.ico", faviconHandler)
	srv.HandleFunc("/404", srv.notFoundHandler)
	srv.HandleFunc("/login", srv.loginHandler)
	srv.HandleFunc("/logout", srv.logoutHandler)
	srv.PathPrefix("/assets/").Handler(http.FileServer(http.FS(htmlContent)))
	srv.HandleFunc("/", srv.indexHandler)
	srv.HandleFunc("/services", srv.indexHandler)
	srv.HandleFunc("/{service}", srv.indexHandler)
	srv.HandleFunc("/{service}/{endpoint}", srv.indexHandler)
	srv.HandleFunc("/{service}/{endpoint}/{method}", srv.indexHandler)

	// create new http server
	server := &http.Server{
		Addr:    Address,
		Handler: srv,
	}

	if err := server.ListenAndServe(); err != nil {
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
			Name:    "login_url",
			EnvVars: []string{"MICRO_WEB_LOGIN_URL"},
			Usage:   "The relative URL where a user can login",
		},
	}
)
