package controllers_test

import (
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	. "../../fascia"
	. "../controllers"
	"../models/db"
	"../models/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("SessionsController", func() {
	var (
		ts *httptest.Server
		currentdb string
	)
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
		testdb := os.Getenv("DB_TEST_NAME")
		currentdb = os.Getenv("DB_NAME")
		os.Setenv("DB_NAME", testdb)
	})
	AfterEach(func() {
		ts.Close()
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Exec("truncate table projects;")
		table.Close()
		os.Setenv("DB_NAME", currentdb)
	})
	Describe("SignIn", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			http.PostForm(ts.URL + "/sign_out", values)
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
		// これの前にログアウト処理をしておかないと
		Context("/", func() {
			It("リダイレクトされること", func() {
				res, err := http.Get(ts.URL + "/")
				Expect(err).To(BeNil())
				Expect(res.StatusCode).To(Equal(200))
				Expect(res.Request.URL.Path).To(Equal("/sign_in"))
			})
		})
	})

	Describe("NewSession", func() {
		JustBeforeEach(func() {
			CheckCSRFToken = func(r *http.Request, token string) bool { return true }
		})
		Context("未登録のとき", func() {
			It("ログインできないこと", func() {
				values := url.Values{}
				values.Add("email", "sign_in@example.com")
				values.Add("password", "hogehoge")
				res, err := http.PostForm(ts.URL + "/sign_in", values)
				Expect(err).To(BeNil())
				Expect(res.Request.URL.Path).To(Equal("/sign_in"))
			})
		})
		Context("登録済みのとき", func() {
			JustBeforeEach(func() {
				id, _ := user.Registration("registration@example.com", "hogehoge")
				LoginRequired = func(r *http.Request) (*user.UserStruct, bool) {
					current_user, _ := user.CurrentUser(id)
					return current_user, true
				}
			})
			It("ログインできること", func() {
				values := url.Values{}
				values.Add("email", "registration@example.com")
				values.Add("password", "hogehoge")
				res, err := http.PostForm(ts.URL + "/sign_in", values)
				Expect(err).To(BeNil())
				fmt.Printf("request: %+v\n", res)
				Expect(res.Request.URL.Path).To(Equal("/"))
			})
		})
	})

	Describe("SignOut", func() {
		JustBeforeEach(func() {
			LoginFaker(ts, "sign_out@example.com", "hogehoge")
		})
		It("ログアウトできること", func() {
			values := url.Values{}
			res, err := http.PostForm(ts.URL + "/sign_out", values)
			Expect(err).To(BeNil())
			Expect(res.Request.URL.Path).To(Equal("/sign_in"))
		})
	})
})
