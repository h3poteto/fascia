package controllers_test

import (
	"github.com/h3poteto/fascia/server"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"

	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegistrationsController", func() {
	var (
		e        *echo.Echo
		rec      *httptest.ResponseRecorder
		database *sql.DB
	)
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		database = db.SharedInstance().Connection
	})

	Describe("SignUp", func() {
		JustBeforeEach(func() {
			GenerateCSRFToken = func(c echo.Context) (string, error) { return "hoge", nil }
			e.Renderer = server.PongoRenderer()
		})
		It("should correctly access", func() {
			req := httptest.NewRequest(echo.GET, "/sign_up", nil)
			c := e.NewContext(req, rec)
			resource := Registrations{}
			err := resource.SignUp(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).NotTo(BeNil())
		})

	})

	Describe("Registration", func() {
		JustBeforeEach(func() {
			CheckCSRFToken = func(c echo.Context, token string) bool { return true }
		})

		Context("パスワードと確認パスワードが一致しているとき", func() {
			It("登録できること", func() {
				f := make(url.Values)
				f.Set("email", "registration@example.com")
				f.Set("password", "hogehoge")
				f.Set("password_confirm", "hogehoge")
				req := httptest.NewRequest(echo.POST, "/sign_up", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				resource := Registrations{}
				err := resource.Registration(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusFound))
				u, _ := rec.Result().Location()
				Expect(u.Path).To(Equal("/sign_in"))
				var id int64
				rows, _ := database.Query("select id from users where email = ?;", "registration@example.com")
				for rows.Next() {
					err := rows.Scan(&id)
					if err != nil {
						panic(err)
					}
				}
				Expect(id).NotTo(Equal(0))
			})
		})
		Context("既に登録されているとき", func() {
			JustBeforeEach(func() {
				handlers.RegistrationUser("registration@example.com", "hogehoge", "hogehoge")
			})
			It("エラーになること", func() {
				f := make(url.Values)
				f.Set("email", "registration@example.com")
				f.Set("password", "hogehoge")
				f.Set("password_confirm", "hogehoge")
				req := httptest.NewRequest(echo.POST, "/sign_up", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				resource := Registrations{}
				err := resource.Registration(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusFound))
				u, _ := rec.Result().Location()
				Expect(u.Path).To(Equal("/sign_up"))
			})
		})
	})
})
