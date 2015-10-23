package controllers_test

import (
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	. "../../fascia"
	. "../controllers"
	"../models/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("RegistrationsController", func() {
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
		table.Close()
		os.Setenv("DB_NAME", currentdb)
	})

	Describe("SignUp", func() {
		It("アクセスできること", func() {
			res, err := http.Get(ts.URL + "/sign_up")
			Expect(err).To(BeNil())
			contents, status := ParseResponse(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
		})

	})

	Describe("Registration", func() {
		JustBeforeEach(func() {
			CheckCSRFToken = func(r *http.Request, token string) bool { return true }
		})

		Context("パスワードと確認パスワードが一致しているとき", func() {
			It("登録できること", func() {
				values := url.Values{}
				values.Add("email", "registration@example.com")
				values.Add("password", "hogehoge")
				values.Add("password-confirm", "hogehoge")
				res, err := http.PostForm(ts.URL + "/sign_up", values)
				Expect(err).To(BeNil())
				Expect(res.Request.URL.Path).To(Equal("/sign_in"))
			})
			It("DBに登録されていること", func() {
				values := url.Values{}
				values.Add("email", "registration@example.com")
				values.Add("password", "hogehoge")
				values.Add("password-confirm", "hogehoge")
				http.PostForm(ts.URL + "/sign_up", values)
				mydb := &db.Database{}
				var database db.DB = mydb
				table := database.Init()
				var id int64
				rows, _ := table.Query("select id from users where email = ?;", "registration@example.com")
				for rows.Next() {
					err := rows.Scan(&id)
					if err != nil {
						panic(err.Error())
					}
				}
				Expect(id).NotTo(Equal(0))
			})
		})
		Context("既に登録されているとき", func() {
			JustBeforeEach(func() {
				values := url.Values{}
				values.Add("email", "registration@example.com")
				values.Add("password", "hogehoge")
				values.Add("password-confirm", "hogehoge")
				http.PostForm(ts.URL + "/sign_up", values)
			})
			It("エラーになること", func() {
				values := url.Values{}
				values.Add("email", "registration@example.com")
				values.Add("password", "hogehoge")
				values.Add("password-confirm", "hogehoge")
				res, err := http.PostForm(ts.URL + "/sign_up", values)
				Expect(err).To(BeNil())
				Expect(res.Request.URL.Path).To(Equal("/sign_up"))
			})
		})
	})
})
