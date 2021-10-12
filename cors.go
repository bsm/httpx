package httpx

import (
	"strconv"
	"strings"

	"github.com/go-chi/cors"
)

func corsDefaults(_ Env) *cors.Options {
	opts := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "CONNECT", "TRACE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	if s := fromEnv("CORS_ALLOWED_ORIGINS", "CORS_ORIGINS"); s != "" {
		opts.AllowedOrigins = strings.Split(s, ",")
	}
	if s := fromEnv("CORS_ALLOWED_METHODS", "CORS_METHODS"); s != "" {
		opts.AllowedMethods = strings.Split(s, ",")
	}
	if s := fromEnv("CORS_ALLOWED_HEADERS", "CORS_HEADERS"); s != "" {
		opts.AllowedHeaders = strings.Split(s, ",")
	}
	if s := fromEnv("CORS_ALLOW_CREDENTIALS"); s == "false" {
		opts.AllowCredentials = false
	}
	if s := fromEnv("CORS_MAX_AGE"); s != "" {
		if n, err := strconv.Atoi(s); err == nil {
			opts.MaxAge = n
		}
	}
	return &opts
}
