package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger is a custom logger.
func Logger(handler func(*http.Request) middleware.LogEntry) middleware.LogFormatter {
	return logFormatter(handler)
}

type logFormatter func(*http.Request) middleware.LogEntry

func (fn logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry { return fn(r) }
