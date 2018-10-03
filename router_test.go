package httpx

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Router", func() {
	var subject chi.Router

	BeforeEach(func() {
		subject = NewRouter(nil)
		subject.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("boom"))
		})
	})

	It("should init new router with defaults", func() {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Add("X-Forwarded-Proto", "https")

		subject.ServeHTTP(w, r)
		Expect(w.Code).To(Equal(200))
		Expect(w.Body.String()).To(Equal("boom"))
		Expect(w.Header()).To(Equal(http.Header{
			"Strict-Transport-Security": {"max-age=25920000; includeSubDomains"},
			"Vary": {"Origin"},
			"X-Content-Type-Options": {"nosniff"},
			"X-Frame-Options":        {"DENY"},
			"X-Xss-Protection":       {"1; mode=block"},
		}))
	})

	It("should mount heartbeat", func() {
		w := httptest.NewRecorder()
		subject.ServeHTTP(w, httptest.NewRequest("GET", "/ping", nil))
		Expect(w.Code).To(Equal(200))
		Expect(w.Body.String()).To(Equal("."))
		Expect(w.Header()).To(Equal(http.Header{
			"Content-Type": {"text/plain"},
		}))
	})
})
