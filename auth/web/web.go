package web

import (
	"net/http"
	"text/template"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

var (
	// Name of the auth web
	Name = "go.micro.web.auth"
	// Address is the auth web address
	Address = ":8012"
)

// Run the micro auth api
func Run(ctx *cli.Context, srvOpts ...micro.Option) {
	log.Init(log.WithFields(map[string]interface{}{"service": "auth"}))

	service := web.NewService(
		web.Name(Name),
		web.Address(Address),
	)

	h := handler{
		auth: service.Options().Service.Options().Auth,
	}

	if h.auth.Options().Provider == nil {
		log.Fatal("Auth provider is not set")
	}

	service.HandleFunc("/", h.indexHandler)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	auth auth.Auth
}

func (h handler) indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		h.createBasicAccountHandler(w, req)
		return
	}

	p := h.auth.Options().Provider
	if len(p.Redirect()) > 0 {
		http.Redirect(w, req, p.Endpoint(), http.StatusFound)
		return
	}

	t, err := template.New("template").Parse(templates)
	if err != nil {
		http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, p.String(), map[string]interface{}{
		"foo": "bar",
	}); err != nil {
		http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
	}
}

func (h handler) createBasicAccountHandler(w http.ResponseWriter, req *http.Request) {
	p := h.auth.Options().Provider
	if p.String() != "basic" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	renderError := func(errMsg string) {
		t, err := template.New("template").Parse(templates)
		if err != nil {
			http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := t.ExecuteTemplate(w, p.String(), map[string]interface{}{
			"error": errMsg,
		}); err != nil {
			http.Error(w, "Error occurred:"+err.Error(), http.StatusInternalServerError)
		}
	}

	email := req.PostFormValue("email")
	if len(email) == 0 {
		renderError("Missing Email")
		return
	}

	pass := req.PostFormValue("password")
	if len(pass) == 0 {
		renderError("Missing Password")
		return
	}

	acc, err := h.auth.Generate(email)
	if err != nil {
		renderError(err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    auth.CookieName,
		Value:   acc.Token,
		Expires: acc.Expiry,
		Secure:  true,
	})

	// TODO: Redirect based on the original request
	http.Redirect(w, req, "/", http.StatusFound)
}
