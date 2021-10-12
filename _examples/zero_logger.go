package main

import (
	"net/http"
	"time"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type zerologEntry struct{ Method, URL, Proto, From, UA, RequestID string }

func (e *zerologEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	level := zerolog.InfoLevel
	if status >= 500 {
		level = zerolog.WarnLevel
	}

	ent := log.WithLevel(level).
		Str("method", e.Method).
		Str("url", e.URL).
		Str("proto", e.Proto).
		Str("from", e.From).
		Str("ua", e.UA).
		Int("status", status).
		Int("bytes", bytes).
		Dur("elapsed", elapsed)
	if e.RequestID != "" {
		ent = ent.Str("request_id", e.RequestID)
	}
	ent.Send()
}

func (e *zerologEntry) Panic(v interface{}, stack []byte) {
	log.Error().Interface("recovered", v).Bytes("stack", stack).Msg("panic")
}

func zerologFormatter(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return &zerologEntry{
		Method:    r.Method,
		URL:       scheme + "://" + r.Host + r.RequestURI,
		Proto:     r.Proto,
		From:      r.RemoteAddr,
		UA:        r.UserAgent(),
		RequestID: r.Header.Get(middleware.RequestIDHeader),
	}
}

func main() {
	mux := httpx.NewMux(&httpx.MuxOptions{
		Logger: httpx.Logger(zerologFormatter),
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panic().Msg(err.Error())
	}
}
