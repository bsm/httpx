package httpx

import (
	"io"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

// RouterOptions support custom router configuration.
type RouterOptions struct {
	Logger     middleware.LogFormatter
	Secure     *secure.Options
	CORS       *cors.Options
	Heartbeat  string // heartbeat path, set to "false" to disable
	NoCompress bool   // disable compression
}

func (o *RouterOptions) norm() {
	if o.Logger == nil {
		var out io.Writer = os.Stdout
		if isTestMode {
			out = io.Discard
		}
		o.Logger = &middleware.DefaultLogFormatter{Logger: log.New(out, "", log.LstdFlags)}
	}

	if o.Secure == nil {
		opts := secureDefaults
		o.Secure = &opts
	}
	if isTestMode {
		o.Secure.IsDevelopment = false
	}

	if o.CORS == nil {
		opts := corsDefaults
		o.CORS = &opts
	}

	if o.Heartbeat == "" {
		o.Heartbeat = coalesce(fromEnv("HTTP_HEARTBEAT"), "/ping")
	}

	o.NoCompress = fromEnv("HTTP_COMPRESS") == "false"
}

// NewRouter inits a new chi.Router with options
func NewRouter(opt *RouterOptions) chi.Router {
	var o RouterOptions
	if opt != nil {
		o = *opt
	}
	o.norm()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(o.Logger))
	r.Use(middleware.Recoverer)

	if o.Heartbeat != "false" {
		r.Use(middleware.Heartbeat(o.Heartbeat))
	}

	r.Use(secure.New(*o.Secure).Handler)
	r.Use(cors.New(*o.CORS).Handler)

	if !o.NoCompress {
		r.Use(middleware.Compress(2))
	}
	return r
}
