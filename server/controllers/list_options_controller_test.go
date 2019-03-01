package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	usecase "github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/views"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListOptionsController", func() {
	var (
		e   *echo.Echo
		rec *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			usecase.RegistrationUser("list_options@example.com", "hogehoge", "hogehoge")
		})
		It("should return", func() {
			req := httptest.NewRequest(echo.GET, "/list_options", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, "list_options@example.com", "hogehoge")
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
