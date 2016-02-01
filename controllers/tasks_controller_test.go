package controllers_test

import (
	. "../../fascia"
	"../controllers"
	seed "../db/seed"
	"../models/db"
	"../models/list"
	"../models/task"
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
		projectId int64
		userId    int64
		listId    int64
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
		table.Exec("truncate table lists;")
		table.Exec("truncate table tasks;")
		table.Exec("truncate table list_options;")
		table.Close()
	})
	JustBeforeEach(func() {
		seed.ListOptions()
		userId = LoginFaker(ts, "tasks@example.com", "hogehoge")
		// projectを作っておく
		values := url.Values{}
		values.Add("title", "projectTitle")
		res, _ := http.PostForm(ts.URL+"/projects", values)
		contents, _ := ParseJson(res)
		parseContents := contents.(map[string]interface{})
		projectId = int64(parseContents["Id"].(float64))

		// listも作っておく
		values.Set("title", "listTitle")
		res, _ = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists", values)
		contents, _ = ParseJson(res)
		parseContents = contents.(map[string]interface{})
		listId = int64(parseContents["Id"].(float64))
	})

	Describe("Create", func() {
		var (
			res *http.Response
			err error
		)
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "taskTitle")
			res, err = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(listId, 10)+"/tasks", values)
		})
		It("can registration", func() {
			Expect(err).To(BeNil())
			contents, status := ParseJson(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("Id"))
		})
		It("should exist in database", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newTask, _ := task.FindTask(listId, int64(parseContents["Id"].(float64)))
			Expect(newTask.Id).To(BeEquivalentTo(parseContents["Id"]))
			Expect(newTask.Title).To(Equal("taskTitle"))
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "task1")
			http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(listId, 10)+"/tasks", values)
			values = url.Values{}
			values.Add("title", "task2")
			http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(listId, 10)+"/tasks", values)
		})
		It("should receive tasks", func() {
			res, err := http.Get(ts.URL + "/projects/" + strconv.FormatInt(projectId, 10) + "/lists/" + strconv.FormatInt(listId, 10) + "/tasks")
			Expect(err).To(BeNil())
			var contents []controllers.TaskJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(contents[0].Title).To(Equal("task1"))
			Expect(contents[1].Title).To(Equal("task2"))
		})
	})

	Describe("MoveTask", func() {
		var (
			newTask *task.TaskStruct
			newList *list.ListStruct
		)
		JustBeforeEach(func() {
			newList = list.NewList(0, projectId, userId, "list2", "", sql.NullInt64{})
			newList.Save(nil, nil)
			newTask = task.NewTask(0, listId, userId, sql.NullInt64{}, "taskTitle", "taskDescription")
			newTask.Save(nil, nil)
		})
		It("should change list the task belongs", func() {
			values := url.Values{}
			values.Add("to_list_id", strconv.FormatInt(newList.Id, 10))
			res, err := http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(listId, 10)+"/tasks/"+strconv.FormatInt(newTask.Id, 10)+"/move_task", values)
			Expect(err).To(BeNil())
			var contents controllers.AllListJSONFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].ListTasks).To(BeEmpty())
			Expect(contents.Lists[4].ListTasks[0].Id).To(Equal(newTask.Id))
		})
	})
})
