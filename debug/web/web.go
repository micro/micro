// Package web provides a dashboard for debugging and introspection of go-micro services
package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/web"
	logpb "github.com/micro/micro/debug/log/proto"
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

	// renders the per service debug dashboard
	r.HandleFunc("/stats/{service}", statsDashboard)
	// endpoint for logs
	r.HandleFunc("/log/{service}", logDashboard)

	wrapper := &netdataWrapper{
		netdataproxy: netdata.ServeHTTP,
	}
	service := web.NewService(opts...)

	// endpoints required for displaying stats and metrics
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

	// endpoints for infrastructure
	service.HandleFunc("/infra", http.RedirectHandler("/debug/infra/", http.StatusTemporaryRedirect).ServeHTTP)
	service.HandleFunc("/infra/", wrapper.proxyNetdata)

	service.Run()
}

func renderDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		dashboardTemplate.Execute(w, nil)
	}
}

func logDashboard(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	service, ok := v["service"]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Service not found\n")
		return
	}
	// get the logs
	c := logpb.NewLogService("go.micro.debug", client.DefaultClient)

	// TODO: ability to stream
	logs, err := c.Read(context.TODO(), &logpb.ReadRequest{
		Service: service,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	t := template.Must(template.New("logs").Parse(logTemplate))

	t.Execute(w, struct {
		Name    string
		Records []*logpb.Record
	}{
		Name:    service,
		Records: logs.Records,
	})
}

// statsDashboard renders the dashboard for services stats
func statsDashboard(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	service, found := v["service"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Service not found\n")
		return
	}

	// execute the dashboad template
	dashboardTemplate.Execute(w, struct{ Name, Service string }{
		Name:    service,
		Service: strings.ReplaceAll(service, ".", "_"),
	})
}

type netdataWrapper struct {
	netdataproxy func(http.ResponseWriter, *http.Request)
}

func (n *netdataWrapper) proxyNetdata(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/infra")
	n.netdataproxy(w, r)
}
