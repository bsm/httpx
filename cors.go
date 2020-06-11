package httpx

import (
	"strings"

	"github.com/bsm/rucksack/v4"
	"github.com/go-chi/cors"
)

var corsDefaults = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"},
	AllowedHeaders:   []string{"*"},
	AllowCredentials: true,
}

func init() {
	if s := rucksack.Env("CORS_ALLOWED_ORIGINS", "CORS_ORIGINS"); s != "" {
		corsDefaults.AllowedOrigins = strings.Split(s, ",")
	}
	if s := rucksack.Env("CORS_ALLOWED_METHODS", "CORS_METHODS"); s != "" {
		corsDefaults.AllowedMethods = strings.Split(s, ",")
	}
	if s := rucksack.Env("CORS_ALLOWED_HEADERS", "CORS_HEADERS"); s != "" {
		corsDefaults.AllowedHeaders = strings.Split(s, ",")
	}
	if s := rucksack.Env("CORS_ALLOW_CREDENTIALS"); s == "false" {
		corsDefaults.AllowCredentials = false
	}
}
