package main

import (
	"net/http"
	"time"

	"github.com/bsm/httpx"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapEntry struct {
	Method, URL, Proto, From, UA, RequestID string

	Logger *zap.Logger
}

func (e *zapEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	fields := append(make([]zapcore.Field, 0, 8),
		zap.String("method", e.Method),
		zap.String("url", e.URL),
		zap.String("proto", e.Proto),
		zap.String("ua", e.UA),
		zap.Int("status", status),
		zap.Int("bytes", bytes),
		zap.Duration("elapsed", elapsed),
	)
	if e.RequestID != "" {
		fields = append(fields, zap.String("request_id", e.RequestID))
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
			Method:    r.Method,
			URL:       scheme + "://" + r.Host + r.RequestURI,
			Proto:     r.Proto,
			From:      r.RemoteAddr,
			UA:        r.UserAgent(),
			RequestID: r.Header.Get(middleware.RequestIDHeader),
			Logger:    logger,
		}
	}

	mux := httpx.NewMux(&httpx.MuxOptions{
		Logger: httpx.Logger(zapper),
	})
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Panic(err.Error())
	}
}
