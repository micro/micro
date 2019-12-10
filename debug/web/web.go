// Package web provides a dashboard for debugging and introspection of go-micro services
package web

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
)

// Run starts go.micro.web.debug
func Run(ctx *cli.Context) {
	dashboardTemplate = template.Must(template.New("dashboard").Parse(dashboardHTML))

	opts := []web.Option{
		web.Name("go.micro.web.debug"),
	}

	address := ctx.GlobalString("server_address")
	if len(address) > 0 {
		opts = append(opts, web.Address(address))
	}

	u, err := url.Parse(ctx.String("netdata_url"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	netdata := httputil.NewSingleHostReverseProxy(u)

	r := mux.NewRouter()
	r.HandleFunc("/", renderDashboard)

	wrapper := &netdataWrapper{
		netdataproxy: netdata.ServeHTTP,
	}
	r.HandleFunc("/infra", wrapper.proxyNetdata)

	// opts = append(opts, web.Handler(r))

	service := web.NewService(opts...)
	service.HandleFunc("/dashboard.js", netdata.ServeHTTP)
	service.HandleFunc("/dashboard.css", netdata.ServeHTTP)
	service.HandleFunc("/dashboard.slate.css", netdata.ServeHTTP)
	service.HandleFunc("/dashboard_info.js", netdata.ServeHTTP)
	service.HandleFunc("/main.css", netdata.ServeHTTP)
	service.HandleFunc("/main.js", netdata.ServeHTTP)
	service.HandleFunc("/images/", netdata.ServeHTTP)
	service.HandleFunc("/lib/", netdata.ServeHTTP)
	service.HandleFunc("/css/", netdata.ServeHTTP)
	service.HandleFunc("/api/", netdata.ServeHTTP)
	service.HandleFunc("/", r.ServeHTTP)

	service.Run()
}

func renderDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		dashboardTemplate.Execute(w, nil)
	}
}

type netdataWrapper struct {
	netdataproxy func(http.ResponseWriter, *http.Request)
}

func (n *netdataWrapper) proxyNetdata(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "/"
	n.netdataproxy(w, r)
}
