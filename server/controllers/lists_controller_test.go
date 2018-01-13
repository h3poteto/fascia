package controllers_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		project, _ = handlers.CreateProject(user.UserEntity.ID, "projectTitle", "", 0, sql.NullString{})
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			j := `{"title":"listTitle","color":"008ed5"}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10))
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
			newList, err := handlers.FindList(project.ProjectEntity.ID, int64(parseContents["ID"].(float64)))
			Expect(err).To(BeNil())
			Expect(newList.ListEntity.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newList.ListEntity.Title.String).To(Equal("listTitle"))
		})
	})

	Describe("Update", func() {
		var (
			err error
		)
		Context("when action is null", func() {
			JustBeforeEach(func() {
				newList := handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "listTitle", "", sql.NullInt64{}, false)
				newList.Save()
				j := `{"title":"newListTitle","color":"008ed5","option_id":"0"}`
				req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(j))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, project)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(newList.ListEntity.ID, 10))
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
				newList := handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "listTitle", "", sql.NullInt64{}, false)
				newList.Save()
				closeListOption, _ := services.FindListOptionByAction("close")
				optionID := strconv.FormatInt(closeListOption.ListOptionEntity.ID, 10)
				j := fmt.Sprintf(`{"title":"newListTitle","color":"008ed5","option_id":"%s"}`, optionID)
				req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(j))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, project)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(newList.ListEntity.ID, 10))
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
			listService := handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "list1", "008ed5", sql.NullInt64{}, false)
			listService.Save()
			listService = handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "list2", "008ed5", sql.NullInt64{}, false)
			listService.Save()
		})
		It("should receive lists", func() {
			req := httptest.NewRequest(echo.GET, "/projects/:project_id/lists", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10))
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
			newList = handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "listTitle", "", sql.NullInt64{}, false)
			newList.Save()
		})
		It("should hide list", func() {
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/hide", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(newList.ListEntity.ID, 10))
			resource := Lists{}
			err := resource.Hide(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeTrue())
			targetList, err := handlers.FindList(project.ProjectEntity.ID, newList.ListEntity.ID)
			Expect(err).To(BeNil())
			Expect(targetList.ListEntity.IsHidden).To(BeTrue())
		})
	})

	Describe("Display", func() {
		var newList *services.List
		JustBeforeEach(func() {
			newList = handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "listTitle", "", sql.NullInt64{}, false)
			newList.Save()
			newList.Hide()
		})
		It("should display list", func() {
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/display", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, project)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(newList.ListEntity.ID, 10))
			resource := Lists{}
			err := resource.Display(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeFalse())
			targetList, err := handlers.FindList(project.ProjectEntity.ID, newList.ListEntity.ID)
			Expect(err).To(BeNil())
			Expect(targetList.ListEntity.IsHidden).To(BeFalse())
		})
	})
})
