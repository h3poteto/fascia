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
	"strconv"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProjectsController", func() {
	var (
		e    *echo.Echo
		rec  *httptest.ResponseRecorder
		user *services.User
	)
	email := "projects@example.com"
	password := "hogehoge"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		user, _ = handlers.RegistrationUser(email, password, password)
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			j := `{"title":"projectTitle"}`
			req := httptest.NewRequest(echo.POST, "/projects", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
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
			Expect(newProject.ProjectEntity.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newProject.ProjectEntity.Title).To(Equal("projectTitle"))
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
			Expect(lists[2].ListOptionID.Int64).To(Equal(closeListOption.ListOptionEntity.ID))
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			handlers.CreateProject(user.UserEntity.ID, "project1", "", 0, sql.NullString{})
			handlers.CreateProject(user.UserEntity.ID, "project2", "", 0, sql.NullString{})
		})
		It("should receive projects", func() {
			req := httptest.NewRequest(echo.GET, "/projects", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
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
			newProject, _ = handlers.CreateProject(user.UserEntity.ID, "title", "desc", 0, sql.NullString{})
		})
		It("should receive project title", func() {
			req := httptest.NewRequest(echo.GET, "/projects/:project_id/show", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, newProject)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ID, 10))
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
			newProject, _ = handlers.CreateProject(user.UserEntity.ID, "title", "desc", 0, sql.NullString{})
		})
		It("should receive new project", func() {
			j := `{"title":"newTitle"}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, newProject)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ID, 10))
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
			newProject, _ = handlers.CreateProject(user.UserEntity.ID, "title", "desc", 0, sql.NullString{})
		})
		It("should update show issues", func() {
			j := `{"show_issues":false,"show_pull_requests":true}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/settings", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, newProject)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ID, 10))
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
			j := `{"show_issues":true,"show_pull_requests":false}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/settings", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, newProject)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(newProject.ProjectEntity.ID, 10))
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
