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
	"github.com/h3poteto/fascia/lib/modules/database"
	. "github.com/h3poteto/fascia/server/controllers"
	"github.com/h3poteto/fascia/server/domains/list"
	"github.com/h3poteto/fascia/server/domains/project"
	"github.com/h3poteto/fascia/server/domains/user"
	userRepo "github.com/h3poteto/fascia/server/infrastructures/user"
	"github.com/h3poteto/fascia/server/usecases/account"
	"github.com/h3poteto/fascia/server/usecases/board"
	"github.com/h3poteto/fascia/server/views"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListsController", func() {
	var (
		e   *echo.Echo
		rec *httptest.ResponseRecorder
		p   *project.Project
		u   *user.User
	)
	email := "lists@example.com"
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
		var err error
		u, err = account.FindUserByEmail(email)
		if err != nil {
			panic(err)
		}
		p, err = board.CreateProject(u.ID, "projectTitle", "", 0, sql.NullString{})
		if err != nil {
			panic(err)
		}
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
			_, c = LoginFaker(c, u)
			c = ProjectContext(c, p)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10))
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
			newList, err := board.FindList(p.ID, int64(parseContents["ID"].(float64)))
			Expect(err).To(BeNil())
			Expect(newList.ID).To(BeEquivalentTo(parseContents["ID"]))
			Expect(newList.Title.String).To(Equal("listTitle"))
		})
	})

	Describe("Update", func() {
		var (
			err error
		)
		Context("when action is null", func() {
			JustBeforeEach(func() {
				newList, _ := board.CreateList(p.ID, u.ID, "listTitle", "", sql.NullInt64{}, false)
				j := `{"title":"newListTitle","color":"008ed5","option_id":"0"}`
				req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(j))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, u)
				c = ProjectContext(c, p)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(newList.ID, 10))
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
				newList, _ := board.CreateList(p.ID, u.ID, "listTitle", "", sql.NullInt64{}, false)
				closeListOption, _ := board.FindListOptionByAction("close")
				optionID := strconv.FormatInt(closeListOption.ID, 10)
				j := fmt.Sprintf(`{"title":"newListTitle","color":"008ed5","option_id":"%s"}`, optionID)
				req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id", strings.NewReader(j))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				_, c = LoginFaker(c, u)
				c = ProjectContext(c, p)
				c = ListContext(c, newList)
				c.SetParamNames("project_id", "list_id")
				c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(newList.ID, 10))
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
			board.CreateList(p.ID, u.ID, "list1", "008ed5", sql.NullInt64{}, false)
			board.CreateList(p.ID, u.ID, "list2", "008ed5", sql.NullInt64{}, false)
		})
		It("should receive lists", func() {
			req := httptest.NewRequest(echo.GET, "/projects/:project_id/lists", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, u)
			c = ProjectContext(c, p)
			c.SetParamNames("project_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10))
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
		var newList *list.List
		JustBeforeEach(func() {
			newList, _ = board.CreateList(p.ID, u.ID, "listTitle", "", sql.NullInt64{}, false)
		})
		It("should hide list", func() {
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/hide", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, u)
			c = ProjectContext(c, p)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(newList.ID, 10))
			resource := Lists{}
			err := resource.Hide(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeTrue())
			targetList, err := board.FindList(p.ID, newList.ID)
			Expect(err).To(BeNil())
			Expect(targetList.IsHidden).To(BeTrue())
		})
	})

	Describe("Display", func() {
		var newList *list.List
		JustBeforeEach(func() {
			newList, _ = board.CreateList(p.ID, u.ID, "listTitle", "", sql.NullInt64{}, false)
			newList.Hide()
		})
		It("should display list", func() {
			req := httptest.NewRequest(echo.POST, "/projects/:project_id/lists/:list_id/display", nil)
			c := e.NewContext(req, rec)
			_, c = LoginFaker(c, u)
			c = ProjectContext(c, p)
			c = ListContext(c, newList)
			c.SetParamNames("project_id", "list_id")
			c.SetParamValues(strconv.FormatInt(p.ID, 10), strconv.FormatInt(newList.ID, 10))
			resource := Lists{}
			err := resource.Display(c)
			Expect(err).To(BeNil())
			var contents views.AllLists
			json.Unmarshal(rec.Body.Bytes(), &contents)
			Expect(contents.Lists[3].IsHidden).To(BeFalse())
			targetList, err := board.FindList(p.ID, newList.ID)
			Expect(err).To(BeNil())
			Expect(targetList.IsHidden).To(BeFalse())
		})
	})
})
