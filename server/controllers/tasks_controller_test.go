package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/server/commands/account"
	"github.com/h3poteto/fascia/server/commands/board"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/handlers"
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

var _ = Describe("TasksController", func() {
	var (
		e       *echo.Echo
		rec     *httptest.ResponseRecorder
		project *board.Project
		user    *account.User
		list    *board.List
	)
	email := "task@example.com"
	password := "hogehoge"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		user, _ = handlers.RegistrationUser(email, password, password)
		// projectを作っておく
		pro, _ = handlers.CreateProject(user.UserEntity.ID, "projectTitle", "", 0, sql.NullString{})

		// listも作っておく
		list = handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "listTitle", "008ed5", sql.NullInt64{}, false)
		list.Save()
	})

	Describe("Create", func() {
		var (
			err error
		)
		JustBeforeEach(func() {
			j := `{"title":"taskTitle","description":"desc"}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, pro)
			c = ListContext(c, list)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10))
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
			newTask, _ := handlers.FindTask(list.ListEntity.ID, int64(contents.Lists[3].ListTasks[0].ID))
			Expect(newTask.TaskEntity.ID).To(BeEquivalentTo(int64(contents.Lists[3].ListTasks[0].ID)))
			Expect(newTask.TaskEntity.Title).To(Equal("taskTitle"))
		})
	})

	Describe("Show", func() {
		var newTask *board.Task
		JustBeforeEach(func() {
			newTask = board.NewTask(0, list.ListEntity.ID, project.ProjectEntity.ID, user.UserEntity.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should receive a task", func() {
			req := httptest.NewRequest(echo.GET, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, pro)
			c = ListContext(c, list)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10), strconv.FormatInt(newTask.TaskEntity.ID, 10))
			resource := Tasks{}
			err := resource.Show(c)
			Expect(err).To(BeNil())
			var contents views.Task
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents.Title).To(Equal(newTask.TaskEntity.Title))
		})
	})

	Describe("MoveTask", func() {
		var (
			newTask *board.Task
			newList *board.List
		)
		JustBeforeEach(func() {
			newList = handlers.NewList(0, project.ProjectEntity.ID, user.UserEntity.ID, "list2", "", sql.NullInt64{}, false)
			newList.Save()
			newTask = board.NewTask(0, list.ListEntity.ID, project.ProjectEntity.ID, user.UserEntity.ID, sql.NullInt64{}, "taskTitle", "taskDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should change list the task belongs", func() {
			listID := strconv.FormatInt(newList.ListEntity.ID, 10)
			j := fmt.Sprintf(`{"to_list_id":%s}`, listID)
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, pro)
			c = ListContext(c, list)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10), strconv.FormatInt(newTask.TaskEntity.ID, 10))
			resource := Tasks{}
			err := resource.MoveTask(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].ListTasks).To(BeEmpty())
			Expect(contents.Lists[4].ListTasks[0].ID).To(Equal(newTask.TaskEntity.ID))
		})
	})

	Describe("Update", func() {
		var newTask *board.Task
		JustBeforeEach(func() {
			newTask = board.NewTask(0, list.ListEntity.ID, project.ProjectEntity.ID, user.UserEntity.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
			newTask.Save()
		})
		It("should update a task", func() {
			j := `{"title":"updateTitle","description":"updateDescription"}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, email, password)
			c = ProjectContext(c, pro)
			c = ListContext(c, list)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10), strconv.FormatInt(newTask.TaskEntity.ID, 10))
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
		var newTask *p.Task
		Context("When a task does not relate issue", func() {
			JustBeforeEach(func() {
				newTask = board.NewTask(0, list.ListEntity.ID, project.ProjectEntity.ID, user.UserEntity.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
				newTask.Save()
			})
			It("should delete a task", func() {
				req := httptest.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, pro)
				c = ListContext(c, list)
				c = TaskContext(c, newTask)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10), strconv.FormatInt(newTask.TaskEntity.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
		Context("When a task relate issue", func() {
			JustBeforeEach(func() {
				newTask = board.NewTask(0, list.ListEntity.ID, project.ProjectEntity.ID, user.UserEntity.ID, sql.NullInt64{Int64: 1, Valid: true}, "sampleTask", "sampleDescription", false, sql.NullString{})
				newTask.Save()
			})
			It("should not delete a task", func() {
				req := httptest.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, email, password)
				c = ProjectContext(c, pro)
				c = ListContext(c, list)
				c = TaskContext(c, newTask)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(project.ProjectEntity.ID, 10), strconv.FormatInt(list.ListEntity.ID, 10), strconv.FormatInt(newTask.TaskEntity.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).NotTo(BeNil())
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
