package controllers_test

import (
	. "../../fascia"
	"../controllers"
	seed "../db/seed"
	"../models/db"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenazn/goji/web"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ListOptionsController", func() {
	var (
		ts     *httptest.Server
		userID int64
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
		table.Exec("truncate table users:")
		table.Exec("truncate table list_options;")
		table.Close()
	})
	JustBeforeEach(func() {
		userID = LoginFaker(ts, "list_options@example.com", "hogehoge")
		seed.ListOptions()
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			res, err := http.Get(ts.URL + "/list_options")
			Expect(err).To(BeNil())
			var contents []controllers.ListOptionJsonFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(contents[0].Action).To(Equal("close"))
			Expect(contents[1].Action).To(Equal("open"))
		})
	})
})
