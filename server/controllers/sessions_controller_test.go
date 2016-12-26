package controllers_test

import (
	. "github.com/h3poteto/fascia/server"
	"github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/models/db"
	"github.com/h3poteto/fascia/server/services"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/PuerkitoBio/goquery"
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
		database := db.SharedInstance().Connection
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
	})
	Describe("SignIn", func() {
		JustBeforeEach(func() {
			controllers.LoginRequired = controllers.CheckLogin
			values := url.Values{}
			http.PostForm(ts.URL+"/sign_out", values)
		})
		Context("/sign_in", func() {
			It("should correctly access", func() {
				res, err := http.Get(ts.URL + "/sign_in")
				Expect(err).To(BeNil())
				contents, status := ParseResponse(res)
				Expect(status).To(Equal(http.StatusOK))
				Expect(contents).NotTo(BeNil())
			})
		})
		Context("/", func() {
			It("should redirect to top page", func() {
				res, err := http.Get(ts.URL + "/")
				Expect(err).To(BeNil())
				Expect(res.StatusCode).To(Equal(200))
				Expect(res.Request.URL.Path).To(Equal("/"))
				topDoc, _ := goquery.NewDocumentFromResponse(res)
				aboutDoc, _ := goquery.NewDocument(ts.URL + "/about")
				topDoc.Find("h1").Each(func(_ int, s *goquery.Selection) {
					aboutDoc.Find("h1").Each(func(_ int, as *goquery.Selection) {
						Expect(as.Text()).To(Equal(s.Text()))
					})
				})
			})
		})
	})

	Describe("NewSession", func() {
		JustBeforeEach(func() {
			controllers.CheckCSRFToken = func(r *http.Request, token string) bool { return true }
		})
		Context("before registration", func() {
			It("should not login", func() {
				values := url.Values{}
				values.Add("email", "sign_in@example.com")
				values.Add("password", "hogehoge")
				res, err := http.PostForm(ts.URL+"/sign_in", values)
				Expect(err).To(BeNil())
				Expect(res.Request.URL.Path).To(Equal("/sign_in"))
			})
		})
		Context("after registration", func() {
			JustBeforeEach(func() {
				user, _ := handlers.RegistrationUser("registration@example.com", "hogehoge", "hogehoge")
				controllers.LoginRequired = func(r *http.Request) (*services.User, error) {
					return user, nil
				}
			})
			Context("when use correctly password", func() {
				It("can login", func() {
					values := url.Values{}
					values.Add("email", "registration@example.com")
					values.Add("password", "hogehoge")
					res, err := http.PostForm(ts.URL+"/sign_in", values)
					Expect(err).To(BeNil())
					Expect(res.Request.URL.Path).To(Equal("/"))
				})
			})
			Context("when use wrong password", func() {
				It("cannot login", func() {
					values := url.Values{}
					values.Add("email", "registration@example.com")
					values.Add("password", "fugafuga")
					res, err := http.PostForm(ts.URL+"/sign_in", values)
					Expect(err).To(BeNil())
					Expect(res.Request.URL.Path).To(Equal("/sign_in"))
				})
			})
		})
	})

	Describe("SignOut", func() {
		JustBeforeEach(func() {
			LoginFaker(ts, "sign_out@example.com", "hogehoge")
		})
		It("can logout", func() {
			values := url.Values{}
			res, err := http.PostForm(ts.URL+"/sign_out", values)
			Expect(err).To(BeNil())
			Expect(res.Request.URL.Path).To(Equal("/sign_in"))
		})
	})

	Describe("Update", func() {
		JustBeforeEach(func() {
			LoginFaker(ts, "update@example.com", "hogehoge")
		})
		It("can update session", func() {
			values := url.Values{}
			res, err := http.PostForm(ts.URL+"/session", values)
			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})
	})
})
