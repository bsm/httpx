package httpx

import (
	"math"
	"net/http"
	"time"

	"github.com/bsm/rucksack/log"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logMiddleware = middleware.RequestLogger(logFormatter{})

type logFormatter struct{}

func (logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry { return logEntry{Request: r} }

type logEntry struct{ *http.Request }

func (e logEntry) Write(status, bytes int, elapsed time.Duration) {
	fields := append(make([]zapcore.Field, 0, 7),
		zap.Int("status", status),
		zap.String("method", e.Method),
		zap.String("uri", e.RequestURI),
		zap.String("remote", e.RemoteAddr),
		zap.Float64("kB", math.Ceil(float64(bytes)/102.4)/10),
		zap.Duration("elapsed", elapsed),
	)
	if reqID := e.Header.Get("X-Request-Id"); reqID != "" {
		fields = append(fields, zap.String("request_id", reqID))
	}

	if status < 500 {
		log.Infow(e.RequestURI, fields...)
	} else {
		log.Warnw(e.RequestURI, fields...)
	}
}

func (e logEntry) Panic(v interface{}, stack []byte) {
	log.Errorw("panic", zap.Any("msg", v), zap.ByteString("stack", stack))
}
