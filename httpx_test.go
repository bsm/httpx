package httpx

import (
	"testing"

	"github.com/bsm/rucksack/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	log.Silence()
}

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "rucksack/httpx")
}
