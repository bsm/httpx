package httpx

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

// MuxOptions allow mux configuration.
type MuxOptions struct {
	Env        Env
	Logger     middleware.LogFormatter
	Secure     *secure.Options
	CORS       *cors.Options
	Heartbeat  string // heartbeat path, set to "false" to disable
	NoCompress bool   // disable compression
}

func (o *MuxOptions) norm() {
	if o.Env == unknownEnv {
		o.Env = guessedEnv
	}

	if o.Logger == nil {
		o.Logger = &middleware.DefaultLogFormatter{Logger: newStdLogger(o.Env)}
	}

	if o.Secure == nil {
		o.Secure = secureDefaults(o.Env)
	}

	if o.CORS == nil {
		o.CORS = corsDefaults(o.Env)
	}

	if o.Heartbeat == "" {
		o.Heartbeat = coalesce(fromEnv("HTTP_HEARTBEAT"), "/ping")
	}

	o.NoCompress = fromEnv("HTTP_COMPRESS") == "false"
}

// NewMux inits a new *chi.Mux with options
func NewMux(opt *MuxOptions) *chi.Mux {
	var o MuxOptions
	if opt != nil {
		o = *opt
	}
	o.norm()

	r := chi.NewMux()
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
