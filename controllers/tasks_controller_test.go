package controllers_test

import (
	"os"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	. "../../fascia"
	"../models/db"
	"../models/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
)

var _ = Describe("TasksController", func() {
	var (
		ts *httptest.Server
		currentdb string
		projectId int64
		listId int64
	)
	BeforeEach(func() {
		m := web.New()
		Routes(m)
		ts = httptest.NewServer(m)
		testdb := os.Getenv("DB_TEST_NAME")
		currentdb = os.Getenv("DB_NAME")
		os.Setenv("DB_NAME", testdb)
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
		table.Close()
		os.Setenv("DB_NAME", currentdb)
	})
	JustBeforeEach(func() {
		LoginFaker(ts, "tasks@example.com", "hogehoge")
		// projectを作っておく
		values := url.Values{}
		values.Add("title", "projectTitle")
		res, _ := http.PostForm(ts.URL + "/projects", values)
		contents, _ := ParseJson(res)
		parseContents := contents.(map[string]interface{})
		projectId = int64(parseContents["Id"].(float64))

		// listも作っておく
		values.Set("title", "listTitle")
		res, _ = http.PostForm(ts.URL + "/projects/" + strconv.FormatInt(projectId, 10) + "/lists", values)
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
			res, err = http.PostForm(ts.URL + "/projects/" + strconv.FormatInt(projectId, 10) + "/lists/" + strconv.FormatInt(listId, 10) + "/tasks", values)
		})
		It("新規登録できること", func() {
			Expect(err).To(BeNil())
			contents, status := ParseJson(res)
			Expect(status).To(Equal(http.StatusOK))
			Expect(contents).NotTo(BeNil())
			Expect(contents).To(HaveKey("Id"))
		})
		It("DBに登録されていること", func() {
			contents, _ := ParseJson(res)
			parseContents := contents.(map[string]interface{})
			newTask := task.FindTask(listId, int64(parseContents["Id"].(float64)))
			Expect(newTask.Id).To(BeEquivalentTo(parseContents["Id"]))
			Expect(newTask.Title.String).To(Equal("taskTitle"))
		})
	})
})
