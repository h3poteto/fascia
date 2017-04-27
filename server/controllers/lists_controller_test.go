package controllers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListsController", func() {
	var (
		e       *echo.Echo
		rec     *httptest.ResponseRecorder
		project *services.Project
		user    *services.User
	)
	email := "lists@example.com"
	password := "hogehoge"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		user, _ = handlers.RegistrationUser(email, password, password)
		project, _ = handlers.CreateProject(user.UserEntity.UserModel.ID, "projectTitle", "", 0, sql.NullString{})
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			f := make(url.Values)
			f.Set("title", "listTitle")
			f.Set("color", "008ed5")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10))
			resource := Lists{}
			err = resource.Create(c)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("ID"))
		})
		It("should exist in database", func() {
			var contents interface{}
			json.Unmarshal(rec.Body.Bytes(), &contents)
			parseContents := contents.(map[string]interface{})
			newList, err := handlers.FindList(project.ProjectEntity.ProjectModel.ID, int64(parseContents["ID"].(float64)))
			Expect(err).To(BeNil())
			Expect(newList.ListEntity.ListModel.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newList.ListEntity.ListModel.Title.String).To(Equal("listTitle"))
		})
	})

	Describe("Update", func() {
		var (
			err error
		)
		Context("when action is null", func() {
			JustBeforeEach(func() {
				newList := handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "listTitle", "", sql.NullInt64{}, false)
				newList.Save()
				f := make(url.Values)
				f.Set("title", "newListTitle")
				f.Set("color", "008ed5")
				f.Set("option_id", "0")
				req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, project)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10), strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
				resource := Lists{}
				err = resource.Update(c)
			})
			It("should update", func() {
				Expect(err).To(BeNil())
				var contents interface{}
				json.Unmarshal(rec.Body.Bytes(), &contents)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(contents).NotTo(BeNil())
				Expect(contents).To(HaveKey("ID"))
			})
		})
		Context("when action is close", func() {
			JustBeforeEach(func() {
				newList := handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "listTitle", "", sql.NullInt64{}, false)
				newList.Save()
				closeListOption, _ := services.FindListOptionByAction("close")
				optionID := strconv.FormatInt(closeListOption.ListOptionEntity.ListOptionModel.ID, 10)
				f := make(url.Values)
				f.Set("title", "newListTitle")
				f.Set("color", "008ed5")
				f.Set("option_id", optionID)
				req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(f.Encode()))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, project)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10), strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
				resource := Lists{}
				err = resource.Update(c)
			})
			It("should update", func() {
				Expect(err).To(BeNil())
				var contents interface{}
				json.Unmarshal(rec.Body.Bytes(), &contents)
				Expect(rec.Code).To(Equal(http.StatusOK))
				Expect(contents).NotTo(BeNil())
				Expect(contents).To(HaveKey("ID"))
			})
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			listService := handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "list1", "008ed5", sql.NullInt64{}, false)
			listService.Save()
			listService = handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "list2", "008ed5", sql.NullInt64{}, false)
			listService.Save()
		})
		It("should receive lists", func() {
			c := e.NewContext(new(http.Request), rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c.SetPath("/projects/:project_id/lists")
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10))
			resource := Lists{}
			err := resource.Index(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].Title).To(Equal("list1"))
			Expect(contents.Lists[4].Title).To(Equal("list2"))
		})
	})

	Describe("Hide", func() {
		var newList *services.List
		JustBeforeEach(func() {
			newList = handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "listTitle", "", sql.NullInt64{}, false)
			newList.Save()
		})
		It("should hide list", func() {
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/hide", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10), strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
			resource := Lists{}
			err := resource.Hide(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeTrue())
			targetList, err := handlers.FindList(project.ProjectEntity.ProjectModel.ID, newList.ListEntity.ListModel.ID)
			Expect(err).To(BeNil())
			Expect(targetList.ListEntity.ListModel.IsHidden).To(BeTrue())
		})
	})

	Describe("Display", func() {
		var newList *services.List
		JustBeforeEach(func() {
			newList = handlers.NewList(0, project.ProjectEntity.ProjectModel.ID, user.UserEntity.UserModel.ID, "listTitle", "", sql.NullInt64{}, false)
			newList.Save()
			newList.Hide()
		})
		It("should display list", func() {
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/display", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ProjectModel.ID, 10), strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
			resource := Lists{}
			err := resource.Display(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeFalse())
			targetList, err := handlers.FindList(project.ProjectEntity.ProjectModel.ID, newList.ListEntity.ListModel.ID)
			Expect(err).To(BeNil())
			Expect(targetList.ListEntity.ListModel.IsHidden).To(BeFalse())
		})
	})
})
