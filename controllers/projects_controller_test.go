package controllers_test

import (
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	"io/ioutil"
	"encoding/json"
	. "../../fascia"
	"../models/db"
	"../models/project"

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
		var (
			res *http.Response
			err error
		)
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "projectTitle")
			res, err = http.PostForm(ts.URL + "/projects", values)
		})
		It("新規登録できること", func() {
			Expect(err).To(BeNil())
			contents, status := ParseResponse(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("Id"))
			Expect(contents).To(HaveKey("UserId"))
			Expect(contents).To(HaveKeyWithValue("Title", "projectTitle"))
		})
		It("DBに登録されていること", func() {
			contents, _ := ParseResponse(res)
			newProject := project.FindProject(int64(contents["Id"].(float64)))
			Expect(newProject.Id).To(BeEquivalentTo(contents["Id"]))
			Expect(newProject.Title).To(Equal("projectTitle"))
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "project1")
			_, _ = http.PostForm(ts.URL + "/projects", values)
			values.Set("title", "project2")
			_, _ = http.PostForm(ts.URL + "/projects", values)
		})
		It("プロジェクト一覧が取得できること", func() {
			res, err := http.Get(ts.URL + "/projects")
			Expect(err).To(BeNil())
			var resp []project.ProjectStruct
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp[0].Title).To(Equal("project1"))
			Expect(resp[1].Title).To(Equal("project2"))
		})
	})
})
