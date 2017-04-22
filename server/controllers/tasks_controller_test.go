package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TasksController", func() {
	var (
		e         *echo.Echo
		rec       *httptest.ResponseRecorder
		projectID int64
		userID    int64
		listID    int64
	)
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		userID = LoginFaker("tasks@example.com", "hogehoge")
		// projectを作っておく
		projectService, _ := handlers.CreateProject(userID, "projectTitle", "", 0, sql.NullString{})
		projectID = projectService.ProjectEntity.ProjectModel.ID

		// listも作っておく
		listService := handlers.NewList(0, projectID, userID, "listTitle", "008ed5", sql.NullInt64{}, false)
		listService.Save()
		listID = listService.ListEntity.ListModel.ID
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			f := make(url.Values)
			f.Set("title", "taskTitle")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10))
			resource := Tasks{}
			err = resource.Create(c)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents.Lists[3].ListTasks[0].Title).To(Equal("taskTitle"))
		})
		It("should exist in database", func() {
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			newTask, _ := handlers.FindTask(listID, int64(contents.Lists[3].ListTasks[0].ID))
			Expect(newTask.TaskEntity.TaskModel.ID).To(BeEquivalentTo(int64(contents.Lists[3].ListTasks[0].ID)))
			Expect(newTask.TaskEntity.TaskModel.Title).To(Equal("taskTitle"))
		})
	})

	Describe("Show", func() {
		var newTask *services.Task
		JustBeforeEach(func() {
			newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should receive a task", func() {
			c := e.NewContext(new(http.Request), rec)
			c.SetPath("/projects/:project_id/lists/:list_id/tasks/:task_id")
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10), strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
			resource := Tasks{}
			err := resource.Show(c)
			Expect(err).To(BeNil())
			var contents views.Task
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents.Title).To(Equal(newTask.TaskEntity.TaskModel.Title))
		})
	})

	Describe("MoveTask", func() {
		var (
			newTask *services.Task
			newList *services.List
		)
		JustBeforeEach(func() {
			newList = handlers.NewList(0, projectID, userID, "list2", "", sql.NullInt64{}, false)
			newList.Save()
			newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{}, "taskTitle", "taskDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should change list the task belongs", func() {
			f := make(url.Values)
			f.Set("to_list_id", strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10), strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
			resource := Tasks{}
			err := resource.MoveTask(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].ListTasks).To(BeEmpty())
			Expect(contents.Lists[4].ListTasks[0].ID).To(Equal(newTask.TaskEntity.TaskModel.ID))
		})
	})

	Describe("Update", func() {
		var newTask *services.Task
		JustBeforeEach(func() {
			newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should update a task", func() {
			f := make(url.Values)
			f.Set("title", "updateTitle")
			f.Set("description", "updateDescription")
			req, _ := http.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id", strings.NewReader(f.Encode()))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			c := e.NewContext(req, rec)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10), strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
			resource := Tasks{}
			err := resource.Update(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents.Lists[3].ListTasks[0].Title).To(Equal("updateTitle"))
			Expect(contents.Lists[3].ListTasks[0].Description).To(Equal("updateDescription"))
		})
	})

	Describe("Delete", func() {
		var newTask *services.Task
		Context("When a task does not relate issue", func() {
			JustBeforeEach(func() {
				newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
				newTask.Save()
			})
			It("should delete a task", func() {
				req, _ := http.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10), strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
		Context("When a task relate issue", func() {
			JustBeforeEach(func() {
				newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{Int64: 1, Valid: true}, "sampleTask", "sampleDescription", false, sql.NullString{})
				newTask.Save()
			})
			It("should not delete a task", func() {
				req, _ := http.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(projectID, 10), strconv.FormatInt(listID, 10), strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).NotTo(BeNil())
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
