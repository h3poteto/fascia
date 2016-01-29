package controllers_test

import (
	. "../../fascia"
	"../controllers"
	seed "../db/seed"
	"../models/db"
	"../models/list"
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

var _ = Describe("ListsController", func() {
	var (
		ts        *httptest.Server
		projectId int64
		userId    int64
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
		table.Exec("truncate table list_options;")
		table.Close()
	})
	JustBeforeEach(func() {
		seed.ListOptions()
		userId = LoginFaker(ts, "lists@example.com", "hogehoge")
		// projectを作っておく
		values := url.Values{}
		values.Add("title", "projectTitle")
		res, _ := http.PostForm(ts.URL+"/projects", values)
		contents, _ := ParseJson(res)
		parseContents := contents.(map[string]interface{})
		projectId = int64(parseContents["Id"].(float64))
	})

	Describe("Create", func() {
		var (
			res *http.Response
			err error
		)
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "listTitle")
			res, err = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists", values)
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
			newList := list.FindList(projectId, int64(parseContents["Id"].(float64)))
			Expect(newList.Id).To(BeEquivalentTo(parseContents["Id"]))
			Expect(newList.Title.String).To(Equal("listTitle"))
		})
	})

	Describe("Update", func() {
		var (
			res *http.Response
			err error
		)
		Context("when action is null", func() {
			JustBeforeEach(func() {
				newList := list.NewList(0, projectId, userId, "listTitle", "", sql.NullInt64{})
				newList.Save(nil, nil)
				values := url.Values{}
				values.Add("title", "newListTitle")
				values.Add("action", "null")
				res, err = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(newList.Id, 10), values)
			})
			It("should update", func() {
				Expect(err).To(BeNil())
				contents, status := ParseJson(res)
				Expect(status).To(Equal(http.StatusOK))
				Expect(contents).NotTo(BeNil())
				Expect(contents).To(HaveKey("Id"))
			})
		})
		Context("when action is close", func() {
			JustBeforeEach(func() {
				newList := list.NewList(0, projectId, userId, "listTitle", "", sql.NullInt64{})
				newList.Save(nil, nil)
				values := url.Values{}
				values.Add("title", "newListTitle")
				values.Add("action", "close")
				res, err = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists/"+strconv.FormatInt(newList.Id, 10), values)
			})
			It("should update", func() {
				Expect(err).To(BeNil())
				contents, status := ParseJson(res)
				Expect(status).To(Equal(http.StatusOK))
				Expect(contents).NotTo(BeNil())
				Expect(contents).To(HaveKey("Id"))
			})
		})
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			values := url.Values{}
			values.Add("title", "list1")
			_, _ = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists", values)
			values.Set("title", "list2")
			_, _ = http.PostForm(ts.URL+"/projects/"+strconv.FormatInt(projectId, 10)+"/lists", values)
		})
		It("should receive lists", func() {
			res, err := http.Get(ts.URL + "/projects/" + strconv.FormatInt(projectId, 10) + "/lists")
			Expect(err).To(BeNil())
			var contents controllers.AllListJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			// 初期リストが入るようになったのでそれ以降
			Expect(contents.Lists[3].Title).To(Equal("list1"))
			Expect(contents.Lists[4].Title).To(Equal("list2"))
		})
	})
})
