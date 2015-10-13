package controllers_test

import (
	"net/http"
	"net/http/httptest"
	. "../../fascia"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("SessionsController", func() {
	var (
		ts *httptest.Server
	)
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
	})
	AfterEach(func() {
		ts.Close()
	})
	Context("/sign_in", func() {
		It("アクセスできること", func() {
			res, err := http.Get(ts.URL + "/sign_in")
			Expect(err).To(BeNil())
			contents, status := ParseResponse(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
		})
	})
	Context("/", func() {
		It("リダイレクトされること", func() {
			res, err := http.Get(ts.URL + "/")
			Expect(err).To(BeNil())
			Expect(res.Request.URL.Path).To(Equal("/sign_in"))
		})
	})
})
