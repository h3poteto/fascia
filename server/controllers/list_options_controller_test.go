package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/views"

	"encoding/json"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ListOptionsController", func() {
	var (
		e      *echo.Echo
		rec    *httptest.ResponseRecorder
		userID int64
	)
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		userID = LoginFaker("list_options@example.com", "hogehoge")
		seed.Seeds()
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			c := e.NewContext(new(http.Request), rec)
			c.SetPath("/list_options")
			resource := ListOptions{}
			err := resource.Index(c)
			Expect(err).To(BeNil())
			var contents []views.ListOption
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents[0].Action).To(Equal("close"))
			Expect(contents[1].Action).To(Equal("open"))
		})
	})
})
