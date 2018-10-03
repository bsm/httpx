package httpx

import (
	"net/http"

	"github.com/bsm/rucksack"
	"github.com/kr/secureheader"
)

var secureheaderDefaults secureheader.Config

func init() {
	secureheaderDefaults = *secureheader.DefaultConfig

	if s := rucksack.Env("SECURE_HEADER_HTTPS_REDIRECT"); s == "false" {
		secureheaderDefaults.HTTPSRedirect = false
	}
	if s := rucksack.Env("SECURE_HEADER_HTTPS_USE_FORWARDED_PROTO"); s != "false" {
		secureheaderDefaults.HTTPSUseForwardedProto = true
	}
}

func secureHeader(cfg *secureheader.Config) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		c := secureheaderDefaults
		if cfg != nil {
			c = *cfg
		}
		c.Next = h
		return &c
	}
}
