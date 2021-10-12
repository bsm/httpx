package httpx

import "os"

var isTestMode bool

func coalesce(vv ...string) string {
	for _, v := range vv {
		if v != "" {
			return v
		}
	}
	return ""
}

func fromEnv(keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return ""
}
