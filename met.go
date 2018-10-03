package httpx

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bsm/rucksack/met"
	"github.com/go-chi/chi/middleware"
)

var metMiddleware = middleware.RequestLogger(metFormatter{})

type metFormatter struct{}

func (metFormatter) NewLogEntry(r *http.Request) middleware.LogEntry { return metEntry{Request: r} }

type metEntry struct{ *http.Request }

func (e metEntry) Write(status, _ int, elapsed time.Duration) {
	tags := append(make([]string, 0, 2),
		"status:"+strconv.Itoa(status),
	)
	if status != 404 {
		tags = append(tags, "path:"+e.URL.Path)
	}

	met.RatePerMin("http.request", tags).Update(1)
	met.Timer("http.request.time", tags).Update(elapsed)
}

func (e metEntry) Panic(_ interface{}, _ []byte) {
	tags := []string{"path:" + e.URL.Path}
	met.RatePerMin("http.request.panic", tags).Update(1)
}
