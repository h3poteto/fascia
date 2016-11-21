package controllers_test

import (
	. "github.com/h3poteto/fascia"
	"github.com/h3poteto/fascia/controllers"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/list_option"
	"github.com/h3poteto/fascia/models/project"

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
		userID int64
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
		database.Exec("truncate table list_options;")
		database.Exec("truncate table lists;")
	})
	JustBeforeEach(func() {
		seed.ListOptions()
		userID = LoginFaker(ts, "projects@example.com", "hogehoge")
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
			Expect(contents).To(HaveKey("ID"))
			Expect(contents).To(HaveKey("UserID"))
			Expect(contents).To(HaveKeyWithValue("Title", "projectTitle"))
		})
		It("should exist in database", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newProject, err := project.FindProject(int64(parseContents["ID"].(float64)))
			Expect(err).To(BeNil())
			Expect(newProject.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newProject.Title).To(Equal("projectTitle"))
		})
		It("should have list which have list_option", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newProject, _ := project.FindProject(int64(parseContents["ID"].(float64)))
			lists, err := newProject.Lists()
			Expect(err).To(BeNil())
			Expect(len(lists)).To(Equal(3))
			closeListOption, err := list_option.FindByAction("close")
			Expect(err).To(BeNil())
			Expect(lists[2].ListOptionID.Int64).To(Equal(closeListOption.ID))
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
			var resp []controllers.ProjectJSONFormat
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
			newProject, _ = project.Create(userID, "title", "desc", 0, "", "", sql.NullString{})
		})
		It("should receive project title", func() {
			res, err := http.Get(ts.URL + "/projects/" + strconv.FormatInt(newProject.ID, 10) + "/show")
			Expect(err).To(BeNil())
			var resp controllers.ProjectJSONFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("title"))
		})
	})

	Describe("Update", func() {
		var newProject *project.ProjectStruct
		JustBeforeEach(func() {
			newProject, _ = project.Create(userID, "title", "desc", 0, "", "", sql.NullString{})
		})
		It("should receive new project", func() {
			values := url.Values{}
			values.Add("title", "newTitle")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.ID, 10), values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJSONFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("newTitle"))
		})
	})

	Describe("Settings", func() {
		var newProject *project.ProjectStruct
		JustBeforeEach(func() {
			newProject, _ = project.Create(userID, "title", "desc", 0, "", "", sql.NullString{})
		})
		It("should update show issues", func() {
			values := url.Values{}
			values.Add("show_issues", "false")
			values.Add("show_pull_requests", "true")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.ID, 10)+"/settings", values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJSONFormat
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
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(newProject.ID, 10)+"/settings", values)
			Expect(err).To(BeNil())
			var resp controllers.ProjectJSONFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &resp)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.ShowIssues).To(BeTrue())
			Expect(resp.ShowPullRequests).To(BeFalse())
		})
	})
})
