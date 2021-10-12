package httpx

import (
	"os"
	"strings"
)

// Env defines the runtime environment
type Env uint8

const (
	unknownEnv Env = iota
	Production
	Test
	Development
)

var guessedEnv = Production

func init() {
	if fromEnv("CI") != "" || (len(os.Args) != 0 && strings.HasSuffix(os.Args[0], ".test")) {
		guessedEnv = Test
	} else if workDir, _ := os.Getwd(); strings.HasPrefix(workDir, "/home/") {
		guessedEnv = Development
	}
}

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
