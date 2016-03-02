package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/micro/internal/handler"
	"github.com/micro/micro/internal/server"
	"github.com/micro/micro/internal/stats"
)

var (
	Address      = ":8080"
	RPCPath      = "/rpc"
	APIPath      = "/"
	Namespace    = "go.micro.api"
	HeaderPrefix = "X-Micro-"
	CORS         = map[string]bool{"*": true}
)

type srv struct {
	*http.ServeMux
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

	s.ServeMux.ServeHTTP(w, r)
}

func run(ctx *cli.Context) {
	// Init API
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

	// create the router
	r := http.NewServeMux()
	log.Infof("Registering RPC Handler at %s", RPCPath)
	r.HandleFunc(RPCPath, handler.RPC)
	log.Infof("Registering API Handler at %s", APIPath)
	r.HandleFunc(APIPath, restHandler)

	s := &srv{r}
	var h http.Handler = s

	if ctx.GlobalBool("enable_stats") {
		st := stats.New()
		s.HandleFunc("/stats", st.StatsHandler)
		h = st.ServeHTTP(s)
		st.Start()
		defer st.Stop()
	}

	// create the server
	api := server.NewServer(Address)
	api.Init(opts...)
	api.Handle("/", h)

	// Initialise Server
	service := micro.NewService(
		micro.Name("go.micro.api"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	// Start API
	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop API
	if err := api.Stop(); err != nil {
		log.Fatal(err)
	}
}

func New(address string) server.Server {
	return server.NewServer(address)
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "api",
			Usage:  "Run the micro API",
			Action: run,
		},
	}
}
