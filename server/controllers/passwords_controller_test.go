package controllers_test

import (
	"github.com/h3poteto/fascia/server"
	. "github.com/h3poteto/fascia/server/controllers"
	resetpassword "github.com/h3poteto/fascia/server/domains/reset_password"
	"github.com/h3poteto/fascia/server/usecases/account"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PasswordsController", func() {
	var (
		e        *echo.Echo
		rec      *httptest.ResponseRecorder
		email    string
		password string
		uid      int64
	)
	BeforeEach(func() {
		e = echo.New()
		e.Renderer = server.PongoRenderer()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		email = "hoge@example.com"
		password = "hogehoge"
		user, _ := account.RegistrationUser(email, password, password)
		uid = user.ID
		GenerateCSRFToken = func(c echo.Context) (string, error) { return "hoge", nil }
	})

	Describe("New", func() {
		It("should correctly access", func() {
			req := httptest.NewRequest(echo.GET, "/passwords/new", nil)
			c := e.NewContext(req, rec)
			resource := Passwords{}
			err := resource.New(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).NotTo(BeNil())
		})
	})

	Describe("Create", func() {
		JustBeforeEach(func() {
			CheckCSRFToken = func(c echo.Context, token string) bool { return true }
		})
		It("should create new reset password", func() {
			f := make(url.Values)
			f.Set("email", "hogehoge@example.com")
			req := httptest.NewRequest(echo.POST, "/passwords/create", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			resource := Passwords{}
			err := resource.Create(c)
			Expect(err).To(BeNil())
			u, _ := rec.Result().Location()
			Expect(u.Path).To(Equal("/sign_in"))
		})
	})

	Describe("Edit", func() {
		var resetPassword *resetpassword.ResetPassword
		JustBeforeEach(func() {
			resetPassword, _ = account.GenerateResetPassword(uid, email)
		})
		Context("token is invalid", func() {
			It("should internal server error", func() {
				q := make(url.Values)
				q.Set("token", "sample")
				req := httptest.NewRequest(echo.GET, "/passwords/:id/edit?"+q.Encode(), nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(strconv.FormatInt(resetPassword.ID, 10))
				resource := Passwords{}
				err := resource.Edit(c)
				Expect(err).NotTo(BeNil())
			})
		})
		Context("token is correct", func() {
			It("should response is ok", func() {
				q := make(url.Values)
				q.Set("token", resetPassword.Token)
				req := httptest.NewRequest(echo.GET, "/passwords/:id/edit?"+q.Encode(), nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(strconv.FormatInt(resetPassword.ID, 10))
				resource := Passwords{}
				err := resource.Edit(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Update", func() {
		var resetPassword *resetpassword.ResetPassword
		JustBeforeEach(func() {
			CheckCSRFToken = func(c echo.Context, token string) bool { return true }
			resetPassword, _ = account.GenerateResetPassword(uid, email)
		})
		Context("token is invalid", func() {
			It("should internal server error", func() {
				f := make(url.Values)
				f.Set("password", "fugafuga")
				f.Set("password_confirm", "fugafuga")
				f.Set("reset_token", "sample")
				req := httptest.NewRequest(echo.POST, "/passwords/:id/update", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(strconv.FormatInt(resetPassword.ID, 10))
				resource := Passwords{}
				err := resource.Update(c)
				fmt.Println(err)
				Expect(err).NotTo(BeNil())
			})
		})
		Context("token is correct", func() {
			It("should response is ok", func() {
				f := make(url.Values)
				f.Set("password", "fugafuga")
				f.Set("password_confirm", "fugafuga")
				f.Set("reset_token", resetPassword.Token)
				req := httptest.NewRequest(echo.POST, "/passwords/:id/update", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				c.SetParamNames("id")
				c.SetParamValues(strconv.FormatInt(resetPassword.ID, 10))
				resource := Passwords{}
				err := resource.Update(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusFound))
				u, _ := rec.Result().Location()
				Expect(u.Path).To(Equal("/sign_in"))
			})
		})
	})
})
