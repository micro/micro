package stats

import (
	"net/http"
	"testing"
)

type testResponseWriter struct {
	code int
}

func (w *testResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *testResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (w *testResponseWriter) WriteHeader(c int) {
	w.code = c
}

func TestWriter(t *testing.T) {
	tw := &testResponseWriter{}
	w := &writer{tw, 0}
	w.WriteHeader(200)

	if w.status != 200 {
		t.Fatalf("Expected status 200 got %d", w.status)
	}
	if tw.code != 200 {
		t.Fatalf("Expected status 200 got %d", tw.code)
	}
}
