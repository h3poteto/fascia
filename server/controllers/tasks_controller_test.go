package controllers_test

import (
	"github.com/h3poteto/fascia/db/seed"
	. "github.com/h3poteto/fascia/server"
	"github.com/h3poteto/fascia/server/handlers"
	"github.com/h3poteto/fascia/server/services"
	"github.com/h3poteto/fascia/server/views"

	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("TasksController", func() {
	var (
		ts        *httptest.Server
		projectID int64
		userID    int64
		listID    int64
	)
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
	})
	AfterEach(func() {
		ts.Close()
	})
	JustBeforeEach(func() {
		seed.Seeds()
		userID = LoginFaker(ts, "tasks@example.com", "hogehoge")
		// projectを作っておく
		values := url.Values{}
		values.Add("title", "projectTitle")
		res, _ := http.PostForm(ts.URL+"/projects", values)
		contents, _ := ParseJson(res)
		parseContents := contents.(map[string]interface{})
		projectID = int64(parseContents["ID"].(float64))

		// listも作っておく
		values.Set("title", "listTitle")
		values.Set("color", "008ed5")
		res, _ = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists", values)
		contents, _ = ParseJson(res)
		parseContents = contents.(map[string]interface{})
		listID = int64(parseContents["ID"].(float64))
	})

	Describe("Create", func() {
		var (
			res *http.Response
			err error
		)
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "taskTitle")
			res, err = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists/"+strconv.FormatInt(listID, 10)+"/tasks", values)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			var contents views.AllLists
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents.Lists[3].ListTasks[0].Title).To(Equal("taskTitle"))
		})
		It("should exist in database", func() {
			var contents views.AllLists
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
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
			res, err := http.Get(ts.URL + "/projects/" + strconv.FormatInt(projectID, 10) + "/lists/" + strconv.FormatInt(listID, 10) + "/tasks/" + strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10))
			Expect(err).To(BeNil())
			var contents views.Task
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
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
			values := url.Values{}
			values.Add("to_list_id", strconv.FormatInt(newList.ListEntity.ListModel.ID, 10))
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists/"+strconv.FormatInt(listID, 10)+"/tasks/"+strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10)+"/move_task", values)
			Expect(err).To(BeNil())
			var contents views.AllLists
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
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
			values := url.Values{}
			values.Add("title", "updateTitle")
			values.Add("description", "updateDescription")
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists/"+strconv.FormatInt(listID, 10)+"/tasks/"+strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10), values)
			Expect(err).To(BeNil())
			var contents views.AllLists
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
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
				req, err := http.NewRequest(http.MethodDelete, ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists/"+strconv.FormatInt(listID, 10)+"/tasks/"+strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10), nil)
				res, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})
		Context("When a task relate issue", func() {
			JustBeforeEach(func() {
				newTask = services.NewTask(0, listID, projectID, userID, sql.NullInt64{Int64: 1, Valid: true}, "sampleTask", "sampleDescription", false, sql.NullString{})
				newTask.Save()
			})
			It("should not delete a task", func() {
				req, err := http.NewRequest(http.MethodDelete, ts.URL+"/projects/"+strconv.FormatInt(projectID, 10)+"/lists/"+strconv.FormatInt(listID, 10)+"/tasks/"+strconv.FormatInt(newTask.TaskEntity.TaskModel.ID, 10), nil)
				res, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())
				Expect(res.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})
})
