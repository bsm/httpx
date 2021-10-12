package httpx

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	isTestMode = true
	os.Exit(m.Run())
}
