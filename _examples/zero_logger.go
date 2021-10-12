package main

import (
	"net"
	"net/http"
	"time"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type zerologEntry struct{ Method, URL, Proto, RemoteAddr, UA, Referer, RequestID string }

// From https://docs.datadoghq.com/logs/log_configuration/attributes_naming_convention/#http-requests
func (e *zerologEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	level := zerolog.InfoLevel
	if status >= 500 {
		level = zerolog.WarnLevel
	}

	ent := log.WithLevel(level).
		Str("http.method", e.Method).
		Str("http.url", e.URL).
		Str("http.version", e.Proto).
		Str("http.useragent", e.UA).
		Int("http.status_code", status).
		Int("network.bytes_written", bytes).
		Dur("duration", elapsed)
	if e.RequestID != "" {
		ent = ent.Str("http.request_id", e.RequestID)
	}
	if e.Referer != "" {
		ent = ent.Str("http.referer", e.Referer)
	}
	if e.RemoteAddr != "" {
		host, port, _ := net.SplitHostPort(e.RemoteAddr)
		ent = ent.
			Str("network.client.ip", host).
			Str("network.client.port", port)
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
		Method:     r.Method,
		URL:        scheme + "://" + r.Host + r.RequestURI,
		Proto:      r.Proto,
		RemoteAddr: r.RemoteAddr,
		UA:         r.UserAgent(),
		Referer:    r.Header.Get("Referer"),
		RequestID:  r.Header.Get(middleware.RequestIDHeader),
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
