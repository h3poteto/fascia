package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectsController", func() {
	var (
		e      *echo.Echo
		rec    *httptest.ResponseRecorder
		userID int64
	)
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		userID = LoginFaker("projects@example.com", "hogehoge")
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			f := make(url.Values)
			f.Set("title", "projectTitle")
			req, _ := http.NewRequest(echo.POST, "/projects", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			resource := Projects{}
			err = resource.Create(c)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("ID"))
			Expect(contents).To(HaveKey("UserID"))
			Expect(contents).To(HaveKeyWithValue("Title", "projectTitle"))
		})
		It("should exist in database", func() {
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			parseContents := contents.(map[string]interface{})
			newProject, err := handlers.FindProject(int64(parseContents["ID"].(float64)))
			Expect(err).To(BeNil())
			Expect(newProject.ProjectEntity.ProjectModel.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newProject.ProjectEntity.ProjectModel.Title).To(Equal("projectTitle"))
		})
		It("should have list which have list_option", func() {
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			parseContents := contents.(map[string]interface{})
			newProject, _ := handlers.FindProject(int64(parseContents["ID"].(float64)))
			lists, err := newProject.ProjectEntity.Lists()
			Expect(err).To(BeNil())
			Expect(len(lists)).To(Equal(3))
			closeListOption, err := services.FindListOptionByAction("close")
			Expect(err).To(BeNil())
			Expect(lists[2].ListModel.ListOptionID.Int64).To(Equal(closeListOption.ListOptionEntity.ListOptionModel.ID))
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			handlers.CreateProject(userID, "project1", "", 0, sql.NullString{})
			handlers.CreateProject(userID, "project2", "", 0, sql.NullString{})
		})
		It("should receive projects", func() {
			c := e.NewContext(new(http.Request), rec)
			c.SetPath("/projects")
			resource := Projects{}
			err := resource.Index(c)
			Expect(err).To(BeNil())
			var resp []views.Project
			json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(resp[0].Title).To(Equal("project1"))
			Expect(resp[1].Title).To(Equal("project2"))
		})
	})

	Describe("Show", func() {
		var newProject *services.Project
		JustBeforeEach(func() {
			newProject, _ = handlers.CreateProject(userID, "title", "desc", 0, sql.NullString{})
		})
		It("should receive project title", func() {
			c := e.NewContext(new(http.Request), rec)
			c.SetPath("/projects/:project_id/show")
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ProjectModel.ID, 10))
			resource := Projects{}
			err := resource.Show(c)
			Expect(err).To(BeNil())
			var resp views.Project
			json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("title"))
		})
	})

	Describe("Update", func() {
		var newProject *services.Project
		JustBeforeEach(func() {
			newProject, _ = handlers.CreateProject(userID, "title", "desc", 0, sql.NullString{})
		})
		It("should receive new project", func() {
			f := make(url.Values)
			f.Set("title", "newTitle")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ProjectModel.ID, 10))
			resource := Projects{}
			err := resource.Update(c)
			Expect(err).To(BeNil())
			var resp views.Project
			json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(resp.Title).To(Equal("newTitle"))
		})
	})

	Describe("Settings", func() {
		var newProject *services.Project
		JustBeforeEach(func() {
			newProject, _ = handlers.CreateProject(userID, "title", "desc", 0, sql.NullString{})
		})
		It("should update show issues", func() {
			f := make(url.Values)
			f.Set("show_issues", "false")
			f.Set("show_pull_requests", "true")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/settings", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ProjectModel.ID, 10))
			resource := Projects{}
			err := resource.Settings(c)
			Expect(err).To(BeNil())
			var resp views.Project
			json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(rec.Code).To(Equal(http.StatusOK))
			fmt.Printf("response: %+v\n", resp)
			Expect(resp.ShowIssues).To(BeFalse())
			Expect(resp.ShowPullRequests).To(BeTrue())
		})
		It("should update show pull requests", func() {
			f := make(url.Values)
			f.Set("show_issues", "true")
			f.Set("show_pull_requests", "false")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/settings", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ProjectModel.ID, 10))
			resource := Projects{}
			err := resource.Settings(c)
			Expect(err).To(BeNil())
			var resp views.Project
			json.Unmarshal(rec.Body.Bytes(), &resp)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(resp.ShowIssues).To(BeTrue())
			Expect(resp.ShowPullRequests).To(BeFalse())
		})
	})
})
