package controllers_test

import (
	"github.com/h3poteto/fascia/server"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"

	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SessionsController", func() {
	var (
		e   *echo.Echo
		rec *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		e = echo.New()
		e.Renderer = server.PongoRenderer()
		rec = httptest.NewRecorder()
	})

	Describe("SignIn", func() {
		Context("/sign_in", func() {
			It("should correctly access", func() {
				c := e.NewContext(new(http.Request), rec)
				c.SetPath("/sign_in")
				resource := Sessions{}
				err := resource.SignIn(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(rec.Body.Len()).NotTo(Equal(0))
			})
		})
		Context("/", func() {
			It("should redirect to top page", func() {
				c := e.NewContext(new(http.Request), rec)
				c.SetPath("/")
				resource := Root{}
				err := resource.Index(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(200))
				topDoc, err := goquery.NewDocumentFromReader(strings.NewReader(rec.Body.String()))
				Expect(err).To(BeNil())
				c.SetPath("/about")
				resource.About(c)
				aboutDoc, err := goquery.NewDocumentFromReader(strings.NewReader(rec.Body.String()))
				Expect(err).To(BeNil())
				topDoc.Selection.Find("h1").Each(func(_ int, s *goquery.Selection) {
					aboutDoc.Selection.Find("h1").Each(func(_ int, as *goquery.Selection) {
						Expect(as.Text()).To(Equal(s.Text()))
					})
				})
			})
		})
	})

	Describe("NewSession", func() {
		JustBeforeEach(func() {
			CSRFFaker()
		})
		Context("before registration", func() {
			It("should not login", func() {
				f := make(url.Values)
				f.Set("email", "sign_in@example.com")
				f.Set("password", "hogehoge")
				req, _ := http.NewRequest(echo.POST, "/sign_in", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				resource := Sessions{}
				err := resource.NewSession(c)
				Expect(err).To(BeNil())
				u, _ := rec.Result().Location()
				Expect(u.Path).To(Equal("/sign_in"))
			})
		})
		Context("after registration", func() {
			JustBeforeEach(func() {
				handlers.RegistrationUser("registration@example.com", "hogehoge", "hogehoge")
			})
			Context("when use correctly password", func() {
				It("can login", func() {
					f := make(url.Values)
					f.Set("email", "registration@example.com")
					f.Set("password", "hogehoge")
					req, _ := http.NewRequest(echo.POST, "/sign_in", strings.NewReader(f.Encode()))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
					c := e.NewContext(req, rec)
					resource := Sessions{}
					err := resource.NewSession(c)
					Expect(err).To(BeNil())
					u, _ := rec.Result().Location()
					Expect(u.Path).To(Equal("/"))
				})
			})
			Context("when use wrong password", func() {
				It("cannot login", func() {
					f := make(url.Values)
					f.Set("email", "registration@example.com")
					f.Set("password", "fugafuga")
					req, _ := http.NewRequest(echo.POST, "/sign_in", strings.NewReader(f.Encode()))
					req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
					c := e.NewContext(req, rec)
					resource := Sessions{}
					err := resource.NewSession(c)
					Expect(err).To(BeNil())
					Expect(rec.Code).To(Equal(http.StatusFound))
					u, _ := rec.Result().Location()
					Expect(u.Path).To(Equal("/sign_in"))
				})
			})
		})
	})

	Describe("SignOut", func() {
		It("can logout", func() {
			req, _ := http.NewRequest(echo.POST, "/sign_out", nil)
			c := e.NewContext(req, rec)
			resource := Sessions{}
			err := resource.SignOut(c)
			Expect(err).To(BeNil())
			u, _ := rec.Result().Location()
			Expect(u.Path).To(Equal("/sign_in"))
		})
	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			handlers.RegistrationUser("update@example.com", "hogehoge", "hogehoge")
		})
		It("can update session", func() {
			req, _ := http.NewRequest(echo.POST, "/update", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, "update@example.com", "hogehoge")
			resource := Sessions{}
			err := resource.Update(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})
})
