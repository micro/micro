package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	log "github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/serenize/snaker"

	"golang.org/x/net/context"
)

var (
	re        = regexp.MustCompile("^[a-zA-Z0-9]+$")
	Address   = ":8082"
	Namespace = "go.micro.web"
)

type server struct {
	*mux.Router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		return
	}

	s.Router.ServeHTTP(w, r)
}

func (s *server) proxy() http.Handler {
	sel := selector.NewSelector(
		selector.Registry(registry.DefaultRegistry),
	)

	director := func(r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 {
			return
		}
		if !re.MatchString(parts[1]) {
			return
		}
		next, err := sel.Select(Namespace + "." + parts[1])
		if err != nil {
			return
		}
		r.URL.Scheme = "http"
		s, err := next()
		if err != nil {
			return
		}
		r.URL.Host = fmt.Sprintf("%s:%d", s.Address, s.Port)
		r.URL.Path = "/" + strings.Join(parts[2:], "/")
	}
	return &httputil.ReverseProxy{
		Director: director,
	}
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, indexTemplate, nil)
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

		render(w, r, serviceTemplate, s)
		return
	}

	services, err := registry.ListServices()
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
		return
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

	render(w, r, registryTemplate, services)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, queryTemplate, nil)
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

func render(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	t, err := template.New("template").Funcs(template.FuncMap{
		"format": format,
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
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), 500)
	}
}

func run() {
	r := mux.NewRouter()
	s := &server{r}
	s.HandleFunc("/", indexHandler)
	s.HandleFunc("/registry", registryHandler)
	s.HandleFunc("/rpc", rpcHandler)
	s.HandleFunc("/query", queryHandler)
	s.PathPrefix("/{service:[a-zA-Z0-9]+}").Handler(s.proxy())

	if err := http.ListenAndServe(Address, s); err != nil {
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
