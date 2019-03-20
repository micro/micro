package v1

import (
	"github.com/gorilla/mux"
	"github.com/micro/go-log"
	"net/http"
)

func (api *API) stats(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	serviceName :=vars["name"]

	log.Logf(serviceName)

	writeJsonData(w, serviceName)
	return
}
