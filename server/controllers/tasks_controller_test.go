package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/lib/modules/database"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/task"
	"github.com/h3poteto/fascia/server/domains/user"
	userRepo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/usecases/board"
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
		e    *echo.Echo
		rec  *httptest.ResponseRecorder
		p    *project.Project
		user *user.User
		l    *list.List
	)
	email := "task@example.com"
	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		db := database.SharedInstance().Connection
		repo := userRepo.New(db)
		repo.Create(
			email,
			"hogehoge",
			sql.NullString{},
			sql.NullString{},
			sql.NullInt64{},
			sql.NullString{},
			sql.NullString{})
		user, _ = account.FindUserByEmail(email)
		// projectを作っておく
		p, _ = board.CreateProject(user.ID, "projectTitle", "", 0, sql.NullString{})

		// listも作っておく
		l, _ = board.CreateList(p.ID, user.ID, "listTitle", "008ed5", sql.NullInt64{}, false)
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
			_, c = LoginFaker(c, user)
			c = ProjectContext(c, p)
			c = ListContext(c, l)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10))
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
			newTask, _ := board.FindTask(int64(contents.Lists[3].ListTasks[0].ID))
			Expect(newTask.ID).To(BeEquivalentTo(int64(contents.Lists[3].ListTasks[0].ID)))
			Expect(newTask.Title).To(Equal("taskTitle"))
		})
	})

	Describe("Show", func() {
		var newTask *task.Task
		JustBeforeEach(func() {
			newTask, _ = board.CreateTask(l.ID, p.ID, user.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
		})
		It("should receive a task", func() {
			req := httptest.NewRequest(echo.GET, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, user)
			c = ProjectContext(c, p)
			c = ListContext(c, l)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10), strconv.FormatInt(newTask.ID, 10))
			resource := Tasks{}
			err := resource.Show(c)
			Expect(err).To(BeNil())
			var contents views.Task
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(contents.Title).To(Equal(newTask.Title))
		})
	})

	Describe("MoveTask", func() {
		var (
			newTask *task.Task
			newList *list.List
		)
		JustBeforeEach(func() {
			newList, _ = board.CreateList(p.ID, user.ID, "list2", "", sql.NullInt64{}, false)
			newTask, _ = board.CreateTask(l.ID, p.ID, user.ID, sql.NullInt64{}, "taskTitle", "taskDescription", false, sql.NullString{})
		})
		It("should change list the task belongs", func() {
			listID := strconv.FormatInt(newList.ID, 10)
			j := fmt.Sprintf(`{"to_list_id":%s}`, listID)
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id/move_task", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, user)
			c = ProjectContext(c, p)
			c = ListContext(c, l)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10), strconv.FormatInt(newTask.ID, 10))
			resource := Tasks{}
			err := resource.MoveTask(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(rec.Code).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].ListTasks).To(BeEmpty())
			Expect(contents.Lists[4].ListTasks[0].ID).To(Equal(newTask.ID))
		})
	})

	Describe("Update", func() {
		var newTask *task.Task
		JustBeforeEach(func() {
			newTask, _ = board.CreateTask(l.ID, p.ID, user.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
		})
		It("should update a task", func() {
			j := `{"title":"updateTitle","description":"updateDescription"}`
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/tasks/:task_id", strings.NewReader(j))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, user)
			c = ProjectContext(c, p)
			c = ListContext(c, l)
			c = TaskContext(c, newTask)
			c.SetParamNames("project_id", "list_id", "task_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10), strconv.FormatInt(newTask.ID, 10))
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
		var newTask *task.Task
		Context("When a task does not relate issue", func() {
			JustBeforeEach(func() {
				newTask, _ = board.CreateTask(l.ID, p.ID, user.ID, sql.NullInt64{}, "sampleTask", "sampleDescription", false, sql.NullString{})
			})
			It("should delete a task", func() {
				req := httptest.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, user)
				c = ProjectContext(c, p)
				c = ListContext(c, l)
				c = TaskContext(c, newTask)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10), strconv.FormatInt(newTask.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).To(BeNil())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
		Context("When a task relate issue", func() {
			JustBeforeEach(func() {
				newTask, _ = board.CreateTask(l.ID, p.ID, user.ID, sql.NullInt64{Int64: 1, Valid: true}, "sampleTask", "sampleDescription", false, sql.NullString{})
			})
			It("should not delete a task", func() {
				req := httptest.NewRequest(echo.DELETE, "/projects/:project_id/lists/:list_id/tasks/:task_id", nil)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, user)
				c = ProjectContext(c, p)
				c = ListContext(c, l)
				c = TaskContext(c, newTask)
				c.SetParamNames("project_id", "list_id", "task_id")
				c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(l.ID, 10), strconv.FormatInt(newTask.ID, 10))
				resource := Tasks{}
				err := resource.Delete(c)
				Expect(err).NotTo(BeNil())
				Expect(rec.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
