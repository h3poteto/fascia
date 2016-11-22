package controllers_test

import (
	"github.com/h3poteto/fascia/controllers"
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/reset_password"
	"github.com/h3poteto/fascia/models/user"
	. "github.com/h3poteto/fascia/server"

	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("PasswordsController", func() {
	var (
		ts       *httptest.Server
		email    string
		password string
		uid      int64
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
		database.Exec("truncate table reset_passwords;")
	})
	JustBeforeEach(func() {
		email = "hoge@example.com"
		password = "hogehoge"
		uid, _ = user.Registration(email, password, password)
	})

	Describe("New", func() {
		It("should correctly access", func() {
			res, err := http.Get(ts.URL + "/passwords/new")
			Expect(err).To(BeNil())
			contents, status := ParseResponse(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
		})
	})

	Describe("Create", func() {
		JustBeforeEach(func() {
			controllers.CheckCSRFToken = func(r *http.Request, token string) bool { return true }
		})
		It("should create new reset password", func() {
			values := url.Values{}
			values.Add("email", "hogehoge@example.com")
			res, err := http.PostForm(ts.URL+"/passwords/create", values)
			Expect(err).To(BeNil())
			Expect(res.Request.URL.Path).To(Equal("/sign_in"))
		})
	})

	Describe("Edit", func() {
		var resetPassword *reset_password.ResetPasswordStruct
		JustBeforeEach(func() {
			resetPassword = reset_password.GenerateResetPassword(uid, email)
			resetPassword.Save()
		})
		Context("token is invalid", func() {
			It("should internal server error", func() {
				res, err := http.Get(ts.URL + "/passwords/" + strconv.FormatInt(resetPassword.ID, 10) + "/edit?token=sample")
				Expect(err).To(BeNil())
				doc, _ := goquery.NewDocumentFromResponse(res)
				doc.Find("h2").Each(func(_ int, s *goquery.Selection) {
					Expect(s.Text()).To(Equal("Internal Server Error."))
				})
			})
		})
		Context("token is correct", func() {
			It("should response is ok", func() {
				res, err := http.Get(ts.URL + "/passwords/" + strconv.FormatInt(resetPassword.ID, 10) + "/edit?token=" + resetPassword.Token)
				Expect(err).To(BeNil())
				_, status := ParseResponse(res)
				Expect(status).To(Equal(http.StatusOK))
			})
		})
	})

	Describe("Update", func() {
		var resetPassword *reset_password.ResetPasswordStruct
		JustBeforeEach(func() {
			controllers.CheckCSRFToken = func(r *http.Request, token string) bool { return true }
			resetPassword = reset_password.GenerateResetPassword(uid, email)
			resetPassword.Save()
		})
		Context("token is invalid", func() {
			It("should internal server error", func() {
				values := url.Values{}
				values.Add("password", "fugafuga")
				values.Add("password_confirm", "fugafuga")
				values.Add("reset_token", "sample")
				res, err := http.PostForm(ts.URL+"/passwords/"+strconv.FormatInt(resetPassword.ID, 10)+"/update", values)
				Expect(err).To(BeNil())
				doc, _ := goquery.NewDocumentFromResponse(res)
				doc.Find("h2").Each(func(_ int, s *goquery.Selection) {
					Expect(s.Text()).To(Equal("Internal Server Error."))
				})
			})
		})
		Context("token is correct", func() {
			It("should response is ok", func() {
				values := url.Values{}
				values.Add("password", "fugafuga")
				values.Add("password_confirm", "fugafuga")
				values.Add("reset_token", resetPassword.Token)
				res, err := http.PostForm(ts.URL+"/passwords/"+strconv.FormatInt(resetPassword.ID, 10)+"/update", values)
				Expect(err).To(BeNil())
				_, status := ParseResponse(res)
				Expect(status).To(Equal(http.StatusOK))
				Expect(res.Request.URL.Path).To(Equal("/sign_in"))
			})
		})
	})
})
