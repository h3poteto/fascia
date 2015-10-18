package controllers_test

import (
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	. "../../fascia"
	"../models/db"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("ProjectsController", func() {
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
	JustBeforeEach(func() {
		LoginFaker(ts, "projects@example.com", "hogehoge")
	})
	Describe("Create", func() {
		It("登録できること", func() {
			values := url.Values{}
			values.Add("title", "projectTitle")
			res, err := http.PostForm(ts.URL + "/projects", values)
			Expect(err).To(BeNil())
			contents, status := ParseResponse(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("Id"))
			Expect(contents).To(HaveKey("UserId"))
			Expect(contents).To(HaveKeyWithValue("Title", "projectTitle"))
		})
	})
})
