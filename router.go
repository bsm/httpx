package httpx

import (
	"github.com/bsm/rucksack"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kr/secureheader"
	"github.com/rs/cors"
)

// RouterOptions support custom router configuration.
type RouterOptions struct {
	CORS         *cors.Options
	SecureHeader *secureheader.Config

	Heartbeat       string // heartbeat path, set to "false" to disable
	DisableLog      bool   // disable log instrumentation
	DisableMet      bool   // disable met instrumentation
	DisableCompress bool   // disable compression
}

func (o *RouterOptions) norm() {
	if o.CORS == nil {
		o.CORS = &corsDefaults
	}
	if o.SecureHeader == nil {
		o.SecureHeader = &secureheaderDefaults
	}
	if o.Heartbeat == "" {
		o.Heartbeat = coalesce(rucksack.Env("ROUTER_HEARTBEAT"), "/ping")
	}
	o.DisableLog = rucksack.Env("ROUTER_LOG") == "false"
	o.DisableMet = rucksack.Env("ROUTER_MET") == "false"
	o.DisableCompress = rucksack.Env("ROUTER_COMPRESS") == "false"
}

// NewRouter inits a new chi.Router with options
func NewRouter(opt *RouterOptions) chi.Router {
	var o RouterOptions
	if opt != nil {
		o = *opt
	}
	o.norm()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	if o.Heartbeat != "false" {
		r.Use(middleware.Heartbeat(o.Heartbeat))
	}
	if !o.DisableLog {
		r.Use(logMiddleware)
	}
	if !o.DisableMet {
		r.Use(metMiddleware)
	}

	r.Use(secureHeader(o.SecureHeader))
	r.Use(cors.New(*o.CORS).Handler)

	if !o.DisableCompress {
		r.Use(middleware.DefaultCompress)
	}
	return r
}
