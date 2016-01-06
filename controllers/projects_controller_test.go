package controllers_test

import (
	. "../../fascia"
	"../controllers"
	seed "../db/seed"
	"../models/db"
	"../models/project"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("ProjectsController", func() {
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
		mydb := &db.Database{}
		var database db.DB = mydb
		table := database.Init()
		table.Exec("truncate table users;")
		table.Exec("truncate table projects;")
		table.Exec("truncate table list_options;")
		table.Exec("truncate table lists;")
		table.Close()
	})
	JustBeforeEach(func() {
		seed.ListOptions()
		LoginFaker(ts, "projects@example.com", "hogehoge")
	})

	// TODO: InitListが作られること，ListOptionが正しく紐づくことを確認
	Describe("Create", func() {
		var (
			res *http.Response
			err error
		)
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "projectTitle")
			res, err = http.PostForm(ts.URL+"/projects", values)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			contents, status := ParseJson(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("Id"))
			Expect(contents).To(HaveKey("UserId"))
			Expect(contents).To(HaveKeyWithValue("Title", "projectTitle"))
		})
		It("should exist in database", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newProject := project.FindProject(int64(parseContents["Id"].(float64)))
			Expect(newProject.Id).To(BeEquivalentTo(parseContents["Id"]))
			Expect(newProject.Title).To(Equal("projectTitle"))
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "project1")
			_, _ = http.PostForm(ts.URL+"/projects", values)
			values.Set("title", "project2")
			_, _ = http.PostForm(ts.URL+"/projects", values)
		})
		It("should receive projects", func() {
			res, err := http.Get(ts.URL + "/projects")
			Expect(err).To(BeNil())
			var resp []controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp[0].Title).To(Equal("project1"))
			Expect(resp[1].Title).To(Equal("project2"))
		})
	})

	Describe("Show", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "project")
			_, _ = http.PostForm(ts.URL+"/projects", values)
		})
		It("should receive project title", func() {
			res, err := http.Get(ts.URL + "/projects")
			Expect(err).To(BeNil())
			var resp []controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp[0].Title).To(Equal("project"))
		})
	})
})
