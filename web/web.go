// Package web is a web dashboard
package web

import (
	"fmt"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/selector"
	"github.com/micro/micro/internal/metadata"
	"github.com/micro/micro/web/api/v1"
	"github.com/micro/micro/web/common"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-api/server"
	httpapi "github.com/micro/go-api/server/http"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/helper"
	"github.com/micro/micro/internal/stats"
	"github.com/micro/micro/plugin"
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

	// path to the html directory
	StaticDir = "web/webapp/dist"

	statsURL string
)

type srv struct {
	*mux.Router
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

		parts := strings.Split(r.URL.Path, "/proxy/")
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

		r.Header.Set(common.BasePathHeader, "/"+parts[1])
		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Path = "/" + strings.Join(parts[2:], "/")
		r.URL.Scheme = "http"
		r.Host = r.URL.Host
	}

	return &common.Proxy{
		Default:  &httputil.ReverseProxy{Director: director},
		Director: director,
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func run(ctx *cli.Context, srvOpts ...micro.Option) {
	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}
	if len(ctx.String("address")) > 0 {
		Address = ctx.String("address")
	}
	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	if len(ctx.String("api_namespace")) > 0 {
		common.APINamespace = ctx.String("api_namespace")
	}

	if len(ctx.String("static_dir")) > 0 {
		// check static-dir existing
		_, err := os.Stat(ctx.String("static_dir"))
		if err == nil {
			StaticDir = ctx.String("static_dir")
		}
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	var h http.Handler
	r := mux.NewRouter()
	s := &srv{r}
	h = s

	if ctx.GlobalBool("enable_stats") {
		statsURL = "/stats"
		st := stats.New()
		s.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	apiV1 := v1.API{}
	apiV1.InitV1Handler(s.Router, Namespace)

	s.HandleFunc("/favicon.ico", faviconHandler)
	s.PathPrefix("/proxy/{service:[a-zA-Z0-9]+}").Handler(s.proxy())
	s.PathPrefix("/").Handler(http.FileServer(http.Dir(StaticDir)))

	var opts []server.Option

	if ctx.GlobalBool("enable_acme") {
		hosts := helper.ACMEHosts(ctx)
		opts = append(opts, server.EnableACME(true))
		opts = append(opts, server.ACMEHosts(hosts...))
	} else if ctx.GlobalBool("enable_tls") {
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
	if i := time.Duration(ctx.GlobalInt("register_ttl")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterTTL(i*time.Second))
	}
	if i := time.Duration(ctx.GlobalInt("register_interval")); i > 0 {
		srvOpts = append(srvOpts, micro.RegisterInterval(i*time.Second))
	}

	serviceMetadata := map[string]string{
		metadata.MetadataFieldNameServerType: metadata.MetadataServiceTypeWebDashboard,
	}
	srvOpts = append(srvOpts, micro.Metadata(serviceMetadata))

	// Initialise Server
	service := micro.NewService(srvOpts...)

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

func Commands(options ...micro.Option) []cli.Command {
	command := cli.Command{
		Name:  "web",
		Usage: "Run the web dashboard",
		Action: func(c *cli.Context) {
			run(c, options...)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "address",
				Usage:  "Set the web UI address e.g 0.0.0.0:8082",
				EnvVar: "MICRO_WEB_ADDRESS",
			},
			cli.StringFlag{
				Name:   "namespace",
				Usage:  "Set the namespace used by the Web proxy e.g. com.example.web",
				EnvVar: "MICRO_WEB_NAMESPACE",
			},
			cli.StringFlag{
				Name:   "static_dir",
				Usage:  "Set the static dir of micro web",
				EnvVar: "MICRO_WEB_STATIC_DIR",
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

	return []cli.Command{command}
}
