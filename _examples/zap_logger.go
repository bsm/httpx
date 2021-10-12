package main

import (
	"net"
	"net/http"
	"time"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapEntry struct {
	Method, URL, Proto, RemoteAddr, UA, Referer, RequestID string

	Logger *zap.Logger
}

// From https://docs.datadoghq.com/logs/log_configuration/attributes_naming_convention/#http-requests
func (e *zapEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	fields := append(make([]zapcore.Field, 0, 11),
		zap.String("http.method", e.Method),
		zap.String("http.url", e.URL),
		zap.String("http.version", e.Proto),
		zap.String("http.useragent", e.UA),
		zap.Int("http.status_code", status),
		zap.Int("network.bytes_written", bytes),
		zap.Duration("duration", elapsed),
	)
	if e.RequestID != "" {
		fields = append(fields, zap.String("http.request_id", e.RequestID))
	}
	if e.Referer != "" {
		fields = append(fields, zap.String("http.referer", e.Referer))
	}
	if e.RemoteAddr != "" {
		host, port, _ := net.SplitHostPort(e.RemoteAddr)
		fields = append(fields,
			zap.String("network.client.ip", host),
			zap.String("network.client.port", port),
		)
	}

	if status < 500 {
		e.Logger.Info("", fields...)
	} else {
		e.Logger.Warn("", fields...)
	}
}

func (e *zapEntry) Panic(v interface{}, stack []byte) {
	e.Logger.Error("panic", zap.Any("recovered", v), zap.ByteString("stack", stack))
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()

	zapper := func(r *http.Request) middleware.LogEntry {
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		return &zapEntry{
			Method:     r.Method,
			URL:        scheme + "://" + r.Host + r.RequestURI,
			Proto:      r.Proto,
			RemoteAddr: r.RemoteAddr,
			UA:         r.UserAgent(),
			Referer:    r.Header.Get("Referer"),
			RequestID:  r.Header.Get(middleware.RequestIDHeader),
			Logger:     logger,
		}
	}

	mux := httpx.NewMux(&httpx.MuxOptions{
		Logger: httpx.Logger(zapper),
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Panic(err.Error())
	}
}
