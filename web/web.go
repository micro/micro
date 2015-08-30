package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/myodc/go-micro/registry"
)

var (
	Address = ":8082"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexTemplate)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {
	services, err := registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	t, err := template.New("registry").Parse(registryTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	if err := t.ExecuteTemplate(w, "T", services); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
}

func run() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/registry", registryHandler)

	if err := http.ListenAndServe(Address, nil); err != nil {
		log.Fatal(err)
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "web",
			Usage: "Run the micro web app",
			Action: func(c *cli.Context) {
				run()
			},
		},
	}
}
