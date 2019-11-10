package controllers_test

import (
	"database/sql"

	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/lib/modules/database"
	. "github.com/h3poteto/fascia/server/controllers"
	userRepo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/views"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
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
			db := database.SharedInstance().Connection
			repo := userRepo.New(db)
			repo.Create(
				"list_options@example.com",
				"hogehoge",
				sql.NullString{},
				sql.NullString{},
				sql.NullInt64{},
				sql.NullString{},
				sql.NullString{})
		})
		It("should return", func() {
			u, _ := account.FindUserByEmail("list_options@example.com")
			req := httptest.NewRequest(echo.GET, "/list_options", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, u)
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
