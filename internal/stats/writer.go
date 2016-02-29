package stats

import (
	"net/http"
)

type writer struct {
	http.ResponseWriter
	status int
}

func (w *writer) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}
