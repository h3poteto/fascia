package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server"
	"github.com/h3poteto/fascia/server/views"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ListOptionsController", func() {
	var (
		ts     *httptest.Server
		userID int64
	)
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
	})
	AfterEach(func() {
		ts.Close()
	})
	JustBeforeEach(func() {
		userID = LoginFaker(ts, "list_options@example.com", "hogehoge")
		seed.Seeds()
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			res, err := http.Get(ts.URL + "/list_options")
			Expect(err).To(BeNil())
			var contents []views.ListOption
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(contents[0].Action).To(Equal("close"))
			Expect(contents[1].Action).To(Equal("open"))
		})
	})
})
