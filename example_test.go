package httpx_test

import (
	"net/http"
	"time"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5/middleware"
)

type myLogEntry struct{ Method, URL, Proto, From, UA string }

func (e *myLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	// instrument request
}

func (e *myLogEntry) Panic(v interface{}, stack []byte) {
	// instrument panic
}

func ExampleLogger() {
	handler := func(r *http.Request) middleware.LogEntry {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		return &myLogEntry{
			Method: r.Method,
			URL:    scheme + "://" + r.Host + r.RequestURI,
			Proto:  r.Proto,
			From:   r.RemoteAddr,
			UA:     r.UserAgent(),
		}
	}

	httpx.NewRouter(&httpx.RouterOptions{
		Logger: httpx.Logger(handler),
	})
}
