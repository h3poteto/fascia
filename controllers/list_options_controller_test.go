package controllers_test

import (
	"github.com/h3poteto/fascia/controllers"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	. "github.com/h3poteto/fascia/server"

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
		database := db.SharedInstance().Connection
		database.Exec("truncate table users:")
		database.Exec("truncate table list_options;")
	})
	JustBeforeEach(func() {
		userID = LoginFaker(ts, "list_options@example.com", "hogehoge")
		seed.ListOptions()
	})

	Describe("Index", func() {
		JustBeforeEach(func() {
			res, err := http.Get(ts.URL + "/list_options")
			Expect(err).To(BeNil())
			var contents []controllers.ListOptionJSONFormat
			con, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(con, &contents)
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(contents[0].Action).To(Equal("close"))
			Expect(contents[1].Action).To(Equal("open"))
		})
	})
})
