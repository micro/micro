package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"

	"golang.org/x/net/context"
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

func queryHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("query").Parse(queryTemplate)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}

	if err := t.ExecuteTemplate(w, "T", nil); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
	}
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var service, method string
	var request interface{}

	// response content type
	w.Header().Set("Content-Type", "application/json")

	switch r.Header.Get("Content-Type") {
	case "application/json":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e := errors.BadRequest("go.micro.api", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		var body map[string]interface{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			e := errors.BadRequest("go.micro.api", err.Error())
			w.WriteHeader(400)
			w.Write([]byte(e.Error()))
			return
		}

		service = body["service"].(string)
		method = body["method"].(string)
		request = body["request"]
	default:
		r.ParseForm()
		service = r.Form.Get("service")
		method = r.Form.Get("method")
		json.Unmarshal([]byte(r.Form.Get("request")), &request)
	}

	var response map[string]interface{}
	req := client.NewJsonRequest(service, method, request)
	err := client.Call(context.Background(), req, &response)
	if err != nil {
		log.Errorf("Error calling %s.%s: %v", service, method, err)
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			w.WriteHeader(500)
		default:
			w.WriteHeader(int(ce.Code))
		}
		w.Write([]byte(ce.Error()))
		return
	}

	b, _ := json.Marshal(response)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}

func run() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/registry", registryHandler)
	http.HandleFunc("/rpc", rpcHandler)
	http.HandleFunc("/query", queryHandler)

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
