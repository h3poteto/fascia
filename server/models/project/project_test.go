package project_test

import (
	"github.com/h3poteto/fascia/config"
	"github.com/h3poteto/fascia/db/seed"
	"github.com/h3poteto/fascia/models/db"
	"github.com/h3poteto/fascia/models/list"
	. "github.com/h3poteto/fascia/models/project"
	"github.com/h3poteto/fascia/models/user"

	"database/sql"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", func() {
	var (
		newProject *ProjectStruct
		uid        int64
		database   *sql.DB
	)

	BeforeEach(func() {
		seed.Seeds()
	})
	AfterEach(func() {
		database.Exec("truncate table users;")
		database.Exec("truncate table projects;")
		database.Exec("truncate table repositories;")
		database.Exec("truncate table lists;")
		database.Exec("truncate table list_options;")
	})

	JustBeforeEach(func() {
		email := "save@example.com"
		password := "hogehoge"
		uid, _ = user.Registration(email, password, password)
		database = db.SharedInstance().Connection
	})

	Describe("Create", func() {
		Context("when did not set repositoryID", func() {
			It("should create new project", func() {
				newProject, err := Create(uid, "new project", "description", 0, sql.NullString{})
				Expect(err).To(BeNil())
				lists, _ := newProject.Lists()
				Expect(len(lists)).To(Equal(3))
				Expect(newProject.NoneList()).NotTo(BeNil())
				Expect(newProject.ShowIssues).To(BeTrue())
				Expect(newProject.ShowPullRequests).To(BeTrue())
			})
			It("should relate user and project", func() {
				newProject, _ = Create(uid, "new project", "description", 0, sql.NullString{})
				rows, _ := database.Query("select id, user_id, title, description from projects where id = ?;", newProject.ID)

				var id int64
				var userID sql.NullInt64
				var title string
				var description string

				for rows.Next() {
					err := rows.Scan(&id, &userID, &title, &description)
					if err != nil {
						panic(err)
					}
				}
				Expect(userID.Valid).To(BeTrue())
				Expect(userID.Int64).To(Equal(uid))
			})
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			newProject, _ = Create(uid, "new project", "description", 0, sql.NullString{})
		})
		It("should set new value", func() {
			err := newProject.Update("newTitle", "newDescription", true, false)
			Expect(err).To(BeNil())
			Expect(newProject.Title).To(Equal("newTitle"))
			Expect(newProject.Description).To(Equal("newDescription"))
			Expect(newProject.RepositoryID.Valid).To(BeFalse())
			Expect(newProject.ShowIssues).To(BeTrue())
			Expect(newProject.ShowPullRequests).To(BeFalse())
		})
	})

	Describe("Lists", func() {
		var (
			newList, noneList *list.ListStruct
			newProject        *ProjectStruct
		)

		BeforeEach(func() {
			newProject, _ = Create(uid, "new project", "description", 0, sql.NullString{})
			newList = list.NewList(0, newProject.ID, newProject.UserID, "list title", "", sql.NullInt64{}, false)
			_ = newList.Save(nil, nil)
			noneList = list.NewList(0, newProject.ID, newProject.UserID, config.Element("init_list").(map[interface{}]interface{})["none"].(string), "", sql.NullInt64{}, false)
			_ = noneList.Save(nil, nil)
		})
		It("should relate project and list", func() {
			lists, err := newProject.Lists()
			Expect(err).To(BeNil())
			Expect(lists).NotTo(BeEmpty())
			Expect(lists[3].ID).To(Equal(newList.ID))
		})
		It("should not take none list", func() {
			lists, err := newProject.Lists()
			Expect(err).To(BeNil())
			Expect(len(lists)).To(Equal(4))
		})

	})

	Describe("NoneList", func() {
		It("should contain only none list", func() {
			newProject, _ := Create(uid, "new project", "description", 0, sql.NullString{})
			noneList, err := newProject.NoneList()
			Expect(err).To(BeNil())
			Expect(noneList.Title.String).To(Equal(config.Element("init_list").(map[interface{}]interface{})["none"].(string)))
		})
	})
})
