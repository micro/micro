package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/myodc/go-micro/registry"
)

var (
	Address = ":8082"
)

func format(v *registry.Value) string {
	if v == nil || len(v.Values) == 0 {
		return "{}"
	}
	return formatEndpoint(v.Values[0], 0)
}

func formatEndpoint(v *registry.Value, r int) string {
	if len(v.Values) == 0 {
		fparts := []string{"\n", "{", "\n", "\t", "%s %s", "\n", "}"}
		if r > 0 {
			fparts = []string{"\n", "\t", "%s %s"}
			for i := 0; i < r; i++ {
				fparts[1] += "\t"
			}
		}
		return fmt.Sprintf(strings.Join(fparts, ""), strings.ToLower(v.Name), v.Type)
	}

	fparts := []string{"\n", "{", "\n", "\t", "%s %s", " {", "\n", "", "\n", "\t", "}", "\n", "}"}
	i := 7

	if r > 0 {
		fparts = []string{"\n", "\t", "%s %s", " {", "\n", "\t", "\n", "\t", "}"}
		i = 5
	}

	var app string
	for j := 0; j < r; j++ {
		if r > 0 {
			fparts[1] += "\t"
			fparts[7] += "\t"
		}
		app += "\t"
	}
	app += "\t%s"

	vals := []interface{}{strings.ToLower(v.Name), v.Type}

	for _, val := range v.Values {
		fparts[i] += app
		vals = append(vals, formatEndpoint(val, r+1))
	}

	return fmt.Sprintf(strings.Join(fparts, ""), vals...)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, indexTemplate)
}

func registryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	svc := r.Form.Get("service")
	if len(svc) > 0 {
		s, err := registry.GetService(svc)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		t, err := template.New("service").Funcs(template.FuncMap{
			"format": format,
		}).Parse(serviceTemplate)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}

		if err := t.ExecuteTemplate(w, "T", s); err != nil {
			http.Error(w, "Error occurred:"+err.Error(), 500)
			return
		}
		return
	}

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
