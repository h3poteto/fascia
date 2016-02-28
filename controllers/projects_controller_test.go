package controllers_test

import (
	. "../../fascia"
	"../controllers"
	seed "../db/seed"
	"../models/db"
	"../models/list_option"
	"../models/project"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("ProjectsController", func() {
	var (
		ts     *httptest.Server
		userId int64
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
		userId = LoginFaker(ts, "projects@example.com", "hogehoge")
	})

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
		It("should have list which have list_option", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newProject := project.FindProject(int64(parseContents["Id"].(float64)))
			lists := newProject.Lists()
			Expect(len(lists)).To(Equal(3))
			closeListOption := list_option.FindByAction("close")
			Expect(lists[2].ListOptionId.Int64).To(Equal(closeListOption.Id))
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
		var newProject *project.ProjectStruct
		JustBeforeEach(func() {
			newProject = project.NewProject(0, userId, "title", "desc", sql.NullInt64{}, true, true)
			newProject.Save()
		})
		It("should receive project title", func() {
			res, err := http.Get(ts.URL + "/projects/" + strconv.FormatInt(newProject.Id, 10) + "/show")
			Expect(err).To(BeNil())
			var resp controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("title"))
		})
	})

	Describe("Update", func() {
		var newProject *project.ProjectStruct
		JustBeforeEach(func() {
			newProject = project.NewProject(0, userId, "title", "desc", sql.NullInt64{}, true, true)
			newProject.Save()
		})
		It("should receive new project", func() {
			values := url.Values{}
			values.Add("title", "newTitle")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.Id, 10), values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("newTitle"))
		})
	})

	Describe("Settings", func() {
		var newProject *project.ProjectStruct
		JustBeforeEach(func() {
			newProject = project.NewProject(0, userId, "title", "description", sql.NullInt64{}, true, true)
			newProject.Save()
		})
		It("should update show issues", func() {
			values := url.Values{}
			values.Add("show_issues", "false")
			values.Add("show_pull_requests", "true")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.Id, 10)+"/settings", values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			fmt.Printf("response: %+v\n", resp)
			Expect(resp.ShowIssues).To(BeFalse())
			Expect(resp.ShowPullRequests).To(BeTrue())
		})
		It("should update show pull requests", func() {
			values := url.Values{}
			values.Add("show_issues", "true")
			values.Add("show_pull_requests", "false")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.Id, 10)+"/settings", values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.ShowIssues).To(BeTrue())
			Expect(resp.ShowPullRequests).To(BeFalse())
		})
	})
})
