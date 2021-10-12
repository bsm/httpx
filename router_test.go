package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5"
)

func seedMux() chi.Router {
	mux := httpx.NewRouter(nil)
	mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("boom"))
	})
	return mux
}

func TestNewRouter_defaults(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	r.Header.Add("X-Forwarded-Proto", "https")

	mux := seedMux()
	mux.ServeHTTP(w, r)
	if exp, got := 200, w.Code; exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := "boom", w.Body.String(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}

	exp := http.Header{
		"Strict-Transport-Security": {"max-age=31104000; includeSubDomains; preload"},
		"Vary":                      {"Origin"},
		"X-Content-Type-Options":    {"nosniff"},
		"X-Frame-Options":           {"SAMEORIGIN"},
		"X-Xss-Protection":          {"1; mode=block"},
	}
	if got := w.Header(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %+v, got %+v", exp, got)
	}
}

func TestNewRouter_heartbeat(t *testing.T) {
	w := httptest.NewRecorder()

	mux := seedMux()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/ping", nil))

	if exp, got := 200, w.Code; exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := ".", w.Body.String(); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
	if exp, got := "text/plain", w.Header().Get("Content-Type"); exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}
