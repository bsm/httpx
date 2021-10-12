package httpx

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5/middleware"
)

// Logger is a custom logger.
func Logger(handler func(*http.Request) middleware.LogEntry) middleware.LogFormatter {
	return logFormatter(handler)
}

func newStdLogger(env Env) *log.Logger {
	var out io.Writer = os.Stdout
	if env == Test {
		out = io.Discard
	}
	return log.New(out, "", log.LstdFlags)
}

type logFormatter func(*http.Request) middleware.LogEntry

func (fn logFormatter) NewLogEntry(r *http.Request) middleware.LogEntry { return fn(r) }
